package csv

import (
	"errors"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestColumnScannerFixture(t *testing.T) {
	gunit.Run(new(ColumnScannerFixture), t)
}

type ColumnScannerFixture struct {
	*gunit.Fixture

	scanner *ColumnScanner
	err     error
	users   []User
}

func (csf *ColumnScannerFixture) Setup() {
	csf.scanner, csf.err = NewColumnScanner(reader(csvCanon))
	csf.So(csf.err, should.BeNil)
	csf.So(csf.scanner.Header(), should.Resemble, []string{"first_name", "last_name", "username"})
}

func (csf *ColumnScannerFixture) ScanAllUsers() {
	for csf.scanner.Scan() {
		csf.users = append(csf.users, csf.scanUser())
	}
}

func (csf *ColumnScannerFixture) TestReadColumns() {
	csf.ScanAllUsers()

	csf.So(csf.scanner.Error(), should.BeNil)
	csf.So(csf.users, should.Resemble, []User{
		{FirstName: "Rob", LastName: "Pike", Username: "rob"},
		{FirstName: "Ken", LastName: "Thompson", Username: "ken"},
		{FirstName: "Robert", LastName: "Griesemer", Username: "gri"},
	})
}

func (csf *ColumnScannerFixture) scanUser() User {
	return User{
		FirstName: csf.scanner.Column(csf.scanner.Header()[0]),
		LastName:  csf.scanner.Column(csf.scanner.Header()[1]),
		Username:  csf.scanner.Column(csf.scanner.Header()[2]),
	}
}

func (csf *ColumnScannerFixture) TestCannotReadHeader() {
	scanner, err := NewColumnScanner(new(ErrorReader))
	csf.So(scanner, should.BeNil)
	csf.So(err, should.NotBeNil)
}

func (csf *ColumnScannerFixture) TestColumnNotFound_Error() {
	csf.scanner.Scan()
	value, err := csf.scanner.ColumnErr("nope")
	csf.So(value, should.BeBlank)
	csf.So(err, should.NotBeNil)
}

func (csf *ColumnScannerFixture) TestColumnNotFound_Panic() {
	csf.scanner.Scan()
	csf.So(func() { csf.scanner.Column("nope") }, should.Panic)
}

type User struct {
	FirstName string
	LastName  string
	Username  string
}

type ErrorReader struct{}

func (csf *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("ERROR")
}
