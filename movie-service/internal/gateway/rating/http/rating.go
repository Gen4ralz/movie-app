package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gen4ralz/movie-app/movie-service/internal/gateway"
	"github.com/gen4ralz/movie-app/rating-service/pkg/model"
)

// Gateway defines an HTTP gateway for a rating service.
type Gateway struct {
	addr string
}

// New creates a new HTTP gateway for a rating service.
func New(addr string) *Gateway {
	return &Gateway{
		addr: addr,
	}
}

// GetAggregatedRating returns the aggregated rating for a record 
// or ErrNotFound if there are no ratings.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()
	
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if response.StatusCode / 100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", response)
	}

	var v float64
	err = json.NewDecoder(response.Body).Decode(&v)
	if err != nil {
		return 0, nil
	}

	return v, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode / 100 != 2 {
		return fmt.Errorf("non-2xx response: %v", response)
	}
	return nil
}