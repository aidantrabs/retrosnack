package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const oEmbedEndpoint = "https://graph.facebook.com/v18.0/instagram_oembed"

type Service interface {
	GetEmbed(ctx context.Context, productID uuid.UUID) (*Link, error)
	RefreshEmbed(ctx context.Context, productID uuid.UUID, postURL string) (*Link, error)
}

type service struct {
	repo   Repository
	client *http.Client
}

func NewService(repo Repository) Service {
	return &service{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *service) GetEmbed(ctx context.Context, productID uuid.UUID) (*Link, error) {
	return s.repo.GetByProductID(ctx, productID)
}

func (s *service) RefreshEmbed(ctx context.Context, productID uuid.UUID, postURL string) (*Link, error) {
	embedHTML, err := s.fetchOEmbed(ctx, postURL)
	if err != nil {
		// Fallback: store the link without embed HTML
		embedHTML = ""
	}

	return s.repo.Upsert(ctx, productID, postURL, embedHTML)
}

func (s *service) fetchOEmbed(ctx context.Context, postURL string) (string, error) {
	endpoint := fmt.Sprintf("%s?url=%s&omitscript=true",
		oEmbedEndpoint, url.QueryEscape(postURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oEmbed API returned %d", resp.StatusCode)
	}

	var result oEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.HTML, nil
}
