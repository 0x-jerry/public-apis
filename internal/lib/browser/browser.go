package browser

import (
	"context"
	"errors"
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

	return m.browser, m.connErr
}

func (m *Manager) connect() (*rod.Browser, error) {
	endpoint := m.cfg.BrowserWS
	if endpoint == "" {
		endpoint = "ws://127.0.0.1:9222"
	}

	browser := rod.New().ControlURL(endpoint)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
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
