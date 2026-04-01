package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strings"

	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"emperror.dev/errors"
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
	"github.com/ocfl-archive/dlza-manager-storage-handler/store"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/ocfl-archive/dlza-manager/models"
	archiveerror "github.com/ocfl-archive/error/pkg/error"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	ublogger "gitlab.switch.ch/ub-unibas/go-ublogger/v2"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"go.ub.unibas.ch/cloud/miniresolverclient/pkg/miniresolverclient"
	"golang.org/x/exp/maps"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/types/known/emptypb"
)

const errorTopic string = "dlza-manager-storage-handler"

var ErrorFactory = archiveerror.NewFactory(errorTopic)

var configFile = flag.String("config", "", "config file in toml format")

var conf *config.Config

var stores map[string]map[string]s3store.S3Store

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

type storeList map[string]map[string]s3store.S3Store

func (sl storeList) String() string {
	str := fmt.Sprintf("StoreList with %d tenants", len(sl))
	for tenant, stores := range sl {
		str += fmt.Sprintf("\nTenant %s has %d stores", tenant, len(stores))
		for storeName, store := range stores {
			str += fmt.Sprintf("\n\tStore %s: %s", storeName, store)
		}
	}
	return str
}

func (sl storeList) Get(tenant, store string) (s3store.S3Store, error) {
	if tenantStores, ok := sl[tenant]; ok {
		if store, ok := tenantStores[store]; ok {
			return store, nil
		}
	}
	return s3store.S3Store{}, fmt.Errorf("store %s not found for tenant %s", store, tenant)
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
		LocalAddr:               "localhost:8443",
		ExternalAddr:            "https://localhost:8443",
		ResolverTimeout:         configutil.Duration(10 * time.Minute),
		ResolverNotFoundTimeout: configutil.Duration(10 * time.Second),
		ServerTLS:               &loader.Config{Type: "DEV"},
		ClientTLS:               &loader.Config{Type: "DEV"},
	}
	if err := config.LoadConfig(cfgFS, cfgFile, conf); err != nil {
		log.Err(err).Msgf("cannot load toml from [%v] %s: %v", cfgFS, cfgFile, err)
	}
	configErrorFactory()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get hostname")
	}

	logger, closers, err := setupLogger(conf, hostname)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot setup logger")
	}
	for _, c := range closers {
		defer c.Close()
	}

	clientTLSConfig, clientLoader, err := loader.CreateClientLoader(conf.ClientTLS, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create client loader")
	}
	defer clientLoader.Close()

	var domainPrefix string
	if conf.Domain != "" {
		domainPrefix = conf.Domain + "."
	}
	certutil.AddDefaultDNSNames(
		domainPrefix+storagehandlerproto.DispatcherStorageHandlerService_ServiceDesc.ServiceName,
		domainPrefix+storagehandlerproto.ClerkStorageHandlerService_ServiceDesc.ServiceName,
	)

	serverTLSConfig, serverLoader, err := loader.CreateServerLoader(true, conf.ServerTLS, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer serverLoader.Close()

	resolverClient, err := miniresolverclient.NewMiniresolverClientNet(
		conf.ResolverAddr, conf.NetName, conf.GRPCClient, clientTLSConfig, serverTLSConfig,
		time.Duration(conf.ResolverTimeout), time.Duration(conf.ResolverNotFoundTimeout), logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create resolver client")
	}
	defer resolverClient.Close()

	grpcServer, err := resolverClient.NewServerAddresses(conf.LocalAddr, conf.Addresses, []string{conf.Domain}, true)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server")
	}
	logger.Info().Msgf("Server address: %s", grpcServer.GetAddr())

	clientStorageHandlerHandler, err := miniresolverclient.NewClient[handlerClientProto.StorageHandlerHandlerServiceClient](
		resolverClient,
		handlerClientProto.NewStorageHandlerHandlerServiceClient,
		handlerClientProto.StorageHandlerHandlerService_ServiceDesc.ServiceName, conf.Domain)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create gRPC client")
	}

	ctx := context.Background()
	storageLocations, err := clientStorageHandlerHandler.GetAllStorageLocations(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot GetAllStorageLocations")
	}

	vfsConfig, err := config.LoadVfsConfig(storageLocations, *conf)
	if err != nil {
		logger.Fatal().Err(err).Msg("error loading VFS config")
	}

	vfs, err := vfsrw.NewFS(vfsConfig, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create vfs")
	}
	defer vfs.Close()

	storagehandlerproto.RegisterDispatcherStorageHandlerServiceServer(grpcServer, &server.DispatcherStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, Logger: logger, Vfs: vfs})
	storagehandlerproto.RegisterCheckerStorageHandlerServiceServer(grpcServer, &server.CheckerStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, Logger: logger, Vfs: vfs})
	storagehandlerproto.RegisterClerkStorageHandlerServiceServer(grpcServer, &server.ClerkStorageHandlerServer{Vfs: vfs})

	uploaderService := service2.UploaderService{
		StorageHandlerHandlerServiceClient: clientStorageHandlerHandler,
		Logger:                             &logger,
		Vfs:                                vfs,
		ConfigObj:                          *conf,
	}

	composer := tusd.NewStoreComposer()
	stores, err = initS3Stores(ctx, clientStorageHandlerHandler, storageLocations, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot init S3 stores")
	}

	storeFunc := func(tenantAlias string, partitionId string) (s3store.S3Store, bool) {
		s, ok := stores[tenantAlias][partitionId]
		return s, ok
	}
	store.NewRoutingStore(storeFunc).UseIn(composer)

	tusHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			filename := hook.HTTPRequest.Header.Get("FileName")
			collectionAlias := hook.HTTPRequest.Header.Get("Collection")
			partitionId := hook.HTTPRequest.Header.Get("PartitionId")
			tenant, err := clientStorageHandlerHandler.FindTenantByCollectionAlias(context.Background(), &pb.Id{Id: collectionAlias})
			if err != nil {
				return tusd.HTTPResponse{StatusCode: 404, Body: err.Error()}, tusd.FileInfoChanges{}, err
			}
			return tusd.HTTPResponse{StatusCode: 0}, tusd.FileInfoChanges{
				ID:       fmt.Sprintf("%s-%s-%s", tenant.Alias, partitionId, filename),
				MetaData: map[string]string{"dlza": filename},
				Storage:  map[string]string{"Path": filename},
			}, nil
		},
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create tus handler")
	}

	go handleUploadEvents(ctx, tusHandler, clientStorageHandlerHandler, uploaderService, *conf, logger)

	tusTLSConfig, tusClosers, err := getTusTLSConfig(conf)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot get tus TLS config")
	}
	for _, c := range tusClosers {
		defer c.Close()
	}

	authCache := cache.New(time.Hour, time.Hour)
	checkAuth := func(c *gin.Context) {
		header := c.Request.Header
		authKey := header.Get("Authorization")
		collection := header.Get("Collection")
		if authKey == "" || collection == "" || header.Get("ObjectJson") == "" || header.Get("StatusId") == "" || header.Get("Checksum") == "" || header.Get("FileName") == "" {
			c.AbortWithStatus(http.StatusExpectationFailed)
			return
		}

		if allowed, ok := authCache.Get(authKey); ok && allowed.(bool) {
			c.Next()
			return
		}

		allowed, err := uploaderService.TenantHasAccess(authKey, collection)
		if err != nil {
			logger.Error().Msgf("access check failed for collection %s: %v", collection, err)
		}
		if allowed {
			authCache.Set(authKey, true, cache.DefaultExpiration)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:    []string{"Authorization", "X-Requested-With", "X-Request-ID", "X-HTTP-Method-Override", "Upload-Length", "Upload-Offset", "Tus-Resumable", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat", "User-Agent", "Referrer", "Origin", "Content-Type", "Content-Length"},
		ExposeHeaders:   []string{"Upload-Offset", "Location", "Upload-Length", "Tus-Version", "Tus-Resumable", "Tus-Max-Size", "Tus-Extension", "Upload-Metadata", "Upload-Defer-Length", "Upload-Concat"},
	}), checkAuth)

	files := router.Group("/files", func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/files")
	})
	files.POST("/", gin.WrapF(tusHandler.PostFile))
	files.HEAD("/*id", gin.WrapF(tusHandler.HeadFile))
	files.PATCH("/*id", gin.WrapF(tusHandler.PatchFile))
	files.GET("/*id", gin.WrapF(tusHandler.GetFile))

	serverTus := http.Server{
		Addr:      conf.TusServer.Addr,
		Handler:   router,
		TLSConfig: tusTLSConfig,
	}
	if err := http2.ConfigureServer(&serverTus, nil); err != nil {
		logger.Fatal().Err(err).Msg("cannot configure http2 server")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info().Msgf("Starting tus server: %s", serverTus.Addr)
		if err := serverTus.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("tus server failed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info().Msgf("Starting grpc server: %s", grpcServer.GetAddr())
		grpcServer.Startup()
	}()
	// Wait for control-c to stop
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logger.Info().Msg("press ctrl+c to stop server")
	s := <-done
	logger.Info().Msgf("got signal: %v", s)

	serverTus.Close()
	grpcServer.GracefulStop()
	logger.Info().Msg("Waiting for server shutdown")
	wg.Wait()
}

func setupLogger(conf *config.Config, hostname string) (zLogger.ZLogger, []io.Closer, error) {
	var closers []io.Closer
	var loggerTLSConfig *tls.Config
	var loggerLoader io.Closer
	var err error
	if conf.Log.Stash.TLS != nil {
		loggerTLSConfig, loggerLoader, err = loader.CreateClientLoader(conf.Log.Stash.TLS, nil)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot create client loader for logger")
		}
		closers = append(closers, loggerLoader)
	}

	_logger, _logstash, _logfile, err := ublogger.CreateUbMultiLoggerTLS(conf.Log.Level, conf.Log.File,
		ublogger.SetDataset(conf.Log.Stash.Dataset),
		ublogger.SetLogStash(conf.Log.Stash.LogstashHost, conf.Log.Stash.LogstashPort, conf.Log.Stash.Namespace, conf.Log.Stash.LogstashTraceLevel),
		ublogger.SetTLS(conf.Log.Stash.TLS != nil),
		ublogger.SetTLSConfig(loggerTLSConfig),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot create logger")
	}

	if _logstash != nil {
		closers = append(closers, _logstash)
	}
	if _logfile != nil {
		closers = append(closers, _logfile)
	}

	l2 := _logger.With().Timestamp().Str("host", hostname).Logger()
	return &l2, closers, nil
}

func getTusTLSConfig(conf *config.Config) (*tls.Config, []io.Closer, error) {
	var closers []io.Closer
	var cert tls.Certificate
	var addCA []*x509.Certificate
	var err error

	if conf.TusServer.TLSCert == "" {
		certBytes, err := fs.ReadFile(certs.CertFS, "ub-log.ub.unibas.ch.cert.pem")
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot read internal cert %v/%s", certs.CertFS, "ub-log.ub.unibas.ch.cert.pem")
		}
		keyBytes, err := fs.ReadFile(certs.CertFS, "ub-log.ub.unibas.ch.key.pem")
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot read internal key %v/%s", certs.CertFS, "ub-log.ub.unibas.ch.key.pem")
		}
		if cert, err = tls.X509KeyPair(certBytes, keyBytes); err != nil {
			return nil, nil, errors.Wrap(err, "cannot create internal cert")
		}
		rootCABytes, err := fs.ReadFile(certs.CertFS, "ca.cert.pem")
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot read root ca %v/%s", certs.CertFS, "ca.cert.pem")
		}
		block, _ := pem.Decode(rootCABytes)
		if block == nil {
			return nil, nil, errors.New("cannot decode root ca")
		}
		rootCA, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot parse root ca")
		}
		addCA = append(addCA, rootCA)
	} else {
		if cert, err = tls.LoadX509KeyPair(conf.TusServer.TLSCert, conf.TusServer.TLSKey); err != nil {
			return nil, nil, errors.Wrapf(err, "cannot load key pair %s - %s", conf.TusServer.TLSCert, conf.TusServer.TLSKey)
		}
		for _, caName := range conf.TusServer.RootCA {
			rootCABytes, err := os.ReadFile(caName)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "cannot read root ca %s", caName)
			}
			block, _ := pem.Decode(rootCABytes)
			if block == nil {
				return nil, nil, errors.Wrapf(err, "cannot decode root ca %s", caName)
			}
			rootCA, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "cannot parse root ca %s", caName)
			}
			addCA = append(addCA, rootCA)
		}
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for _, ca := range addCA {
		rootCAs.AddCert(ca)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
	}, closers, nil
}

func initS3Stores(ctx context.Context, clientStorageHandlerHandler handlerClientProto.StorageHandlerHandlerServiceClient, storageLocations *pb.StorageLocations, logger zLogger.ZLogger) (storeList, error) {
	tenants, err := clientStorageHandlerHandler.FindAllTenants(ctx, &pb.NoParam{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get tenants")
	}

	sl := make(storeList)
	for _, tenant := range tenants.Tenants {
		storesForPartitions := make(map[string]s3store.S3Store)
		for _, storageLocation := range storageLocations.StorageLocations {
			if !storageLocation.FillFirst || tenant.Id != storageLocation.TenantId {
				continue
			}

			connection := models.Connection{}
			if err = json.Unmarshal([]byte(storageLocation.Connection), &connection); err != nil {
				return nil, errors.Wrapf(err, "error mapping storageLocation json for storageLocation ID: %s", storageLocation.Id)
			}

			vfsS3 := maps.Values(connection.VFS)
			var s3config *vfsrw.S3
			for _, v := range vfsS3 {
				if v.S3 != nil {
					s3config = v.S3
					break
				}
			}
			if s3config == nil {
				continue
			}

			credentialsS3 := credentials.NewStaticCredentialsProvider(string(s3config.AccessKeyID), string(s3config.SecretAccessKey), "")
			cnf, err := configuration.LoadDefaultConfig(ctx,
				configuration.WithCredentialsProvider(credentialsS3),
				configuration.WithRegion("us-east-1"),
			)
			if err != nil {
				return nil, errors.Wrap(err, "cannot load AWS default config")
			}

			svc := s3.NewFromConfig(cnf, func(o *s3.Options) {
				o.EndpointResolverV2 = staticResolver{url: "https://" + string(s3config.Endpoint)}
			})

			partitions, err := clientStorageHandlerHandler.GetStoragePartitionsByStorageLocationId(ctx, &pb.Id{Id: storageLocation.Id})
			if err != nil {
				logger.Error().Msgf("cannot get storage partitions for location %s: %v", storageLocation.Id, err)
				continue
			}

			for _, partition := range partitions.StoragePartitions {
				alias := strings.Split(partition.Alias, "/")
				if len(alias) < 2 {
					logger.Error().Msgf("invalid partition alias: %s", partition.Alias)
					continue
				}
				s3Store := s3store.New(alias[0], &service.S3Service{Client: svc, AddDisableEndpointPrefix: addDisableEndpointPrefix})
				s3Store.ObjectPrefix = alias[1] + "/"
				storesForPartitions[partition.Id] = s3Store
			}
		}
		sl[tenant.Alias] = storesForPartitions
	}
	return sl, nil
}

func handleUploadEvents(ctx context.Context, handler *tusd.Handler, clientStorageHandlerHandler handlerClientProto.StorageHandlerHandlerServiceClient, uploaderService service2.UploaderService, conf config.Config, logger zLogger.ZLogger) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-handler.CompleteUploads:
			logger.Info().Msgf("Upload %s finished", event.Upload.ID)

			header := event.HTTPRequest.Header
			filename := header.Get("FileName")
			objectJson := header.Get("ObjectJson")
			collection := header.Get("Collection")
			statusId := header.Get("StatusId")
			severalObjects := header.Get("SeveralObjects")
			partitionId := header.Get("PartitionId")

			if severalObjects == "0" {
				continue
			}

			objectInstance, err := clientStorageHandlerHandler.GetObjectInstanceByFileNameAndPartitionId(ctx, &pb.ObjectAndFile{StatusId: partitionId, FileName: filename})
			if err != nil {
				logger.Error().Msgf("could not GetObjectInstanceByFileNameAndPartitionId for file %s and partitionId %s: %v", filename, partitionId, err)
				continue
			}

			objectInstance.Status = "new"
			storageLocation, err := clientStorageHandlerHandler.GetStorageLocationByObjectInstanceId(ctx, &pb.Id{Id: objectInstance.Id})
			if err != nil {
				logger.Error().Msgf("could not GetStorageLocationByObjectInstanceId for file %s and partitionId %s: %v", filename, partitionId, err)
				continue
			}

			connection := models.Connection{}
			if err = json.Unmarshal([]byte(storageLocation.Connection), &connection); err != nil {
				logger.Error().Msgf("error mapping storageLocation json for ID %s: %v", storageLocation.Id, err)
				continue
			}

			if _, err = clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: statusId, Status: "copied to temp storage"}); err != nil {
				logger.Error().Msgf("could not AlterStatus with status id %s: %v", statusId, err)
			}

			basePathString := strings.TrimSuffix(objectInstance.Path, filename)
			object := models.Object{}
			if err = json.Unmarshal([]byte(objectJson), &object); err != nil {
				logger.Error().Msgf("cannot unmarshal object: %v", err)
				continue
			}

			objectAndFiles, err := uploaderService.CreateObjectAndFiles(filename, object, collection, basePathString, severalObjects, connection, conf, ErrorFactory)
			if err != nil {
				if _, errStatus := clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: statusId, Status: "error"}); errStatus != nil {
					logger.Error().Msgf("could not AlterStatus to error for %s: %v", statusId, errStatus)
				}
				logger.Error().Msgf("could not CreateObjectAndFiles for file %s: %v", filename, err)
				continue
			}

			objectAndFiles.ObjectInstance = objectInstance
			objectAndFiles.Object.Id = objectInstance.ObjectId
			if object.Head == "v+" {
				objectAndFiles.NewVersion = true
			}

			order := &pb.IncomingOrder{
				CollectionAlias: collection,
				StatusId:        statusId,
				FilePath:        basePathString + filename,
				ObjectAndFiles:  objectAndFiles,
				FileName:        filename,
			}
			if err = uploaderService.StoringFiles(order, partitionId, severalObjects); err != nil {
				if _, errStatus := clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: statusId, Status: "error"}); errStatus != nil {
					logger.Error().Msgf("could not AlterStatus to error for %s: %v", statusId, errStatus)
				}
				logger.Error().Msgf("could not StoringFiles for file %s: %v", filename, err)
			}

		case event := <-handler.CreatedUploads:
			logger.Info().Msgf("Upload %s created", event.Upload.ID)
		case event := <-handler.TerminatedUploads:
			logger.Info().Msgf("Upload %s terminated", event.Upload.ID)
		case event := <-handler.UploadProgress:
			if event.Upload.Size > 0 {
				logger.Info().Msgf("Upload %s progress: %v%%", event.Upload.ID, event.Upload.Offset*100/event.Upload.Size)
			}
		}
	}
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
