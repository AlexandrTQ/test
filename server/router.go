package server

import (
	"TransactionServer/database"
	"TransactionServer/model/dto"
	"TransactionServer/model/enum"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"TransactionServer/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type key string

var userKey key = "user"

func InitRouter() http.Handler {
	cors := setCors()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Compress(5))
	router.Use(timeout)
	//free methods
	router.Group(func(router chi.Router) {
		router.Post("/api/transaction/run", handleTransaction())
	})

	return router
}

func handleTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.NewTransactionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("error on unmarshall req, err: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		txId, err := database.GetDatabase().GetTransactionDao().CreateNewTransaction(req.UserId, req.TxType, req.Amount, enum.TxResultFail)
		if err != nil {
			log.Printf("error on try find client by id %d, err: %s", req.UserId, err.Error())
			handleResult(w, r, http.StatusNotFound, nil)
			return
		}

		client, err := service.GetClient(req.UserId)
		if err != nil {
			log.Printf("error on try find client by id %d, err: %s", req.UserId, err.Error())
			handleResult(w, r, http.StatusNotFound, nil)
			return
		}

		newBalance, code, err := client.HandleNewTransaction(req.TxType, req.Amount, req.UserId, txId)
		if err != nil {
			log.Printf("error on tx handling, err: %s", err.Error())
			result := err.Error()
			handleResult(w, r, code, &result)
			return
		}

		handleResult(w, r, http.StatusOK, &newBalance)
	}
}

func timeout(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerDone := make(chan bool)
		go runHandler(handlerDone, h, w, r)
		if err := timeOutOrNot(r, handlerDone); err != nil {
			log.Println(r.Context().Err().Error())
		}
	})
}

func handleResult(w http.ResponseWriter, r *http.Request, status int, result interface{}) {
	if result == nil || reflect.ValueOf(result).IsNil() {
		w.WriteHeader(status)
		return
	}

	if r.Context().Err() != nil {
		return
	}

	setBodyContentType(w, r)
	w.WriteHeader(status)
	if result != nil {
		json.NewEncoder(w).Encode(result)
	}
}

func setCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func setBodyContentType(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil {
		w.Header().Set("Content-Type", http.DetectContentType(bodyBytes))
	}
}

func runHandler(handlerDone chan bool, h http.Handler,
	w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
	handlerDone <- true
}

func timeOutOrNot(r *http.Request, handlerDone chan bool) error {
	select {
	case <-r.Context().Done():
		return errors.New("Timeout")
	case <-handlerDone:
		return nil
	}
}
