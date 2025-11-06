package services

import (
	"context"
	"file-sharing/config"
	"file-sharing/ent"
	"file-sharing/ent/file"
	"file-sharing/internal/lib/crypto"
	"file-sharing/internal/lib/filelib"
	"file-sharing/internal/lib/reply"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type File struct {
	dc *ent.Client
}

type AttachedGinFile struct {
	dc  *ent.Client
	c   *gin.Context
	ctx context.Context
}

// INIT

func NewFile(client *ent.Client) *File {
	return &File{dc: client}
}

func (s *File) AttachGin(c *gin.Context) *AttachedGinFile {
	return &AttachedGinFile{s.dc, c, c.Request.Context()}
}

// PRIVATE UTIL

func (s *AttachedGinFile) replyDbError(err error) {
	rp := reply.New(s.c)

	if ent.IsNotFound(err) {
		rp.Error(reply.CodeNotFound, "File not found. This could be happen because file sharing was expired").Fail()
		return
	}
	rp.Error(reply.CodeBadGateWay, err.Error()).Fail()
}

// SERVICES

func (s *AttachedGinFile) GetMany(offset int) ([]*ent.File, error) {
	return s.dc.File.Query().Where().Offset(offset).Limit(config.PAGINATION_LIMIT).All(s.ctx)
}

func (s *AttachedGinFile) GetOne(token string, allowReply bool) (*ent.File, error) {
	f, err := s.dc.File.Query().Where(file.Token(token)).First(s.ctx)

	if allowReply && err != nil {
		s.replyDbError(err)
		return nil, err
	}

	return f, err
}

func (s *AttachedGinFile) DeleteOne(token string, allowReply bool) (int, error) {
	f, err := s.dc.File.Delete().Where(file.Token(token)).Exec(s.ctx)

	if allowReply && err != nil {
		s.replyDbError(err)
		return 0, err
	}

	return f, err
}

func (s *AttachedGinFile) DeleteOneFile(file *ent.File, allowReply bool) error {
	path := filelib.GetPathname(file)
	err := os.Remove(path)
	if err != nil {
		reply.New(s.c).Error(reply.CodeServerError, "Error cannot remove file").Fail()
		return err
	}
	return nil
}

func (s *AttachedGinFile) SendToDownload(file *ent.File) error {
	err := s.dc.File.UpdateOneID(file.ID).AddDownloadCount(1).Exec(s.ctx)
	s.c.FileAttachment(
		filelib.GetPathname(file),
		file.FileName,
	)
	file.DownloadCount++
	return err
}

func (s *AttachedGinFile) ProcessUpload(allowReply bool) (*ent.File, error) {
	rp := reply.New(s.c)

	// Validate max size
	contentLength := s.c.Request.ContentLength
	if contentLength > config.MAX_UPLOAD*config.MB {
		if allowReply {
			rp.Error(
				reply.CodeBadRequest,
				fmt.Sprintf("Max uploaded file is %vMB", config.MAX_UPLOAD),
				fmt.Sprintf("File size: %.2fMB", float64(contentLength)/float64(config.MB)),
			).Fail()
		}
		return nil, fmt.Errorf("file too large")
	}

	// Hard validate max size
	s.c.Request.Body = http.MaxBytesReader(s.c.Writer, s.c.Request.Body, (config.MAX_UPLOAD*config.MB)+(10*config.MB))

	// Get file and optional parameters from form
	u, err := s.c.FormFile("file")
	p := s.c.Request.FormValue("password")
	maxDownloads := s.c.Request.FormValue("max-downloads")

	if err != nil {
		s.c.Request.Body.Close()
		// Validate max size
		if strings.Contains(err.Error(), "http: request body too large") {
			if allowReply {
				rp.Error(
					reply.CodeBadRequest,
					fmt.Sprintf("Max file to upload is %vMB", config.MAX_UPLOAD),
				).Fail()
			}
			return nil, err
		}

		if allowReply {
			rp.Error(reply.CodeBadRequest, "Please add 'file' field in form and make sure it's file formatted", err.Error()).Fail()
		}
		return nil, err
	}
	defer s.c.Request.Body.Close()

	// Validate max size
	if u.Size > config.MAX_UPLOAD*config.MB {
		if allowReply {
			rp.Error(
				reply.CodeBadRequest,
				fmt.Sprintf("Max uploaded file is %vMB", config.MAX_UPLOAD),
				fmt.Sprintf("File size: %.2fMB", float64(u.Size)/float64(config.MB)),
			).Fail()
		}
		return nil, fmt.Errorf("file too large")
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
		q.SetPassword(crypto.HashPassword(p))
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
	pathname := filelib.GetPathname(file)

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
