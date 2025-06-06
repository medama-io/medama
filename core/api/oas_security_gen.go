// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

// SecurityHandler is handler for security parameters.
type SecurityHandler interface {
	// HandleCookieAuth handles CookieAuth security.
	// Session token for authentication.
	HandleCookieAuth(ctx context.Context, operationName OperationName, t CookieAuth) (context.Context, error)
}

func findAuthorization(h http.Header, prefix string) (string, bool) {
	v, ok := h["Authorization"]
	if !ok {
		return "", false
	}
	for _, vv := range v {
		scheme, value, ok := strings.Cut(vv, " ")
		if !ok || !strings.EqualFold(scheme, prefix) {
			continue
		}
		return value, true
	}
	return "", false
}

var operationRolesCookieAuth = map[string][]string{
	DeleteUserOperation:             []string{},
	DeleteWebsitesIDOperation:       []string{},
	GetUserOperation:                []string{},
	GetUserUsageOperation:           []string{},
	GetWebsiteIDBrowsersOperation:   []string{},
	GetWebsiteIDCampaignsOperation:  []string{},
	GetWebsiteIDCountryOperation:    []string{},
	GetWebsiteIDDeviceOperation:     []string{},
	GetWebsiteIDLanguageOperation:   []string{},
	GetWebsiteIDMediumsOperation:    []string{},
	GetWebsiteIDOsOperation:         []string{},
	GetWebsiteIDPagesOperation:      []string{},
	GetWebsiteIDPropertiesOperation: []string{},
	GetWebsiteIDReferrersOperation:  []string{},
	GetWebsiteIDSourcesOperation:    []string{},
	GetWebsiteIDSummaryOperation:    []string{},
	GetWebsiteIDTimeOperation:       []string{},
	GetWebsitesOperation:            []string{},
	GetWebsitesIDOperation:          []string{},
	PatchUserOperation:              []string{},
	PatchWebsitesIDOperation:        []string{},
	PostWebsitesOperation:           []string{},
}

func (s *Server) securityCookieAuth(ctx context.Context, operationName OperationName, req *http.Request) (context.Context, bool, error) {
	var t CookieAuth
	const parameterName = "_me_sess"
	var value string
	switch cookie, err := req.Cookie(parameterName); {
	case err == nil: // if NO error
		value = cookie.Value
	case errors.Is(err, http.ErrNoCookie):
		return ctx, false, nil
	default:
		return nil, false, errors.Wrap(err, "get cookie value")
	}
	t.APIKey = value
	t.Roles = operationRolesCookieAuth[operationName]
	rctx, err := s.sec.HandleCookieAuth(ctx, operationName, t)
	if errors.Is(err, ogenerrors.ErrSkipServerSecurity) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return rctx, true, err
}
