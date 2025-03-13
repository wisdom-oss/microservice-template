package router

import (
	"github.com/gin-gonic/gin"

	internal "microservice/internal/router"
	v1Routes "microservice/routes/v1"
)

func Configure() (*gin.Engine, error) {
	r, err := internal.GenerateRouter()
	if err != nil {
		return nil, err
	}

	/* Define and import your routes here */
	v1 := r.Group("/v1")
	{
		v1.GET("/", v1Routes.BasicHandler)
	}

	return r, nil
}
