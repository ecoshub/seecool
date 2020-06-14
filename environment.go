package seecool

import (
	"errors"
	"penman"
	"strings"
)

var (
	malformedEnv error = errors.New("Malformed .env file. line format must be 'key' = 'value'.")
)

func GetEnv(dir string) (map[string]string, error) {
	dir = penman.PreProcess(dir)
	rl, err := penman.Reader(dir)
	if err != nil {
		return nil, err
	}
	defer rl.Close()
	envMap := make(map[string]string)
	line := rl.Next()
	for line != nil {
		tokens := wordSplit(string(line))
		// blank line check
		// comment sysmbol check '//'
		if tokens == nil || startsWith(tokens[0], "//") {
			line = rl.Next()
			continue
		}
		lent := len(tokens)
		// missing key or value
		if lent < 3 {
			return nil, malformedEnv
		}
		// middle char control.
		if tokens[1] != "=" {
			return nil, malformedEnv
		}
		if lent > 3 {
			// end comment control
			if startsWith(tokens[3], "//") {
				tokens = tokens[:3]
			}
		}
		envMap[tokens[0]] = tokens[2]
		line = rl.Next()
	}
	return envMap, nil
}

func GetEnvString(file string) (map[string]string, error) {
	lines := strings.Split(file, "\n")
	envMap := make(map[string]string)
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		tokens := wordSplit(string(line))
		// blank line check
		// comment sysmbol check '//'
		if tokens == nil || startsWith(tokens[0], "//") {
			continue
		}
		lent := len(tokens)
		// missing key or value
		if lent < 3 {
			return nil, malformedEnv
		}
		// middle char control.
		if tokens[1] != "=" {
			return nil, malformedEnv
		}
		if lent > 3 {
			// end comment control
			if startsWith(tokens[3], "//") {
				tokens = tokens[:3]
			}
		}
		envMap[tokens[0]] = tokens[2]
	}
	return envMap, nil
}

func startsWith(word, prefix string) bool {
	lenp := len(prefix)
	if len(word) >= lenp {
		if string(word[:lenp]) == prefix {
			return true
		}
	}
	return false
}
