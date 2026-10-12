package admin

import (
	"context"
	"errors"

	adminmodel "github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/idtoken"
)

type googleTokenVerifier struct {
	audience string
}

func NewGoogleTokenVerifier(conf *config.Config) adminmodel.GoogleTokenVerifier {
	return &googleTokenVerifier{audience: conf.GoogleClientID}
}

func (v *googleTokenVerifier) Verify(ctx context.Context, idToken string) (*adminmodel.GoogleClaims, error) {
	ctx, span := otel.Tracer("googleTokenVerifier").Start(ctx, "Verify")
	defer span.End()

	payload, err := idtoken.Validate(ctx, idToken, v.audience)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	if email == "" {
		return nil, errors.New("email claim not found in google id token")
	}
	if verified, ok := payload.Claims["email_verified"].(bool); ok && !verified {
		return nil, errors.New("google email is not verified")
	}
	name, _ := payload.Claims["name"].(string)

	return &adminmodel.GoogleClaims{Email: email, Name: name}, nil
}
