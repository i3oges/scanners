package csv

import (
	"fmt"
	"io"
	"log"
)

type ColumnScanner struct {
	*Scanner
	headerRecord []string
	columnIndex  map[string]int
}

func NewColumnScanner(reader io.Reader, options ...Option) (*ColumnScanner, error) {
	inner := NewScanner(reader, append(options, FieldsPerRecord(0))...)
	if !inner.Scan() {
		return nil, inner.Error()
	}
	scanner := &ColumnScanner{
		Scanner:      inner,
		headerRecord: inner.Record(),
		columnIndex:  make(map[string]int),
	}
	scanner.readHeader()
	return scanner, nil
}

func (cs *ColumnScanner) readHeader() {
	for i, value := range cs.headerRecord {
		cs.columnIndex[value] = i
	}
}

func (cs *ColumnScanner) Header() []string {
	return cs.headerRecord
}

func (cs *ColumnScanner) ColumnErr(column string) (string, error) {
	index, ok := cs.columnIndex[column]
	if !ok {
		return "", fmt.Errorf("Column [%s] not present in header record: %#v\n", column, cs.headerRecord)
	}
	return cs.Record()[index], nil
}

func (cs *ColumnScanner) Column(column string) string {
	value, err := cs.ColumnErr(column)
	if err != nil {
		log.Panic(err)
	}
	return value
}
