package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"encoding/pem"
	"flag"
	"fmt"
	configuration "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/je4/filesystem/v3/pkg/vfsrw"
	"github.com/je4/trustutil/v2/pkg/certutil"
	configutil "github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerClientProto "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"github.com/ocfl-archive/dlza-manager-storage-handler/certs"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/internal"
	"github.com/ocfl-archive/dlza-manager-storage-handler/server"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	service2 "github.com/ocfl-archive/dlza-manager-storage-handler/service"
	"github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	archiveerror "github.com/ocfl-archive/error/pkg/error"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	ublogger "gitlab.switch.ch/ub-unibas/go-ublogger/v2"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"go.ub.unibas.ch/cloud/miniresolver/v2/pkg/resolver"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const errorTopic string = "dlza-manager-storage-handler"

var ErrorFactory = archiveerror.NewFactory(errorTopic)

var configFile = flag.String("config", "", "config file in toml format")

var conf *config.Config

const separator = "+"

// disableEndpointPrefix applies the flag that will prevent any
// operation-specific host prefix from being applied
type disableEndpointPrefix struct{}

func (disableEndpointPrefix) ID() string { return "disableEndpointPrefix" }

func (disableEndpointPrefix) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (middleware.InitializeOutput, middleware.Metadata, error) {
	ctx = smithyhttp.SetHostnameImmutable(ctx, true)
	return next.HandleInitialize(ctx, in)
}

func addDisableEndpointPrefix(o *s3.Options) {
	o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
		return stack.Initialize.Add(disableEndpointPrefix{}, middleware.After)
	})
}

type staticResolver struct {
	url string
}

func (s staticResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(s.url)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{URI: *u}, nil
}

func main() {

	flag.Parse()

	var cfgFS fs.FS
	var cfgFile string
	if *configFile != "" {
		cfgFS = os.DirFS(filepath.Dir(*configFile))
		cfgFile = filepath.Base(*configFile)
	} else {
		cfgFS = config.ConfigFS
		cfgFile = "storagehandler.toml"
	}

	conf = &config.Config{
		LocalAddr: "localhost:8443",
		//ResolverTimeout: config.Duration(10 * time.Minute),
		ExternalAddr:            "https://localhost:8443",
		ResolverTimeout:         configutil.Duration(10 * time.Minute),
		ResolverNotFoundTimeout: configutil.Duration(10 * time.Second),
		ServerTLS: &loader.Config{
			Type: "DEV",
		},
		ClientTLS: &loader.Config{
			Type: "DEV",
		},
	}
	if err := config.LoadConfig(cfgFS, cfgFile, conf); err != nil {
		log.Err(err).Msgf("cannot load toml from [%v] %s: %v", cfgFS, cfgFile, err)
	}
	configErrorFactory()

	// create logger instance
	hostname, err := os.Hostname()
	if err != nil {
		log.Err(err).Msgf("cannot get hostname: %v", err)
	}

	var loggerTLSConfig *tls.Config
	var loggerLoader io.Closer
	if conf.Log.Stash.TLS != nil {
		loggerTLSConfig, loggerLoader, err = loader.CreateClientLoader(conf.Log.Stash.TLS, nil)
		if err != nil {
			log.Err(err).Msgf("cannot create client loader: %v", err)
		}
		defer loggerLoader.Close()
	}

	_logger, _logstash, _logfile, err := ublogger.CreateUbMultiLoggerTLS(conf.Log.Level, conf.Log.File,
		ublogger.SetDataset(conf.Log.Stash.Dataset),
		ublogger.SetLogStash(conf.Log.Stash.LogstashHost, conf.Log.Stash.LogstashPort, conf.Log.Stash.Namespace, conf.Log.Stash.LogstashTraceLevel),
		ublogger.SetTLS(conf.Log.Stash.TLS != nil),
		ublogger.SetTLSConfig(loggerTLSConfig),
	)
	if err != nil {
		log.Err(err).Msgf("cannot create logger: %v", err)
	}
	if _logstash != nil {
		defer _logstash.Close()
	}

	if _logfile != nil {
		defer _logfile.Close()
	}

	l2 := _logger.With().Timestamp().Str("host", hostname).Logger() //.Output(output)
	var logger zLogger.ZLogger = &l2

	clientTLSConfig, clientLoader, err := loader.CreateClientLoader(conf.ClientTLS, logger)
	if err != nil {
		logger.Panic().Msgf("cannot create client loader: %v", err)
	}
	defer clientLoader.Close()

	// create TLS Certificate.
	// the certificate MUST contain <package>.<service> as DNS name

	var domainPrefix string
	if conf.Domain != "" {
		domainPrefix = conf.Domain + "."
	}
	certutil.AddDefaultDNSNames(domainPrefix+storagehandlerproto.DispatcherStorageHandlerService_ServiceDesc.ServiceName, domainPrefix+storagehandlerproto.ClerkStorageHandlerService_ServiceDesc.ServiceName)

	serverTLSConfig, serverLoader, err := loader.CreateServerLoader(true, conf.ServerTLS, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer serverLoader.Close()

	logger.Info().Msgf("resolver address is %s", conf.ResolverAddr)
	resolverClient, err := resolver.NewMiniresolverClient(conf.ResolverAddr, conf.GRPCClient, clientTLSConfig, serverTLSConfig, time.Duration(conf.ResolverTimeout), time.Duration(conf.ResolverNotFoundTimeout), logger)
	if err != nil {
		logger.Fatal().Msgf("cannot create resolver client: %v", err)
	}
	defer resolverClient.Close()

	// create grpc server with resolver for name resolution
	grpcServer, err := resolverClient.NewServer(conf.LocalAddr, []string{conf.Domain}, true)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server")
	}
	addr := grpcServer.GetAddr()
	l2 = _logger.With().Timestamp().Str("addr", addr).Logger() //.Output(output)
	logger = &l2

	clientStorageHandlerHandler, err := resolver.NewClient[handlerClientProto.StorageHandlerHandlerServiceClient](
		resolverClient,
		handlerClientProto.NewStorageHandlerHandlerServiceClient,
		handlerClientProto.StorageHandlerHandlerService_ServiceDesc.ServiceName, conf.Domain)
	if err != nil {
		logger.Panic().Msgf("cannot create clientStorageHandlerHandler grpc client: %v", err)
	}

	storageLocations, err := clientStorageHandlerHandler.GetAllStorageLocations(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Panic().Msgf("cannot GetAllStorageLocations: %v", err)
	}

	vfsConfig, err := config.LoadVfsConfig(storageLocations, *conf)
	if err != nil {
		logger.Panic().Msgf("error mapping json for storage location connection field: %v", err)
	}

	vfs, err := vfsrw.NewFS(vfsConfig, &l2)
	if err != nil {
		logger.Panic().Err(err).Msg("cannot create vfs")
	}

	defer func() {
		if err := vfs.Close(); err != nil {
			logger.Error().Err(err).Msg("cannot close vfs")
		}
	}()

	storagehandlerproto.RegisterDispatcherStorageHandlerServiceServer(grpcServer, &server.DispatcherStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, Logger: logger, Vfs: vfs})
	storagehandlerproto.RegisterCheckerStorageHandlerServiceServer(grpcServer, &server.CheckerStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, Logger: logger, Vfs: vfs})

	uploaderService := service2.UploaderService{StorageHandlerHandlerServiceClient: clientStorageHandlerHandler, Logger: &logger, Vfs: vfs, ConfigObj: *conf}
	ctx := context.Background()
	cs := cache.New(60*time.Minute, 60*time.Minute)
	credentialsS3 := credentials.NewStaticCredentialsProvider(conf.S3TempStorage.Key, conf.S3TempStorage.Secret, "")

	// Obtaining the S3 SDK client configuration based on the passed parameters.
	cnf, err := configuration.LoadDefaultConfig(
		ctx,
		configuration.WithCredentialsProvider(credentialsS3),
		configuration.WithRegion("us-east-1"),
		//configuration.WithRequestChecksumCalculation(0),
		//configuration.WithResponseChecksumValidation(0),
	)
	if err != nil {
		panic(err.Error())
	}
	// Create a new S3 SDK client instance.
	svc := s3.NewFromConfig(cnf, func(o *s3.Options) {
		o.EndpointResolverV2 = staticResolver{url: conf.S3TempStorage.ApiUrlValue}
	})

	// Create a new S3 SDK client instance.
	composer := tusd.NewStoreComposer()

	s3Store := s3store.New(conf.S3TempStorage.Bucket, &service.S3Service{Client: svc, AddDisableEndpointPrefix: addDisableEndpointPrefix})

	s3Store.UseIn(composer)
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contain details about the upload
	// itself and the relevant HTTP request.

	go func() {
		for {
			select {
			case event := <-handler.CompleteUploads:
				fmt.Printf("Upload %s finished\n", event.Upload.ID)
				basePathString := conf.S3TempStorage.UploadFolder + "/" + conf.S3TempStorage.Bucket + "/"
				uploadId := strings.Split(event.Upload.ID, separator)[0]
				filename := event.HTTPRequest.Header.Get("FileName")
				objectJson := event.HTTPRequest.Header.Get("ObjectJson")
				collection := event.HTTPRequest.Header.Get("Collection")
				statusId := event.HTTPRequest.Header.Get("StatusId")
				_, err = clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: statusId, Status: "copied to temp storage"})
				if err != nil {
					log.Printf("could not AlterStatus with status id %s:  to copied to temp storage", statusId)
				}
				objectAndFiles, err := uploaderService.CreateObjectAndFiles(uploadId, objectJson, collection, *conf, ErrorFactory)
				if err != nil {
					log.Printf("could not CreateObjectAndFiles for upload id %s: %v", event.Upload.ID, err)
				} else {
					order := &pb.IncomingOrder{CollectionAlias: collection, StatusId: statusId,
						FilePath: basePathString + uploadId, ObjectAndFiles: objectAndFiles, FileName: filename}
					err = uploaderService.CopyFiles(order)
					if err != nil {
						log.Printf("could not copy file with upload id %s:", event.Upload.ID)
					}
				}
			case event := <-handler.CreatedUploads:
				fmt.Printf("Upload %s created\n", event.Upload.ID)
			case event := <-handler.TerminatedUploads:
				fmt.Printf("Upload %s terminated\n", event.Upload.ID)
			case event := <-handler.UploadProgress:
				fmt.Printf("Upload %s progress: %v\n", event.Upload.ID, event.Upload.Offset*100/event.Upload.Size)
			}
		}
	}()

	var cert tls.Certificate
	var addCA = []*x509.Certificate{}
	if conf.TusServer.TLSCert == "" {
		certBytes, err := fs.ReadFile(certs.CertFS, "ub-log.ub.unibas.ch.cert.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal cert %v/%s", certs.CertFS, "ub-log.ub.unibas.ch.cert.pem"))
		}
		keyBytes, err := fs.ReadFile(certs.CertFS, "ub-log.ub.unibas.ch.key.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal key %v/%s", certs.CertFS, "ub-log.ub.unibas.ch.key.pem"))
		}
		if cert, err = tls.X509KeyPair(certBytes, keyBytes); err != nil {
			emperror.Panic(errors.Wrap(err, "cannot create internal cert"))
		}
		rootCABytes, err := fs.ReadFile(certs.CertFS, "ca.cert.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read root ca %v/%s", certs.CertFS, "ca.cert.pem"))
		}
		block, _ := pem.Decode(rootCABytes)
		if block == nil {
			emperror.Panic(errors.Wrapf(err, "cannot decode root ca"))
		}
		rootCA, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			emperror.Panic(errors.Wrap(err, "cannot parse root ca"))
		}
		addCA = append(addCA, rootCA)
	} else {
		if cert, err = tls.LoadX509KeyPair(conf.TusServer.TLSCert, conf.TusServer.TLSKey); err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot load key pair %s - %s", conf.TusServer.TLSCert, conf.TusServer.TLSKey))
		}
		if conf.TusServer.RootCA != nil {
			for _, caName := range conf.TusServer.RootCA {
				rootCABytes, err := os.ReadFile(caName)
				if err != nil {
					emperror.Panic(errors.Wrapf(err, "cannot read root ca %s", caName))
				}
				block, _ := pem.Decode(rootCABytes)
				if block == nil {
					emperror.Panic(errors.Wrapf(err, "cannot decode root ca"))
				}
				rootCA, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					emperror.Panic(errors.Wrap(err, "cannot parse root ca"))
				}
				addCA = append(addCA, rootCA)
			}
		}
	}
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for _, ca := range addCA {
		rootCAs.AddCert(ca)
	}

	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
	}
	corsV := cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:  []string{"http://example.com"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "X-Requested-With", "X-Request-ID", "X-HTTP-Method-Override", "Upload-Length", "Upload-Offset", "Tus-Resumable", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat", "User-Agent", "Referrer", "Origin", "Content-Type", "Content-Length"},
		ExposeHeaders: []string{"Upload-Offset", "Location", "Upload-Length", "Tus-Version", "Tus-Resumable", "Tus-Max-Size", "Tus-Extension", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat", "Location", "Upload-Offset", "Upload-Length"},
	})
	checkAuth := func(c *gin.Context) {
		authKey := c.Request.Header.Get("Authorization")
		collection := c.Request.Header.Get("Collection")
		objectJson := c.Request.Header.Get("ObjectJson")
		statusId := c.Request.Header.Get("StatusId")
		checksum := c.Request.Header.Get("Checksum")
		fileName := c.Request.Header.Get("FileName")

		if authKey == "" || collection == "" || objectJson == "" || statusId == "" || checksum == "" || fileName == "" {
			c.AbortWithStatus(http.StatusExpectationFailed)
			return
		}

		allowed := false
		allowedCache, hasCache := cs.Get(authKey)

		if !hasCache {
			allowedDb, err := uploaderService.TenantHasAccess(authKey, collection)
			if err != nil {
				log.Printf("could not get tenant access status for collection %s:", collection)
			}
			if allowedDb == true {
				cs.Set(authKey, allowedDb, cache.DefaultExpiration)
			}
			allowed = allowedDb
		} else {
			allowed = allowedCache.(bool)
		}
		if !allowed {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Continue down the chain to handler etc
		c.Next()
	}
	router := gin.Default()
	router.Use(corsV, checkAuth)
	router.POST("/files/", gin.WrapF(handler.PostFile))
	router.HEAD("/files/:id", gin.WrapF(handler.HeadFile))
	router.PATCH("/files/:id", gin.WrapF(handler.PatchFile))
	router.GET("/files/:id", gin.WrapF(handler.GetFile))

	serverTus := http.Server{
		Addr:      conf.TusServer.Addr,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	var wg = sync.WaitGroup{}
	if err := http2.ConfigureServer(&serverTus, nil); err != nil {
		emperror.Panic(errors.Wrap(err, "cannot configure http2 server"))
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info().Msgf("Starting tus server: %s", serverTus.Addr)
		if err := serverTus.ListenAndServeTLS("", ""); err != nil {
			emperror.Panic(errors.Wrap(err, "cannot start http2 server"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info().Msgf("Starting grpc server: %s", grpcServer.GetAddr())
		grpcServer.Startup()
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	fmt.Println("press ctrl+c to stop server")
	s := <-done
	fmt.Println("got signal:", s)

	serverTus.Close()
	grpcServer.GracefulStop()
	fmt.Println("Waiting for server shutdown")
	wg.Wait()
}

func configErrorFactory() {
	var archiveErrs []*archiveerror.Error
	if conf.ErrorConfig != "" {
		errorExt := filepath.Ext(conf.ErrorConfig)
		var err error
		switch errorExt {
		case ".toml":
			archiveErrs, err = archiveerror.LoadTOMLFile(conf.ErrorConfig)
		case ".yaml":
			archiveErrs, err = archiveerror.LoadYAMLFile(conf.ErrorConfig)
		default:
			err = errors.Errorf("unknown error config file extension %s", errorExt)
		}
		if err != nil {
			log.Fatal().Err(err).Msgf("cannot load error config file %s", conf.ErrorConfig)
		}
	} else {
		var err error
		const errorsEmbedToml string = "errors.toml"
		archiveErrs, err = archiveerror.LoadTOMLFileFS(internal.InternalFS, errorsEmbedToml)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot load error config file")
		}
	}
	if err := ErrorFactory.RegisterErrors(archiveErrs); err != nil {
		log.Fatal().Err(err).Msg("cannot register errors")
	}
}
