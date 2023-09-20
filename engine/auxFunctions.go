package engine

import (
	"strconv"
	"strings"
)

// DrawNumber Get drawNumber
func DrawNumber(c Controller) int {
	draw, err := strconv.Atoi(c.GetString("draw"))
	if err != nil {
		return 0
	}
	return draw
}

// ParamToBool get a param and parse it to bool
func ParamToBool(c Controller, paramName string) (requestRegex bool) {
	requestRegex, err := strconv.ParseBool(c.GetString(paramName))
	if err != nil {
		requestRegex = false
	}
	return
}

// CheckReserved Skip reserved words
func CheckReserved(columnName string) string {
	if isReserved(columnName) {
		return "\"" + columnName + "\""
	}
	return columnName
}

func isReserved(text string) bool {
	for _, item := range ReservedWords {
		if item == strings.ToUpper(text) {
			return true
		}
	}
	return false
}
