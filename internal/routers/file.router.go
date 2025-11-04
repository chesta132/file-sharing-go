package routers

import (
	"file-sharing/internal/handlers"
	"file-sharing/internal/services"

	"github.com/gin-gonic/gin"
)

func (r *Router) RegisterFile(router *gin.Engine) {
	fs := services.NewFile(r.dc)
	fh := handlers.NewFile(fs)

	router.POST("/files", fh.CreateOne)

	router.GET("/files", fh.GetMany)
	router.GET("/files/:token", fh.GetOne)
}
