package csv

import (
	"encoding/csv"
	"io"
)

// Scanner wraps a csv.Reader via an API similar to that of bufio.Scanner.
type Scanner struct {
	reader *csv.Reader
	record []string
	err    error

	continueOnError bool
}

// NewScanner returns a scanner configured with the provided options.
func NewScanner(reader io.Reader, options ...Option) *Scanner {
	return new(Scanner).initialize(reader).configure(options)
}
func (s *Scanner) initialize(reader io.Reader) *Scanner {
	s.reader = csv.NewReader(reader)
	return s
}
func (s *Scanner) configure(options []Option) *Scanner {
	for _, configure := range options {
		configure(s)
	}
	return s
}

// Scan advances the Scanner to the next record, which will then be available
// through the Record method. It returns false when the scan stops, either by
// reaching the end of the input or an error. After Scan returns false, the
// Error method will return any error that occurred during scanning, except
// that if it was io.EOF, Error will return nil.
func (s *Scanner) Scan() bool {
	if s.eof() {
		return false
	}
	s.record, s.err = s.reader.Read()
	return !s.eof()
}

func (s *Scanner) eof() bool {
	if s.err == io.EOF {
		return true
	}
	if s.err == nil {
		return false
	}
	return !s.continueOnError
}

// Record returns the most recent record generated by a call to Scan as a
// []string. See *csv.Reader.ReuseRecord for details on the strategy for
// reusing the underlying array: https://golang.org/pkg/encoding/csv/#Reader
func (s *Scanner) Record() []string {
	return s.record
}

// Error returns the last non-nil error produced by Scan (if there was one).
// It will not ever return io.EOF. s method may be called at any point
// during or after scanning but the underlying err will be reset by each call
// to Scan.
func (s *Scanner) Error() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
