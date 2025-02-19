package liblscp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parses an integer value.
// @throws LscpException If the string does not contain valid integer value.
func parseInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not int: %s %w", s, err)
	}
	return i, nil
}

// Parses a float value.
func parseFloat(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("not float: %s %w", s, err)
	}
	return i, nil
}

// Parses a comma separated list with boolean values
func ParceBoolList(list string) ([]bool, error) {
	ar := strings.Split(list, ",")
	bar := make([]bool, len(ar))
	for i, v := range ar {
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return nil, err
		}
		bar[i] = b
	}
	return bar, nil
}

// Parses a comma separated list with integer values
func parseIntList(list string) ([]int, error) {
	ar := strings.Split(list, ",")
	bar := make([]int, len(ar))
	for i, v := range ar {
		b, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return nil, fmt.Errorf("not int: %s %w", v, err)
		}
		bar[i] = b
	}
	return bar, nil
}

// Parses a comma separated list with float values.
func parseFloatList(list string) ([]float64, error) {
	ar := strings.Split(list, ",")
	bar := make([]float64, len(ar))
	for i, v := range ar {
		b, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return nil, fmt.Errorf("not float: %s %w", v, err)
		}
		bar[i] = b
	}
	return bar, nil
}

// Parses a comma separated list whose items are encapsulated into curly braces.
func ParseArray(list string) ([]string, error) {
	pattern := regexp.MustCompile(`\{([^}]*)\}`)
	matches := pattern.FindAllString(list, -1)

	bp := regexp.MustCompile(`[{}]`)
	for i, v := range matches {
		matches[i] = bp.ReplaceAllString(v, "")
	}
	return matches, nil

}

// Parses a comma separated string list, which elements contains escaped sequences.
func parseEscapedStringListComma(list string) ([]string, error) {
	return ParseEscapedStringList(list, ",")
}

// Parses a string list, which elements contains escaped sequences.
func ParseEscapedStringList(list string, sep string) ([]string, error) {
	unescaped, err := strconv.Unquote(list)
	if err != nil {
		return nil, fmt.Errorf("can't unescape: %s %w", list, err)
	}
	return strings.Split(unescaped, sep), nil
}

// Parses a list whose items are encapsulated into apostrophes.
func ParseStringList(list string, sep string) ([]string, error) {
	return ParseEscapedStringList(list, sep)
}

// Gets the type of the parameter represented by the specified result set.
// resultSet A string array containing the information categories of a multi-line result set.
func parseType(resultSet []string) (ParameterType, error) {
	if resultSet == nil || len(resultSet) == 0 {
		return Unknown, nil
	}
	for _, s := range resultSet {
		if strings.HasPrefix(s, "TYPE: ") {
			tname, _ := strings.CutPrefix(s, "TYPE: ")
			return ParameterToType[tname], nil

		}
	}
	return Unknown, nil
}
