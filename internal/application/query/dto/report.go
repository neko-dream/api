package dto

import (
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/presentation/oas"
	"github.com/neko-dream/server/pkg/utils"
)

type ReportDetail struct {
	Opinion     Opinion
	User        User
	Reasons     []ReportDetailReason
	ReportCount int
	Status      string
}

type ReportDetailReason struct {
	ReportID shared.UUID[opinion.Report]
	Reason   string
	Content  *string
}

// ToResponse converts ReportDetail DTO to OAS ReportDetail response
func (r *ReportDetail) ToResponse() oas.ReportDetail {
	// Convert reasons
	reasons := make([]oas.ReportDetailReasonsItem, 0, len(r.Reasons))
	for _, reason := range r.Reasons {
		reasons = append(reasons, oas.ReportDetailReasonsItem{
			Reason:  reason.Reason,
			Content: utils.ToOptNil[oas.OptNilString](reason.Content),
		})
	}

	return oas.ReportDetail{
		Opinion: r.Opinion.ToResponse(),
		User: oas.User{
			DisplayID:   r.User.DisplayID,
			DisplayName: r.User.DisplayName,
			IconURL:     utils.ToOptNil[oas.OptNilString](r.User.IconURL),
		},
		Status:      oas.ReportStatus(r.Status),
		Reasons:     reasons,
		ReportCount: r.ReportCount,
	}
}
