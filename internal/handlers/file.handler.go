package handlers

import (
	"file-sharing/internal/lib/reply"
	"file-sharing/internal/services"

	"github.com/gin-gonic/gin"
)

type File struct {
	s *services.File
}

func NewFile(service *services.File) *File {
	return &File{service}
}

func (h *File) CreateOne(c *gin.Context) {
	s := h.s.AttachGin(c)
	rp := reply.New(c)

	file, err := s.ProcessUpload(true)

	if err != nil {
		return
	}

	rp.Success(file).SetInfo("File successfully uploaded").Created()
}
