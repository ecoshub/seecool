package seecool

import (
	"errors"
	"fmt"
	"jin"
	"penman"
	"strings"
)

var (
	emptyFile error = errors.New("File is empty.")
)

func CsvToJson(file string) ([]byte, error) {
	rl, err := penman.Reader(file)
	if err != nil {
		return nil, err
	}
	defer rl.Close()

	line := rl.Next()

	if line == nil {
		return nil, emptyFile
	}

	columns := strings.Split(string(line), ",")
	columScheme := jin.MakeScheme(columns...)

	arr := jin.MakeEmptyArray()
	line = rl.Next()
	for line != nil {
		cols := strings.Split(string(line), ",")
		json := columScheme.MakeJsonString(cols...)
		arr, err = jin.Add(arr, json)
		if err != nil {
			return nil, err
		}
		line = rl.Next()
	}
	return arr, nil
}

func CsvToJsonNoHeader(file string) ([]byte, error) {
	rl, err := penman.Reader(file)
	if err != nil {
		return nil, err
	}
	defer rl.Close()
	line := rl.Next()

	if line == nil {
		return nil, emptyFile
	}

	firstLine := strings.Split(string(line), ",")

	columns := make([]string, len(firstLine))
	temp := "column_"
	for i := 0; i < len(firstLine); i++ {
		columns[i] = fmt.Sprintf("%v%v", temp, i+1)
	}
	columScheme := jin.MakeScheme(columns...)
	arr := jin.MakeEmptyArray()

	cols := strings.Split(string(line), ",")
	json := columScheme.MakeJsonString(cols...)
	arr, err = jin.Add(arr, json)

	line = rl.Next()
	for line != nil {
		cols := strings.Split(string(line), ",")
		json := columScheme.MakeJsonString(cols...)
		arr, err = jin.Add(arr, json)
		if err != nil {
			return nil, err
		}
		line = rl.Next()
	}
	return arr, nil
}
