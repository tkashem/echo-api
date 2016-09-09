package service

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"goji.io/pat"
)

type httpHandler interface {
	echo(req *http.Request) (int, string)
}

type httpHandlerImpl struct{}

func (h httpHandlerImpl) echo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := pat.Param(ctx, "name")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!", name)
}
