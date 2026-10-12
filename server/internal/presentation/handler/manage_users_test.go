package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userquery "github.com/neko-dream/api/internal/application/query/user"
	"github.com/neko-dream/api/internal/domain/messages"
	"github.com/neko-dream/api/internal/domain/model/admin"
	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/user"
	adminauth "github.com/neko-dream/api/internal/infrastructure/auth/admin"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/internal/presentation/oas"
	http_utils "github.com/neko-dream/api/pkg/http"
)

func TestManagedUsersHTTPAuthorization(t *testing.T) {
	id := shared.NewUUID[admin.AdminUser]()
	mgr := adminauth.NewAdminTokenManager(&config.Config{AdminJWTSecret: "operator-http-test-signing-secret-32bytes"})
	token, err := mgr.Generate(context.Background(), admin.AdminUser{AdminUserID: id})
	if err != nil {
		t.Fatal(err)
	}
	q := &managedUsersStub{output: &userquery.ManagementListOutput{}}
	srv, err := oas.NewServer(&handlers{ManageHandler: &manageHandler{userManagement: q}}, NewSecurityHandler(nil, nil, mgr, adminSecurityUsers{}), oas.WithErrorHandler(CustomErrorHandler))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, path, cookie string
		status             int
	}{
		{"anonymous", "/v1/manage/users/list", "", 403},
		{"ordinary user cookie", "/v1/manage/users/list", "SessionId=" + token, 403},
		{"organization user cookie", "/v1/manage/users/list", "SessionId=organization-session", 403},
		{"tampered admin cookie", "/v1/manage/users/list", "AdminSessionId=invalid", 403},
		{"operator list", "/v1/manage/users/list", "AdminSessionId=" + token, 200},
		{"invalid pagination", "/v1/manage/users/list?limit=101", "AdminSessionId=" + token, 400},
		{"invalid user ID", "/v1/manage/users/invalid", "AdminSessionId=" + token, 400},
		{"missing user", "/v1/manage/users/" + shared.NewUUID[user.User]().String(), "AdminSessionId=" + token, 404},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.cookie != "" {
				r.Header.Set("Cookie", tc.cookie)
			}
			w := httptest.NewRecorder()
			r = r.WithContext(http_utils.WithHTTPResReqContext(r.Context(), r, w))
			srv.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status=%d want=%d body=%s", w.Code, tc.status, w.Body)
			}
			for _, secret := range []string{"password_hash", "token_hash", "subject", "salt"} {
				if strings.Contains(w.Body.String(), secret) {
					t.Fatal("credential field exposed")
				}
			}
		})
	}
}

type managedUsersStub struct {
	input  userquery.ManagementListInput
	output *userquery.ManagementListOutput
	user   *userquery.ManagedUser
	err    error
	calls  int
}

func (q *managedUsersStub) List(_ context.Context, input userquery.ManagementListInput) (*userquery.ManagementListOutput, error) {
	q.calls++
	q.input = input
	return q.output, q.err
}

func (q *managedUsersStub) FindByID(context.Context, shared.UUID[user.User]) (*userquery.ManagedUser, error) {
	q.calls++
	return q.user, q.err
}

func TestManagedUsersRequireAdmin(t *testing.T) {
	q := &managedUsersStub{}
	h := &manageHandler{userManagement: q}
	if _, err := h.GetUserListManage(context.Background(), oas.GetUserListManageParams{}); !errors.Is(err, messages.ForbiddenError) {
		t.Fatal(err)
	}
	if _, err := h.GetUserDetailManage(context.Background(), oas.GetUserDetailManageParams{}); !errors.Is(err, messages.ForbiddenError) {
		t.Fatal(err)
	}
	if q.calls != 0 {
		t.Fatal("unauthenticated request reached user data")
	}
}

func TestManagedUserList(t *testing.T) {
	ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{})
	q := &managedUsersStub{output: &userquery.ManagementListOutput{TotalCount: 7}}
	h := &manageHandler{userManagement: q}
	res, err := h.GetUserListManage(ctx, oas.GetUserListManageParams{Search: oas.NewOptString("  city  ")})
	if err != nil || q.input.Limit != 20 || q.input.Offset != 0 || q.input.Search != "city" || res.TotalCount != 7 || res.Users == nil {
		t.Fatalf("unexpected list %v %v %+v", res, err, q.input)
	}
	for _, params := range []oas.GetUserListManageParams{{Limit: oas.NewOptInt32(0)}, {Limit: oas.NewOptInt32(101)}, {Offset: oas.NewOptInt32(-1)}} {
		if _, err := h.GetUserListManage(ctx, params); !errors.Is(err, messages.BadRequestError) {
			t.Fatal("invalid pagination accepted")
		}
	}
	q.err = errors.New("database failed")
	if _, err := h.GetUserListManage(ctx, oas.GetUserListManageParams{}); !errors.Is(err, q.err) {
		t.Fatal("database failure hidden")
	}
}

func TestManagedUserDetail(t *testing.T) {
	ctx := admin.SetAdminSession(context.Background(), &admin.AdminClaim{})
	q := &managedUsersStub{}
	h := &manageHandler{userManagement: q}
	if _, err := h.GetUserDetailManage(ctx, oas.GetUserDetailManageParams{UserID: "bad-id"}); !errors.Is(err, messages.BadRequestError) {
		t.Fatal("invalid UUID accepted")
	}
	id := shared.NewUUID[user.User]()
	res, err := h.GetUserDetailManage(ctx, oas.GetUserDetailManageParams{UserID: id.String()})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.(*oas.GetUserDetailManageNotFound); !ok {
		t.Fatalf("got %T", res)
	}
	email := "city@example.com"
	q.user = &userquery.ManagedUser{UserID: id, Email: &email, Providers: []string{"GOOGLE", "PASSWORD"}}
	res, err = h.GetUserDetailManage(ctx, oas.GetUserDetailManageParams{UserID: id.String()})
	if err != nil {
		t.Fatal(err)
	}
	if u := res.(*oas.ManagedUser); u.UserID != id.String() || u.Email.Or("") != email || len(u.Providers) != 2 {
		t.Fatal("detail projection mismatch")
	}
	q.err = errors.New("database failed")
	if _, err := h.GetUserDetailManage(ctx, oas.GetUserDetailManageParams{UserID: id.String()}); !errors.Is(err, q.err) {
		t.Fatal("database failure became 404")
	}
}
