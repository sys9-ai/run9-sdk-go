package run9

import (
	"fmt"
	"io"
	"net/http"
)

// ArchiveReader reads a downloaded tar archive and verifies that the runtime
// completed it. Callers must read through EOF before using the downloaded files;
// the tar end marker alone does not confirm transfer completion.
type ArchiveReader struct {
	response *http.Response
	terminal error
}

// Read returns archive bytes. Missing completion confirmation is an error even
// when an intermediary has ended the HTTP response cleanly.
func (r *ArchiveReader) Read(p []byte) (int, error) {
	if r.terminal != nil {
		return 0, r.terminal
	}
	n, err := r.response.Body.Read(p)
	if err == io.EOF && r.response.Trailer.Get("X-Run9-Archive-Complete") != "true" {
		err = fmt.Errorf("archive download completion was not confirmed: %w", io.ErrUnexpectedEOF)
	}
	if err != nil {
		r.terminal = err
	}
	return n, err
}

// Close releases the download connection. It does not confirm completion.
func (r *ArchiveReader) Close() error { return r.response.Body.Close() }
