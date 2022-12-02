package api

import mux "gitlab.ddos-guard.net/dma/gorilla"

type Handlers interface {
	Register(router *mux.Router)
}
