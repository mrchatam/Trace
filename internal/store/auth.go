package store

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AccessTokenEnv is the environment variable Open checks when an access token
// file is present under .trace/.
const AccessTokenEnv = "TRACE_ACCESS_TOKEN"

const accessTokenFileName = "access.token"

// ErrUnauthorized is returned by Open when .trace/access.token is non-empty
// and TRACE_ACCESS_TOKEN is missing or does not match.
var ErrUnauthorized = errors.New("store: unauthorized (set TRACE_ACCESS_TOKEN to match .trace/access.token)")

// SetAccessToken writes <projectRoot>/.trace/access.token with mode 0600.
// Does not require Open (so the first set works without a matching env).
func SetAccessToken(projectRoot, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("store: access token must be non-empty")
	}
	absRoot, traceDir, err := resolveTraceDir(projectRoot)
	if err != nil {
		return err
	}
	_ = absRoot
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return fmt.Errorf("store: mkdir .trace: %w", err)
	}
	path := filepath.Join(traceDir, accessTokenFileName)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(token+"\n"), 0o600); err != nil {
		return fmt.Errorf("store: write access token: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("store: install access token: %w", err)
	}
	return nil
}

// ClearAccessToken removes .trace/access.token if present.
func ClearAccessToken(projectRoot string) error {
	_, traceDir, err := resolveTraceDir(projectRoot)
	if err != nil {
		return err
	}
	path := filepath.Join(traceDir, accessTokenFileName)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("store: clear access token: %w", err)
	}
	return nil
}

// AccessTokenEnabled reports whether a non-empty .trace/access.token exists.
func AccessTokenEnabled(projectRoot string) (bool, error) {
	_, traceDir, err := resolveTraceDir(projectRoot)
	if err != nil {
		return false, err
	}
	tok, err := readAccessTokenFile(filepath.Join(traceDir, accessTokenFileName))
	if err != nil {
		return false, err
	}
	return tok != "", nil
}

func resolveTraceDir(projectRoot string) (absRoot, traceDir string, err error) {
	absRoot, err = filepath.Abs(projectRoot)
	if err != nil {
		return "", "", fmt.Errorf("store: resolve project root: %w", err)
	}
	return absRoot, filepath.Join(absRoot, traceDirName), nil
}

func readAccessTokenFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("store: read access token: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

// checkAccessToken gates Open: if the token file is non-empty, TRACE_ACCESS_TOKEN
// must match via constant-time compare.
func checkAccessToken(traceDir string) error {
	tok, err := readAccessTokenFile(filepath.Join(traceDir, accessTokenFileName))
	if err != nil {
		return err
	}
	if tok == "" {
		return nil
	}
	env := os.Getenv(AccessTokenEnv)
	if subtle.ConstantTimeCompare([]byte(env), []byte(tok)) != 1 {
		return ErrUnauthorized
	}
	return nil
}
