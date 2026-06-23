package run9

import (
	"bufio"
	"encoding/json"
	"io"
)

// ExecStream reads one inline foreground exec event stream.
type ExecStream struct {
	// ExecID is the durable exec identifier assigned by the control plane.
	ExecID  string
	body    io.ReadCloser
	decoder *json.Decoder
}

func newExecStream(execID string, body io.ReadCloser) *ExecStream {
	return &ExecStream{
		ExecID:  execID,
		body:    body,
		decoder: json.NewDecoder(bufio.NewReader(body)),
	}
}

// ReadEvent reads the next inline foreground exec event.
func (s *ExecStream) ReadEvent() (ExecStreamEvent, error) {
	var event ExecStreamEvent
	if err := s.decoder.Decode(&event); err != nil {
		return ExecStreamEvent{}, err
	}
	return event, nil
}

// Close closes the underlying stream body.
func (s *ExecStream) Close() error {
	if s == nil || s.body == nil {
		return nil
	}
	return s.body.Close()
}
