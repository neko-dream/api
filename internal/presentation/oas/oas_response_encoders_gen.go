// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"io"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/uri"
)

func encodeAuthorizeResponse(response AuthorizeRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *AuthorizeFoundHeaders:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Encoding response headers.
		{
			h := uri.NewHeaderEncoder(w.Header())
			// Encode "Location" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Location",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					return e.EncodeValue(conv.StringToString(response.Location))
				}); err != nil {
					return errors.Wrap(err, "encode Location header")
				}
			}
			// Encode "Set-Cookie" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Set-Cookie",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					return e.EncodeArray(func(e uri.Encoder) error {
						for i, item := range response.SetCookie {
							if err := func() error {
								return e.EncodeValue(conv.StringToString(item))
							}(); err != nil {
								return errors.Wrapf(err, "[%d]", i)
							}
						}
						return nil
					})
				}); err != nil {
					return errors.Wrap(err, "encode Set-Cookie header")
				}
			}
		}
		w.WriteHeader(302)
		span.SetStatus(codes.Ok, http.StatusText(302))

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *AuthorizeBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *AuthorizeInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeCreateTalkSessionResponse(response CreateTalkSessionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *CreateTalkSessionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *CreateTalkSessionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeDevAuthorizeResponse(response DevAuthorizeRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *DevAuthorizeFoundHeaders:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Encoding response headers.
		{
			h := uri.NewHeaderEncoder(w.Header())
			// Encode "Location" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Location",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					if val, ok := response.Location.Get(); ok {
						return e.EncodeValue(conv.StringToString(val))
					}
					return nil
				}); err != nil {
					return errors.Wrap(err, "encode Location header")
				}
			}
			// Encode "Set-Cookie" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Set-Cookie",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					if response.SetCookie != nil {
						return e.EncodeArray(func(e uri.Encoder) error {
							for i, item := range response.SetCookie {
								if err := func() error {
									return e.EncodeValue(conv.StringToString(item))
								}(); err != nil {
									return errors.Wrapf(err, "[%d]", i)
								}
							}
							return nil
						})
					}
					return nil
				}); err != nil {
					return errors.Wrap(err, "encode Set-Cookie header")
				}
			}
		}
		w.WriteHeader(302)
		span.SetStatus(codes.Ok, http.StatusText(302))

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *DevAuthorizeBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *DevAuthorizeInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeDummiInitResponse(response DummiInitRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *DummiInitOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *DummiInitBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *DummiInitInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeEditTalkSessionResponse(response EditTalkSessionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *EditTalkSessionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditTalkSessionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditTalkSessionInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeEditTimeLineResponse(response EditTimeLineRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *EditTimeLineOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditTimeLineBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditTimeLineInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeEditUserProfileResponse(response EditUserProfileRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *EditUserProfileOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditUserProfileBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *EditUserProfileInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetConclusionResponse(response GetConclusionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetConclusionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetConclusionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetConclusionInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetOpenedTalkSessionResponse(response GetOpenedTalkSessionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetOpenedTalkSessionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetOpenedTalkSessionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetOpenedTalkSessionInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetOpinionDetailResponse(response GetOpinionDetailRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetOpinionDetailOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetOpinionDetailInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetOpinionsForTalkSessionResponse(response GetOpinionsForTalkSessionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetOpinionsForTalkSessionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetOpinionsForTalkSessionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetOpinionsForTalkSessionInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetTalkSessionDetailResponse(response GetTalkSessionDetailRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetTalkSessionDetailOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTalkSessionDetailBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetTalkSessionListResponse(response GetTalkSessionListRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetTalkSessionListOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTalkSessionListBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTalkSessionListInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetTalkSessionReportResponse(response GetTalkSessionReportRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetTalkSessionReportOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTalkSessionReportInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetTimeLineResponse(response GetTimeLineRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetTimeLineOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTimeLineBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetTimeLineInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetUserInfoResponse(response GetUserInfoRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *GetUserInfoOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetUserInfoInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeManageIndexResponse(response ManageIndexOK, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	writer := w
	if _, err := io.Copy(writer, response); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeManageRegenerateResponse(response *ManageRegenerateOK, w http.ResponseWriter, span trace.Span) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	span.SetStatus(codes.Ok, http.StatusText(200))

	e := new(jx.Encoder)
	response.Encode(e)
	if _, err := e.WriteTo(w); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func encodeOAuthCallbackResponse(response OAuthCallbackRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OAuthCallbackFoundHeaders:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Encoding response headers.
		{
			h := uri.NewHeaderEncoder(w.Header())
			// Encode "Location" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Location",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					return e.EncodeValue(conv.StringToString(response.Location))
				}); err != nil {
					return errors.Wrap(err, "encode Location header")
				}
			}
			// Encode "Set-Cookie" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Set-Cookie",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					return e.EncodeValue(conv.StringToString(response.SetCookie))
				}); err != nil {
					return errors.Wrap(err, "encode Set-Cookie header")
				}
			}
		}
		w.WriteHeader(302)
		span.SetStatus(codes.Ok, http.StatusText(302))

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthCallbackBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthCallbackInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeOAuthTokenInfoResponse(response OAuthTokenInfoRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OAuthTokenInfoOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthTokenInfoBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthTokenInfoInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeOAuthTokenRevokeResponse(response OAuthTokenRevokeRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OAuthTokenRevokeNoContentHeaders:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Encoding response headers.
		{
			h := uri.NewHeaderEncoder(w.Header())
			// Encode "Set-Cookie" header.
			{
				cfg := uri.HeaderParameterEncodingConfig{
					Name:    "Set-Cookie",
					Explode: false,
				}
				if err := h.EncodeParam(cfg, func(e uri.Encoder) error {
					return e.EncodeArray(func(e uri.Encoder) error {
						for i, item := range response.SetCookie {
							if err := func() error {
								return e.EncodeValue(conv.StringToString(item))
							}(); err != nil {
								return errors.Wrapf(err, "[%d]", i)
							}
						}
						return nil
					})
				}); err != nil {
					return errors.Wrap(err, "encode Set-Cookie header")
				}
			}
		}
		w.WriteHeader(204)
		span.SetStatus(codes.Ok, http.StatusText(204))

		e := new(jx.Encoder)
		response.Response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthTokenRevokeBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OAuthTokenRevokeInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeOpinionCommentsResponse(response OpinionCommentsRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OpinionCommentsOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OpinionCommentsBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OpinionCommentsInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeOpinionsHistoryResponse(response OpinionsHistoryRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *OpinionsHistoryOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OpinionsHistoryBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *OpinionsHistoryInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodePostConclusionResponse(response PostConclusionRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *PostConclusionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostConclusionBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostConclusionInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodePostOpinionPostResponse(response PostOpinionPostRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *PostOpinionPostOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostOpinionPostBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostOpinionPostInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodePostTimeLineItemResponse(response PostTimeLineItemRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *PostTimeLineItemOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostTimeLineItemBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *PostTimeLineItemInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeRegisterUserResponse(response RegisterUserRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *RegisterUserOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *RegisterUserBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *RegisterUserInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeSessionsHistoryResponse(response SessionsHistoryRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *SessionsHistoryOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *SessionsHistoryBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *SessionsHistoryInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeSwipeOpinionsResponse(response SwipeOpinionsRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *SwipeOpinionsOKApplicationJSON:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *SwipeOpinionsBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *SwipeOpinionsInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeTalkSessionAnalysisResponse(response TalkSessionAnalysisRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *TalkSessionAnalysisOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *TalkSessionAnalysisBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *TalkSessionAnalysisInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeTestResponse(response TestRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *TestOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *TestBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *TestInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeVoteResponse(response VoteRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *VoteOKApplicationJSON:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *VoteBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *VoteInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)
		span.SetStatus(codes.Error, http.StatusText(500))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}
