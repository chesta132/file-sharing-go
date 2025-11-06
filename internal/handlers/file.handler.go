package handlers

import (
	"file-sharing/config"
	"file-sharing/ent"
	"file-sharing/internal/lib/filelib"
	"file-sharing/internal/lib/reply"
	"file-sharing/internal/services"
	"strconv"

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

func (h *File) GetMany(c *gin.Context) {
	rp := reply.New(c)
	s := h.s.AttachGin(c)
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	f, err := s.GetMany(offset)
	if err != nil {
		rp.Error(reply.CodeBadGateWay, err.Error()).Fail()
		return
	}
	rp.Success(f).Ok()
}

func (h *File) GetOne(c *gin.Context) {
	rp := reply.New(c)
	s := h.s.AttachGin(c)
	token := c.Param("token")

	file, err := s.GetOne(token, true)
	if err != nil {
		return
	}

	rp.Success(file).Ok()
}

func (h *File) Download(c *gin.Context) {
	rp := reply.New(c)
	s := h.s.AttachGin(c)
	token := c.Param("token")
	pw := c.Query("password")

	file, err := s.GetOne(token, true)
	if err != nil {
		return
	}

	if a, c := filelib.IsDownloadable(file, pw); !a {
		message := "Max download reached"
		if c == "PASSWORD" {
			message = "Wrong password"
		}
		rp.Error(reply.CodeBadRequest, message).Fail()
		return
	}

	s.SendToDownload(file)
}

func (h *File) DeleteOne(c *gin.Context) {
	rp := reply.New(c)
	s := h.s.AttachGin(c)
	token := c.Param("token")
	pw := c.Query("password")

	file, err := s.GetOne(token, true)
	if err != nil {
		return
	}
	if !filelib.IsPasswordCorrect(file, pw) {
		rp.Error(reply.CodeBadRequest, "Wrong password").Fail()
		return
	}

	err = filelib.CreateDeleteLog([]*ent.File{file}, config.REQUEST_DELETE_LOG_PATH)
	if err != nil {
		rp.Error(reply.CodeServerError, err.Error()).Fail()
		return
	}

	err = s.DeleteOneFile(file, true)
	if err != nil {
		return
	}

	_, err = s.DeleteOne(token, true)
	if err != nil {
		return
	}

	rp.Success(file).SetInfo(file.FileName + " successfully deleted").Ok()
}
