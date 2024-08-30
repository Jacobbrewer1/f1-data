package pagefilter

import (
	"fmt"
	"testing"
)

func TestParamParser(t *testing.T) {
	testEq := ParamParser("test", Equal)
	test2All := ParamParser("test2", Equal, GreaterThan, LessThan)
	test3NoDefault := ParamParser("test3", "", Equal)
	test4Impossible := ParamParser("test4", "")
	testLike := ParamParser("test", Like)

	tests := []struct {
		name    string
		parser  ParamParserFunc
		param   string
		op      string
		value   string
		wantErr error
	}{
		{
			name:    "invalid op",
			parser:  testEq,
			param:   "invalid:value",
			wantErr: InvalidOpError("invalid"),
		},
		{
			name:   "valid op",
			parser: testEq,
			param:  "eq:value",
			op:     Equal,
			value:  "value",
		},
		{
			name:    "missing value",
			parser:  testEq,
			param:   "",
			wantErr: MissingValueError("test"),
		},
		{
			name:   "valid default",
			parser: testEq,
			param:  "value",
			op:     Equal,
			value:  "value",
		},
		{
			name:   "valid non-default op",
			parser: test2All,
			param:  "gt:value",
			op:     GreaterThan,
			value:  "value",
		},
		{
			name:    "invalid no-default",
			parser:  test3NoDefault,
			param:   "value",
			wantErr: InvalidOpError(""),
		},
		{
			name:   "valid no-default",
			parser: test3NoDefault,
			param:  "eq:value",
			op:     Equal,
			value:  "value",
		},
		{
			name:    "impossible parser",
			parser:  test4Impossible,
			wantErr: InvalidParamError("test4"),
		},
		{
			name:   "test like",
			parser: testLike,
			param:  "like:192.168",
			op:     Like,
			value:  "192.168",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, value, err := tt.parser(tt.param)
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no errors, got %v", err)
				return
			} else if tt.wantErr != nil && err == nil {
				t.Error("expected an error but didn't get one")
				return
			} else if tt.wantErr != err {
				t.Errorf("expected an error of %v but got %v", tt.wantErr, err)
				return
			}

			if op != tt.op {
				t.Errorf("expected op to equal %v, got %v", tt.op, op)
			}

			if value != tt.value {
				t.Errorf("expected value to equal %v, got %v", tt.value, value)
			}
		})
	}
}

func ExampleParamParser() {
	// A parser for a param called 'count' which defaults to the Equal operator if none specified, and also accepts the
	// GreaterThan operator.
	parser := ParamParser("count", Equal, GreaterThan)

	for _, example := range []string{
		"gt:123",
		"eq:456",
		"789",
		"lt:0",
	} {
		op, query, err := parser(example)
		fmt.Printf("op=%q, query=%q, err=%v\n", op, query, err)
	}

	// Output:
	// op="gt", query="123", err=<nil>
	// op="eq", query="456", err=<nil>
	// op="eq", query="789", err=<nil>
	// op="", query="", err=op not implemented: lt
}
