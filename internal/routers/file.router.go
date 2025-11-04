package routers

import (
	"file-sharing/internal/handlers"
	"file-sharing/internal/lib/reply"
	"file-sharing/internal/services"

	"github.com/gin-gonic/gin"
)

func (r *Router) RegisterFile(router *gin.Engine) {
	fs := services.NewFile(r.dc)
	fh := handlers.NewFile(fs)

	router.GET("/files", func(c *gin.Context) {
		rp := reply.New(c)
		f, _ := r.dc.File.Query().Where().All(c.Request.Context())
		rp.Success(f).Ok()
	})

	router.POST("/files", fh.CreateOne)
}
