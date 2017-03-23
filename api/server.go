package api

import (
	"net/http"
	"fmt"
)

func NewGateway() *Gateway {
	return &Gateway{}
}

type Gateway struct {
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Serving........")

	
}
