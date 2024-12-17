//go:build docker

// This file contains a default configuration for local development.
// It automatically disables the requirement for authenticated calls and sets
// the default logging level to Debug to make debugging easier.
// Furthermore, the default listen address is set to localhost only as the
// authentication has been turned off to minimize the risk of data leaks

package config

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"

	errorHandler "github.com/wisdom-oss/common-go/v3/middleware/gin/error-handler"
	"github.com/wisdom-oss/common-go/v3/middleware/gin/jwt"
	"github.com/wisdom-oss/common-go/v3/middleware/gin/recoverer"
	"github.com/wisdom-oss/common-go/v3/types"
)

const ListenAddress = "0.0.0.0:8000"

// Middlewares configures and outputs the middlewares used in the configuration.
// The contained middlewares are the following:
//   - gin.Logger
func Middlewares() []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc

	middlewares = append(middlewares,
		logger.SetLogger(
			logger.WithDefaultLevel(zerolog.DebugLevel),
			logger.WithUTC(false),
		))

	// read the OIDC authority from the environment
	oidcAuthority, isSet := os.LookupEnv("OIDC_AUTHORITY")
	if !isSet {
		oidcAuthority = "http://backend/api/auth/"
	}

	validator := jwt.Validator{}
	err := validator.Discover(oidcAuthority)
	if err != nil {
		panic(err)
	}

	// TODO: Remove if authorization is required for service on all routes
	validator.EnableOptional()

	middlewares = append(middlewares, requestid.New())
	middlewares = append(middlewares, errorHandler.Handler)
	middlewares = append(middlewares, gin.CustomRecovery(recoverer.RecoveryHandler))
	return middlewares
}

func PrepareRouter() *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(Middlewares()...)

	router.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, types.ServiceError{
			Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.6",
			Status: http.StatusMethodNotAllowed,
			Title:  "Method Not Allowed",
			Detail: "The used HTTP method is not allowed on this route. Please check the documentation and your request",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, types.ServiceError{
			Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.5",
			Status: http.StatusNotFound,
			Title:  "Route Not Found",
			Detail: "The requested path does not exist in this microservice. Please check the documentation and your request",
		})

	})

	return router
}
