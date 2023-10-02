package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gen4ralz/movie-app/metadata-service/pkg/model"
	"github.com/gen4ralz/movie-app/movie-service/internal/gateway"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	addr	string
}

// New creates a new HTTP gateway for a movie metadata service
func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
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