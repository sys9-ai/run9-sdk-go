package run9

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadArchiveAcceptsCompletionTrailer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trailer", "X-Run9-Archive-Complete")
		_, _ = io.WriteString(w, "archive bytes")
		w.Header().Set("X-Run9-Archive-Complete", "true")
	}))
	defer server.Close()
	reader, err := newProjectTestClient(t, server.URL, "default").DownloadArchive(t.Context(), "box", "/file")
	require.NoError(t, err)
	defer reader.Close()
	data, err := io.ReadAll(reader)
	require.Equal(t, "archive bytes", string(data))
	require.NoError(t, err)
}

func TestDownloadArchiveRejectsMissingCompletionTrailer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trailer", "X-Run9-Archive-Complete")
		_, _ = io.WriteString(w, "archive bytes")

	}))
	defer server.Close()
	reader, err := newProjectTestClient(t, server.URL, "default").DownloadArchive(t.Context(), "box", "/file")
	require.NoError(t, err)
	defer reader.Close()
	data, err := io.ReadAll(reader)
	require.Equal(t, "archive bytes", string(data))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestDownloadArchiveRejectsCompletionHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Run9-Archive-Complete", "true")
		_, _ = io.WriteString(w, "archive bytes")

	}))
	defer server.Close()
	reader, err := newProjectTestClient(t, server.URL, "default").DownloadArchive(t.Context(), "box", "/file")
	require.NoError(t, err)
	defer reader.Close()
	data, err := io.ReadAll(reader)
	require.Equal(t, "archive bytes", string(data))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestDownloadArchiveRejectsFailedCompletionTrailer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trailer", "X-Run9-Archive-Complete")
		_, _ = io.WriteString(w, "archive bytes")
		w.Header().Set("X-Run9-Archive-Complete", "false")
	}))
	defer server.Close()
	reader, err := newProjectTestClient(t, server.URL, "default").DownloadArchive(t.Context(), "box", "/file")
	require.NoError(t, err)
	defer reader.Close()
	data, err := io.ReadAll(reader)
	require.Equal(t, "archive bytes", string(data))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestDownloadArchivePreservesTransportFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trailer", "X-Run9-Archive-Complete")
		_, _ = io.WriteString(w, "archive bytes")
		w.(http.Flusher).Flush()
		panic(http.ErrAbortHandler)
	}))
	defer server.Close()
	reader, err := newProjectTestClient(t, server.URL, "default").DownloadArchive(t.Context(), "box", "/file")
	require.NoError(t, err)
	defer reader.Close()
	data, err := io.ReadAll(reader)
	require.Equal(t, "archive bytes", string(data))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
