package browser

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"public-apis/internal/config"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Manager struct {
	cfg     *config.Config
	mu      sync.Mutex
	browser *rod.Browser
	once    sync.Once
	connErr error
}

func New(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) isEnabled() bool {
	return m.cfg.BrowserWSEnabled
}

func (m *Manager) Get() (*rod.Browser, error) {
	if !m.isEnabled() {
		log.Printf("browser: not enabled")
		return nil, nil
	}

	m.mu.Lock()
	if m.browser != nil {
		b := m.browser
		m.mu.Unlock()
		return b, nil
	}
	m.mu.Unlock()

	m.once.Do(func() {
		m.browser, m.connErr = m.connect()
	})

	if m.connErr != nil {
		log.Printf("browser: unavailable: %v", m.connErr)
	}

	return m.browser, m.connErr
}

func (m *Manager) connect() (*rod.Browser, error) {
	endpoint := m.cfg.BrowserWS
	if endpoint == "" {
		endpoint = "ws://127.0.0.1:9222"
	}

	log.Printf("browser: connecting to %s", endpoint)
	browser := rod.New().ControlURL(endpoint)
	if err := browser.Connect(); err != nil {
		log.Printf("browser: connect to %s failed: %v", endpoint, err)
		return nil, err
	}
	log.Printf("browser: connected to %s", endpoint)
	return browser, nil
}

func (m *Manager) FetchHTML(ctx context.Context, url string) (string, error) {
	b, err := m.Get()
	if err != nil {
		return "", err
	}
	if b == nil {
		return "", errors.New("browser not available")
	}

	page, err := b.Page(proto.TargetCreateTarget{})
	if err != nil {
		return "", err
	}
	defer page.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	page = page.Context(ctx)
	if err := page.Navigate(url); err != nil {
		return "", err
	}

	page.WaitStable(1 * time.Second)

	return page.HTML()
}

// FetchResult is the raw result of a browser navigation.
type FetchResult struct {
	Body        []byte
	ContentType string
	StatusCode  int
}

// Fetch navigates to url with the headless browser and returns the raw response
// body (for example HTML or XML) along with its content type and status code.
// Unlike FetchHTML, it does not render or otherwise transform the document.
func (m *Manager) Fetch(ctx context.Context, url string) (*FetchResult, error) {
	b, err := m.Get()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("browser not available")
	}

	page, err := b.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}
	defer page.Close()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	page = page.Context(ctx)

	if err := (proto.NetworkEnable{}).Call(page); err != nil {
		return nil, err
	}

	var (
		requestID   proto.NetworkRequestID
		contentType string
		statusCode  int
	)

	// Capture the final document response, skipping redirect hops.
	wait := page.EachEvent(func(e *proto.NetworkResponseReceived) bool {
		if e.Type != proto.NetworkResourceTypeDocument {
			return false
		}
		if e.Response.Status >= 300 && e.Response.Status < 400 {
			return false
		}

		requestID = e.RequestID
		statusCode = e.Response.Status
		contentType = e.Response.MIMEType

		for key, value := range e.Response.Headers {
			if strings.EqualFold(key, "content-type") {
				contentType = value.Str()
				break
			}
		}

		return true
	})

	if err := page.Navigate(url); err != nil {
		return nil, err
	}
	wait()

	if requestID == "" {
		return nil, errors.New("no document response received")
	}

	res, err := proto.NetworkGetResponseBody{RequestID: requestID}.Call(page)
	if err != nil {
		return nil, err
	}

	body := []byte(res.Body)
	if res.Base64Encoded {
		body, err = base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			return nil, err
		}
	}

	return &FetchResult{
		Body:        body,
		ContentType: contentType,
		StatusCode:  statusCode,
	}, nil
}
