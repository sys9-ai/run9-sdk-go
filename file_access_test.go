package run9

import (
	"context"
	"encoding/json"
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

func TestFileSystemGlobFilesUsesOneBoundedServerRequest(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/access/work", r.URL.Path)
		require.Equal(t, "**/*srch*", r.URL.Query().Get("glob"))
		require.Equal(t, "25", r.URL.Query().Get("limit"))
		require.Equal(t, []string{".git", "node_modules"}, r.URL.Query()["exclude_dir"])
		require.Equal(t, "/workspace", r.Header.Get(fileRootHeader))
		writeJSONResponse(t, w, http.StatusOK, GlobFilesResult{
			Matches:   []FileGlobMatch{{Path: "src/search.go"}},
			Truncated: true,
		})
	}))
	defer server.Close()

	files, err := newFileSystem(server.URL+"/access/", server.Client())
	require.NoError(t, err)
	files, err = files.RootedAt("/workspace")
	require.NoError(t, err)
	result, err := files.GlobFiles(context.Background(), "/work", GlobFilesRequest{
		Pattern:            "**/*srch*",
		Limit:              25,
		ExcludeDirectories: []string{".git", "node_modules"},
	})
	require.NoError(t, err)
	require.Equal(t, []FileGlobMatch{{Path: "src/search.go"}}, result.Matches)
	require.True(t, result.Truncated)
	require.Equal(t, int64(1), requests.Load())
}

func TestFileSystemGlobFilesAllowsEscapedLiteralBrace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Contains(t, []string{`\{foo\}.go`, `[{]foo[}].go`}, r.URL.Query().Get("glob"))
		writeJSONResponse(t, w, http.StatusOK, GlobFilesResult{
			Matches: []FileGlobMatch{{Path: "{foo}.go"}},
		})
	}))
	defer server.Close()

	files, err := newFileSystem(server.URL+"/access/", server.Client())
	require.NoError(t, err)
	result, err := files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: `\{foo\}.go`})
	require.NoError(t, err)
	require.Equal(t, []FileGlobMatch{{Path: "{foo}.go"}}, result.Matches)
	result, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: `[{]foo[}].go`})
	require.NoError(t, err)
	require.Equal(t, []FileGlobMatch{{Path: "{foo}.go"}}, result.Matches)
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
	_, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{})
	require.EqualError(t, err, "file glob pattern is required")
	_, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: "/**/*.go"})
	require.EqualError(t, err, "file glob pattern must be relative to the requested directory")
	_, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: "{src,lib}/**/*.go"})
	require.EqualError(t, err, "file glob brace alternatives are not supported")
	_, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: "**/*.go", Limit: 201})
	require.EqualError(t, err, "file glob limit must be between 1 and 200, or zero for the default")
	_, err = files.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: "**/*.go", ExcludeDirectories: []string{"src/generated"}})
	require.EqualError(t, err, "excluded directory names must be single path components")

	targetURL, canonicalPath, err := files.requestURL("/work/ leading and trailing ", nil)
	require.NoError(t, err)
	require.Equal(t, "/work/ leading and trailing ", canonicalPath)
	require.Equal(t, "https://static.run.sys9.ai/file-token/work/%20leading%20and%20trailing%20", targetURL)
}

func TestFileSystemRootedAtBindsEveryRequestToOneRoot(t *testing.T) {
	var roots []string
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		roots = append(roots, request.Header.Get(fileRootHeader))
		switch request.Method {
		case http.MethodHead:
			response.Header().Set("X-Run9-File-Type", "file")
		default:
			if request.URL.Query().Has("glob") {
				_ = json.NewEncoder(response).Encode(GlobFilesResult{Matches: []FileGlobMatch{}})
				return
			}
			if request.URL.Query().Get("list") == "1" {
				_ = json.NewEncoder(response).Encode(ReadDirPage{Entries: []FileEntry{}})
				return
			}
			_, _ = response.Write([]byte("content"))
		}
	}))
	defer server.Close()

	files, err := newFileSystem(server.URL+"/access/", server.Client())
	require.NoError(t, err)
	workspace, err := files.RootedAt("/workspace")
	require.NoError(t, err)

	reader, err := workspace.Open(context.Background(), "/README.md", OpenFileOptions{})
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	_, err = workspace.Stat(context.Background(), "/README.md")
	require.NoError(t, err)
	_, err = workspace.ReadDir(context.Background(), "/", ReadDirRequest{})
	require.NoError(t, err)
	_, err = workspace.GlobFiles(context.Background(), "/", GlobFilesRequest{Pattern: "**/*"})
	require.NoError(t, err)

	require.Equal(t, []string{"/workspace", "/workspace", "/workspace", "/workspace"}, roots)
}

func TestFileSystemRootedAtRejectsInvalidRootAndKeepsOriginalRoot(t *testing.T) {
	files, err := newFileSystem("https://static.run.sys9.ai/file-token/", http.DefaultClient)
	require.NoError(t, err)

	_, err = files.RootedAt("/workspace/../etc")
	require.EqualError(t, err, "file system root: file path must be canonical")

	workspace, err := files.RootedAt("/workspace")
	require.NoError(t, err)
	require.Equal(t, "/", files.root)
	require.Equal(t, "/workspace", workspace.root)
}

func basicAuthUser(request *http.Request) string {
	username, _, _ := request.BasicAuth()
	return username
}
