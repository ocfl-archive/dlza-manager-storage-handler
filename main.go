package main

import (
	"emperror.dev/errors"
	"flag"
	"github.com/je4/utils/v2/pkg/zLogger"
	handlerClient "github.com/ocfl-archive/dlza-manager-handler/client"
	"github.com/ocfl-archive/dlza-manager-storage-handler/config"
	"github.com/ocfl-archive/dlza-manager-storage-handler/server"
	storageHandlerPb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var configParam = flag.String("config", "", "config file in toml format, no need for filetype for this param")

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
	storageHandlerPb.RegisterUploaderStorageHandlerServiceServer(grpcServerStorageHandler, &server.UploaderStorageHandlerServer{ClientStorageHandlerHandler: clientStorageHandlerHandler, ConfigObj: cfg, Logger: daLogger})
	log.Printf("server started at %v", lisDispatcher.Addr())
	if err := grpcServerStorageHandler.Serve(lisDispatcher); err != nil {
		panic(errors.Wrapf(err, "Failed to serve gRPC server on port: %v", cfg.StorageHandler.Port))
	}

}
