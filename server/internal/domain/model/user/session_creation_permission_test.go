package user

import "testing"

func TestCanCreateTalkSession(t *testing.T) {
	for _, tc := range []struct{ granted, guest, withdrawn, want bool }{
		{true, false, false, true}, {false, false, false, false},
		{true, true, false, false}, {true, false, true, false}, {true, true, true, false},
	} {
		if got := CanCreateTalkSession(tc.granted, tc.guest, tc.withdrawn); got != tc.want {
			t.Fatalf("%+v got=%v", tc, got)
		}
	}
}
