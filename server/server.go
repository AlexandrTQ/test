package server

import (
	"TransactionServer/database"
	"TransactionServer/database/postgres"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func StartServer() error {
	host := viper.GetString("server_settings.host")
	port := viper.GetString("server_settings.port")
	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{Addr: addr, Handler: InitRouter()}
	serverCtx, serverCancel := context.WithCancel(context.Background())
	withGracefulShutdown(serverCtx, serverCancel, server)
	log.Printf("successful server start on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-serverCtx.Done()
	return nil
}

func StartDb() error {
	dbType := viper.GetString("data_base_settings.type")
	addr := viper.GetString("data_base_settings.addr")
	dbName := viper.GetString("data_base_settings.db_name")
	user := viper.GetString("data_base_settings.user")
	password := viper.GetString("data_base_settings.password")
	switch dbType {
	case "postgres":
		database.InitDataBase(&postgres.PostgresImpl{})
		if err := database.GetDatabase().Start(addr, user, password, dbName); err != nil {
			return err
		}
	default:
		return fmt.Errorf("wrong DB connection type %s" + dbType)
	}
	log.Printf("successful access with database %s on addr: %s", dbName, addr)
	database.DbIsInit = true
	return nil
}

func withGracefulShutdown(serverCtx context.Context, serverCancel context.CancelFunc, server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer serverCancel()
		<-c
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		go shutdownByTimeOut(shutdownCancel, shutdownCtx)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()
}

func shutdownByTimeOut(shutdownCancel context.CancelFunc, shutdownCtx context.Context) {
	defer shutdownCancel()
	<-shutdownCtx.Done()
	if shutdownCtx.Err() == context.DeadlineExceeded {
		log.Fatal("graceful shutdown timed out.. forcing exit.")
	}
}
