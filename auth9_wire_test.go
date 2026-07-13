package run9

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// The structs below mirror the portal's exact serialization types for the
// auth9 endpoints (portal/api/auth_auth9.go, identity_model_views.go). They
// pin the generated SDK models to the real wire shape produced by the portal.

type portalMeView struct {
	UserID             string    `json:"user_id"`
	PrimaryEmail       string    `json:"primary_email"`
	DisplayName        string    `json:"display_name,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	IsSystemManager    bool      `json:"is_system_manager"`
	CanCreateSharedOrg bool      `json:"can_create_shared_org"`
}

type portalAuth9SessionView struct {
	User             portalMeView `json:"user"`
	SessionToken     string       `json:"session_token"`
	SessionExpiresAt time.Time    `json:"session_expires_at"`
	RefreshToken     string       `json:"refresh_token"`
	RefreshExpiresAt *time.Time   `json:"refresh_expires_at,omitempty"`
}

type portalAuth9ConfigView struct {
	Enabled      bool   `json:"enabled"`
	AuthorizeURL string `json:"authorize_url,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
}

type portalAuth9SignInRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	RedirectURI  string `json:"redirect_uri"`
}

type portalAuth9RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func portalSessionFixture(t *testing.T) (portalAuth9SessionView, []byte) {
	t.Helper()
	refreshExpiry := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	view := portalAuth9SessionView{
		User: portalMeView{
			UserID:             "user-123",
			PrimaryEmail:       "dev@example.com",
			DisplayName:        "Dev",
			CreatedAt:          time.Date(2026, 7, 1, 8, 30, 0, 0, time.UTC),
			IsSystemManager:    true,
			CanCreateSharedOrg: true,
		},
		SessionToken:     "session-jwt",
		SessionExpiresAt: time.Date(2026, 7, 13, 12, 15, 0, 0, time.UTC),
		RefreshToken:     "refresh-1",
		RefreshExpiresAt: &refreshExpiry,
	}
	raw, err := json.Marshal(view)
	require.NoError(t, err)
	return view, raw
}

func TestAuth9SessionViewDecodesPortalWire(t *testing.T) {
	_, raw := portalSessionFixture(t)

	var got genmodels.APIAuth9SessionView
	require.NoError(t, got.UnmarshalBinary(raw))
	require.NotNil(t, got.User)
	require.Equal(t, "user-123", got.User.UserID)
	require.Equal(t, "dev@example.com", got.User.PrimaryEmail)
	require.Equal(t, "Dev", got.User.DisplayName)
	require.Equal(t, "2026-07-01T08:30:00Z", got.User.CreatedAt)
	require.True(t, got.User.IsSystemManager)
	require.True(t, got.User.CanCreateSharedOrg)
	require.Equal(t, "session-jwt", got.SessionToken)
	require.Equal(t, "2026-07-13T12:15:00Z", got.SessionExpiresAt)
	require.Equal(t, "refresh-1", got.RefreshToken)
	require.Equal(t, "2026-08-10T12:00:00Z", got.RefreshExpiresAt)
	require.NoError(t, got.Validate(strfmt.Default))

	// The refresh endpoint omits refresh_expires_at (nil pointer on the
	// portal side); the SDK must decode that as an empty string.
	refreshWire, err := json.Marshal(portalAuth9SessionView{
		User:             portalMeView{UserID: "user-123", PrimaryEmail: "dev@example.com", CreatedAt: time.Date(2026, 7, 1, 8, 30, 0, 0, time.UTC)},
		SessionToken:     "session-jwt-2",
		SessionExpiresAt: time.Date(2026, 7, 13, 12, 30, 0, 0, time.UTC),
		RefreshToken:     "refresh-2",
	})
	require.NoError(t, err)
	var rotated genmodels.APIAuth9SessionView
	require.NoError(t, rotated.UnmarshalBinary(refreshWire))
	require.Empty(t, rotated.RefreshExpiresAt)
	require.Equal(t, "refresh-2", rotated.RefreshToken)
}

func TestAuth9ConfigViewDecodesPortalWire(t *testing.T) {
	enabledWire, err := json.Marshal(portalAuth9ConfigView{
		Enabled:      true,
		AuthorizeURL: "https://auth9.example.com/v1/oauth/authorize",
		ClientID:     "client-1",
		RedirectURI:  "https://run9.example.com/auth/auth9/callback",
	})
	require.NoError(t, err)
	var enabled genmodels.APIAuth9ConfigView
	require.NoError(t, enabled.UnmarshalBinary(enabledWire))
	require.True(t, enabled.Enabled)
	require.Equal(t, "https://auth9.example.com/v1/oauth/authorize", enabled.AuthorizeURL)
	require.Equal(t, "client-1", enabled.ClientID)
	require.Equal(t, "https://run9.example.com/auth/auth9/callback", enabled.RedirectURI)

	// Disabled deployments serve the bare {"enabled":false} bootstrap.
	var disabled genmodels.APIAuth9ConfigView
	require.NoError(t, disabled.UnmarshalBinary([]byte(`{"enabled":false}`)))
	require.False(t, disabled.Enabled)
	require.Empty(t, disabled.AuthorizeURL)
	require.Empty(t, disabled.ClientID)
	require.Empty(t, disabled.RedirectURI)
}

func TestAuth9RequestsEncodePortalWire(t *testing.T) {
	signIn := genmodels.APIAuth9SignInRequest{
		Code:         "authz-code",
		CodeVerifier: "pkce-verifier",
		RedirectURI:  "https://run9.example.com/auth/auth9/callback",
	}
	raw, err := signIn.MarshalBinary()
	require.NoError(t, err)
	var gotSignIn portalAuth9SignInRequest
	require.NoError(t, json.Unmarshal(raw, &gotSignIn))
	require.Equal(t, "authz-code", gotSignIn.Code)
	require.Equal(t, "pkce-verifier", gotSignIn.CodeVerifier)
	require.Equal(t, "https://run9.example.com/auth/auth9/callback", gotSignIn.RedirectURI)

	refresh := genmodels.APIAuth9RefreshRequest{RefreshToken: "refresh-1"}
	raw, err = refresh.MarshalBinary()
	require.NoError(t, err)
	var gotRefresh portalAuth9RefreshRequest
	require.NoError(t, json.Unmarshal(raw, &gotRefresh))
	require.Equal(t, "refresh-1", gotRefresh.RefreshToken)
}

func TestAuth9ModelsRejectMalformedResponses(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		decode  func([]byte) error
	}{
		{
			name:    "session token wrong type",
			payload: `{"session_token":123}`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9SessionView; return m.UnmarshalBinary(b) },
		},
		{
			name:    "user wrong type",
			payload: `{"user":"not-an-object"}`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9SessionView; return m.UnmarshalBinary(b) },
		},
		{
			name:    "top level array",
			payload: `[{"session_token":"x"}]`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9SessionView; return m.UnmarshalBinary(b) },
		},
		{
			name:    "truncated json",
			payload: `{"session_token":"x"`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9SessionView; return m.UnmarshalBinary(b) },
		},
		{
			name:    "enabled wrong type",
			payload: `{"enabled":"yes"}`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9ConfigView; return m.UnmarshalBinary(b) },
		},
		{
			name:    "refresh token wrong type",
			payload: `{"refresh_token":{}}`,
			decode:  func(b []byte) error { var m genmodels.APIAuth9RefreshRequest; return m.UnmarshalBinary(b) },
		},
		{
			name:    "empty body",
			payload: ``,
			decode:  func(b []byte) error { var m genmodels.APIAuth9SessionView; return m.UnmarshalBinary(b) },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				require.Error(t, tc.decode([]byte(tc.payload)))
			})
		})
	}
}

func TestAuth9SessionViewHugePayload(t *testing.T) {
	hugeToken := strings.Repeat("a", 8<<20)
	raw, err := json.Marshal(map[string]any{
		"session_token": hugeToken,
		"refresh_token": "refresh-1",
	})
	require.NoError(t, err)

	var got genmodels.APIAuth9SessionView
	require.NoError(t, got.UnmarshalBinary(raw))
	require.Len(t, got.SessionToken, 8<<20)
	require.Equal(t, "refresh-1", got.RefreshToken)
}

// TestAuth9EndpointsConcurrentAgainstPortalStub drives the three auth9
// portal routes from many goroutines sharing one HTTP client, decoding every
// response through the generated models. Run with -race it also stresses the
// shared strfmt.Default registry used by Validate.
func TestAuth9EndpointsConcurrentAgainstPortalStub(t *testing.T) {
	_, sessionWire := portalSessionFixture(t)
	configWire, err := json.Marshal(portalAuth9ConfigView{
		Enabled:      true,
		AuthorizeURL: "https://auth9.example.com/v1/oauth/authorize",
		ClientID:     "client-1",
		RedirectURI:  "https://run9.example.com/auth/auth9/callback",
	})
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth/auth9/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(configWire)
	})
	writeSessionOr400 := func(w http.ResponseWriter, r *http.Request, requiredField string) {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body[requiredField] == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing " + requiredField})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(sessionWire)
	}
	mux.HandleFunc("POST /auth/auth9/exchange", func(w http.ResponseWriter, r *http.Request) {
		writeSessionOr400(w, r, "code")
	})
	mux.HandleFunc("POST /auth/auth9/refresh", func(w http.ResponseWriter, r *http.Request) {
		writeSessionOr400(w, r, "refresh_token")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	post := func(path string, reqBody []byte) (int, []byte, error) {
		resp, err := client.Post(server.URL+path, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return 0, nil, err
		}
		return resp.StatusCode, buf.Bytes(), nil
	}

	const goroutines = 24
	const iterations = 25
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				resp, err := client.Get(server.URL + "/auth/auth9/config")
				if err != nil {
					t.Error(err)
					return
				}
				var cfg genmodels.APIAuth9ConfigView
				decodeErr := json.NewDecoder(resp.Body).Decode(&cfg)
				resp.Body.Close()
				if decodeErr != nil || !cfg.Enabled || cfg.ClientID != "client-1" {
					t.Errorf("bad config response: err=%v view=%+v", decodeErr, cfg)
					return
				}

				signIn := genmodels.APIAuth9SignInRequest{Code: "authz-code", CodeVerifier: "pkce", RedirectURI: "https://run9.example.com/auth/auth9/callback"}
				reqRaw, _ := signIn.MarshalBinary()
				status, raw, err := post("/auth/auth9/exchange", reqRaw)
				if err != nil {
					t.Error(err)
					return
				}
				var session genmodels.APIAuth9SessionView
				if status != http.StatusOK || session.UnmarshalBinary(raw) != nil {
					t.Errorf("bad exchange response: status=%d body=%s", status, raw)
					return
				}
				if session.SessionToken != "session-jwt" || session.User == nil || session.User.UserID != "user-123" {
					t.Errorf("unexpected session view: %+v", session)
					return
				}
				if err := session.Validate(strfmt.Default); err != nil {
					t.Errorf("session validate: %v", err)
					return
				}

				refresh := genmodels.APIAuth9RefreshRequest{RefreshToken: "refresh-1"}
				reqRaw, _ = refresh.MarshalBinary()
				status, raw, err = post("/auth/auth9/refresh", reqRaw)
				if err != nil || status != http.StatusOK {
					t.Errorf("bad refresh response: err=%v status=%d body=%s", err, status, raw)
					return
				}

				// Empty refresh token must surface the portal's Error shape.
				status, raw, err = post("/auth/auth9/refresh", []byte(`{}`))
				if err != nil {
					t.Error(err)
					return
				}
				var apiErr genmodels.Error
				if status != http.StatusBadRequest || apiErr.UnmarshalBinary(raw) != nil || apiErr.Error != "missing refresh_token" {
					t.Errorf("bad refresh error response: status=%d body=%s", status, raw)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestAuth9StubCancellationPropagates(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/auth/auth9/config", nil)
	require.NoError(t, err)
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err = server.Client().Do(req)
	require.ErrorIs(t, err, context.Canceled)

	// ContextValidate on a canceled context must not hang or panic.
	view := genmodels.APIAuth9SessionView{User: &genmodels.APIMeView{UserID: "user-123"}}
	canceled, cancelNow := context.WithCancel(context.Background())
	cancelNow()
	require.NotPanics(t, func() {
		_ = view.ContextValidate(canceled, strfmt.Default)
	})
}
