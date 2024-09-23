// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/otelogen"
)

// handleAuthLoginRequest handles auth_login operation.
//
// OAuthログイン.
//
// GET /auth/{provider}/login
func (s *Server) handleAuthLoginRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("auth_login"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/auth/{provider}/login"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "AuthLogin",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "AuthLogin",
			ID:   "auth_login",
		}
	)
	params, err := decodeAuthLoginParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response *AuthLoginFound
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "AuthLogin",
			OperationSummary: "OAuthログイン",
			OperationID:      "auth_login",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "provider",
					In:   "path",
				}: params.Provider,
				{
					Name: "redirect_url",
					In:   "query",
				}: params.RedirectURL,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = AuthLoginParams
			Response = *AuthLoginFound
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackAuthLoginParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.AuthLogin(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.AuthLogin(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeAuthLoginResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleCreateTalkSessionRequest handles createTalkSession operation.
//
// トークセッション作成.
//
// POST /api/talksessions
func (s *Server) handleCreateTalkSessionRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("createTalkSession"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/api/talksessions"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "CreateTalkSession",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err error
	)

	var response *CreateTalkSessionOK
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CreateTalkSession",
			OperationSummary: "トークセッション作成",
			OperationID:      "createTalkSession",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *CreateTalkSessionOK
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
				response, err = s.h.CreateTalkSession(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.CreateTalkSession(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeCreateTalkSessionResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleEditUserProfileRequest handles editUserProfile operation.
//
// ユーザー情報の変更.
//
// PUT /api/user
func (s *Server) handleEditUserProfileRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("editUserProfile"),
		semconv.HTTPRequestMethodKey.String("PUT"),
		semconv.HTTPRouteKey.String("/api/user"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "EditUserProfile",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err error
	)

	var response *EditUserProfileOK
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "EditUserProfile",
			OperationSummary: "ユーザー情報の変更",
			OperationID:      "editUserProfile",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *EditUserProfileOK
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
				response, err = s.h.EditUserProfile(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.EditUserProfile(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeEditUserProfileResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetTalkSessionDetailRequest handles getTalkSessionDetail operation.
//
// トークセッションの詳細.
//
// GET /api/talksessions/{talkSessionId}
func (s *Server) handleGetTalkSessionDetailRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getTalkSessionDetail"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/api/talksessions/{talkSessionId}"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "GetTalkSessionDetail",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetTalkSessionDetail",
			ID:   "getTalkSessionDetail",
		}
	)
	params, err := decodeGetTalkSessionDetailParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response *GetTalkSessionDetailOK
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetTalkSessionDetail",
			OperationSummary: "トークセッションの詳細",
			OperationID:      "getTalkSessionDetail",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "talkSessionId",
					In:   "path",
				}: params.TalkSessionId,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetTalkSessionDetailParams
			Response = *GetTalkSessionDetailOK
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetTalkSessionDetailParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetTalkSessionDetail(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetTalkSessionDetail(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetTalkSessionDetailResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetTalkSessionsRequest handles getTalkSessions operation.
//
// トークセッションリスト.
//
// GET /api/talksessions
func (s *Server) handleGetTalkSessionsRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getTalkSessions"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/api/talksessions"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "GetTalkSessions",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err error
	)

	var response *GetTalkSessionsOK
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetTalkSessions",
			OperationSummary: "トークセッションリスト",
			OperationID:      "getTalkSessions",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *GetTalkSessionsOK
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
				err = s.h.GetTalkSessions(ctx)
				return response, err
			},
		)
	} else {
		err = s.h.GetTalkSessions(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetTalkSessionsResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetUserProfileRequest handles getUserProfile operation.
//
// ユーザー情報の取得.
//
// GET /api/user
func (s *Server) handleGetUserProfileRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getUserProfile"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/api/user"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "GetUserProfile",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err error
	)

	var response *GetUserProfileOK
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetUserProfile",
			OperationSummary: "ユーザー情報の取得",
			OperationID:      "getUserProfile",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = *GetUserProfileOK
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
				response, err = s.h.GetUserProfile(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetUserProfile(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetUserProfileResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleIndicateIntentionRequest handles indicateIntention operation.
//
// 意思表明API.
//
// POST /api/talksessions/{talkSessionID}/opinions/{opinionID}/intentions
func (s *Server) handleIndicateIntentionRequest(args [2]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("indicateIntention"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/api/talksessions/{talkSessionID}/opinions/{opinionID}/intentions"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "IndicateIntention",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "IndicateIntention",
			ID:   "indicateIntention",
		}
	)
	params, err := decodeIndicateIntentionParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response IndicateIntentionRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "IndicateIntention",
			OperationSummary: "意思表明API",
			OperationID:      "indicateIntention",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "talkSessionID",
					In:   "path",
				}: params.TalkSessionID,
				{
					Name: "opinionID",
					In:   "path",
				}: params.OpinionID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = IndicateIntentionParams
			Response = IndicateIntentionRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackIndicateIntentionParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.IndicateIntention(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.IndicateIntention(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeIndicateIntentionResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleListOpinionsRequest handles listOpinions operation.
//
// セッションの意見一覧.
//
// GET /api/talksession/{talkSessionID}/opinions
func (s *Server) handleListOpinionsRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("listOpinions"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/api/talksession/{talkSessionID}/opinions"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "ListOpinions",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "ListOpinions",
			ID:   "listOpinions",
		}
	)
	params, err := decodeListOpinionsParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response ListOpinionsRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "ListOpinions",
			OperationSummary: "セッションの意見一覧",
			OperationID:      "listOpinions",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "talkSessionID",
					In:   "path",
				}: params.TalkSessionID,
				{
					Name: "parentOpinionID",
					In:   "query",
				}: params.ParentOpinionID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = ListOpinionsParams
			Response = ListOpinionsRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackListOpinionsParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.ListOpinions(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.ListOpinions(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeListOpinionsResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleOAuthCallbackRequest handles oauth_callback operation.
//
// Auth Callback.
//
// GET /auth/{provider}/callback
func (s *Server) handleOAuthCallbackRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("oauth_callback"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/auth/{provider}/callback"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "OAuthCallback",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "OAuthCallback",
			ID:   "oauth_callback",
		}
	)
	params, err := decodeOAuthCallbackParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response *OAuthCallbackFound
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "OAuthCallback",
			OperationSummary: "Auth Callback",
			OperationID:      "oauth_callback",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "state",
					In:   "cookie",
				}: params.CookieState,
				{
					Name: "redirect_url",
					In:   "cookie",
				}: params.RedirectURL,
				{
					Name: "provider",
					In:   "path",
				}: params.Provider,
				{
					Name: "code",
					In:   "query",
				}: params.Code,
				{
					Name: "state",
					In:   "query",
				}: params.QueryState,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = OAuthCallbackParams
			Response = *OAuthCallbackFound
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackOAuthCallbackParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.OAuthCallback(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.OAuthCallback(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeOAuthCallbackResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostOpinionPostRequest handles postOpinionPost operation.
//
// セッションに対して意見投稿.
//
// POST /api/talksessions/{talkSessionID}/opinions
func (s *Server) handlePostOpinionPostRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("postOpinionPost"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/api/talksessions/{talkSessionID}/opinions"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "PostOpinionPost",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "PostOpinionPost",
			ID:   "postOpinionPost",
		}
	)
	params, err := decodePostOpinionPostParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response PostOpinionPostRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "PostOpinionPost",
			OperationSummary: "セッションに対して意見投稿",
			OperationID:      "postOpinionPost",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "talkSessionID",
					In:   "path",
				}: params.TalkSessionID,
				{
					Name: "opinionContent",
					In:   "query",
				}: params.OpinionContent,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = PostOpinionPostParams
			Response = PostOpinionPostRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackPostOpinionPostParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostOpinionPost(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostOpinionPost(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostOpinionPostResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleRegisterUserRequest handles registerUser operation.
//
// ユーザー作成.
//
// POST /api/user
func (s *Server) handleRegisterUserRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("registerUser"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/api/user"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), "RegisterUser",
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)
		attrOpt := metric.WithAttributeSet(labeler.AttributeSet())

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(float64(elapsedDuration)/float64(time.Millisecond)), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)
			span.SetStatus(codes.Error, stage)
			s.errors.Add(ctx, 1, metric.WithAttributeSet(labeler.AttributeSet()))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "RegisterUser",
			ID:   "registerUser",
		}
	)
	{
		type bitset = [1]uint8
		var satisfied bitset
		{
			sctx, ok, err := s.securitySessionId(ctx, "RegisterUser", r)
			if err != nil {
				err = &ogenerrors.SecurityError{
					OperationContext: opErrContext,
					Security:         "SessionId",
					Err:              err,
				}
				defer recordError("Security:SessionId", err)
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
			defer recordError("Security", err)
			s.cfg.ErrorHandler(ctx, w, r, err)
			return
		}
	}
	params, err := decodeRegisterUserParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response RegisterUserRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "RegisterUser",
			OperationSummary: "ユーザー作成",
			OperationID:      "registerUser",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "displayName",
					In:   "query",
				}: params.DisplayName,
				{
					Name: "displayID",
					In:   "query",
				}: params.DisplayID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = RegisterUserParams
			Response = RegisterUserRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackRegisterUserParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.RegisterUser(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.RegisterUser(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeRegisterUserResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}
