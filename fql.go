package main

import (
	"os"
	"fmt"
)

func main() {
	query := ParseQuery(os.Args[1])
	//fmt.Println(listener.result)

	file, err := os.Open(os.Args[2])
	if err != nil {
		panic("Failed to open file")
	}

	defer file.Close()

	reader := NewJSONReader(file)
	for  {
		row := reader.Read()
		if row == nil {
			break
		}

		if query.Eval(row) {
			fmt.Println()
			fmt.Println(row)
		}
	}
}
