// Code generated by ogen, DO NOT EDIT.

package oas

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