package vps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// DormantRequest represents the request payload for a dormant operation.
type DormantRequest struct {
	Dormant bool   `json:"dormant"`
	Product string `json:"product,omitempty"`
}

// DormantResponse represents the response from a dormant operation.
type DormantResponse struct {
	Message string `json:"message"`
}

// MakeDormant decommissions the VPS but retains its storage and IP addresses.
// The transition is a forced power off and discards whatever is in RAM, so a
// running server should be shut down gracefully first with ShutdownWithGrace
// or SetPower using PowerActionShutdown.
// Returns ErrEmptyIdentifier if the identifier is blank.
func (s *Service) MakeDormant(ctx context.Context, identifier string) (DormantResponse, error) {
	if strings.TrimSpace(identifier) == "" {
		return DormantResponse{}, ErrEmptyIdentifier
	}

	return s.setDormant(ctx, identifier, DormantRequest{Dormant: true})
}

// Reactivate returns a dormant VPS to service under the given product code.
// Valid product codes are returned by GetProducts.
// Returns ErrEmptyIdentifier if the identifier is blank.
func (s *Service) Reactivate(ctx context.Context, identifier, product string) (DormantResponse, error) {
	if strings.TrimSpace(identifier) == "" {
		return DormantResponse{}, ErrEmptyIdentifier
	}

	if strings.TrimSpace(product) == "" {
		return DormantResponse{}, errors.New("product is required to re-activate a dormant server")
	}

	return s.setDormant(ctx, identifier, DormantRequest{Dormant: false, Product: product})
}

// setDormant issues the dormant request for the given payload.
func (s *Service) setDormant(ctx context.Context, identifier string, payload DormantRequest) (DormantResponse, error) {
	url := fmt.Sprintf("/vps/servers/%s/dormant", identifier)

	var result DormantResponse
	if _, _, err := s.DoJSON(ctx, http.MethodPut, url, payload, &result, http.StatusOK); err != nil {
		return DormantResponse{}, err
	}

	return result, nil
}
