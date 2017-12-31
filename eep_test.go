package eep

import (
	"testing"
)

func TestEval(t *testing.T) {
	testCases := []struct {
		expression string
		val        interface{}
	}{
		{
			"(1)",
			float64(1),
		},
		{
			`("abc"+"def")`,
			"abcdef",
		},
		{
			"2-1",
			float64(1),
		},
		{
			"2*(10-1)",
			float64(18),
		},
		{
			"true and true",
			true,
		},
		{
			"false or true",
			true,
		},
		{
			"true and false",
			false,
		},
		{
			"1==1",
			true,
		},
		{
			"100!=100",
			false,
		},
		{
			`"abc"=="abc"`,
			true,
		},
		{
			`(!(1 != 1)) == ("x" == "x")`,
			true,
		},
	}

	for _, tc := range testCases {
		val, err := Eval(tc.expression)
		if err != nil {
			t.Errorf("eval error: %v\n", err)
		}
		if val != tc.val {
			t.Errorf("expected: %v,get: %v\n", tc.val, val)
		}
	}
}

func TestEvalWithEnv(t *testing.T) {
	f := func(args ...interface{}) interface{} {
		s1, s2 := args[0].(string), args[1].(string)
		return s1 + s2
	}

	add := func(args ...interface{}) interface{} {
		n1, n2 := args[0].(int), args[1].(int)
		return n1 + n2
	}

	exponent := func(args ...interface{}) interface{} {
		n1, n2 := args[0].(float64), args[1].(float64)
		var n float64 = 1
		var i int64 = 1
		for i = 1; i <= int64(n2); i++ {
			n *= n1
		}
		return n
	}

	len := func(args ...interface{}) interface{} {
		return len(args[0].([]int))
	}

	testCases := []struct {
		expression string
		val        interface{}
		env        map[string]interface{}
	}{
		{
			"Concat(x, y)",
			"123456",
			map[string]interface{}{"x": "123", "y": "456", "Concat": f},
		},
		{
			`Concat("x", "y")`,
			"xy",
			map[string]interface{}{"Concat": f},
		},
		{
			"add(a, b)",
			100,
			map[string]interface{}{"add": add, "a": 99, "b": 1},
		},
		{
			"expn(2, 10)==1024",
			true,
			map[string]interface{}{"expn": exponent},
		},
		{
			"len(x)",
			6,
			map[string]interface{}{"len": len, "x": []int{1, 3, 2, 8, 6, 100}},
		},
	}

	for _, tc := range testCases {
		val, err := EvalWithEnv(tc.expression, tc.env)
		if err != nil {
			t.Errorf("eval error: %v\n", err)
		}
		if val != tc.val {
			t.Errorf("expected: %v,get: %v\n", tc.val, val)
		}
	}
}
