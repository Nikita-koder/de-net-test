package main

import (
	"de-net/libs/4_common/smart_context"

	"github.com/go-chi/chi/v5"
)

// TODO: вынести все руты в данную функцию
func initRoutes(sctx smart_context.ISmartContext) (*chi.Mux, error) {
	r := chi.NewRouter()

	return r, nil
}
