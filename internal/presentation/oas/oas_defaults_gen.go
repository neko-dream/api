// Code generated by ogen, DO NOT EDIT.

package oas

// setDefaults set default value of fields.
func (s *EditUserProfileReq) setDefaults() {
	{
		val := bool(false)
		s.DeleteIcon.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *GetOpinionDetailOKOpinion) setDefaults() {
	{
		val := GetOpinionDetailOKOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *GetOpinionsForTalkSessionOKOpinionsItem) setDefaults() {
	{
		val := GetOpinionsForTalkSessionOKOpinionsItemMyVoteType("unvote")
		s.MyVoteType = val
	}
}

// setDefaults set default value of fields.
func (s *GetOpinionsForTalkSessionOKOpinionsItemOpinion) setDefaults() {
	{
		val := GetOpinionsForTalkSessionOKOpinionsItemOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *OpinionCommentsOKOpinionsItem) setDefaults() {
	{
		val := OpinionCommentsOKOpinionsItemMyVoteType("unvote")
		s.MyVoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *OpinionCommentsOKOpinionsItemOpinion) setDefaults() {
	{
		val := OpinionCommentsOKOpinionsItemOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *OpinionCommentsOKRootOpinion) setDefaults() {
	{
		val := OpinionCommentsOKRootOpinionMyVoteType("unvote")
		s.MyVoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *OpinionCommentsOKRootOpinionOpinion) setDefaults() {
	{
		val := OpinionCommentsOKRootOpinionOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *OpinionsHistoryOKOpinionsItemOpinion) setDefaults() {
	{
		val := OpinionsHistoryOKOpinionsItemOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *RegisterUserReq) setDefaults() {
	{
		val := int(0)
		s.YearOfBirth.SetTo(val)
	}
	{
		val := RegisterUserReqGender("preferNotToSay")
		s.Gender.SetTo(val)
	}
	{
		val := RegisterUserReqOccupation("無回答")
		s.Occupation.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *SwipeOpinionsOKItemOpinion) setDefaults() {
	{
		val := SwipeOpinionsOKItemOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *TalkSessionAnalysisOKGroupOpinionsItemOpinionsItemOpinion) setDefaults() {
	{
		val := TalkSessionAnalysisOKGroupOpinionsItemOpinionsItemOpinionVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}

// setDefaults set default value of fields.
func (s *VoteOKItem) setDefaults() {
	{
		val := VoteOKItemVoteType("unvote")
		s.VoteType.SetTo(val)
	}
}
