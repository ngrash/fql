package main

import (
	"os"
	"fmt"
)

func main() {
	query := ParseQuery(os.Args[1])

	file, err := os.Open(os.Args[2])
	if err != nil {
		panic("Failed to open file")
	}

	defer file.Close()

	reader := NewJSONReader(file)
	output := NewPivotOutput()
	for  {
		row := reader.Read()
		if row == nil {
			break
		}

		if query.Eval(row) {
			fmt.Print(output.Format(reader, row))
		}
	}
}
