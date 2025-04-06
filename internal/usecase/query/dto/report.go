package dto

import (
	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
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
