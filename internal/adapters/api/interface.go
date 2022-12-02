package api

import mux "github.com/gorilla"

type Handlers interface {
	Register(router *mux.Router)
}
