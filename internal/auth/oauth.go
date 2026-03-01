package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
)

const (
	// OAuthRedirectURI is the redirect URI for the OAuth flow.
	OAuthRedirectURI = "http://localhost:8080/callback"
	// OAuthPort is the local server port for the OAuth callback.
	OAuthPort = ":8080"
)

// OAuthFlow handles the OAuth authorization flow.
type OAuthFlow struct {
	AppID       string
	AppSecret   string
	Client      api.Client
	Keychain    KeychainStore
	OpenBrowser func(url string) error // injectable; defaults to system browser
}

// OAuthResult contains the result of a successful OAuth flow.
type OAuthResult struct {
	Username  string
	UserID    string
	Token     string
	ExpiresIn int64
}

// Run executes the OAuth flow: starts local server, opens browser, waits for callback.
func (o *OAuthFlow) Run() (*OAuthResult, error) {
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			if errMsg == "" {
				errMsg = "no authorization code received"
			}
			fmt.Fprintf(w, "<html><body><h2>Authentication Failed</h2><p>%s</p><p>You can close this window.</p></body></html>", errMsg)
			errChan <- fmt.Errorf("OAuth callback error: %s", errMsg)
			return
		}

		fmt.Fprint(w, "<html><body><h2>Authentication Successful!</h2><p>You can close this window and return to the terminal.</p></body></html>")
		codeChan <- code
	})

	server := &http.Server{
		Addr:    OAuthPort,
		Handler: mux,
	}

	go func() {
		slog.Debug("Starting OAuth callback server", "addr", OAuthPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("OAuth server error: %w", err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Open browser for authorization
	authURL := fmt.Sprintf(
		"https://www.instagram.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=instagram_business_basic,instagram_business_manage_insights",
		o.AppID, OAuthRedirectURI,
	)

	slog.Debug("Opening browser", "url", authURL)
	fmt.Printf("Opening browser for Instagram authorization...\n")
	fmt.Printf("If the browser doesn't open, visit this URL:\n%s\n\n", authURL)

	opener := o.OpenBrowser
	if opener == nil {
		opener = openBrowser
	}
	if err := opener(authURL); err != nil {
		slog.Warn("Failed to open browser", "error", err)
	}

	fmt.Println("Waiting for authorization...")

	// Wait for callback or error
	var code string
	select {
	case code = <-codeChan:
		slog.Debug("Received authorization code")
	case err := <-errChan:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("OAuth flow timed out after 5 minutes")
	}

	// Exchange code for short-lived token
	fmt.Println("Exchanging authorization code for token...")
	shortToken, err := o.Client.ExchangeCodeForToken(o.AppID, o.AppSecret, OAuthRedirectURI, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}

	// Exchange for long-lived token
	fmt.Println("Obtaining long-lived token...")
	longToken, err := o.Client.ExchangeForLongLivedToken(o.AppID, o.AppSecret, shortToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("getting long-lived token: %w", err)
	}

	// Get user profile
	fmt.Println("Fetching user profile...")
	profile, err := o.Client.GetUserProfile(longToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}

	return &OAuthResult{
		Username:  profile.Username,
		UserID:    profile.ID,
		Token:     longToken.AccessToken,
		ExpiresIn: longToken.ExpiresIn,
	}, nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
