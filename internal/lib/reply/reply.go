package reply

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// =======================
// == Struct Definitions ==
// =======================
//

// Meta represents metadata information of the API reply.
type Meta struct {
	Status      string `json:"status"`                // Overall reply status (e.g. "SUCCESS" or "ERROR")
	Information string `json:"information,omitempty"` // Action information
}

// ReplyEnvelope wraps the general structure of an API reply.
type ReplyEnvelope struct {
	Meta Meta `json:"meta"` // Contains status metadata
	Data any  `json:"data"` // Holds the actual reply payload
}

// ReplyError defines the structure of an error reply payload.
type ReplyError struct {
	Code    string `json:"code"`              // Short human-readable error code
	Message string `json:"message"`           // Human-readable message describing the error
	Details string `json:"details,omitempty"` // Optional detailed context or debug info
}

// Reply represents a unified HTTP reply writer with utility methods.
type Reply struct {
	c       *gin.Context
	Payload ReplyEnvelope // Payload to reply
}

//
// =======================
// == Constructor & Utils ==
// =======================
//

// New creates a new Reply instance with the JSON content-type header already set.
//
//	rp := reply.New(c)
func New(c *gin.Context) *Reply {
	return &Reply{c: c}
}

//
// =======================
// == Setter Methods ==
// =======================
//

// SetStatus sets the "status" in meta.
//
// PLEASE USE "SUCCESS" OR "ERROR"
func (r *Reply) SetStatus(status string) *Reply {
	r.Payload.Meta.Status = status
	return r
}

// SetData assigns data to the "data" field.
func (r *Reply) SetData(data any) *Reply {
	r.Payload.Data = data
	return r
}

func (r *Reply) SetInfo(information string) *Reply {
	r.Payload.Meta.Information = information
	return r
}

//
// =======================
// == High-level Helpers ==
// =======================
//

// Success marks the reply as successful and attaches data.
func (r *Reply) Success(data any) *Reply {
	r.SetStatus("SUCCESS")
	r.SetData(data)
	return r
}

// Error sets reply status to "ERROR" and attaches an error payload.
func (r *Reply) Error(code, message string, details ...string) *Reply {
	r.SetStatus("ERROR")
	d := ""
	if len(details) > 0 {
		d = details[0]
	}
	r.SetData(ReplyError{code, message, d})
	return r
}

//
// =======================
// == Senders ==
// =======================
//

// Reply writes the full JSON reply to the client.
func (r *Reply) Reply(code int) {
	r.c.JSON(code, r.Payload)
}

// Reply data only (without meta) in payload to the client.
func (r *Reply) ReplyData(code int) {
	r.c.JSON(code, r.Payload.Data)
}

// Ok sends a 200 OK reply.
func (r *Reply) Ok() {
	r.Reply(http.StatusOK)
}

// NoContent sends a 204 No Content reply.
func (r *Reply) NoContent() {
	r.Reply(http.StatusNoContent)
}

// Created sends a 201 Created reply.
func (r *Reply) Created() {
	r.Reply(http.StatusCreated)
}

// Fail sends a failure reply with a specific HTTP status code.
func (r *Reply) Fail(code ...int) {
	c := codeAlias[r.Payload.Data.(ReplyError).Code]
	if len(code) > 0 {
		c = code[0]
	}
	r.Reply(c)
}

//
// =======================
// == Templates ==
// =======================
//

var (
	CodeNotFound    = "NOT_FOUND"
	CodeServerError = "SERVER_ERROR"
	CodeBadRequest  = "CLIENT_ERROR"
	CodeBadGateWay  = "BAD_GATEWAY"
)

var codeAlias = map[string]int{
	"NOT_FOUND":    http.StatusNotFound,
	"SERVER_ERROR": http.StatusInternalServerError,
	"CLIENT_ERROR": http.StatusBadRequest,
	"BAD_GATEWAY":  http.StatusBadGateway,
}
