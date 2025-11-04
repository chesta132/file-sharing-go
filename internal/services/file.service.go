package services

import (
	"context"
	"file-sharing/config"
	"file-sharing/ent"
	"file-sharing/internal/lib/filelib"
	"file-sharing/internal/lib/reply"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type File struct {
	dc *ent.Client
}

type AttachedGinFile struct {
	*File
	c   *gin.Context
	ctx context.Context
}

// INIT

func NewFile(client *ent.Client) *File {
	return &File{dc: client}
}

func (s *File) AttachGin(c *gin.Context) *AttachedGinFile {
	return &AttachedGinFile{s, c, c.Request.Context()}
}

// SERVICES

func (s *AttachedGinFile) ProcessUpload(allowReply bool) (*ent.File, error) {
	rp := reply.New(s.c)

	// Parse multipart form with max size limit
	if err := s.c.Request.ParseMultipartForm(config.MAX_UPLOAD * config.MB); err != nil {
		if allowReply {
			rp.Error(reply.CodeBadRequest, fmt.Sprintf("Max uploaded file is %vMB", config.MAX_UPLOAD), err.Error()).Fail()
		}
		return nil, err
	}

	// Get file and optional parameters from form
	u, err := s.c.FormFile("file")
	p := s.c.Request.FormValue("password")
	maxDownloads := s.c.Request.FormValue("max-downloads")

	if err != nil {
		if allowReply {
			rp.Error(reply.CodeBadRequest, "Please add 'file' field in form and make sure it's file formatted", err.Error()).Fail()
		}
		return nil, err
	}

	// Ensure upload directories exist
	if err := filelib.CreateDir(); err != nil {
		if allowReply {
			rp.Error(reply.CodeServerError, "Error creating directories for upload file", err.Error()).Fail()
		}
		return nil, err
	}

	// Detect MIME type from header, fallback to "unknown"
	mime := u.Header.Get("Content-Type")
	if mime == "" {
		mime = "unknown"
	}

	// Build database query with file metadata
	q := s.dc.File.Create().SetFileName(u.Filename).SetFileSize(u.Size).SetMime(mime)

	// Set optional password if provided
	if p != "" {
		q.SetPassword(filelib.HashPassword(p))
	}

	// Set max downloads limit if valid number provided
	if md, err := strconv.Atoi(maxDownloads); err == nil {
		q.SetMaxDownloads(md)
	}

	// Save metadata to database first to get generated ID
	file, err := q.Save(s.ctx)
	if err != nil {
		if allowReply {
			rp.Error(reply.CodeServerError, "Error while saving file metadata", err.Error()).Fail()
		}
		return nil, err
	}

	// Generate file path using ID from database
	pathname := filelib.GetPathname(file.FileSize, file.ID, file.FileName)

	// Save physical file to disk
	if err := s.c.SaveUploadedFile(u, pathname); err != nil {
		// Rollback: delete database record if file save fails
		s.dc.File.DeleteOneID(file.ID).ExecX(s.ctx)

		if allowReply {
			rp.Error(reply.CodeServerError, "Error while saving file to disk", err.Error()).Fail()
		}
		return nil, err
	}

	return file, nil
}
