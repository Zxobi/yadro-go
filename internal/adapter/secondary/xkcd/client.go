package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
	"yadro-go/internal/adapter/secondary"
	"yadro-go/internal/core/domain"
	"yadro-go/pkg/logger"
)

type HttpClient struct {
	log *slog.Logger
	c   *http.Client
	url string
}

func NewHttpClient(log *slog.Logger, url string, timeout time.Duration) *HttpClient {
	c := &http.Client{Timeout: timeout}
	return &HttpClient{log: log, c: c, url: url}
}

func (xc *HttpClient) GetById(id int) (*domain.Comic, error) {
	const op = "xkcd.GetById"
	log := xc.log.With(slog.String("op", op))

	resp, err := xc.doGet(xc.makeComicUrl(id))
	if err != nil {
		log.Error("failed to make a request", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", op, secondary.ErrComicNotFound)
		}

		log.With(slog.Int("status", resp.StatusCode)).Error("failed to fetch comic")
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read body", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	comic, err := parseBody(body)
	if err != nil {
		log.Error("failed to parse body", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, secondary.ErrInternal)
	}

	return comic, nil
}

func (xc *HttpClient) doGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := xc.c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (xc *HttpClient) makeComicUrl(id int) string {
	return fmt.Sprintf("%s/%d/info.0.json", xc.url, id)
}

func parseBody(b []byte) (*domain.Comic, error) {
	comic := &domain.Comic{}
	return comic, json.Unmarshal(b, comic)
}
