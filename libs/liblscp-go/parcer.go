package liblscp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type EmptyError struct {
	Message string
}

func (e *EmptyError) Error() string {
	return fmt.Sprintf("data is empty: %s", e.Message)
}

type LscpError struct {
	Code    int
	Message string
}

func (e *LscpError) Error() string {
	return fmt.Sprintf("code: %d message: '%s'", e.Code, e.Message)
}

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
func parseFloat(s string) (float32, error) {
	i, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, fmt.Errorf("not float: %s %w", s, err)
	}
	return float32(i), nil
}

// Parses a comma separated list with boolean values
func parceBoolList(list string) ([]bool, error) {
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
func parseFloatList(list string) ([]float32, error) {
	ar := strings.Split(list, ",")
	bar := make([]float32, len(ar))
	for i, v := range ar {
		b, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
		if err != nil {
			return nil, fmt.Errorf("not float: %s %w", v, err)
		}
		bar[i] = float32(b)
	}
	return bar, nil
}

// Parses a comma separated list whose items are encapsulated into curly braces.
func parseArray(list string) ([]string, error) {
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
	return parseEscapedStringList(list, ",")
}

// Parses a string list, which elements contains escaped sequences.
func parseEscapedStringList(list string, sep string) ([]string, error) {
	unescaped, err := strconv.Unquote(list)
	if err != nil {
		return nil, fmt.Errorf("can't unescape: %s %w", list, err)
	}
	return strings.Split(unescaped, sep), nil
}

// Parses a list whose items are encapsulated into apostrophes.
func parseStringList(list string, sep string) ([]string, error) {
	return parseEscapedStringList(list, sep)
}

// Gets the type of the parameter represented by the specified result set.
// resultSet A string array containing the information categories of a multi-line result set.
func parseType(resultSet []string) (ParameterType, error) {
	if resultSet == nil || len(resultSet) == 0 {
		return ptUnknown, nil
	}
	for _, s := range resultSet {
		if strings.HasPrefix(s, "TYPE: ") {
			tname, _ := strings.CutPrefix(s, "TYPE: ")
			return ParameterToType[tname], nil

		}
	}
	return ptUnknown, nil
}

// Determines whether the parameter represented by the specified result set allows only one value or a list of values.
// resultSet a String array containing the information categories of a multi-line result set.
// return false if the parameter represented by the specified result set allows only one value and true if allows a list of values.
func parseMultiplicity(resultSet []string) (bool, error) {
	if resultSet == nil || len(resultSet) == 0 {
		return false, &EmptyError{"resultSet"}
	}
	for _, s := range resultSet {
		if strings.HasPrefix(s, "MULTIPLICITY: ") {
			b, _ := strings.CutPrefix(s, "MULTIPLICITY: ")
			return strconv.ParseBool(strings.TrimSpace(b))
		}
	}
	return false, nil
}

// Parses an empty result set and returns an appropriate ResultSet object.
// Notice that the result set may be of type warning or error.
// ln A <code>String</code> representing the single line result set to be parsed.
func parseError(ln string, rs *ResultSet) error {
	m, f := strings.CutPrefix(ln, "ERR:")
	if !f {
		return fmt.Errorf("not an error result: '%s'", ln)
	}
	code, msg, f := strings.Cut(m, ":")
	if !f {
		return fmt.Errorf("cant fine error code: '%s'", ln)
	}
	i, err := strconv.Atoi(code)
	if err != nil {
		return fmt.Errorf("code not int: '%s' %w", code, err)
	}
	rs.Type = ResultType.Error
	rs.Code = i
	rs.Message = msg
	return nil
}

func cutIndex(ln string) (index int, found bool, msg string, err error) {
	index = -1
	found = false
	ind, f := strings.CutPrefix(ln, "[")
	if f {
		ind, msg, f = strings.Cut(ind, "]")
		if !f {
			return 0, false, "", fmt.Errorf("cant parse index: '%s'", ln)
		}
		idx, err := strconv.Atoi(ind)
		if err != nil {
			return 0, false, "", fmt.Errorf("index not int: '%s' %w", ind, err)
		}
		index = idx
		found = true
	} else {
		msg = ln
	}
	return
}

// Parses warning message.
// ln The warning message to be parsed.
// rs A <code>ResultSet</code> instance where the warning must be stored.
func parseWarning(ln string, rs *ResultSet) error {
	_, msg, f := strings.Cut(ln, "WRN")
	if !f {
		return fmt.Errorf("not a warning result: '%s'", ln)
	}

	idx, f, msg, err := cutIndex(msg)
	if err != nil {
		return err
	}
	if f {
		rs.Index = idx
	}

	msg, f = strings.CutPrefix(msg, ":")
	if !f {
		return fmt.Errorf("cant parse code: '%s'", ln)
	}
	code, msg, f := strings.Cut(msg, ":")
	if !f {
		return fmt.Errorf("cant fine error code: '%s'", ln)
	}
	i, err := strconv.Atoi(code)
	if err != nil {
		return fmt.Errorf("code not int: '%s' %w", code, err)
	}
	rs.Type = ResultType.Warning
	rs.Code = i
	rs.Message = msg
	return nil

}

func parseOk(ln string, rs *ResultSet) error {
	_, msg, f := strings.Cut(ln, "OK")
	if !f {
		return fmt.Errorf("not an OK result: '%s'", ln)
	}
	if len(msg) == 0 {
		// it's empty "OK" result
		return nil
	}
	idx, f, msg, err := cutIndex(msg)
	if err != nil {
		return err
	}
	if f {
		rs.Index = idx
	}
	rs.Type = ResultType.Ok
	rs.Message = msg
	// it's "OK" result with index
	return nil
}
