// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

func encodeCreateTalkSessionRequest(
	req OptCreateTalkSessionReq,
	r *http.Request,
) error {
	const contentType = "application/x-www-form-urlencoded"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	{
		// Encode "theme" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "theme",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.Theme.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	encoded := q.Values().Encode()
	ht.SetBody(r, strings.NewReader(encoded), contentType)
	return nil
}

func encodeIntentionRequest(
	req OptIntentionReq,
	r *http.Request,
) error {
	const contentType = "application/x-www-form-urlencoded"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	{
		// Encode "intentionStatus" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "intentionStatus",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(request.IntentionStatus))
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	encoded := q.Values().Encode()
	ht.SetBody(r, strings.NewReader(encoded), contentType)
	return nil
}

func encodePostOpinionPostRequest(
	req OptPostOpinionPostReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	{
		// Encode "parentOpinionID" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "parentOpinionID",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.ParentOpinionID.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "intentionStatus" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "intentionStatus",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(request.IntentionStatus))
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "title" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "title",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.Title.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "opinion" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "opinion",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(request.Opinion))
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "referenceURL" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "referenceURL",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.ReferenceURL.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if val, ok := request.Picture.Get(); ok {
			if err := val.WriteMultipart("picture", w); err != nil {
				return errors.Wrap(err, "write \"picture\"")
			}
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}

func encodeRegisterUserRequest(
	req OptRegisterUserReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	{
		// Encode "displayName" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "displayName",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(request.DisplayName))
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "displayID" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "displayID",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			return e.EncodeValue(conv.StringToString(request.DisplayID))
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "yearOfBirth" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "yearOfBirth",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.YearOfBirth.Get(); ok {
				return e.EncodeValue(conv.IntToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "gender" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "gender",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.Gender.Get(); ok {
				return e.EncodeValue(conv.StringToString(string(val)))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "municipality" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "municipality",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.Municipality.Get(); ok {
				return e.EncodeValue(conv.StringToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "occupation" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "occupation",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.Occupation.Get(); ok {
				return e.EncodeValue(conv.StringToString(string(val)))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	{
		// Encode "householdSize" form field.
		cfg := uri.QueryParameterEncodingConfig{
			Name:    "householdSize",
			Style:   uri.QueryStyleForm,
			Explode: true,
		}
		if err := q.EncodeParam(cfg, func(e uri.Encoder) error {
			if val, ok := request.HouseholdSize.Get(); ok {
				return e.EncodeValue(conv.IntToString(val))
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "encode query")
		}
	}
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if val, ok := request.Icon.Get(); ok {
			if err := val.WriteMultipart("icon", w); err != nil {
				return errors.Wrap(err, "write \"icon\"")
			}
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}

func encodeTestRequest(
	req OptTestReq,
	r *http.Request,
) error {
	const contentType = "multipart/form-data"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	request := req.Value

	q := uri.NewFormEncoder(map[string]string{})
	body, boundary := ht.CreateMultipartBody(func(w *multipart.Writer) error {
		if err := request.File.WriteMultipart("file", w); err != nil {
			return errors.Wrap(err, "write \"file\"")
		}
		if err := q.WriteMultipart(w); err != nil {
			return errors.Wrap(err, "write multipart")
		}
		return nil
	})
	ht.SetCloserBody(r, body, mime.FormatMediaType(contentType, map[string]string{"boundary": boundary}))
	return nil
}
