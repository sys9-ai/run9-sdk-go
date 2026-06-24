package run9

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	backgroundExecFrameStarted byte = iota + 1
	backgroundExecFrameStdout
	backgroundExecFrameStderr
	backgroundExecFrameGap
	backgroundExecFrameTruncated
	backgroundExecFrameExit
	backgroundExecFrameCancelled
	backgroundExecFrameError
)

const backgroundExecFrameHeaderSize = 1 + 8 + 4

const (
	backgroundExecGapNoticeFormat  = "[run9] background exec omitted %d bytes around this cursor\n"
	backgroundExecTruncatedMessage = "[run9] background exec output was truncated at the service limit\n"
)

// BackgroundExecOutputEventType identifies one event decoded from a background exec output window.
type BackgroundExecOutputEventType string

const (
	// BackgroundExecOutputEventStarted reports the durable exec start marker.
	BackgroundExecOutputEventStarted BackgroundExecOutputEventType = "started"
	// BackgroundExecOutputEventStdout reports one stdout chunk.
	BackgroundExecOutputEventStdout BackgroundExecOutputEventType = "stdout"
	// BackgroundExecOutputEventStderr reports one stderr chunk.
	BackgroundExecOutputEventStderr BackgroundExecOutputEventType = "stderr"
	// BackgroundExecOutputEventGap reports that older bytes were omitted around the current cursor.
	BackgroundExecOutputEventGap BackgroundExecOutputEventType = "gap"
	// BackgroundExecOutputEventTruncated reports that the service truncated retained output at its hard limit.
	BackgroundExecOutputEventTruncated BackgroundExecOutputEventType = "truncated"
	// BackgroundExecOutputEventExit reports a terminal process exit.
	BackgroundExecOutputEventExit BackgroundExecOutputEventType = "exit"
	// BackgroundExecOutputEventCancelled reports a terminal cancellation.
	BackgroundExecOutputEventCancelled BackgroundExecOutputEventType = "cancelled"
	// BackgroundExecOutputEventError reports a terminal runtime or control-plane failure.
	BackgroundExecOutputEventError BackgroundExecOutputEventType = "error"
)

// BackgroundExecOutputEvent describes one decoded background exec output event.
type BackgroundExecOutputEvent struct {
	// Seq is the durable event sequence number for incremental replay.
	Seq uint64
	// Type identifies whether this event is stdout, stderr, gap, truncated, or terminal.
	Type BackgroundExecOutputEventType
	// Data carries stdout or stderr bytes on output events.
	Data []byte
	// GapBytes reports how many bytes were omitted on gap events.
	GapBytes uint64
	// ExitCode is set on exit events.
	ExitCode *int
	// Reason is set on cancelled and error events.
	Reason string
}

func decodeBackgroundExecOutputEvents(body []byte) ([]BackgroundExecOutputEvent, error) {
	reader := bytes.NewReader(body)
	events := make([]BackgroundExecOutputEvent, 0)
	for reader.Len() > 0 {
		frameType, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		var seq uint64
		if err := binary.Read(reader, binary.BigEndian, &seq); err != nil {
			return nil, err
		}

		var payloadLen uint32
		if err := binary.Read(reader, binary.BigEndian, &payloadLen); err != nil {
			return nil, err
		}

		payload := make([]byte, payloadLen)
		if payloadLen > 0 {
			if _, err := io.ReadFull(reader, payload); err != nil {
				return nil, err
			}
		}

		event, err := decodeBackgroundExecOutputEvent(frameType, seq, payload)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func decodeBackgroundExecOutputEvent(frameType byte, seq uint64, payload []byte) (BackgroundExecOutputEvent, error) {
	event := BackgroundExecOutputEvent{Seq: seq}
	switch frameType {
	case backgroundExecFrameStarted:
		if len(payload) != 0 {
			return BackgroundExecOutputEvent{}, errors.New("started frame expects empty payload")
		}
		event.Type = BackgroundExecOutputEventStarted
	case backgroundExecFrameStdout:
		event.Type = BackgroundExecOutputEventStdout
		event.Data = append([]byte(nil), payload...)
	case backgroundExecFrameStderr:
		event.Type = BackgroundExecOutputEventStderr
		event.Data = append([]byte(nil), payload...)
	case backgroundExecFrameGap:
		if len(payload) != 8 {
			return BackgroundExecOutputEvent{}, errors.New("gap frame expects 8-byte payload")
		}
		event.Type = BackgroundExecOutputEventGap
		event.GapBytes = binary.BigEndian.Uint64(payload)
	case backgroundExecFrameTruncated:
		if len(payload) != 0 {
			return BackgroundExecOutputEvent{}, errors.New("truncated frame expects empty payload")
		}
		event.Type = BackgroundExecOutputEventTruncated
	case backgroundExecFrameExit:
		if len(payload) != 4 {
			return BackgroundExecOutputEvent{}, errors.New("exit frame expects 4-byte payload")
		}
		exitCode := int(int32(binary.BigEndian.Uint32(payload)))
		event.Type = BackgroundExecOutputEventExit
		event.ExitCode = &exitCode
	case backgroundExecFrameCancelled:
		event.Type = BackgroundExecOutputEventCancelled
		event.Reason = string(payload)
	case backgroundExecFrameError:
		event.Type = BackgroundExecOutputEventError
		event.Reason = string(payload)
	default:
		return BackgroundExecOutputEvent{}, fmt.Errorf("unsupported background exec frame type %d", frameType)
	}
	return event, nil
}

func encodeBackgroundExecOutputEvents(events []BackgroundExecOutputEvent) ([]byte, error) {
	total := 0
	payloads := make([][]byte, 0, len(events))
	types := make([]byte, 0, len(events))
	for _, event := range events {
		frameType, payload, err := encodeBackgroundExecOutputEvent(event)
		if err != nil {
			return nil, err
		}
		types = append(types, frameType)
		payloads = append(payloads, payload)
		total += backgroundExecFrameHeaderSize + len(payload)
	}

	buf := bytes.NewBuffer(make([]byte, 0, total))
	for i, event := range events {
		buf.WriteByte(types[i])
		if err := binary.Write(buf, binary.BigEndian, event.Seq); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, uint32(len(payloads[i]))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(payloads[i]); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func encodeBackgroundExecOutputEvent(event BackgroundExecOutputEvent) (byte, []byte, error) {
	switch event.Type {
	case BackgroundExecOutputEventStarted:
		return backgroundExecFrameStarted, nil, nil
	case BackgroundExecOutputEventStdout:
		return backgroundExecFrameStdout, append([]byte(nil), event.Data...), nil
	case BackgroundExecOutputEventStderr:
		return backgroundExecFrameStderr, append([]byte(nil), event.Data...), nil
	case BackgroundExecOutputEventGap:
		payload := make([]byte, 8)
		binary.BigEndian.PutUint64(payload, event.GapBytes)
		return backgroundExecFrameGap, payload, nil
	case BackgroundExecOutputEventTruncated:
		return backgroundExecFrameTruncated, nil, nil
	case BackgroundExecOutputEventExit:
		payload := make([]byte, 4)
		exitCode := 0
		if event.ExitCode != nil {
			exitCode = *event.ExitCode
		}
		binary.BigEndian.PutUint32(payload, uint32(int32(exitCode)))
		return backgroundExecFrameExit, payload, nil
	case BackgroundExecOutputEventCancelled:
		return backgroundExecFrameCancelled, []byte(event.Reason), nil
	case BackgroundExecOutputEventError:
		return backgroundExecFrameError, []byte(event.Reason), nil
	default:
		return 0, nil, fmt.Errorf("unsupported background exec output event type %q", event.Type)
	}
}

func normalizeWriter(writer io.Writer) io.Writer {
	if writer == nil {
		return io.Discard
	}
	return writer
}

func writeBackgroundExecOutput(events []BackgroundExecOutputEvent, stdout io.Writer, stderr io.Writer, notices io.Writer) error {
	stdout = normalizeWriter(stdout)
	stderr = normalizeWriter(stderr)
	notices = normalizeWriter(notices)

	for _, event := range events {
		switch event.Type {
		case BackgroundExecOutputEventStarted:
			continue
		case BackgroundExecOutputEventStdout:
			if _, err := stdout.Write(event.Data); err != nil {
				return err
			}
		case BackgroundExecOutputEventStderr:
			if _, err := stderr.Write(event.Data); err != nil {
				return err
			}
		case BackgroundExecOutputEventGap:
			if _, err := fmt.Fprintf(notices, backgroundExecGapNoticeFormat, event.GapBytes); err != nil {
				return err
			}
		case BackgroundExecOutputEventTruncated:
			if _, err := io.WriteString(notices, backgroundExecTruncatedMessage); err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteOutput writes stdout events to Stdout, stderr events to Stderr, and gap or truncation notices to Stderr.
func (r BackgroundExecPullOutput) WriteOutput(writers ExecOutputWriters) error {
	writers = normalizeExecOutputWriters(writers)
	return writeBackgroundExecOutput(r.Events, writers.Stdout, writers.Stderr, writers.Stderr)
}

// WriteMergedOutput writes stdout and stderr chunks to output in event order and writes gap or truncation notices to notices.
func (r BackgroundExecPullOutput) WriteMergedOutput(output io.Writer, notices io.Writer) error {
	return writeBackgroundExecOutput(r.Events, output, output, notices)
}

// TerminalResult returns the terminal outcome carried by this output window when one is available.
func (r BackgroundExecPullOutput) TerminalResult() *ExecTerminalResult {
	for i := len(r.Events) - 1; i >= 0; i-- {
		event := r.Events[i]
		switch event.Type {
		case BackgroundExecOutputEventExit:
			return &ExecTerminalResult{
				Status:   ExecTerminalStatusExited,
				ExitCode: cloneOptionalInt(event.ExitCode),
			}
		case BackgroundExecOutputEventCancelled:
			return &ExecTerminalResult{
				Status: ExecTerminalStatusCancelled,
				Reason: event.Reason,
			}
		case BackgroundExecOutputEventError:
			return &ExecTerminalResult{
				Status: ExecTerminalStatusError,
				Reason: event.Reason,
			}
		}
	}

	state := strings.TrimSpace(r.State)
	switch state {
	case "succeeded":
		exitCode := 0
		if r.ExitCode != nil {
			exitCode = *r.ExitCode
		}
		return &ExecTerminalResult{
			Status:   ExecTerminalStatusExited,
			ExitCode: &exitCode,
		}
	case "failed":
		if r.ExitCode != nil {
			return &ExecTerminalResult{
				Status:   ExecTerminalStatusExited,
				ExitCode: cloneOptionalInt(r.ExitCode),
			}
		}
		if strings.TrimSpace(r.Reason) != "" {
			return &ExecTerminalResult{
				Status: ExecTerminalStatusError,
				Reason: r.Reason,
			}
		}
	case "cancelled":
		return &ExecTerminalResult{
			Status: ExecTerminalStatusCancelled,
			Reason: r.Reason,
		}
	case "error":
		return &ExecTerminalResult{
			Status:   ExecTerminalStatusError,
			ExitCode: cloneOptionalInt(r.ExitCode),
			Reason:   r.Reason,
		}
	}
	return nil
}

func cloneOptionalInt(value *int) *int {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}
