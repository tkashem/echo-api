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

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	output := echo(name)
	fmt.Fprintf(w, "%s", output)
}
