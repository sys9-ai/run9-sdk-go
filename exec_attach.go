package run9

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ExecAttachSocket wraps one exec attach websocket connection.
type ExecAttachSocket struct {
	mu   sync.Mutex
	conn *websocket.Conn
}

// ReadEvent reads the next exec stream event from the websocket.
func (s *ExecAttachSocket) ReadEvent() (ExecStreamEvent, error) {
	var event ExecStreamEvent
	if err := s.conn.ReadJSON(&event); err != nil {
		return ExecStreamEvent{}, err
	}
	return event, nil
}

// WriteInput writes one exec attach input frame to the websocket.
func (s *ExecAttachSocket) WriteInput(input ExecAttachInput) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteJSON(input)
}

// Close closes the underlying websocket connection.
func (s *ExecAttachSocket) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

// OpenExecAttach dials one exec attach websocket returned by the control plane.
func (c *Client) OpenExecAttach(ctx context.Context, attachURL string) (*ExecAttachSocket, error) {
	resolvedURL, err := resolveHTTPURL(c.baseURL, attachURL)
	if err != nil {
		return nil, err
	}
	wsURL, err := websocketURL(resolvedURL, "")
	if err != nil {
		return nil, err
	}
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 15 * time.Second
	conn, resp, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		if resp != nil {
			return nil, responseError(resp)
		}
		return nil, err
	}
	return &ExecAttachSocket{conn: conn}, nil
}

func websocketURL(baseURL string, path string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	switch parsed.Scheme {
	case "http":
		parsed.Scheme = "ws"
	case "https":
		parsed.Scheme = "wss"
	default:
		return "", errors.New("expected http or https endpoint")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + path
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func resolveHTTPURL(baseURL string, target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", errors.New("missing attach url")
	}
	parsedTarget, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	if parsedTarget.IsAbs() {
		return parsedTarget.String(), nil
	}
	base, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	return base.ResolveReference(parsedTarget).String(), nil
}
