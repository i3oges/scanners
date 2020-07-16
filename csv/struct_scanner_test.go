package csv

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStructScannerFixture(t *testing.T) {
	gunit.Run(new(StructScannerFixture), t)
}

type StructScannerFixture struct {
	*gunit.Fixture
	scanner *StructScanner
	numScanner *StructScanner
	err     error
	users   []TaggedUser
	numUsers []NumTaggedUser
}

func (s *StructScannerFixture) Setup() {
	s.scanner, s.err = NewStructScanner(reader(csvCanon))
	s.So(s.err, should.BeNil)
	s.numScanner, s.err = NewStructScanner(reader(csvNums))
	s.So(s.err, should.BeNil)
}

func (s *StructScannerFixture) ScanAll() {
	for s.scanner.Scan() {
		var user TaggedUser
		s.scanner.Populate(&user)
		s.users = append(s.users, user)
	}
	for s.numScanner.Scan() {
		var numUser NumTaggedUser
		s.numScanner.Populate(&numUser)
		s.numUsers = append(s.numUsers, numUser)
	}
}

func (s *StructScannerFixture) Test() {
	s.ScanAll()

	s.So(s.scanner.Error(), should.BeNil)
	s.So(s.users, should.Resemble, []TaggedUser{
		{FirstName: "Rob", LastName: "Pike", Username: "rob"},
		{FirstName: "Ken", LastName: "Thompson", Username: "ken"},
		{FirstName: "Robert", LastName: "Griesemer", Username: "gri"},
	})
	s.So(s.numUsers, should.Resemble, []NumTaggedUser{
		{FirstName: "Jim", Age: 18, Height: 4.20},
		{FirstName: "Steve", Age: 9, Height: 1.1},
		{FirstName: "Bart", Age: 80, Height: 1.0},
	})
}

type TaggedUser struct {
	FirstName string `csv:"first_name"`
	LastName  string `csv:"last_name"`
	Username  string `csv:"username"`
}

type NumTaggedUser struct {
	FirstName string `csv:"first_name"`
	Height  float64 `csv:"height"`
	Age  int64 `csv:"age"`
}

func (s *StructScannerFixture) TestCannotReadHeader() {
	scanner, err := NewStructScanner(new(ErrorReader))
	s.So(scanner, should.BeNil)
	s.So(err, should.NotBeNil)
}

func (s *StructScannerFixture) TestScanIntoLessCompatibleType() {
	s.scanner.Scan()

	var nonPointer User
	s.So(s.scanner.Populate(nonPointer), should.NotBeNil)

	var nilPointer *User
	s.So(s.scanner.Populate(nilPointer), should.NotBeNil)
}
