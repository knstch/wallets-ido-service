package users_test

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	knlog "github.com/knstch/knstch-libs/log"
	"github.com/stretchr/testify/require"

	"wallets-service/internal/wallets"
)

func envOr(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustExtractState(t *require.Assertions, loginURL string) string {
	u, err := url.Parse(loginURL)
	t.NoError(err)
	state := u.Query().Get("state")
	t.NotEmpty(state)
	return state
}

func mustBase64Decode(t *require.Assertions, raw string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	t.NoError(err)
	return b
}

func mustDecodeOAuthState(t *require.Assertions, state string) wallets.OAuthState {
	var st wallets.OAuthState
	err := json.Unmarshal(mustBase64Decode(t, state), &st)
	t.NoError(err)
	return st
}

func mustEncodeOAuthState(t *require.Assertions, st wallets.OAuthState) string {
	b, err := json.Marshal(&st)
	t.NoError(err)
	return base64.RawURLEncoding.EncodeToString(b)
}

func makeFakeIDToken(t *require.Assertions, claims map[string]any) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payloadBytes, err := json.Marshal(claims)
	t.NoError(err)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payload + ".sig"
}

func mustFindRepoRoot(t *require.Assertions) string {
	_, thisFile, _, ok := runtime.Caller(0)
	t.True(ok)
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.FailNow("could not find repo root (go.mod)")
	return ""
}

func newTestLogger(serviceName string) *knlog.Logger {
	if serviceName == "" {
		serviceName = "test"
	}
	return knlog.NewLogger(serviceName, knlog.InfoLevel)
}
