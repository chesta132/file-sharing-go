package routers

import "file-sharing/ent"

type Router struct {
	dc *ent.Client
}

func New(dbClient *ent.Client) *Router {
	return &Router{dbClient}
}
