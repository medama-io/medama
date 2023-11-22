// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
)

func recordError(string, error) {}

// handleDeleteUserRequest handles delete-user operation.
//
// Delete a user account.
//
// DELETE /user
func (s *Server) handleDeleteUserRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "DeleteUser",
			ID:   "delete-user",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "DeleteUser", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeDeleteUserParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response DeleteUserRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "DeleteUser",
			OperationSummary: "Delete User.",
			OperationID:      "delete-user",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = DeleteUserParams
			Response = DeleteUserRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackDeleteUserParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.DeleteUser(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.DeleteUser(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeDeleteUserResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleDeleteWebsitesIDRequest handles delete-websites-id operation.
//
// Delete a website.
//
// DELETE /websites/{hostname}
func (s *Server) handleDeleteWebsitesIDRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "DeleteWebsitesID",
			ID:   "delete-websites-id",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "DeleteWebsitesID", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeDeleteWebsitesIDParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response DeleteWebsitesIDRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "DeleteWebsitesID",
			OperationSummary: "Delete Website.",
			OperationID:      "delete-websites-id",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
				{
					Name: "hostname",
					In:   "path",
				}: params.Hostname,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = DeleteWebsitesIDParams
			Response = DeleteWebsitesIDRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackDeleteWebsitesIDParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.DeleteWebsitesID(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.DeleteWebsitesID(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeDeleteWebsitesIDResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetEventPingRequest handles get-event-ping operation.
//
// This is a ping endpoint to determine if the user is unique or not.
//
// GET /event/ping
func (s *Server) handleGetEventPingRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetEventPing",
			ID:   "get-event-ping",
		}
	)
	params, err := decodeGetEventPingParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetEventPingRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetEventPing",
			OperationSummary: "Unique User Check.",
			OperationID:      "get-event-ping",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "If-Modified-Since",
					In:   "header",
				}: params.IfModifiedSince,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetEventPingParams
			Response = GetEventPingRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetEventPingParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetEventPing(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetEventPing(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetEventPingResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetUserRequest handles get-user operation.
//
// Retrieve the information of the user with the matching user ID.
//
// GET /user
func (s *Server) handleGetUserRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetUser",
			ID:   "get-user",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "GetUser", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeGetUserParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetUserRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetUser",
			OperationSummary: "Get User Info.",
			OperationID:      "get-user",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetUserParams
			Response = GetUserRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetUserParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetUser(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetUser(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetUserResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetWebsiteIDSummaryRequest handles get-website-id-summary operation.
//
// Get a summary of the website's stats.
//
// GET /website/{hostname}/summary
func (s *Server) handleGetWebsiteIDSummaryRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetWebsiteIDSummary",
			ID:   "get-website-id-summary",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "GetWebsiteIDSummary", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeGetWebsiteIDSummaryParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetWebsiteIDSummaryRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetWebsiteIDSummary",
			OperationSummary: "Get Stat Summary.",
			OperationID:      "get-website-id-summary",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
				{
					Name: "start",
					In:   "query",
				}: params.Start,
				{
					Name: "end",
					In:   "query",
				}: params.End,
				{
					Name: "hostname",
					In:   "path",
				}: params.Hostname,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetWebsiteIDSummaryParams
			Response = GetWebsiteIDSummaryRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetWebsiteIDSummaryParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetWebsiteIDSummary(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetWebsiteIDSummary(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetWebsiteIDSummaryResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetWebsitesRequest handles get-websites operation.
//
// Get a list of all websites from the user.
//
// GET /websites
func (s *Server) handleGetWebsitesRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetWebsites",
			ID:   "get-websites",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "GetWebsites", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeGetWebsitesParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetWebsitesRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetWebsites",
			OperationSummary: "List Websites.",
			OperationID:      "get-websites",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetWebsitesParams
			Response = GetWebsitesRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetWebsitesParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetWebsites(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetWebsites(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetWebsitesResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetWebsitesIDRequest handles get-websites-id operation.
//
// Get website details for an individual website.
//
// GET /websites/{hostname}
func (s *Server) handleGetWebsitesIDRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetWebsitesID",
			ID:   "get-websites-id",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "GetWebsitesID", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeGetWebsitesIDParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetWebsitesIDRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetWebsitesID",
			OperationSummary: "Get Website.",
			OperationID:      "get-websites-id",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
				{
					Name: "hostname",
					In:   "path",
				}: params.Hostname,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetWebsitesIDParams
			Response = GetWebsitesIDRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetWebsitesIDParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetWebsitesID(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetWebsitesID(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetWebsitesIDResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetWebsitesIDActiveRequest handles get-websites-id-active operation.
//
// Return the number of active users who triggered a pageview in the past 5 minutes.
//
// GET /websites/{hostname}/active
func (s *Server) handleGetWebsitesIDActiveRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetWebsitesIDActive",
			ID:   "get-websites-id-active",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "GetWebsitesIDActive", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeGetWebsitesIDActiveParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetWebsitesIDActiveRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetWebsitesIDActive",
			OperationSummary: "Get Active Users.",
			OperationID:      "get-websites-id-active",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
				{
					Name: "hostname",
					In:   "path",
				}: params.Hostname,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetWebsitesIDActiveParams
			Response = GetWebsitesIDActiveRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetWebsitesIDActiveParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetWebsitesIDActive(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetWebsitesIDActive(ctx, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetWebsitesIDActiveResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePatchUserRequest handles patch-user operation.
//
// Update a user account's details.
//
// PATCH /user
func (s *Server) handlePatchUserRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PatchUser",
			ID:   "patch-user",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "PatchUser", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodePatchUserParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodePatchUserRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PatchUserRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PatchUser",
			OperationSummary: "Update User Info.",
			OperationID:      "patch-user",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
			},
			Raw: r,
		}

		type (
			Request  = *UserPatch
			Params   = PatchUserParams
			Response = PatchUserRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackPatchUserParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PatchUser(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.PatchUser(ctx, request, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePatchUserResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePatchWebsitesIDRequest handles patch-websites-id operation.
//
// Update a website's information.
//
// PATCH /websites/{hostname}
func (s *Server) handlePatchWebsitesIDRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PatchWebsitesID",
			ID:   "patch-websites-id",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "PatchWebsitesID", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodePatchWebsitesIDParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodePatchWebsitesIDRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PatchWebsitesIDRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PatchWebsitesID",
			OperationSummary: "Update Website.",
			OperationID:      "patch-websites-id",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "_me_sess",
					In:   "cookie",
				}: params.MeSess,
				{
					Name: "hostname",
					In:   "path",
				}: params.Hostname,
			},
			Raw: r,
		}

		type (
			Request  = *WebsitePatch
			Params   = PatchWebsitesIDParams
			Response = PatchWebsitesIDRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackPatchWebsitesIDParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PatchWebsitesID(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.PatchWebsitesID(ctx, request, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePatchWebsitesIDResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostAuthLoginRequest handles post-auth-login operation.
//
// Login to the service and retrieve a session token for authentication.
//
// POST /auth/login
func (s *Server) handlePostAuthLoginRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PostAuthLogin",
			ID:   "post-auth-login",
		}
	)
	request, close, err := s.decodePostAuthLoginRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PostAuthLoginRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PostAuthLogin",
			OperationSummary: "Session Token Authentication.",
			OperationID:      "post-auth-login",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *AuthLogin
			Params   = struct{}
			Response = PostAuthLoginRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostAuthLogin(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostAuthLogin(ctx, request)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostAuthLoginResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostEventHitRequest handles post-event-hit operation.
//
// Send a hit event to register a user view.
//
// POST /event/hit
func (s *Server) handlePostEventHitRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PostEventHit",
			ID:   "post-event-hit",
		}
	)
	params, err := decodePostEventHitParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodePostEventHitRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PostEventHitRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PostEventHit",
			OperationSummary: "Send Hit Event.",
			OperationID:      "post-event-hit",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "User-Agent",
					In:   "header",
				}: params.UserAgent,
				{
					Name: "Accept-Language",
					In:   "header",
				}: params.AcceptLanguage,
			},
			Raw: r,
		}

		type (
			Request  = *EventHit
			Params   = PostEventHitParams
			Response = PostEventHitRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackPostEventHitParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostEventHit(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostEventHit(ctx, request, params)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostEventHitResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostUserRequest handles post-user operation.
//
// Create a new user.
//
// POST /user
func (s *Server) handlePostUserRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PostUser",
			ID:   "post-user",
		}
	)
	request, close, err := s.decodePostUserRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PostUserRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PostUser",
			OperationSummary: "Create New User.",
			OperationID:      "post-user",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *UserCreate
			Params   = struct{}
			Response = PostUserRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostUser(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostUser(ctx, request)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostUserResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostWebsitesRequest handles post-websites operation.
//
// Add a new website.
//
// POST /websites
func (s *Server) handlePostWebsitesRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PostWebsites",
			ID:   "post-websites",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securityCookieAuth(ctx, "PostWebsites", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "CookieAuth",
					Err:              err,
				}
				recordError("Security:CookieAuth", err)
				s.cfg.ErrorHandler(ctx, w, r, err)
				return
			}
			if ok {
				satisfied[0] |= 1 << 0
				ctx = sctx
			}
		}

		if ok := func() bool {
		nextRequirement:
			for _, requirement := range []bitset{
				{0b00000001},
			} {
				for i, mask := range requirement {
					if satisfied[i]&mask != mask {
						continue nextRequirement
					}
				}
				return true
			}
			return false
		}(); !ok {
			err = &ogenerrors.SecurityError{
				OperationContext: opErrContext,
				Err:              ogenerrors.ErrSecurityRequirementIsNotSatisfied,
			}
			recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	request, close, err := s.decodePostWebsitesRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PostWebsitesRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PostWebsites",
			OperationSummary: "Add Website.",
			OperationID:      "post-websites",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *WebsiteCreate
			Params   = struct{}
			Response = PostWebsitesRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostWebsites(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostWebsites(ctx, request)
	}
	if err != nil {
		recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostWebsitesResponse(response, w); err != nil {
		recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}
