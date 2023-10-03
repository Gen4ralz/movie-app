package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gen4ralz/movie-app/metadata-service/pkg/model"
	"github.com/gen4ralz/movie-app/movie-service/internal/gateway"
	"github.com/gen4ralz/movie-app/pkg/discovery"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	registry	discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"

	log.Printf("Calling metadata service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	
	if response.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", response)
	}

	var v *model.Metadata
	err = json.NewDecoder(response.Body).Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}