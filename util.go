package seecool

import (
	"strconv"
	"strings"
	"unsafe"
)

func formatType(val string) string {
	if len(val) > 0 {
		if isBool(val) {
			return val
		}
		if isInt(val) {
			if val[0] == 48 && len(val) > 1 {
				return `"` + val + `"`
			}
			return val
		}
		if isFloat(val) {
			return val
		}
		if val == "null" {
			return val
		}
		start := val[0]
		end := val[len(val)-1]
		if (start == 34 && end == 34) || (start == 91 && end == 93) || (start == 123 && end == 125) {
			return val
		}
		return `"` + val + `"`
	}
	return `""`
}

func isBool(val string) bool {
	return val == "true" || val == "false"
}

func isFloat(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}
	return true
}

func isInt(val string) bool {
	_, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false
	}
	return true
}

func stringToByteArray(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
func wordSplit(line string) []string {
	lenl := len(line)
	if lenl < 5 {
		return nil
	}
	tokens := make([]string, 0, 3)
	onWord := false
	temp := ""
	for i := 0; i < lenl; i++ {
		curr := line[i]
		if !space(curr) {
			temp += string(curr)
			if !onWord {
				onWord = true
				continue
			}
		} else {
			if onWord {
				tokens = append(tokens, temp)
				temp = ""
				onWord = false
				continue
			}
		}
	}
	if temp != "" {
		tokens = append(tokens, temp)
	}
	return tokens
}

func space(curr byte) bool {
	// space
	if curr == 32 {
		return true
	}
	// tab
	if curr == 9 {
		return true
	}
	// new line NL
	if curr == 10 {
		return true
	}
	// return CR
	if curr == 13 {
		return true
	}
	return false
}

func arrStr(arr []string) string {
	lena := len(arr)
	switch lena {
	case 0:
		return "*"
	case 1:
		return arr[0]
	default:
		str := ""
		for i := 0; i < lena-1; i++ {
			str += arr[i] + ", "
		}
		str += arr[lena-1]
		return str
	}
}

func inQuote(arr []string) []string {
	for i := 0; i < len(arr); i++ {
		if arr[i][0] != '(' {
			arr[i] = `'` + arr[i] + `'`
		}
	}
	return arr
}

func columnCheck(cols []string) []string {
	for i, c := range cols {
		prefix, rest := astrixCheck(c)
		switch prefix {
		case "":
		case "c":
			cols[i] = "COUNT(" + rest + ")"
		case "a":
			cols[i] = "AVG(" + rest + ")"
		case "mi":
			cols[i] = "MIN(" + rest + ")"
		case "ma":
			cols[i] = "MAX(" + rest + ")"
		case "s":
			cols[i] = "SUM(" + rest + ")"
		default:
			cols[i] = strings.ToUpper(prefix) + "(" + rest + ")"
		}
	}
	return cols
}

func astrixCheck(str string) (string, string) {
	for i := 0; i < len(str); i++ {
		if str[i] == 42 {
			return str[:i], str[i+1:]
		}
	}
	return "", str
}

func EscapeQuote(str string) string {
	newStr := make([]byte, 0, len(str))
	for i := 0; i < len(str); i++ {
		// double quote
		if str[i] == 34 {
			// single quote
			newStr = append(newStr, 39)
		} else {
			newStr = append(newStr, str[i])
		}
	}
	return string(newStr)
}
