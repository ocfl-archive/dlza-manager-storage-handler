package main

import (
	"context"
	"crypto/tls"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	configuration "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerClient "github.com/ocfl-archive/dlza-manager-handler/client"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/data/certs"
	"github.com/ocfl-archive/dlza-manager-storage-handler/server"
	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
	service2 "github.com/ocfl-archive/dlza-manager-storage-handler/service"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var configParam = flag.String("config", "", "config file in toml format, no need for filetype for this param")

const separator = "+"

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func main() {

	flag.Parse()
	cfg := config.GetConfig(*configParam)

	clientStorageHandlerHandler, connectionStorageHandlerHandler, err := handlerClient.NewStorageHandlerHandlerClient(cfg.Handler.Host+":"+strconv.Itoa(cfg.Handler.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connectionStorageHandlerHandler.Close()

	// create logger instance
	var out io.Writer = os.Stdout
	if string(cfg.Logging.LogFile) != "" {
		fp, err := os.OpenFile(string(cfg.Logging.LogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("cannot open logfile %s: %v", string(cfg.Logging.LogFile), err)
		}
		defer fp.Close()
		out = fp
	}

	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	_logger := zerolog.New(output).With().Timestamp().Logger()
	_logger.Level(zLogger.LogLevel(string(cfg.Logging.LogLevel)))
	var logger zLogger.ZLogger = &_logger
	daLogger := zLogger.NewZWrapper(logger)

	//Listen Clerk and Dispatcher
	lisDispatcher, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.StorageHandler.Port))
	if err != nil {
		panic(errors.Wrapf(err, "Failed to listen gRPC server"))
	}
	grpcServerStorageHandler := grpc.NewServer()
	storageHandlerPb.RegisterDispatcherStorageHandlerServiceServer(grpcServerStorageHandler, &server.DispatcherStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, Logger: daLogger})
	storageHandlerPb.RegisterClerkStorageHandlerServiceServer(grpcServerStorageHandler, &server.ClerkStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler})
	log.Printf("server started at %v", lisDispatcher.Addr())
	go func() {
		if err := grpcServerStorageHandler.Serve(lisDispatcher); err != nil {
			panic(errors.Wrapf(err, "Failed to serve gRPC server on port: %v", cfg.StorageHandler.Port))
		}
	}()
	uploaderService := service2.UploaderService{StorageHandlerHandlerServiceClient: clientStorageHandlerHandler, Logger: &logger, ConfigObj: cfg}
	ctx := context.Background()
	cs := cache.New(60*time.Minute, 60*time.Minute)
	credentialsS3 := credentials.NewStaticCredentialsProvider("AKIAFEDBDB2704C24D21", "0jmsjtQd0ka66thzFDJn6ESUeiLii4dIHHHgTPU6", "")

	// Obtaining the S3 SDK client configuration based on the passed parameters.
	cnf, err := configuration.LoadDefaultConfig(
		ctx,
		configuration.WithCredentialsProvider(credentialsS3),
		configuration.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err.Error())
	}
	// Create a new S3 SDK client instance.
	s3Client := s3.NewFromConfig(cnf, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://vip-ecs-ub.storage.p.unibas.ch")
		o.EndpointResolverV2 = &resolverV2{}
	})

	// Create a new S3 SDK client instance.
	composer := tusd.NewStoreComposer()

	s3Store := s3store.New("ubbasel-test", &service.S3Service{Client: s3Client})

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

				basePathString := "vfs:/temp_switch_ch" + "/" + "ubbasel-test" + "/"
				uploadId := strings.Split(event.Upload.ID, separator)[0]
				filename := event.HTTPRequest.Header.Get("FileName")
				objectJson := event.HTTPRequest.Header.Get("ObjectJson")
				collection := event.HTTPRequest.Header.Get("Collection")
				statusId := event.HTTPRequest.Header.Get("StatusId")
				_, err = clientStorageHandlerHandler.AlterStatus(ctx, &pb.StatusObject{Id: statusId, Status: "copied to temp storage"})
				if err != nil {
					log.Printf("could not AlterStatus with status id %s:  to copied to temp storage", statusId)
				}
				objectAndFiles, err := uploaderService.CreateObjectAndFiles(uploadId, objectJson, collection, cfg)
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
	if cfg.ServerConfig.TLSCert == "" {
		certBytes, err := fs.ReadFile(certs.CertFS, "localhost.cert.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal cert %v/%s", certs.CertFS, "localhost.cert.pem"))
		}
		keyBytes, err := fs.ReadFile(certs.CertFS, "localhost.key.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal key %v/%s", certs.CertFS, "localhost.key.pem"))
		}
		if cert, err = tls.X509KeyPair(certBytes, keyBytes); err != nil {
			emperror.Panic(errors.Wrap(err, "cannot create internal cert"))
		}
	} else {
		if cert, err = tls.LoadX509KeyPair(cfg.ServerConfig.TLSCert, cfg.ServerConfig.TLSKey); err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot load key pair %s - %s", cfg.ServerConfig.TLSCert, cfg.ServerConfig.TLSKey))
		}
	}

	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
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

	server := http.Server{
		Addr:      "localhost:8085",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	if err := http2.ConfigureServer(&server, nil); err != nil {
		emperror.Panic(errors.Wrap(err, "cannot configure http2 server"))
	}
	if err := server.ListenAndServeTLS("", ""); err != nil {
		emperror.Panic(errors.Wrap(err, "cannot start http2 server"))
	}
	defer server.Close()
}
