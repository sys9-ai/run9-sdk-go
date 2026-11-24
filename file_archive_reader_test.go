package run9

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

// testTarArchive is a complete gateway-style archive: PAX headers, a
// directory, a file whose content ends mid-block, a symlink and a hard link.
func testTarArchive(t *testing.T) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := tar.NewWriter(&buffer)
	for _, header := range []*tar.Header{
		{Typeflag: tar.TypeDir, Name: "out/", Mode: 0o755},
		{Typeflag: tar.TypeReg, Name: "out/file", Mode: 0o644, Size: 1000},
		{Typeflag: tar.TypeSymlink, Name: "out/link", Linkname: "file"},
		{Typeflag: tar.TypeLink, Name: "out/hard", Linkname: "out/file"},
	} {
		header.Format = tar.FormatPAX
		require.NoError(t, writer.WriteHeader(header))
		_, err := writer.Write(bytes.Repeat([]byte{'x'}, int(header.Size)))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	return buffer.Bytes()
}

func readFileArchive(data []byte) ([]byte, error) {
	return io.ReadAll(newFileArchiveReader(io.NopCloser(bytes.NewReader(data))))
}

func TestFileArchiveReaderRequiresEndOfArchiveMarker(t *testing.T) {
	archive := testTarArchive(t)
	data, err := readFileArchive(archive)
	require.NoError(t, err)
	require.Equal(t, archive, data, "the archive passes through unchanged")

	// Cut cleanly after the last entry: archive/tar reads this as complete.
	atBoundary := archive[:len(archive)-1024]
	entries := tar.NewReader(bytes.NewReader(atBoundary))
	for {
		_, err := entries.Next()
		if err != nil {
			require.ErrorIs(t, err, io.EOF)
			break
		}
	}
	_, err = readFileArchive(atBoundary)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)

	_, err = readFileArchive(archive[:len(archive)-512])
	require.ErrorIs(t, err, io.ErrUnexpectedEOF, "one zero block is not the marker")
	_, err = readFileArchive(archive[:700])
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	_, err = readFileArchive(nil)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// A file of 8 GiB or more keeps its size only in its PAX header; the framing
// must follow that size, not the zero in the ustar field.
func TestFileArchiveReaderFollowsPAXSize(t *testing.T) {
	block := func(name string, typeflag byte, size int) []byte {
		header := make([]byte, 512)
		copy(header, name)
		copy(header[124:], fmt.Sprintf("%011o", size))
		header[156] = typeflag
		return header
	}
	records := "12 size=600\n"
	var archive []byte
	archive = append(archive, block("PaxHeaders/big", 'x', len(records))...)
	archive = append(archive, append([]byte(records), make([]byte, 512-len(records))...)...)
	archive = append(archive, block("big", '0', 0)...)
	// Content whose first block is zero: read as headers, it would be a zero
	// block followed by an entry, which fails.
	archive = append(archive, make([]byte, 512)...)
	archive = append(archive, append(bytes.Repeat([]byte{'x'}, 88), make([]byte, 512-88)...)...)
	archive = append(archive, make([]byte, 1024)...)

	_, err := readFileArchive(archive)
	require.NoError(t, err)
	_, err = readFileArchive(archive[:len(archive)-1024])
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// A framing error that arrives with the end of the body is reported, not
// replaced by io.EOF.
func TestFileArchiveReaderReportsFramingErrorAtEOF(t *testing.T) {
	header := make([]byte, 512)
	copy(header, "out")
	copy(header[124:], "zz")
	reader := newFileArchiveReader(io.NopCloser(iotest.DataErrReader(bytes.NewReader(header))))
	_, err := io.ReadAll(reader)
	require.EqualError(t, err, `archive entry size "zz" is invalid`)
}
