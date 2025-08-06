package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice/internal/db"
)

// HealthcheckHandler is a function that is called if a call to /_/health is
// being made.
var HealthcheckHandler = func(c *gin.Context) {
	if err := db.Pool().Ping(c); err != nil {
		c.String(http.StatusInternalServerError, "service health failed: %s", err)
		return
	}

	c.String(http.StatusOK, "aaaaaaaa")
}
