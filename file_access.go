package run9

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const fileRootHeader = "X-Run9-File-Root"

// FileType identifies the kind of one filesystem entry.
type FileType string

const (
	// FileTypeFile identifies a regular file.
	FileTypeFile FileType = "file"
	// FileTypeDirectory identifies a directory.
	FileTypeDirectory FileType = "directory"
	// FileTypeSymlink identifies a symbolic link in a directory listing.
	FileTypeSymlink FileType = "symlink"
)

// FileInfo describes one path returned by FileSystem.Stat.
type FileInfo struct {
	Path        string    `json:"path"`
	Type        FileType  `json:"type"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time,omitempty"`
	ETag        string    `json:"etag,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
}

// FileEntry describes one child returned by FileSystem.ReadDir.
type FileEntry struct {
	Name    string    `json:"name"`
	Type    FileType  `json:"type"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	ETag    string    `json:"etag,omitempty"`
}

// ReadDirRequest controls one bounded directory page.
type ReadDirRequest struct {
	// Limit is the maximum number of entries. Zero uses the gateway default.
	Limit int
	// Cursor continues a previous page for the same directory and filesystem view.
	Cursor string
}

// ReadDirPage is one name-ordered directory page.
type ReadDirPage struct {
	Entries []FileEntry `json:"entries"`
	Cursor  string      `json:"cursor,omitempty"`
}

// GlobFilesRequest controls one bounded recursive file glob.
type GlobFilesRequest struct {
	// Pattern is a case-sensitive doublestar glob relative to the requested directory.
	Pattern string
	// Limit is the maximum number of matches. Zero uses the gateway default.
	Limit int
	// ExcludeDirectories skips directories with these exact names at any depth.
	ExcludeDirectories []string
}

// FileGlobMatch identifies one regular file relative to the requested directory.
type FileGlobMatch struct {
	Path string `json:"path"`
}

// GlobFilesResult contains lexical matches from one bounded filesystem view.
type GlobFilesResult struct {
	Matches   []FileGlobMatch `json:"matches"`
	Truncated bool            `json:"truncated"`
}

// OpenFileOptions controls one streamed file read.
type OpenFileOptions struct {
	// Range is an optional HTTP byte range, for example "bytes=0-1023".
	// Single and multipart ranges supported by the gateway are accepted.
	Range string
}

// FileReader streams one file response and exposes its metadata.
type FileReader struct {
	body io.ReadCloser
	info FileInfo
}

// Read reads file content from the gateway response.
func (reader *FileReader) Read(data []byte) (int, error) {
	return reader.body.Read(data)
}

// Close releases the underlying HTTP response body.
func (reader *FileReader) Close() error {
	return reader.body.Close()
}

// Info returns metadata from the file response headers.
func (reader *FileReader) Info() FileInfo {
	return reader.info
}

// FileSystem is a reusable read-only Box or Snap filesystem capability.
// It reuses the parent Client's HTTP transport and does not contact the control
// plane again after resolution.
type FileSystem struct {
	baseURL *url.URL
	http    *http.Client
	root    string
}

// BoxFileSystem resolves one Box filesystem capability.
func (c *Client) BoxFileSystem(ctx context.Context, boxID string) (*FileSystem, error) {
	view, err := c.GetBox(ctx, boxID)
	if err != nil {
		return nil, err
	}
	return newFileSystem(view.FileAccessURL, c.http)
}

// SnapFileSystem resolves one Snap filesystem capability. Detached snaps use
// their immutable read view; attached snaps resolve to their owning Box view.
func (c *Client) SnapFileSystem(ctx context.Context, snapID string) (*FileSystem, error) {
	view, err := c.GetSnap(ctx, snapID)
	if err != nil {
		return nil, err
	}
	return newFileSystem(view.FileAccessURL, c.http)
}

func newFileSystem(accessURL string, client *http.Client) (*FileSystem, error) {
	if strings.TrimSpace(accessURL) == "" {
		return nil, errors.New("file access is unavailable for this resource")
	}
	parsed, err := url.Parse(accessURL)
	if err != nil {
		return nil, fmt.Errorf("parse file access URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("file access URL must use HTTP or HTTPS")
	}
	if parsed.Host == "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("invalid file access URL")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/"
	return &FileSystem{baseURL: parsed, http: client, root: "/"}, nil
}

// RootedAt returns a filesystem whose "/" is rooted at root inside the Box or
// Snap. The gateway also resolves symlinks inside that root.
func (fileSystem *FileSystem) RootedAt(root string) (*FileSystem, error) {
	if fileSystem == nil || fileSystem.baseURL == nil || fileSystem.http == nil {
		return nil, errors.New("file system is not initialized")
	}
	canonicalRoot, err := canonicalFilePath(root)
	if err != nil {
		return nil, fmt.Errorf("file system root: %w", err)
	}
	rooted := *fileSystem
	rooted.root = canonicalRoot
	return &rooted, nil
}

// Open starts a streamed read of one regular file.
func (fileSystem *FileSystem) Open(ctx context.Context, filePath string, options OpenFileOptions) (*FileReader, error) {
	targetURL, canonicalPath, err := fileSystem.requestURL(filePath, nil)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set(fileRootHeader, fileSystem.root)
	if byteRange := strings.TrimSpace(options.Range); byteRange != "" {
		if !strings.HasPrefix(byteRange, "bytes=") {
			return nil, errors.New("file range must start with bytes=")
		}
		request.Header.Set("Range", byteRange)
	}
	response, err := fileSystem.http.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= http.StatusBadRequest {
		return nil, responseError(response)
	}
	return &FileReader{body: response.Body, info: fileInfoFromResponse(canonicalPath, response)}, nil
}

// Stat loads metadata for one file or directory without reading its content.
func (fileSystem *FileSystem) Stat(ctx context.Context, filePath string) (FileInfo, error) {
	targetURL, canonicalPath, err := fileSystem.requestURL(filePath, nil)
	if err != nil {
		return FileInfo{}, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, targetURL, nil)
	if err != nil {
		return FileInfo{}, err
	}
	request.Header.Set(fileRootHeader, fileSystem.root)
	response, err := fileSystem.http.Do(request)
	if err != nil {
		return FileInfo{}, err
	}
	if response.StatusCode >= http.StatusBadRequest {
		return FileInfo{}, responseError(response)
	}
	defer response.Body.Close()
	return fileInfoFromResponse(canonicalPath, response), nil
}

// ReadDir loads one bounded, name-ordered directory page.
func (fileSystem *FileSystem) ReadDir(ctx context.Context, directoryPath string, request ReadDirRequest) (ReadDirPage, error) {
	if request.Limit < 0 || request.Limit > 1000 {
		return ReadDirPage{}, errors.New("file list limit must be between 1 and 1000, or zero for the default")
	}
	query := url.Values{"list": []string{"1"}}
	if request.Limit != 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}
	if cursor := strings.TrimSpace(request.Cursor); cursor != "" {
		query.Set("cursor", cursor)
	}
	response, err := fileSystem.getResponse(ctx, directoryPath, query)
	if err != nil {
		return ReadDirPage{}, err
	}
	defer response.Body.Close()
	var page ReadDirPage
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		return ReadDirPage{}, fmt.Errorf("decode file list: %w", err)
	}
	if page.Entries == nil {
		page.Entries = []FileEntry{}
	}
	return page, nil
}

// GlobFiles matches regular files in one server-side operation. Pattern uses
// doublestar glob syntax and is relative to directoryPath. Results are
// case-sensitive and lexically ordered. Symlinks are not returned or followed.
// Truncated reports additional matches or a scan stopped by safety bounds.
func (fileSystem *FileSystem) GlobFiles(ctx context.Context, directoryPath string, request GlobFilesRequest) (GlobFilesResult, error) {
	if request.Pattern == "" {
		return GlobFilesResult{}, errors.New("file glob pattern is required")
	}
	if len(request.Pattern) > 256 {
		return GlobFilesResult{}, errors.New("file glob pattern must be at most 256 bytes")
	}
	if strings.HasPrefix(request.Pattern, "/") {
		return GlobFilesResult{}, errors.New("file glob pattern must be relative to the requested directory")
	}
	if strings.Contains(request.Pattern, "{") {
		return GlobFilesResult{}, errors.New("file glob brace alternatives are not supported")
	}
	if request.Limit < 0 || request.Limit > 200 {
		return GlobFilesResult{}, errors.New("file glob limit must be between 1 and 200, or zero for the default")
	}
	if len(request.ExcludeDirectories) > 32 {
		return GlobFilesResult{}, errors.New("file glob accepts at most 32 excluded directory names")
	}
	query := url.Values{"glob": []string{request.Pattern}}
	if request.Limit != 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}
	for _, name := range request.ExcludeDirectories {
		if name == "" || name == "." || name == ".." || len(name) > 255 || strings.ContainsAny(name, "/\\\x00") {
			return GlobFilesResult{}, errors.New("excluded directory names must be single path components")
		}
		query.Add("exclude_dir", name)
	}
	response, err := fileSystem.getResponse(ctx, directoryPath, query)
	if err != nil {
		return GlobFilesResult{}, err
	}
	defer response.Body.Close()
	var result GlobFilesResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return GlobFilesResult{}, fmt.Errorf("decode file glob: %w", err)
	}
	if result.Matches == nil {
		result.Matches = []FileGlobMatch{}
	}
	return result, nil
}

func (fileSystem *FileSystem) getResponse(ctx context.Context, filePath string, query url.Values) (*http.Response, error) {
	targetURL, _, err := fileSystem.requestURL(filePath, query)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set(fileRootHeader, fileSystem.root)
	response, err := fileSystem.http.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= http.StatusBadRequest {
		return nil, responseError(response)
	}
	return response, nil
}

func (fileSystem *FileSystem) requestURL(filePath string, query url.Values) (string, string, error) {
	if fileSystem == nil || fileSystem.baseURL == nil || fileSystem.http == nil {
		return "", "", errors.New("file system is not initialized")
	}
	canonicalPath, err := canonicalFilePath(filePath)
	if err != nil {
		return "", "", err
	}
	target := *fileSystem.baseURL
	target.Path = strings.TrimRight(target.Path, "/") + canonicalPath
	target.RawPath = ""
	target.RawQuery = query.Encode()
	return target.String(), canonicalPath, nil
}

func canonicalFilePath(value string) (string, error) {
	if value == "" || value == "/" {
		return "/", nil
	}
	if strings.IndexByte(value, 0) >= 0 || strings.Contains(value, "\\") {
		return "", errors.New("file path contains an invalid character")
	}
	segments := strings.Split(strings.Trim(value, "/"), "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", errors.New("file path must be canonical")
		}
	}
	return path.Clean("/" + strings.Join(segments, "/")), nil
}

func fileInfoFromResponse(filePath string, response *http.Response) FileInfo {
	size := response.ContentLength
	if size < 0 {
		size = 0
	}
	modTime, _ := http.ParseTime(response.Header.Get("Last-Modified"))
	return FileInfo{
		Path:        filePath,
		Type:        FileType(response.Header.Get("X-Run9-File-Type")),
		Size:        size,
		ModTime:     modTime,
		ETag:        response.Header.Get("ETag"),
		ContentType: response.Header.Get("Content-Type"),
	}
}
