package vertexai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// gcpAuth manages GCP access tokens from Application Default Credentials.
// Tokens are cached and auto-refreshed.
type gcpAuth struct {
	mu          sync.Mutex
	tokenSource oauth2.TokenSource
}

// newGCPAuth creates a new GCP auth helper using Application Default Credentials.
func newGCPAuth(ctx context.Context) (*gcpAuth, error) {
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("vertex-ai: failed to find GCP credentials (run 'gcloud auth application-default login'): %w", err)
	}

	return &gcpAuth{
		tokenSource: creds.TokenSource,
	}, nil
}

// token returns a valid access token, refreshing if needed.
func (a *gcpAuth) token(ctx context.Context) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	tok, err := a.tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("vertex-ai: failed to get access token: %w", err)
	}

	if tok.Expiry.Before(time.Now().Add(30 * time.Second)) {
		// Token is about to expire, force refresh
		tok, err = a.tokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("vertex-ai: failed to refresh access token: %w", err)
		}
	}

	return tok.AccessToken, nil
}
