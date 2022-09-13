package register

import (
	"context"
	"net/http"
)

type Registration interface {
	Register(context.Context) ([]string, error)
}

type registration struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) Registration {
	return &registration{httpClient: httpClient}
}
