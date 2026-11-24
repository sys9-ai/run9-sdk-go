package run9

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	tarBlockSize = 512
	// tarMaxPAXHeader bounds one buffered PAX extended header, as archive/tar does.
	tarMaxPAXHeader = 1 << 20
)

// fileArchiveReader passes a gateway archive through unchanged and reports
// io.ErrUnexpectedEOF if the stream ends before tar's end-of-archive marker.
//
// The gateway writes the marker only after a complete archive. Readers such as
// archive/tar accept a stream that ends right after an entry as complete, and
// a hop that fails, for example a proxy or reader process that exits, can end
// the response cleanly at such a boundary. Following the block framing tells
// the two apart: each 512-byte header is followed by its content rounded up to
// whole blocks, where a PAX extended header may replace the next entry's size.
type fileArchiveReader struct {
	body      io.ReadCloser
	header    [tarBlockSize]byte
	filled    int    // Bytes of the current header block received.
	skip      int64  // Content and padding bytes left before the next header.
	pax       []byte // Content of the PAX extended header being received.
	paxLeft   int64  // Bytes of that content not yet received.
	paxPad    int64  // Padding after that content.
	nextSize  int64  // Size a PAX header gave the next entry, or -1.
	zeroes    int    // Consecutive all-zero header blocks.
	completed bool   // Two zero blocks: the end-of-archive marker.
	err       error  // A framing error, reported once the data before it is read.
}

func newFileArchiveReader(body io.ReadCloser) *fileArchiveReader {
	return &fileArchiveReader{body: body, nextSize: -1}
}

// Read reads the archive, failing with io.ErrUnexpectedEOF in place of an
// io.EOF that comes before the end-of-archive marker.
func (reader *fileArchiveReader) Read(buffer []byte) (int, error) {
	if reader.err != nil {
		return 0, reader.err
	}
	n, err := reader.body.Read(buffer)
	reader.err = reader.follow(buffer[:n])
	if err == io.EOF {
		switch {
		case reader.err != nil:
			err = reader.err
		case !reader.completed:
			err = io.ErrUnexpectedEOF
		}
	}
	return n, err
}

// Close releases the response.
func (reader *fileArchiveReader) Close() error {
	return reader.body.Close()
}

func (reader *fileArchiveReader) follow(data []byte) error {
	for len(data) > 0 && !reader.completed {
		switch {
		case reader.paxLeft > 0:
			received := min(reader.paxLeft, int64(len(data)))
			reader.pax = append(reader.pax, data[:received]...)
			reader.paxLeft -= received
			data = data[received:]
			if reader.paxLeft == 0 {
				size, err := parsePAXSize(reader.pax)
				if err != nil {
					return err
				}
				reader.nextSize, reader.skip = size, reader.paxPad
			}
		case reader.skip > 0:
			skipped := min(reader.skip, int64(len(data)))
			reader.skip -= skipped
			data = data[skipped:]
		default:
			copied := copy(reader.header[reader.filled:], data)
			reader.filled += copied
			data = data[copied:]
			if reader.filled == tarBlockSize {
				reader.filled = 0
				if err := reader.endHeader(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// endHeader interprets one complete header block.
func (reader *fileArchiveReader) endHeader() error {
	if bytes.Equal(reader.header[:], make([]byte, tarBlockSize)) {
		reader.zeroes++
		reader.completed = reader.zeroes == 2
		return nil
	}
	if reader.zeroes != 0 {
		return fmt.Errorf("archive has an entry after a zero block: %w", io.ErrUnexpectedEOF)
	}
	size, err := parseTarSize(reader.header[124:136])
	if err != nil {
		return err
	}
	padded := func(size int64) int64 { return (size + tarBlockSize - 1) / tarBlockSize * tarBlockSize }
	switch typeflag := reader.header[156]; typeflag {
	case 'x':
		if size > tarMaxPAXHeader {
			return fmt.Errorf("archive PAX header of %d bytes exceeds %d", size, tarMaxPAXHeader)
		}
		reader.pax, reader.paxLeft, reader.paxPad = reader.pax[:0], size, padded(size)-size
	case 'g':
		reader.skip = padded(size)
	default:
		if reader.nextSize >= 0 {
			size, reader.nextSize = reader.nextSize, -1
		}
		// As in archive/tar, these entries carry no content whatever their size.
		if strings.IndexByte("123456", typeflag) >= 0 {
			size = 0
		}
		reader.skip = padded(size)
	}
	return nil
}

// parseTarSize decodes a header's size field: octal digits, or a big-endian
// base-256 number when the high bit of the first byte is set.
func parseTarSize(field []byte) (int64, error) {
	if field[0]&0x80 != 0 {
		// The 11 value bytes must fit a non-negative int64.
		if field[0] != 0x80 || field[1]|field[2]|field[3] != 0 || field[4]&0x80 != 0 {
			return 0, fmt.Errorf("archive entry size %x is out of range", field)
		}
		var size int64
		for _, value := range field[4:] {
			size = size<<8 | int64(value)
		}
		return size, nil
	}
	text := strings.Trim(string(field), " \x00")
	if text == "" {
		return 0, nil
	}
	size, err := strconv.ParseInt(text, 8, 64)
	if err != nil || size < 0 {
		return 0, fmt.Errorf("archive entry size %q is invalid", text)
	}
	return size, nil
}

// parsePAXSize returns the size record of PAX extended header content made of
// "<length> <key>=<value>\n" records, or -1 when it has none.
func parsePAXSize(records []byte) (int64, error) {
	size := int64(-1)
	for len(records) > 0 {
		lengthText, _, found := bytes.Cut(records, []byte(" "))
		length, err := strconv.Atoi(string(lengthText))
		if !found || err != nil || length <= len(lengthText)+1 || length > len(records) || records[length-1] != '\n' {
			return 0, fmt.Errorf("archive PAX header has an invalid record")
		}
		key, value, found := strings.Cut(string(records[len(lengthText)+1:length-1]), "=")
		if !found {
			return 0, fmt.Errorf("archive PAX header has an invalid record")
		}
		if key == "size" {
			size, err = strconv.ParseInt(value, 10, 64)
			if err != nil || size < 0 {
				return 0, fmt.Errorf("archive PAX size %q is invalid", value)
			}
		}
		records = records[length:]
	}
	return size, nil
}
