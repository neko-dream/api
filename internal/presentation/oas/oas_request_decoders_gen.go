// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"mime"
	"net/http"
	"net/url"

	"github.com/go-faster/errors"
	"go.uber.org/multierr"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

func (s *Server) decodeCreateTalkSessionRequest(r *http.Request) (
	req OptCreateTalkSessionReq,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	if _, ok := r.Header["Content-Type"]; !ok && r.ContentLength == 0 {
		return req, close, nil
	}
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "application/x-www-form-urlencoded":
		if r.ContentLength == 0 {
			return req, close, nil
		}
		form, err := ht.ParseForm(r)
		if err != nil {
			return req, close, errors.Wrap(err, "parse form")
		}

		var request OptCreateTalkSessionReq
		{
			var optForm CreateTalkSessionReq
			q := uri.NewQueryDecoder(form)
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "theme",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotThemeVal string
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotThemeVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.Theme.SetTo(optFormDotThemeVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"theme\"")
					}
				}
			}
			request = OptCreateTalkSessionReq{
				Value: optForm,
				Set:   true,
			}
		}
		return request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}

func (s *Server) decodeIntentionRequest(r *http.Request) (
	req OptIntentionReq,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	if _, ok := r.Header["Content-Type"]; !ok && r.ContentLength == 0 {
		return req, close, nil
	}
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "application/x-www-form-urlencoded":
		if r.ContentLength == 0 {
			return req, close, nil
		}
		form, err := ht.ParseForm(r)
		if err != nil {
			return req, close, errors.Wrap(err, "parse form")
		}

		var request OptIntentionReq
		{
			var optForm IntentionReq
			q := uri.NewQueryDecoder(form)
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "intentionStatus",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToString(val)
						if err != nil {
							return err
						}

						optForm.IntentionStatus = c
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"intentionStatus\"")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			request = OptIntentionReq{
				Value: optForm,
				Set:   true,
			}
		}
		return request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}

func (s *Server) decodePostOpinionPostRequest(r *http.Request) (
	req OptPostOpinionPostReq,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	if _, ok := r.Header["Content-Type"]; !ok && r.ContentLength == 0 {
		return req, close, nil
	}
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "multipart/form-data":
		if r.ContentLength == 0 {
			return req, close, nil
		}
		if err := r.ParseMultipartForm(s.cfg.MaxMultipartMemory); err != nil {
			return req, close, errors.Wrap(err, "parse multipart form")
		}
		// Remove all temporary files created by ParseMultipartForm when the request is done.
		//
		// Notice that the closers are called in reverse order, to match defer behavior, so
		// any opened file will be closed before RemoveAll call.
		closers = append(closers, r.MultipartForm.RemoveAll)
		// Form values may be unused.
		form := url.Values(r.MultipartForm.Value)
		_ = form

		var request OptPostOpinionPostReq
		{
			var optForm PostOpinionPostReq
			q := uri.NewQueryDecoder(form)
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "parentOpinionID",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotParentOpinionIDVal string
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotParentOpinionIDVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.ParentOpinionID.SetTo(optFormDotParentOpinionIDVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"parentOpinionID\"")
					}
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "intentionStatus",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToString(val)
						if err != nil {
							return err
						}

						optForm.IntentionStatus = c
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"intentionStatus\"")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "title",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotTitleVal string
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotTitleVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.Title.SetTo(optFormDotTitleVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"title\"")
					}
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "opinion",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToString(val)
						if err != nil {
							return err
						}

						optForm.Opinion = c
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"opinion\"")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "referenceURL",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotReferenceURLVal string
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotReferenceURLVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.ReferenceURL.SetTo(optFormDotReferenceURLVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"referenceURL\"")
					}
				}
			}
			{
				if err := func() error {
					files, ok := r.MultipartForm.File["picture"]
					if !ok || len(files) < 1 {
						return nil
					}
					fh := files[0]

					f, err := fh.Open()
					if err != nil {
						return errors.Wrap(err, "open")
					}
					closers = append(closers, f.Close)
					optForm.Picture.SetTo(ht.MultipartFile{
						Name:   fh.Filename,
						File:   f,
						Size:   fh.Size,
						Header: fh.Header,
					})
					return nil
				}(); err != nil {
					return req, close, errors.Wrap(err, "decode \"picture\"")
				}
			}
			request = OptPostOpinionPostReq{
				Value: optForm,
				Set:   true,
			}
		}
		return request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}

func (s *Server) decodeRegisterUserRequest(r *http.Request) (
	req OptRegisterUserReq,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	if _, ok := r.Header["Content-Type"]; !ok && r.ContentLength == 0 {
		return req, close, nil
	}
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "multipart/form-data":
		if r.ContentLength == 0 {
			return req, close, nil
		}
		if err := r.ParseMultipartForm(s.cfg.MaxMultipartMemory); err != nil {
			return req, close, errors.Wrap(err, "parse multipart form")
		}
		// Remove all temporary files created by ParseMultipartForm when the request is done.
		//
		// Notice that the closers are called in reverse order, to match defer behavior, so
		// any opened file will be closed before RemoveAll call.
		closers = append(closers, r.MultipartForm.RemoveAll)
		// Form values may be unused.
		form := url.Values(r.MultipartForm.Value)
		_ = form

		var request OptRegisterUserReq
		{
			var optForm RegisterUserReq
			q := uri.NewQueryDecoder(form)
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "displayName",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToString(val)
						if err != nil {
							return err
						}

						optForm.DisplayName = c
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"displayName\"")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "displayID",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						val, err := d.DecodeValue()
						if err != nil {
							return err
						}

						c, err := conv.ToString(val)
						if err != nil {
							return err
						}

						optForm.DisplayID = c
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"displayID\"")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "yearOfBirth",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotYearOfBirthVal int
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToInt(val)
							if err != nil {
								return err
							}

							optFormDotYearOfBirthVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.YearOfBirth.SetTo(optFormDotYearOfBirthVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"yearOfBirth\"")
					}
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "gender",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotGenderVal RegisterUserReqGender
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotGenderVal = RegisterUserReqGender(c)
							return nil
						}(); err != nil {
							return err
						}
						optForm.Gender.SetTo(optFormDotGenderVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"gender\"")
					}
					if err := func() error {
						if value, ok := optForm.Gender.Get(); ok {
							if err := func() error {
								if err := value.Validate(); err != nil {
									return err
								}
								return nil
							}(); err != nil {
								return err
							}
						}
						return nil
					}(); err != nil {
						return req, close, errors.Wrap(err, "validate")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "municipality",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotMunicipalityVal string
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotMunicipalityVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.Municipality.SetTo(optFormDotMunicipalityVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"municipality\"")
					}
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "occupation",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotOccupationVal RegisterUserReqOccupation
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToString(val)
							if err != nil {
								return err
							}

							optFormDotOccupationVal = RegisterUserReqOccupation(c)
							return nil
						}(); err != nil {
							return err
						}
						optForm.Occupation.SetTo(optFormDotOccupationVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"occupation\"")
					}
					if err := func() error {
						if value, ok := optForm.Occupation.Get(); ok {
							if err := func() error {
								if err := value.Validate(); err != nil {
									return err
								}
								return nil
							}(); err != nil {
								return err
							}
						}
						return nil
					}(); err != nil {
						return req, close, errors.Wrap(err, "validate")
					}
				} else {
					return req, close, errors.Wrap(err, "query")
				}
			}
			{
				cfg := uri.QueryParameterDecodingConfig{
					Name:    "householdSize",
					Style:   uri.QueryStyleForm,
					Explode: true,
				}
				if err := q.HasParam(cfg); err == nil {
					if err := q.DecodeParam(cfg, func(d uri.Decoder) error {
						var optFormDotHouseholdSizeVal int
						if err := func() error {
							val, err := d.DecodeValue()
							if err != nil {
								return err
							}

							c, err := conv.ToInt(val)
							if err != nil {
								return err
							}

							optFormDotHouseholdSizeVal = c
							return nil
						}(); err != nil {
							return err
						}
						optForm.HouseholdSize.SetTo(optFormDotHouseholdSizeVal)
						return nil
					}); err != nil {
						return req, close, errors.Wrap(err, "decode \"householdSize\"")
					}
					if err := func() error {
						if value, ok := optForm.HouseholdSize.Get(); ok {
							if err := func() error {
								if err := (validate.Int{
									MinSet:        true,
									Min:           0,
									MaxSet:        false,
									Max:           0,
									MinExclusive:  false,
									MaxExclusive:  false,
									MultipleOfSet: false,
									MultipleOf:    0,
								}).Validate(int64(value)); err != nil {
									return errors.Wrap(err, "int")
								}
								return nil
							}(); err != nil {
								return err
							}
						}
						return nil
					}(); err != nil {
						return req, close, errors.Wrap(err, "validate")
					}
				}
			}
			{
				if err := func() error {
					files, ok := r.MultipartForm.File["icon"]
					if !ok || len(files) < 1 {
						return nil
					}
					fh := files[0]

					f, err := fh.Open()
					if err != nil {
						return errors.Wrap(err, "open")
					}
					closers = append(closers, f.Close)
					optForm.Icon.SetTo(ht.MultipartFile{
						Name:   fh.Filename,
						File:   f,
						Size:   fh.Size,
						Header: fh.Header,
					})
					return nil
				}(); err != nil {
					return req, close, errors.Wrap(err, "decode \"icon\"")
				}
			}
			request = OptRegisterUserReq{
				Value: optForm,
				Set:   true,
			}
		}
		return request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}

func (s *Server) decodeTestRequest(r *http.Request) (
	req OptTestReq,
	close func() error,
	rerr error,
) {
	var closers []func() error
	close = func() error {
		var merr error
		// Close in reverse order, to match defer behavior.
		for i := len(closers) - 1; i >= 0; i-- {
			c := closers[i]
			merr = multierr.Append(merr, c())
		}
		return merr
	}
	defer func() {
		if rerr != nil {
			rerr = multierr.Append(rerr, close())
		}
	}()
	if _, ok := r.Header["Content-Type"]; !ok && r.ContentLength == 0 {
		return req, close, nil
	}
	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return req, close, errors.Wrap(err, "parse media type")
	}
	switch {
	case ct == "multipart/form-data":
		if r.ContentLength == 0 {
			return req, close, nil
		}
		if err := r.ParseMultipartForm(s.cfg.MaxMultipartMemory); err != nil {
			return req, close, errors.Wrap(err, "parse multipart form")
		}
		// Remove all temporary files created by ParseMultipartForm when the request is done.
		//
		// Notice that the closers are called in reverse order, to match defer behavior, so
		// any opened file will be closed before RemoveAll call.
		closers = append(closers, r.MultipartForm.RemoveAll)
		// Form values may be unused.
		form := url.Values(r.MultipartForm.Value)
		_ = form

		var request OptTestReq
		{
			var optForm TestReq
			{
				if err := func() error {
					files, ok := r.MultipartForm.File["file"]
					if !ok || len(files) < 1 {
						return validate.ErrFieldRequired
					}
					fh := files[0]

					f, err := fh.Open()
					if err != nil {
						return errors.Wrap(err, "open")
					}
					closers = append(closers, f.Close)
					optForm.File = ht.MultipartFile{
						Name:   fh.Filename,
						File:   f,
						Size:   fh.Size,
						Header: fh.Header,
					}
					return nil
				}(); err != nil {
					return req, close, errors.Wrap(err, "decode \"file\"")
				}
			}
			request = OptTestReq{
				Value: optForm,
				Set:   true,
			}
		}
		return request, close, nil
	default:
		return req, close, validate.InvalidContentType(ct)
	}
}
