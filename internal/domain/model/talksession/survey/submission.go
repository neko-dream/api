package survey

import (
	"time"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
)

type Submission struct {
	submissionID shared.UUID[Submission]
	surveyID     shared.UUID[Survey]
	userID       shared.UUID[user.User]
	submittedAt  time.Time
	answers      []*Answer
}

func NewSubmission(
	submissionID shared.UUID[Submission],
	surveyID shared.UUID[Survey],
	userID shared.UUID[user.User],
	answers []*Answer,
) *Submission {
	return &Submission{
		submissionID: submissionID,
		surveyID:     surveyID,
		userID:       userID,
		submittedAt:  time.Now(),
		answers:      answers,
	}
}

func ReconstructSubmission(
	submissionID shared.UUID[Submission],
	surveyID shared.UUID[Survey],
	userID shared.UUID[user.User],
	submittedAt time.Time,
	answers []*Answer,
) *Submission {
	return &Submission{
		submissionID: submissionID,
		surveyID:     surveyID,
		userID:       userID,
		submittedAt:  submittedAt,
		answers:      answers,
	}
}

func (s *Submission) SubmissionID() shared.UUID[Submission] { return s.submissionID }
func (s *Submission) SurveyID() shared.UUID[Survey]         { return s.surveyID }
func (s *Submission) UserID() shared.UUID[user.User]        { return s.userID }
func (s *Submission) SubmittedAt() time.Time                { return s.submittedAt }
func (s *Submission) Answers() []*Answer                    { return s.answers }
