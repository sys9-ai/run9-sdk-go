package run9

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBoxFileSystemResolvesOnceAndStreamsFiles(t *testing.T) {
	var boxLookups atomic.Int64
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-a":
			boxLookups.Add(1)
			require.Equal(t, "ak-1", basicAuthUser(r))
			writeJSONResponse(t, w, http.StatusOK, BoxView{BoxID: "box-a", FileAccessURL: server.URL + "/file-token/"})
		case "/file-token/work/one.txt":
			require.Empty(t, r.Header.Get("Authorization"))
			require.Equal(t, "bytes=1-2", r.Header.Get("Range"))
			w.Header().Set("X-Run9-File-Type", "file")
			w.Header().Set("ETag", `"one"`)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("ne"))
		case "/file-token/work/two.txt":
			w.Header().Set("X-Run9-File-Type", "file")
			_, _ = w.Write([]byte("two"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, Credentials{AK: "ak-1", SK: "sk-1"})
	require.NoError(t, err)
	files, err := client.WithProject("default").BoxFileSystem(context.Background(), "box-a")
	require.NoError(t, err)

	one, err := files.Open(context.Background(), "/work/one.txt", OpenFileOptions{Range: "bytes=1-2"})
	require.NoError(t, err)
	require.Equal(t, FileTypeFile, one.Info().Type)
	require.Equal(t, `"one"`, one.Info().ETag)
	data, err := io.ReadAll(one)
	require.NoError(t, err)
	require.NoError(t, one.Close())
	require.Equal(t, "ne", string(data))

	two, err := files.Open(context.Background(), "work/two.txt", OpenFileOptions{})
	require.NoError(t, err)
	data, err = io.ReadAll(two)
	require.NoError(t, err)
	require.NoError(t, two.Close())
	require.Equal(t, "two", string(data))
	require.Equal(t, int64(1), boxLookups.Load())
}

func TestSnapFileSystemStatAndReadDir(t *testing.T) {
	modified := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/snaps/snap-a":
			writeJSONResponse(t, w, http.StatusOK, SnapView{SnapID: "snap-a", FileAccessURL: server.URL + "/snap-token/"})
		case "/snap-token/work/site/index.html":
			require.Equal(t, http.MethodHead, r.Method)
			w.Header().Set("X-Run9-File-Type", "file")
			w.Header().Set("Content-Length", "42")
			w.Header().Set("Last-Modified", modified.Format(http.TimeFormat))
			w.Header().Set("ETag", `"index"`)
		case "/snap-token/work/site":
			require.Equal(t, "1", r.URL.Query().Get("list"))
			require.Equal(t, "25", r.URL.Query().Get("limit"))
			require.Equal(t, "next-a", r.URL.Query().Get("cursor"))
			writeJSONResponse(t, w, http.StatusOK, ReadDirPage{
				Entries: []FileEntry{{Name: "index.html", Type: FileTypeFile, Size: 42, ModTime: modified, ETag: `"index"`}},
				Cursor:  "next-b",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, Credentials{AK: "ak-1", SK: "sk-1"})
	require.NoError(t, err)
	files, err := client.WithProject("default").SnapFileSystem(context.Background(), "snap-a")
	require.NoError(t, err)

	info, err := files.Stat(context.Background(), "/work/site/index.html")
	require.NoError(t, err)
	require.Equal(t, FileInfo{Path: "/work/site/index.html", Type: FileTypeFile, Size: 42, ModTime: modified, ETag: `"index"`}, info)

	page, err := files.ReadDir(context.Background(), "/work/site", ReadDirRequest{Limit: 25, Cursor: "next-a"})
	require.NoError(t, err)
	require.Len(t, page.Entries, 1)
	require.Equal(t, "index.html", page.Entries[0].Name)
	require.Equal(t, "next-b", page.Cursor)
}

func TestFileSystemRejectsInvalidPathsAndRanges(t *testing.T) {
	files, err := newFileSystem("https://static.run.sys9.ai/file-token/", http.DefaultClient)
	require.NoError(t, err)

	_, err = files.Open(context.Background(), "../secret", OpenFileOptions{})
	require.EqualError(t, err, "file path must be canonical")
	_, err = files.Open(context.Background(), "/work/file", OpenFileOptions{Range: "0-10"})
	require.EqualError(t, err, "file range must start with bytes=")
	_, err = files.ReadDir(context.Background(), "/", ReadDirRequest{Limit: 1001})
	require.EqualError(t, err, "file list limit must be between 1 and 1000, or zero for the default")

	targetURL, canonicalPath, err := files.requestURL("/work/ leading and trailing ", nil)
	require.NoError(t, err)
	require.Equal(t, "/work/ leading and trailing ", canonicalPath)
	require.Equal(t, "https://static.run.sys9.ai/file-token/work/%20leading%20and%20trailing%20", targetURL)
}

func basicAuthUser(request *http.Request) string {
	username, _, _ := request.BasicAuth()
	return username
}
