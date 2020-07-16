package csv

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestScanAllFixture(t *testing.T) {
	gunit.Run(new(ScanAllFixture), t)
}

type ScanAllFixture struct {
	*gunit.Fixture
}

func (saf *ScanAllFixture) scanAll(inputs []string, options ...Option) (scanned []Record) {
	scanner := NewScanner(reader(inputs), options...)
	line := 1
	for ; scanner.Scan(); line++ {
		scanned = append(scanned, Record{
			line:   line,
			record: scanner.Record(),
			err:    scanner.Error(),
		})
	}
	if err := scanner.Error(); err != nil {
		scanned = append(scanned, Record{
			line: line,
			err:  err,
		})
	}
	return scanned
}

func (saf *ScanAllFixture) TestCanonical() {
	scanned := saf.scanAll(csvCanon, Comma(','), FieldsPerRecord(3))
	saf.So(scanned, should.Resemble, expectedScannedOutput)
}

func (saf *ScanAllFixture) TestCanonicalWithOptions() {
	scanned := saf.scanAll(csvCanonRequiringConfigOptions, Comma(';'), Comment('#'))
	saf.So(scanned, should.Resemble, expectedScannedOutput)
}

func (saf *ScanAllFixture) TestOptions() {
	scanner := NewScanner(nil, ReuseRecord(true), TrimLeadingSpace(true), LazyQuotes(true))
	saf.So(scanner.reader.ReuseRecord, should.BeTrue)
	saf.So(scanner.reader.LazyQuotes, should.BeTrue)
	saf.So(scanner.reader.TrimLeadingSpace, should.BeTrue)
}

func (saf *ScanAllFixture) TestInconsistentFieldCounts_ContinueOnError() {
	scanned := saf.scanAll(csvCanonInconsistentFieldCounts, ContinueOnError(true))
	saf.So(scanned, should.Resemble, []Record{
		{line: 1, record: []string{"1", "2", "3"}, err: nil},
		{line: 2, record: []string{"1", "2", "3", "4"}, err: &csv.ParseError{StartLine: 2, Line: 2, Column: 0, Err: csv.ErrFieldCount}},
		{line: 3, record: []string{"1", "2", "3"}, err: nil},
	})
}

func (saf *ScanAllFixture) TestInconsistentFieldCounts_HaltOnError() {
	scanned := saf.scanAll(csvCanonInconsistentFieldCounts)
	saf.So(scanned, should.Resemble, []Record{
		{line: 1, record: []string{"1", "2", "3"}, err: nil},
		{line: 2, record: nil, err: &csv.ParseError{StartLine: 2, Line: 2, Column: 0, Err: csv.ErrFieldCount}},
	})
}

func (saf *ScanAllFixture) TestCallsToScanAfterEOFReturnFalse() {
	scanner := NewScanner(strings.NewReader("1,2,3"), Comma(','))

	saf.So(scanner.Scan(), should.BeTrue)
	saf.So(scanner.Record(), should.Resemble, []string{"1", "2", "3"})
	saf.So(scanner.Error(), should.BeNil)

	for x := 0; x < 100; x++ {
		saf.So(scanner.Scan(), should.BeFalse)
		saf.So(scanner.Record(), should.BeNil)
		saf.So(scanner.Error(), should.BeNil)
	}
}

func (saf *ScanAllFixture) TestSkipHeader() {
	scanned := saf.scanAll(csvCanon, Comma(','), SkipHeaderRecord())
	saf.So(scanned, should.Resemble, []Record{
		{line: 1, record: []string{"Rob", "Pike", "rob"}},
		{line: 2, record: []string{"Ken", "Thompson", "ken"}},
		{line: 3, record: []string{"Robert", "Griesemer", "gri"}},
	})
}

func (saf *ScanAllFixture) TestRecords() {
	scanned := saf.scanAll(csvCanon, Comma(','), SkipRecords(3))
	saf.So(scanned, should.Resemble, []Record{
		{line: 1, record: []string{"Robert", "Griesemer", "gri"}},
	})
}

func reader(lines []string) io.Reader {
	return strings.NewReader(strings.Join(lines, "\n"))
}

var ( // https://golang.org/pkg/encoding/csv/#example_Reader
	csvCanon = []string{
		"first_name,last_name,username",
		`"Rob","Pike",rob`,
		`Ken,Thompson,ken`,
		`"Robert","Griesemer","gri"`,
	}
	csvNums = []string{
		"first_name,height,age",
		`Jim,4.2,18`,
		`Steve,1.1,9`,
		`Bart,1.0,80`,
	}
	csvCanonRequiringConfigOptions = []string{
		`first_name;last_name;username`,
		`"Rob";"Pike";rob`,
		`# lines beginning with a # character are ignored`,
		`Ken;Thompson;ken`,
		`"Robert";"Griesemer";"gri"`,
	}
	csvCanonInconsistentFieldCounts = []string{
		`1,2,3`,
		`1,2,3,4`,
		`1,2,3`,
	}
	expectedScannedOutput = []Record{
		{1, []string{"first_name", "last_name", "username"}, nil},
		{2, []string{"Rob", "Pike", "rob"}, nil},
		{3, []string{"Ken", "Thompson", "ken"}, nil},
		{4, []string{"Robert", "Griesemer", "gri"}, nil},
	}
)

type Record struct {
	line   int
	record []string
	err    error
}
