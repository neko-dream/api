package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/neko-dream/api/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/session"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	"github.com/neko-dream/api/internal/presentation/oas"
	httputils "github.com/neko-dream/api/pkg/http"
)

type creationErrorCommand struct {
	err   error
	calls int
}

func (s *creationErrorCommand) Execute(context.Context, talksession_usecase.StartTalkSessionUseCaseInput) (talksession_usecase.StartTalkSessionUseCaseOutput, error) {
	s.calls++
	return talksession_usecase.StartTalkSessionUseCaseOutput{}, s.err
}

func TestCreationHTTPSeparatesAuthenticationAndAuthorizationErrors(t *testing.T) {
	id := shared.NewUUID[user.User]()
	sid := shared.NewUUID[session.Session]()
	claim := &session.Claim{Sub: id.String(), Jti: sid.String(), IsRegistered: true, Exp: time.Now().Add(time.Hour).Unix()}
	for _, tc := range []struct {
		name   string
		cookie string
		want   *messages.APIError
	}{
		{"unauthenticated", "", messages.ForbiddenError},
		{"not granted or revoked", "SessionId=test-user-session", messages.TalkSessionCreationNotAllowed},
		{"not an organization member", "SessionId=test-user-session", messages.TalkSessionOrganizationMembershipRequired},
		{"guest", "SessionId=test-user-session", messages.GuestCannotCreateSessionError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			command := &creationErrorCommand{err: tc.want}
			srv, err := oas.NewServer(&handlers{TalkSessionHandler: &talkSessionHandler{authorizationService: currentAuthStub{id: id}, startTalkSessionCommand: command}}, NewSecurityHandler(meSessionToken{claim: claim}, meSessions{value: session.NewSession(sid, id, shared.ProviderGoogle, session.SESSION_ACTIVE, time.Now().Add(time.Hour), time.Now())}, nil, nil), oas.WithErrorHandler(CustomErrorHandler))
			if err != nil {
				t.Fatal(err)
			}
			var body bytes.Buffer
			form := multipart.NewWriter(&body)
			for key, value := range map[string]string{"theme": "creation permission", "scheduledEndTime": "2100-12-31T23:59:59Z"} {
				if err := form.WriteField(key, value); err != nil {
					t.Fatal(err)
				}
			}
			if err := form.Close(); err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRequest(http.MethodPost, "/talksessions", &body)
			r.Header.Set("Content-Type", form.FormDataContentType())
			r.Header.Set("Cookie", tc.cookie)
			w := httptest.NewRecorder()
			r = r.WithContext(httputils.WithHTTPResReqContext(r.Context(), r, w))
			srv.ServeHTTP(w, r)
			var got messages.APIError
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if w.Code != http.StatusForbidden || got.Code != tc.want.Code || got.Message != tc.want.Message {
				t.Fatalf("status=%d body=%s want=%+v", w.Code, w.Body, tc.want)
			}
			wantCalls := 1
			if tc.cookie == "" {
				wantCalls = 0
			}
			if command.calls != wantCalls {
				t.Fatalf("usecase calls=%d want=%d", command.calls, wantCalls)
			}
		})
	}
}
