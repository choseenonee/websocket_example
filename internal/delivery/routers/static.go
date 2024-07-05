package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	"websockets/internal/delivery/docs"
)

func RegisterStatic(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/static", "../static")
	r.LoadHTMLFiles("../static/asyncapi/asyncapi.html", "../static/frontend/index.html")

	r.GET("/asyncapi", func(c *gin.Context) {
		c.HTML(http.StatusOK, "asyncapi.html", gin.H{})
	})
	r.GET("/frontend", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
