package main

import (
	"fmt"
	"bufio"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"strings"
)

type Row interface {
	Value(key string) string
	Values() []string
}

type Reader interface {
	Read() Row
}

/*
 * Generic Row
 */

type MemoryRow struct {
	columns *ColumnData
	values []string
}

func (r *MemoryRow) String() string {
	fields := make([]string, len(r.columns.indices))
	for column, index := range r.columns.indices {
		fields[index] = fmt.Sprintf("%v=%v", column, r.values[index])
	}
	return fmt.Sprintf("{%v}", strings.Join(fields, ", "))
}

func (r *MemoryRow) Value(key string) string {
	i := r.columns.Index(key)
	v := r.values[i]
	return v
}

func (r *MemoryRow) Values() []string {
	return r.values
}

type ColumnData struct {
	indices map[string]uint8
}

func (c *ColumnData) Index(key string) uint8 {
	index, ok := c.indices[key]
	if !ok {
		panic(fmt.Sprintf("No such column: %v", key))
	}

	return index
}

/*
 * CSV Reader
 */

type CSVReader struct {
	columns ColumnData
	scanner *bufio.Scanner
}

func NewCSVReader(reader io.Reader) Reader {
	scanner := bufio.NewScanner(reader)

	// first line contains columns
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		panic("Failed to read columns")
	}

	columnsLine := scanner.Text()
	columns := strings.Split(columnsLine, ";")

	csvReader := CSVReader{
		columns: ColumnData{indices: make(map[string]uint8)},
		scanner: scanner}

	for i, column := range columns {
		csvReader.columns.indices[column] = uint8(i)
	}

	return &csvReader
}

func (r *CSVReader) Read() Row {
	if ok := r.scanner.Scan(); ok {
		line := r.scanner.Text()
		return &MemoryRow{
			values: strings.Split(line, ";"),
			columns: &r.columns}
	}

	if err := r.scanner.Err(); err != nil {
		panic("Failed to read line")
	}

	return nil // EOF
}

/*
 * JSON Reader
 */

type JSONReader struct {
	data []map[string]string
	index int
}

type JSONRow struct {
	values map[string]string
}

func (r JSONRow) Value(key string) string {
	return r.values[key]
}

func (r JSONRow) Values() []string {
	values := make([]string, len(r.values))
	i := 0
	for _, value := range r.values {
		values[i] = value
		i++
	}

	return values
}

func NewJSONReader(reader io.Reader) Reader {
	blob, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Panicf("Failed to read JSON: %v\n", err)
	}

	var array []map[string]string
	if err := json.Unmarshal(blob, &array); err != nil {
		log.Panicln(err)
	}

	return &JSONReader{array, 0}
}

func (r *JSONReader) Read() Row {
	if r.index >= len(r.data) {
		return nil
	}

	values := r.data[r.index]
	r.index++
	return JSONRow{values}
}
