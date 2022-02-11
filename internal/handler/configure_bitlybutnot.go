// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/ineverbee/bitlybutnot/internal/handler/operations"
	"github.com/ineverbee/bitlybutnot/internal/shortener"
	"github.com/ineverbee/bitlybutnot/internal/store"
	"github.com/ineverbee/bitlybutnot/internal/store/db"
	"github.com/ineverbee/bitlybutnot/internal/store/mapstore"
	"github.com/ineverbee/bitlybutnot/models"
)

//go:generate swagger generate server --target ../../../bitlybutnot --name Bitlybutnot --spec ../../api/api.yml --server-package internal/handler --principal interface{}

func configureFlags(api *operations.BitlybutnotAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.BitlybutnotAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error
	var storage store.Store
	if storeEnv := os.Getenv("STORE"); storeEnv == "postgres" {
		storage, err = db.NewLinkDB(ctx, api.Logger, &db.Options{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     5432,
			Timeout:  5,
		})
		if err != nil {
			panic(err)
		}
	} else {
		c, err := strconv.Atoi(os.Getenv("STORE_CAP"))
		if err != nil {
			panic(err)
		}
		storage = mapstore.NewLinkMap(c)
	}

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetLinkHandler = operations.GetLinkHandlerFunc(func(params operations.GetLinkParams) middleware.Responder {
		id, err := shortener.DecodeString(params.ShortLink)
		if err != nil {
			message := err.Error()
			return operations.NewGetLinkDefault(http.StatusBadRequest).WithPayload(&models.Error{
				Code:    http.StatusBadRequest,
				Message: &message,
			})
		}
		longLink, err := storage.Get(id)
		if err != nil {
			message := err.Error()
			return operations.NewGetLinkDefault(http.StatusBadRequest).WithPayload(&models.Error{
				Code:    http.StatusBadRequest,
				Message: &message,
			})
		}
		return operations.NewGetLinkOK().WithPayload(&models.Item{
			LongLink:  longLink,
			ShortLink: params.ShortLink,
		})
	})

	api.PostLinkHandler = operations.PostLinkHandlerFunc(func(params operations.PostLinkParams) middleware.Responder {
		id := shortener.HashURL(params.OriginalLink.LongLink)
		shortLink := shortener.UniqueString(id)
		err := storage.Set(id, shortLink, params.OriginalLink.LongLink)
		if err != nil {
			message := err.Error()
			return operations.NewGetLinkDefault(http.StatusBadRequest).WithPayload(&models.Error{
				Code:    http.StatusBadRequest,
				Message: &message,
			})
		}
		return operations.NewPostLinkCreated().WithPayload(&models.Item{
			LongLink:  params.OriginalLink.LongLink,
			ShortLink: shortLink,
		})
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
