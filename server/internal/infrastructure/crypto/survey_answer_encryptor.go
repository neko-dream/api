package crypto

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/crypto"
	"github.com/neko-dream/api/internal/domain/model/talksession/survey"
	"github.com/neko-dream/api/pkg/utils"
	"go.opentelemetry.io/otel"
)

// EncryptAnswerValue は回答値全体を暗号化する。
// 平文（SelectedChoiceIDs / Other）を JSON 化し、AES-GCM で暗号化した文字列を
// Encrypted フィールドにのみ格納した AnswerValue を返す。集計に必要な選択肢 ID も
// 含めてすべて暗号化されるため、集計時はサーバ側で復号する必要がある。
func EncryptAnswerValue(
	ctx context.Context,
	encryptor crypto.Encryptor,
	v survey.AnswerValue,
) (survey.AnswerValue, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "EncryptAnswerValue")
	defer span.End()

	// Encrypted を除いた平文部分のみを JSON 化する。
	plain := survey.AnswerValue{
		SelectedChoiceIDs: v.SelectedChoiceIDs,
		Other:             v.Other,
	}
	data, err := plain.Marshal()
	if err != nil {
		utils.HandleError(ctx, err, "AnswerValue.Marshal")
		return survey.AnswerValue{}, err
	}

	enc, err := encryptor.EncryptString(ctx, string(data))
	if err != nil {
		utils.HandleError(ctx, err, "EncryptString AnswerValue")
		return survey.AnswerValue{}, err
	}

	return survey.AnswerValue{Encrypted: &enc}, nil
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

	if v.Encrypted == nil || *v.Encrypted == "" {
		return v, nil
	}

	data, err := encryptor.DecryptString(ctx, *v.Encrypted)
	if err != nil {
		utils.HandleError(ctx, err, "DecryptString AnswerValue")
		return survey.AnswerValue{}, err
	}

	out, err := survey.UnmarshalAnswerValue([]byte(data))
	if err != nil {
		utils.HandleError(ctx, err, "UnmarshalAnswerValue")
		return survey.AnswerValue{}, err
	}
	return out, nil
}
