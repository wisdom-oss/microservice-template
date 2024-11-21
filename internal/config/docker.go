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

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/wisdom-oss/common-go/v2/middleware"
	"github.com/wisdom-oss/common-go/v2/types"
)

const ListenAddress = "0.0.0.0:8000"

func init() {
	// set gin to the production mode (aka release mode)
	gin.SetMode(gin.ReleaseMode)
}

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

	middlewares = append(middlewares, requestid.New())
	middlewares = append(middlewares, middleware.ErrorHandler{}.Gin)
	middlewares = append(middlewares, gin.CustomRecovery(middleware.RecoveryHandler))

	// read the OpenID Connect issuer from the environment
	oidcIssuer := os.Getenv("OIDC_AUTHORITY")
	jwtValidator := middleware.JWTValidator{}
	err := jwtValidator.DiscoverAndConfigure(oidcIssuer)
	if err != nil {
		panic(err)
	}

	middlewares = append(middlewares, jwtValidator.GinHandler)
	return middlewares
}

func PrepareRouter() *gin.Engine {
	router := gin.New()
	router.ForwardedByClientIP = true
	_ = router.SetTrustedProxies([]string{"172.16.31.4"})
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
