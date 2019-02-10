package postal

import (
	"testing"
)

func TestValidatePostalCode(t *testing.T) {
	cases := []struct {
		Value  string
		Expect bool
	}{
		{"08773380", true},
		{"08773-380", true},
		{"8773380", false},
		{"8773-380", false},
		{"CEP08773380", true},
		{"08773.380", true},
		{"08773 380", true},
		{"O8773-380", false},
		{"O8773380", false},
	}

	for i, tc := range cases {
		t.Logf("Testing case #%v", i)

		if _, valid := ValidatePostalCode(tc.Value); valid != tc.Expect {
			t.Errorf("Unexpected validation: %t of '%s', expect %t", valid, tc.Value, tc.Expect)
		} else {
			t.Logf("Value:'%s' is valid:%t", tc.Value, tc.Expect)
		}
	}
}
