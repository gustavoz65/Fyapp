package firebase

import (
	"context"
	"encoding/base64"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gustavoz65/Cashing-go/internal/config"
	"google.golang.org/api/option"
)

type Client struct {
	authClient *auth.Client
	config     *config.FirebaseConfig
}

// NewClient inicializa Firebase Admin SDK
func NewClient(cfg *config.FirebaseConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	var opt option.ClientOption

	if cfg.ServiceAccountJSON != "" {
		// Produção: base64 encoded JSON
		decoded, err := base64.StdEncoding.DecodeString(cfg.ServiceAccountJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to decode service account JSON: %w", err)
		}
		opt = option.WithCredentialsJSON(decoded)
	} else if cfg.ServiceAccountPath != "" {
		// Local: arquivo JSON
		opt = option.WithCredentialsFile(cfg.ServiceAccountPath)
	} else {
		return nil, fmt.Errorf("firebase service account not configured")
	}

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: cfg.ProjectID,
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}

	return &Client{
		authClient: authClient,
		config:     cfg,
	}, nil
}

// VerifyIDToken valida Firebase ID token e retorna claims
func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := c.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}
	return token, nil
}

// GetUser busca informações do usuário no Firebase
func (c *Client) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := c.authClient.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
