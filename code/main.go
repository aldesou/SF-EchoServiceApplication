package main

import (
	"context"
	"echoservice/logger"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type echoHandler struct {
	logger logger.Logger
}

func (eh *echoHandler) echo(c *gin.Context) {
	correlationId := uuid.New().String()

	if data, err := io.ReadAll(c.Request.Body); err != nil {
		eh.logger.LogError(correlationId, fmt.Sprintf("Unable to read request body: %v", err))
		c.Data(http.StatusInternalServerError, "text/plain", []byte("error: unable to read request body"))
	} else {
		eh.logger.LogInfo(correlationId, fmt.Sprintf("Echo request: %s", data))
		c.Data(http.StatusOK, "text/plain", data)
	}
}

func getEnvironmentVariable(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment variable \"%s\" not defined", key))
	}

	return value
}

func main() {
	var logFolderPath string
	flag.StringVar(&logFolderPath, "logFolderPath", "", "Path to folder where logs will be written")

	var host string
	flag.StringVar(&host, "host", "", "Server host used for binding")

	var port string
	flag.StringVar(&port, "port", "", "Server port used for binding")

	flag.Parse()

	if len(host) == 0 {
		host = getEnvironmentVariable("Fabric_Endpoint_IPOrFQDN_ServiceEndpoint")
	}

	if len(port) == 0 {
		port = getEnvironmentVariable("Fabric_Endpoint_ServiceEndpoint")
	}

	if len(logFolderPath) == 0 {
		logFolderPath = getEnvironmentVariable("Fabric_Folder_App_Log")
	}

	logFilePath := fmt.Sprintf("%s\\service.log", logFolderPath)

	correlationId := uuid.New().String()

	logger := logger.NewFileLogger(logFilePath)
	defer logger.Close()

	handler := echoHandler{logger}

	router := gin.Default()
	router.POST("/echo", handler.echo)

	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{Addr: addr, Handler: router}

	go func() {
		logger.LogInfo(correlationId, fmt.Sprintf("Starting EchoService server on address %s", addr))

		if err := server.ListenAndServe(); err != nil {
			logger.LogInfo(correlationId, fmt.Sprintf("EchoService server exited with error \"%v\"", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	// Wait for ctrl+c
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.LogInfo(correlationId, "Stopping EchoService server")

	if err := server.Shutdown(ctx); err != nil {
		logger.LogInfo(correlationId, fmt.Sprintf("EchoService server shutdown error \"%v\"", err))
	}
}
