package opinion_command

import (
	"context"
	"mime/multipart"

	"github.com/neko-dream/server/internal/domain/messages"
	"github.com/neko-dream/server/internal/domain/model/clock"
	"github.com/neko-dream/server/internal/domain/model/image"
	"github.com/neko-dream/server/internal/domain/model/image/meta"
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/talksession"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/neko-dream/server/internal/domain/service"
	"github.com/neko-dream/server/internal/infrastructure/persistence/db"
	"github.com/neko-dream/server/pkg/utils"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
)

type (
	SubmitOpinion interface {
		Execute(context.Context, SubmitOpinionInput) error
	}

	SubmitOpinionInput struct {
		// TalkSessionID
		TalkSessionID *shared.UUID[talksession.TalkSession]
		// ParentOpinionID 親意見IDがない場合はTalkSessionIDが必須
		ParentOpinionID *shared.UUID[opinion.Opinion]
		OwnerID         shared.UUID[user.User]
		UserID          shared.UUID[user.User]
		Title           *string
		Content         string
		ReferenceURL    *string
		Picture         *multipart.FileHeader
	}

	submitOpinionHandler struct {
		opinion.OpinionRepository
		opinion.OpinionService
		vote.VoteRepository
		service.TalkSessionAccessControl
		image.ImageStorage
		image.ImageRepository
		*db.DBManager
	}
)

func NewSubmitOpinionHandler(
	opinionRepository opinion.OpinionRepository,
	opinionService opinion.OpinionService,
	voteRepository vote.VoteRepository,
	talkSessionAccessControl service.TalkSessionAccessControl,
	dbManager *db.DBManager,
	imageRepository image.ImageRepository,
	imageStorage image.ImageStorage,
) SubmitOpinion {
	return &submitOpinionHandler{
		DBManager:                dbManager,
		OpinionService:           opinionService,
		OpinionRepository:        opinionRepository,
		VoteRepository:           voteRepository,
		TalkSessionAccessControl: talkSessionAccessControl,
		ImageStorage:             imageStorage,
		ImageRepository:          imageRepository,
	}
}

func (h *submitOpinionHandler) Execute(ctx context.Context, input SubmitOpinionInput) error {
	ctx, span := otel.Tracer("opinion_command").Start(ctx, "submitOpinionHandler.Execute")
	defer span.End()

	// TalkSessionIDが指定されているのならば、そのTalkSessionを取得する
	// ParentOpinionIDが指定されている場合は、その意見のTalkSessionを取得する
	// どちらも指定されていない場合はエラーを返す
	var talkSessionID shared.UUID[talksession.TalkSession]
	if input.TalkSessionID != nil {
		talkSessionID = *input.TalkSessionID
	} else if input.ParentOpinionID != nil {
		parentOpinion, err := h.OpinionRepository.FindByID(ctx, *input.ParentOpinionID)
		if err != nil {
			utils.HandleError(ctx, err, "OpinionRepository.FindByID")
			return messages.OpinionCreateFailed
		}
		talkSessionID = parentOpinion.TalkSessionID()
	} else {
		return messages.OpinionCreateFailed
	}

	// 参加制限を満たしているか確認。満たしていない場合はエラーを返す
	if _, err := h.TalkSessionAccessControl.CanUserJoin(ctx, talkSessionID, lo.ToPtr(input.UserID)); err != nil {
		utils.HandleError(ctx, err, "TalkSessionAccessControl.CanUserJoin")
		return err
	}

	if err := h.ExecTx(ctx, func(ctx context.Context) error {
		opinion, err := opinion.NewOpinion(
			shared.NewUUID[opinion.Opinion](),
			talkSessionID,
			input.OwnerID,
			input.ParentOpinionID,
			input.Title,
			input.Content,
			clock.Now(ctx),
			input.ReferenceURL,
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewOpinion")
			return err
		}
		if input.Picture != nil {
			file, err := input.Picture.Open()
			if err != nil {
				utils.HandleError(ctx, err, "input.Icon.Open")
				return messages.OpinionReferenceImageUploadFailed
			}
			defer file.Close()

			imageMeta, err := meta.NewImageForReference(ctx, opinion.OpinionID(), file)
			if err != nil {
				utils.HandleError(ctx, err, "meta.NewImageForProfile")
				return messages.OpinionReferenceImageUploadFailed
			}
			if err := imageMeta.Validate(ctx, meta.ReferenceImageValidationRule); err != nil {
				utils.HandleError(ctx, err, "ImageMeta.Validate")
				msg := messages.OpinionReferenceImageUploadFailed
				msg.Message = err.Error()
				return msg
			}

			// 画像をアップロード
			url, err := h.ImageStorage.Upload(ctx, *imageMeta, input.Picture)
			if err != nil {
				utils.HandleError(ctx, err, "ImageRepository.Upload")
				return messages.OpinionReferenceImageUploadFailed
			}
			if err := h.ImageRepository.Create(ctx, image.NewUserImage(
				ctx,
				shared.NewUUID[image.UserImage](),
				input.OwnerID,
				*imageMeta,
				*url,
			)); err != nil {
				utils.HandleError(ctx, err, "ImageRepository.Create")
				return messages.OpinionReferenceImageUploadFailed
			}

			opinion.ChangeReferenceImageURL(url)
		}

		if err := h.OpinionRepository.Create(ctx, *opinion); err != nil {
			utils.HandleError(ctx, err, "OpinionRepository.Create")
			return messages.OpinionCreateFailed
		}

		// 自分の意見には必ず投票を紐付ける
		v, err := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			opinion.OpinionID(),
			talkSessionID,
			input.OwnerID,
			vote.Agree,
			clock.Now(ctx),
		)
		if err != nil {
			utils.HandleError(ctx, err, "NewVote")
			return err
		}
		if err := h.VoteRepository.Create(ctx, *v); err != nil {
			return messages.VoteFailed
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
