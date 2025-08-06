//go:build no_db

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthcheckHandler is a function that is called if a call to /_/health is
// being made.
var HealthcheckHandler = func(c *gin.Context) {
	c.Status(http.StatusOK)
}
