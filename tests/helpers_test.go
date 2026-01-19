package wallets_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	knlog "github.com/knstch/knstch-libs/log"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"

	"github.com/knstch/knstch-libs/svcerrs"
)

func envOr(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustBase64Decode(t *require.Assertions, raw string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	t.NoError(err)
	return b
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

func mustGenerateSolanaKeypair(t *require.Assertions) (pubkeyBase58 string, priv ed25519.PrivateKey) {
	// Solana pubkeys are base58-encoded 32-byte ed25519 public keys.
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	t.NoError(err)
	return base58.Encode(pub), priv
}

func mustSignBase64(priv ed25519.PrivateKey, msg string) string {
	sig := ed25519.Sign(priv, []byte(msg))
	return base64.StdEncoding.EncodeToString(sig)
}

func requireSvcErrIs(t *testing.T, err error, target error) {
	t.Helper()
	switch target {
	case svcerrs.ErrDataNotFound, svcerrs.ErrConflict, svcerrs.ErrInvalidData:
		// ok
	default:
		// also ok; keep generic
	}
	require.Error(t, err)
	require.ErrorIs(t, err, target)
}
