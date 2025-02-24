package multibytes

import "io"

var (
	_ io.Reader = (*MultiBytes)(nil)
	_ io.Writer = (*MultiBytes)(nil)
)

type MultiBytes struct {
	data  [][]byte
	index int
	pos   int
}

func NewMultiBytes(data [][]byte) *MultiBytes {
	return &MultiBytes{data: data}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call. If some data is available but not len(p) bytes, Read
// conventionally returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after successfully reading n > 0 bytes, it returns the number of bytes read. It
// may return the (non-nil) error from the same call or return the error (and n == 0) from a subsequent call. An instance of this general
// case is that a Reader returning a non-zero number of bytes at the end of the input stream may return either err == EOF or err == nil.
// The next Read should return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before considering the error err. Doing so correctly handles I/O errors that happen
// after reading some bytes and also both of the allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a zero byte count with a nil error, except when len(p) == 0. Callers should treat a
// return of 0 and nil as indicating that nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func (m *MultiBytes) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if m.index >= len(m.data) {
		return 0, io.EOF
	}

	for n < len(p) {
		if m.pos >= len(m.data[m.index]) {
			m.index++
			m.pos = 0
			if m.index >= len(m.data) {
				break
			}
		}

		bs := m.data[m.index]
		cnt := copy(p[n:], bs[m.pos:])
		m.pos += cnt
		n += cnt
	}

	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

// Write writes len(p) bytes from p to the underlying data stream. It returns the number of bytes written from p (0 <= n <= len(p)) and any
// error encountered that caused the write to stop early. Write must return a non-nil error if it returns n < len(p). Write must not modify
// the slice data, even temporarily.
//
// Implementations of Write must not retain p.
func (m *MultiBytes) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	clone := make([]byte, len(p))
	copy(clone, p)
	m.data = append(m.data, clone)
	return len(p), nil
}
