// Package util
package util

import "testing"

func TestStringJoin(t *testing.T) {
	var testCase = []struct{
		Input []string
		Separator, LastSeparator string
		Expected string
	}{
		{
			Input: []string{"id", "name", "address"},
			Separator: ",",
			LastSeparator: "",
			Expected: "id,name,address",
		},
		{
			Input: []string{"id", "name", "address"},
			Separator: "=? AND ",
			LastSeparator: "=?",
			Expected: "id=? AND name=? AND address=?",
		},
		{
			Input: []string{},
			Separator: "=? AND ",
			LastSeparator: "=?",
			Expected: "",
		},
		{
			Input: []string{"id"},
			Separator: "=? AND ",
			LastSeparator: "=?",
			Expected: "id=?",
		},
	}

	for _, x := range testCase{

		result := StringJoin(x.Input,x.Separator,x.LastSeparator)

		if x.Expected == result {
			t.Logf("expected '%s', got '%s", x.Expected, result)
		} else {
			t.Errorf("expected '%s', got '%s'", x.Expected, result)
		}
	}
}
