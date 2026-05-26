package cotahist

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://bvmf.bmfbovespa.com.br/InstDados/SerHist"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type Option func(*Client)

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.httpClient = h }
}

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Minute},
		baseURL:    defaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// AnnualURL returns the COTAHIST URL for a full calendar year (COTAHIST_AYYYY.ZIP).
func (c *Client) AnnualURL(year int) string {
	return fmt.Sprintf("%s/COTAHIST_A%d.ZIP", c.baseURL, year)
}

// MonthlyURL returns the COTAHIST URL for a single month (COTAHIST_MMMYYYY.ZIP).
func (c *Client) MonthlyURL(year int, month time.Month) string {
	return fmt.Sprintf("%s/COTAHIST_M%02d%d.ZIP", c.baseURL, int(month), year)
}

// DailyURL returns the COTAHIST URL for a single trading day (COTAHIST_DDDMMYYYY.ZIP).
func (c *Client) DailyURL(date time.Time) string {
	return fmt.Sprintf("%s/COTAHIST_D%s.ZIP", c.baseURL, date.Format("02012006"))
}

// FetchAnnual downloads and parses a full-year COTAHIST file.
func (c *Client) FetchAnnual(ctx context.Context, year int) ([]Record, error) {
	return c.fetch(ctx, c.AnnualURL(year))
}

// FetchMonthly downloads and parses a single-month COTAHIST file.
func (c *Client) FetchMonthly(ctx context.Context, year int, month time.Month) ([]Record, error) {
	return c.fetch(ctx, c.MonthlyURL(year, month))
}

// FetchDaily downloads and parses a single-day COTAHIST file.
func (c *Client) FetchDaily(ctx context.Context, date time.Time) ([]Record, error) {
	return c.fetch(ctx, c.DailyURL(date))
}

// WalkAnnual streams a full-year COTAHIST file record by record.
// Use this instead of FetchAnnual to avoid loading millions of rows into memory.
func (c *Client) WalkAnnual(ctx context.Context, year int, fn func(Record) error) error {
	return c.walk(ctx, c.AnnualURL(year), fn)
}

// WalkMonthly streams a monthly COTAHIST file record by record.
func (c *Client) WalkMonthly(ctx context.Context, year int, month time.Month, fn func(Record) error) error {
	return c.walk(ctx, c.MonthlyURL(year, month), fn)
}

// WalkDaily streams a daily COTAHIST file record by record.
func (c *Client) WalkDaily(ctx context.Context, date time.Time, fn func(Record) error) error {
	return c.walk(ctx, c.DailyURL(date), fn)
}

func (c *Client) fetch(ctx context.Context, url string) ([]Record, error) {
	rc, err := c.open(ctx, url)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return Parse(rc)
}

func (c *Client) walk(ctx context.Context, url string, fn func(Record) error) error {
	rc, err := c.open(ctx, url)
	if err != nil {
		return err
	}
	defer rc.Close()
	return Walk(rc, fn)
}

// open downloads the COTAHIST ZIP at url and returns a reader for the inner text file.
func (c *Client) open(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	if len(zr.File) == 0 {
		return nil, errors.New("empty zip")
	}

	return zr.File[0].Open()
}
