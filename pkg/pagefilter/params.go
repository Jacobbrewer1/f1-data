package pagefilter

import (
	"strings"
)

// ParamParserFunc is a function that takes a parameter value,
// validates it and parses it into an operation and query pair.
type ParamParserFunc func(string) (string, string, error)

// ParamParser returns a ParamParseFunc that will parse and validate a param value based on the provided options.
//
// Using ParamParser you can create basic validators and parser for pagefilter params. As an example, if
// you wanted to accept a `count` parameter that by default would do an equality check, but also
// supported greater than and less than checks, you could do the following:
//
//	parser := pagefilter.ParamParser("count", pagefilter.Equal, pagefilter.GreaterThan, pagefilter.LessThan)
//
// You could then use the parser to get the results, like so:
//
//	op, query, err := parser("eq:123") // op == "eq", query == "123"
func ParamParser(paramName, defaultOp string, validOps ...string) ParamParserFunc {
	vOps := make(map[string]struct{}, len(validOps))
	for _, op := range validOps {
		vOps[op] = struct{}{}
	}

	// The default op is always valid
	if defaultOp != "" {
		vOps[defaultOp] = struct{}{}
	} else if len(validOps) == 0 {
		return func(_ string) (string, string, error) {
			return "", "", InvalidParamError(paramName)
		}
	}

	return func(p string) (op, query string, err error) {
		v := strings.SplitN(p, ":", 2)
		switch len(v) {
		case 1:
			op = defaultOp
			query = v[0]
		case 2:
			op = v[0]
			query = v[1]
		}

		if query == "" {
			return "", "", MissingValueError(paramName)
		}

		if _, ok := vOps[op]; !ok {
			return "", "", InvalidOpError(op)
		}

		return op, query, nil
	}
}
