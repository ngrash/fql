package main

import (
	"encoding/json"
	"log"
	"strings"
	"fmt"
)

type Output interface {
	Format(reader Reader, row Row) string
}

/*
 * JSON Output
 */

type JSONOutput struct {}

func NewJSONOutput() Output {
	return &JSONOutput{}
}

func (o *JSONOutput) Format(reader Reader, row Row) string {
	bytes, err := json.Marshal(reader.Fields(row))
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

/*
 * CSV Output
 */

type CSVOutput struct {}

func NewCSVOutput() Output {
	return &CSVOutput{}
}

func (o *CSVOutput) Format(reader Reader, row Row) string {
	// TODO: Converting JSON to CSV might result in CSV rows with missing
	// columns as JSON structures might have different fields present
	fields := reader.Fields(row)
	values := make([]string, len(fields))

	i := 0
	for _, value := range fields {
		values[i] = value
		i++
	}

	return strings.Join(values, ";")
}

/*
 * Pivot Table Output
 */

type PivotOutput struct {}

func NewPivotOutput() Output {
	return &PivotOutput{}
}

func (o *PivotOutput) Format(reader Reader, row Row) string {
	const maxLineLen = 79

	fields := reader.Fields(row)

	var str strings.Builder

	for key, value := range fields {
		str.WriteString(fmt.Sprintf("%v: ", key))
		indent := len(key) + 2 // +2 for ": "

		for i, paragraph := range strings.Split(value, "\n") {
			lineLen := indent
			if i > 0 {
				str.WriteString(strings.Repeat(" ", indent))
			}

			for _, word := range strings.SplitAfter(paragraph, " ") {
				wordLen := len(word)
				if lineLen + wordLen > maxLineLen {
					str.WriteString("\n")
					str.WriteString(strings.Repeat(" ", indent))
					lineLen = indent
				}
				str.WriteString(word)
				lineLen += wordLen
			}
			str.WriteString("\n")
		}
	}

	str.WriteString("\n")

	return str.String()
}
