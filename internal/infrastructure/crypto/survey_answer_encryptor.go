package crypto

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

// EncryptAnswerValue は PII フラグ付き設問への回答を暗号化する。
// 暗号化対象は Text / Other / Date のみ。SelectedChoiceIDs と Rating は集計可能性のため平文のまま。
func EncryptAnswerValue(
	ctx context.Context,
	encryptor crypto.Encryptor,
	v survey.AnswerValue,
) (survey.AnswerValue, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "EncryptAnswerValue")
	defer span.End()

	out := survey.AnswerValue{
		SelectedChoiceIDs: v.SelectedChoiceIDs,
		Rating:            v.Rating,
	}

	if v.Text != nil && *v.Text != "" {
		enc, err := encryptor.EncryptString(ctx, *v.Text)
		if err != nil {
			utils.HandleError(ctx, err, "EncryptString Text")
			return survey.AnswerValue{}, err
		}
		out.Text = &enc
	} else {
		out.Text = v.Text
	}

	if v.Other != nil && *v.Other != "" {
		enc, err := encryptor.EncryptString(ctx, *v.Other)
		if err != nil {
			utils.HandleError(ctx, err, "EncryptString Other")
			return survey.AnswerValue{}, err
		}
		out.Other = &enc
	} else {
		out.Other = v.Other
	}

	if v.Date != nil && *v.Date != "" {
		enc, err := encryptor.EncryptString(ctx, *v.Date)
		if err != nil {
			utils.HandleError(ctx, err, "EncryptString Date")
			return survey.AnswerValue{}, err
		}
		out.Date = &enc
	} else {
		out.Date = v.Date
	}

	return out, nil
}

// DecryptAnswerValue は EncryptAnswerValue で暗号化した値を復号する。
// is_encrypted=true の Answer に対してのみ呼び出すこと。
func DecryptAnswerValue(
	ctx context.Context,
	encryptor crypto.Encryptor,
	v survey.AnswerValue,
) (survey.AnswerValue, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "DecryptAnswerValue")
	defer span.End()

	out := survey.AnswerValue{
		SelectedChoiceIDs: v.SelectedChoiceIDs,
		Rating:            v.Rating,
	}

	if v.Text != nil && *v.Text != "" {
		dec, err := encryptor.DecryptString(ctx, *v.Text)
		if err != nil {
			utils.HandleError(ctx, err, "DecryptString Text")
			return survey.AnswerValue{}, err
		}
		out.Text = &dec
	} else {
		out.Text = v.Text
	}

	if v.Other != nil && *v.Other != "" {
		dec, err := encryptor.DecryptString(ctx, *v.Other)
		if err != nil {
			utils.HandleError(ctx, err, "DecryptString Other")
			return survey.AnswerValue{}, err
		}
		out.Other = &dec
	} else {
		out.Other = v.Other
	}

	if v.Date != nil && *v.Date != "" {
		dec, err := encryptor.DecryptString(ctx, *v.Date)
		if err != nil {
			utils.HandleError(ctx, err, "DecryptString Date")
			return survey.AnswerValue{}, err
		}
		out.Date = &dec
	} else {
		out.Date = v.Date
	}

	return out, nil
}
