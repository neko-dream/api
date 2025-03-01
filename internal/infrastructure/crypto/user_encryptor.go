package crypto

import (
	"context"
	"database/sql"

	"github.com/neko-dream/server/internal/domain/model/crypto"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/user"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/internal/usecase/query/dto"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

// Encrypt時はuser.UserDemographicの各フィールドを暗号化し、model.UserDemographicの各フィールドにセットする
// QueryでもRepositoryでも使える。
func EncryptUserDemographics(
	ctx context.Context,
	encryptor crypto.Encryptor,
	userID shared.UUID[user.User],
	userDemographic *user.UserDemographic,
) (*model.UserDemographic, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "EncryptUserDemographics")
	defer span.End()

	_ = ctx

	var city, prefecture, yearOfBirth, gender sql.NullString
	if userDemographic.City() != nil {
		encryptedCity, err := encryptor.EncryptString(ctx, userDemographic.City().String())
		if err != nil {
			utils.HandleError(ctx, err, "encryptor.EncryptString City")
			return nil, err
		}
		city = sql.NullString{String: encryptedCity, Valid: true}
	}

	if userDemographic.Prefecture() != nil {
		encryptedPrefecture, err := encryptor.EncryptString(ctx, *userDemographic.Prefecture())
		if err != nil {
			utils.HandleError(ctx, err, "encryptor.EncryptString Prefecture")
			return nil, err
		}
		prefecture = sql.NullString{String: encryptedPrefecture, Valid: true}
	}

	if userDemographic.YearOfBirth() != nil {
		encryptedYear, err := encryptor.EncryptInt(ctx, int64(*userDemographic.YearOfBirth()))
		if err != nil {
			utils.HandleError(ctx, err, "encryptor.EncryptInt YearOfBirth")
			return nil, err
		}
		yearOfBirth = sql.NullString{String: encryptedYear, Valid: true}
	}

	if userDemographic.Gender() != nil {
		encryptedGender, err := encryptor.EncryptInt(ctx, int64(*userDemographic.Gender()))
		if err != nil {
			utils.HandleError(ctx, err, "encrypt.EncryptInt Gender")
			return nil, err
		}
		gender = sql.NullString{String: encryptedGender, Valid: true}
	}

	// 2つには消えてもらうので一旦暗号化しない
	var householdSize sql.NullInt16
	if userDemographic.HouseholdSize() != nil {
		householdSize = sql.NullInt16{Int16: int16(*userDemographic.HouseholdSize()), Valid: true}
	}
	var occupation sql.NullInt16
	if userDemographic.Occupation() != nil {
		occupation = sql.NullInt16{Int16: int16(*userDemographic.Occupation()), Valid: true}
	}

	return &model.UserDemographic{
		UserID:             userID.UUID(),
		UserDemographicsID: userDemographic.ID().UUID(),
		City:               city,
		YearOfBirth:        yearOfBirth,
		HouseholdSize:      householdSize,
		Prefecture:         prefecture,
		Occupation:         occupation,
		Gender:             gender,
	}, nil
}

// DecryptUserDemographicsはmodel.UserDemographicの各フィールドを復号化し、user.UserDemographicの各フィールドにセットする
func DecryptUserDemographics(
	ctx context.Context,
	decryptor crypto.Encryptor,
	userDemographic *model.UserDemographic,
) (*user.UserDemographic, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "DecryptUserDemographics")
	defer span.End()

	var city, prefecture, gender *string
	var household, yearOfBirth, occupation *int

	if userDemographic.City.Valid {
		decryptedCity, err := decryptor.DecryptString(ctx, userDemographic.City.String)
		if err != nil {
			utils.HandleError(ctx, err, "decryptor.DecryptString City")
			return nil, err
		}
		city = &decryptedCity
	}
	if userDemographic.Prefecture.Valid {
		decryptedPrefecture, err := decryptor.DecryptString(ctx, userDemographic.Prefecture.String)
		if err != nil {
			utils.HandleError(ctx, err, "decryptor.DecryptString Prefecture")
			return nil, err
		}
		prefecture = &decryptedPrefecture
	}
	if userDemographic.YearOfBirth.Valid {
		decryptedYear, err := decryptor.DecryptInt(ctx, userDemographic.YearOfBirth.String)
		if err != nil {
			utils.HandleError(ctx, err, "decryptor.DecryptInt YearOfBirth")
			return nil, err
		}
		yearOfBirth = lo.ToPtr(int(decryptedYear))
	}
	if userDemographic.Gender.Valid {
		decryptedGender, err := decryptor.DecryptInt(ctx, userDemographic.Gender.String)
		if err != nil {
			utils.HandleError(ctx, err, "decrypt.DecryptInt Gender")
			return nil, err
		}

		gender = lo.ToPtr(user.Gender(decryptedGender).String())

	}

	if userDemographic.Occupation.Valid {
		occupationInt := int(userDemographic.Occupation.Int16)
		occupation = &occupationInt
	}

	var occupationStr *string
	if occupation != nil {
		str := user.Occupation(*occupation).String()
		occupationStr = &str
	}
	if userDemographic.HouseholdSize.Valid {
		householdSize := int(userDemographic.HouseholdSize.Int16)
		household = &householdSize
	}

	demo := user.NewUserDemographic(
		ctx,
		shared.UUID[user.UserDemographic](userDemographic.UserDemographicsID),
		yearOfBirth,
		occupationStr,
		gender,
		city,
		household,
		prefecture,
	)

	return &demo, nil
}

// DecryptUserDemographicsDTO model.UserDemographicをdto.UserDemographicに変換する
func DecryptUserDemographicsDTO(
	ctx context.Context,
	decryptor crypto.Encryptor,
	userDemographic *model.UserDemographic,
) (*dto.UserDemographic, error) {
	ctx, span := otel.Tracer("crypto").Start(ctx, "DecryptUserDemographicsDTO")
	defer span.End()

	decrypted, err := DecryptUserDemographics(ctx, decryptor, userDemographic)
	if err != nil {
		return nil, err
	}
	udDTO := dto.UserDemographic{
		UserDemographicID: userDemographic.UserDemographicsID,
		UserID:            shared.UUID[user.User](userDemographic.UserID),
	}

	if decrypted.YearOfBirth() != nil {
		udDTO.YearOfBirth = lo.ToPtr(int(*decrypted.YearOfBirth()))
	}
	if decrypted.Occupation() != nil {
		udDTO.Occupation = lo.ToPtr(int(*decrypted.Occupation()))
	}
	if decrypted.Gender() != nil {
		udDTO.Gender = lo.ToPtr(int(*decrypted.Gender()))
	}
	if decrypted.City() != nil {
		udDTO.City = lo.ToPtr(string(*decrypted.City()))
	}
	if decrypted.Prefecture() != nil {
		udDTO.Prefecture = lo.ToPtr(*decrypted.Prefecture())
	}
	if decrypted.HouseholdSize() != nil {
		udDTO.HouseholdSize = lo.ToPtr(int(*decrypted.HouseholdSize()))
	}

	return &udDTO, nil
}
