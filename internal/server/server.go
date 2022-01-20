package server

import (
	"context"
	"employe/internal/server/api/rest"
	"employe/internal/server/db"
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const restPort = 1020

var (
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
)

func Launch() error {
	log.Info("Starting Server")
	dbStruct, err := db.NewDBStruct()
	if err != nil {
		return err
	}
	ctx, cancel = context.WithCancel(context.Background())
	rtr, err := rest.NewServer(dbStruct)
	if err != nil {
		return err
	}
	go func() {
		server = &http.Server{Addr: fmt.Sprintf(":%d", restPort), Handler: rtr}
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorln("HTTP endpoint returned error: ", err)
			os.Exit(1)
		} else {
			log.Info("Server closed gracefully")
		}

	}()
	log.Info("Server started")
	select {
	case <-ctx.Done():
		log.Info("Server is stopped")
		return nil
	}
}

func ShutDown() {
	fmt.Println("Shutting down the HTTP server...")
	server.Shutdown(ctx)
	cancel()
}
