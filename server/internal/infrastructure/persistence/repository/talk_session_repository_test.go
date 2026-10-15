package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/domain/model/clock"
	"github.com/neko-dream/api/internal/domain/model/organization"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/infrastructure/crypto"
	ci "github.com/neko-dream/api/internal/infrastructure/crypto"
	"github.com/neko-dream/api/internal/infrastructure/di"
	"github.com/neko-dream/api/internal/infrastructure/persistence/db"
	"github.com/neko-dream/api/internal/infrastructure/persistence/repository"
	"github.com/neko-dream/api/internal/test/txtest"
	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestTalkSessionRepository_Create(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	encryptor, _ := ci.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012", // テスト用の32バイトキー
	}))
	type TestData struct {
		TsRepo      talksession.TalkSessionRepository
		TalkSession *talksession.TalkSession
	}

	eventStore := repository.NewEventStoreMock()
	initData := TestData{
		TsRepo: repository.NewTalkSessionRepository(dbManager, eventStore),
	}
	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	ownerUserID := shared.NewUUID[user.User]()
	userRepo := repository.NewUserRepository(
		dbManager,
		repository.NewImageRepositoryMock(),
		encryptor,
	)
	userRepo.Create(context.Background(), user.NewUser(
		ownerUserID,
		lo.ToPtr("test"),
		lo.ToPtr("test"),
		"test",
		"GOOGLE",
		nil,
	))

	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "トークセッション作成ができる",
			SetupFn: func(ctx context.Context, data *TestData) error {
				data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					nil,
					lo.ToPtr("https://example.com/test.jpg"),
					ownerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					nil,
					nil,
					true,
					nil, // organizationID
					nil, // organizationAliasID
					nil,
				)
				return nil
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				if err := data.TsRepo.Create(ctx, data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSession.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}

				return nil
			},
			WantErr: false,
		},
		{
			Name: "トークセッション作成ができ、Locationも保存される",
			SetupFn: func(ctx context.Context, data *TestData) error {
				data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					nil,
					lo.ToPtr("https://example.com/test.jpg"),
					ownerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					talksession.NewLocation(
						talkSessionID,
						30.0,
						30.0,
					),
					nil, nil,
					true,
					nil, // organizationID
					nil, // organizationAliasID
					nil,
				)
				return nil
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				if err := data.TsRepo.Create(ctx, data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSession.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}
				location := talksession.NewLocation(
					talkSessionID,
					ts.Latitude,
					ts.Longitude,
				)
				assert.Equal(t, location.ToGeographyText(), "POINT(30.000000 30.000000)")

				return nil
			},
			WantErr: false,
		},
		{
			Name: "トークセッション作成ができ、ThumbnailURLも保存される",
			SetupFn: func(ctx context.Context, data *TestData) error {
				data.TalkSession = talksession.NewTalkSession(
					talkSessionID,
					"test",
					nil,
					lo.ToPtr("https://example.com/thumbnail.jpg"),
					ownerUserID,
					clock.Now(ctx),
					clock.Now(ctx).Add(time.Hour*24),
					nil,
					nil,
					nil,
					true,
					nil, // organizationID
					nil, // organizationAliasID
					nil,
				)
				return nil
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				if err := data.TsRepo.Create(ctx, data.TalkSession); err != nil {
					return err
				}

				ts, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				if ts.TalkSession.TalkSessionID != talkSessionID.UUID() {
					return errors.New("トークセッションIDが一致しません")
				}
				assert.Equal(t, *data.TalkSession.ThumbnailURL(), ts.TalkSession.ThumbnailUrl.String)

				return nil
			},
			WantErr: false,
		},
	}

	txtest.RunTransactionalTests(t, dbManager, &initData, testCases)
}

func TestTalkSessionRepository_UpdateKeepsOrganization(t *testing.T) {
	container := di.BuildContainer()
	dbManager := di.Invoke[*db.DBManager](container)
	encryptor, _ := ci.NewEncryptor(lo.ToPtr(config.Config{
		ENCRYPTION_VERSION: crypto.Version1,
		ENCRYPTION_SECRET:  "12345678901234567890123456789012", // テスト用の32バイトキー
	}))
	type TestData struct {
		TsRepo talksession.TalkSessionRepository
	}
	initData := TestData{
		TsRepo: repository.NewTalkSessionRepository(dbManager, repository.NewEventStoreMock()),
	}
	talkSessionID := shared.NewUUID[talksession.TalkSession]()
	ownerUserID := shared.NewUUID[user.User]()
	organizationID := shared.NewUUID[organization.Organization]()
	var aliasID shared.UUID[organization.OrganizationAlias]

	testCases := []*txtest.TransactionalTestCase[TestData]{
		{
			Name: "編集しても組織と組織エイリアスの紐づけが残る",
			SetupFn: func(ctx context.Context, data *TestData) error {
				userRepo := repository.NewUserRepository(dbManager, repository.NewImageRepositoryMock(), encryptor)
				if err := userRepo.Create(ctx, user.NewUser(ownerUserID, lo.ToPtr("owner"), lo.ToPtr("owner"), "org-owner-"+ownerUserID.String(), "GOOGLE", nil)); err != nil {
					return err
				}
				if err := repository.NewOrganizationRepository(dbManager).Create(ctx, organization.NewOrganization(
					organizationID, organization.OrganizationTypeNormal, "テスト組織", "org"+organizationID.String()[:8], nil, ownerUserID,
				)); err != nil {
					return err
				}
				alias, err := organization.NewOrganizationAlias("テストエイリアス", organizationID, ownerUserID)
				if err != nil {
					return err
				}
				aliasID = alias.AliasID()
				if err := repository.NewOrganizationAliasRepository(dbManager).Create(ctx, alias); err != nil {
					return err
				}
				return data.TsRepo.Create(ctx, talksession.NewTalkSession(
					talkSessionID, "編集前のテーマ", nil, nil, ownerUserID,
					clock.Now(ctx), clock.Now(ctx).Add(time.Hour*24),
					nil, nil, nil, false,
					&organizationID, &aliasID, nil,
				))
			},
			TestFn: func(ctx context.Context, data *TestData) error {
				ts, err := data.TsRepo.FindByID(ctx, talkSessionID)
				if err != nil {
					return err
				}
				if ts.OrganizationID() == nil || *ts.OrganizationID() != organizationID {
					return errors.New("FindByID で組織IDが読み込まれていません")
				}
				if ts.OrganizationAliasID() == nil || *ts.OrganizationAliasID() != aliasID {
					return errors.New("FindByID で組織エイリアスIDが読み込まれていません")
				}

				ts.ChangeTheme("編集後のテーマ")
				if err := data.TsRepo.Update(ctx, ts); err != nil {
					return err
				}

				row, err := dbManager.GetQueries(ctx).GetTalkSessionByID(ctx, talkSessionID.UUID())
				if err != nil {
					return err
				}
				assert.Equal(t, "編集後のテーマ", row.TalkSession.Theme)
				assert.Equal(t, organizationID.UUID(), row.TalkSession.OrganizationID.UUID)
				assert.True(t, row.TalkSession.OrganizationID.Valid, "組織IDが消えています")
				assert.Equal(t, aliasID.UUID(), row.TalkSession.OrganizationAliasID.UUID)
				assert.True(t, row.TalkSession.OrganizationAliasID.Valid, "組織エイリアスIDが消えています")
				return nil
			},
		},
	}

	txtest.RunTransactionalTests(t, dbManager, &initData, testCases)
}
