package lng

import "testing"

func TestParseKey(t *testing.T) {
	var tests = []struct{
		s       string
		flag    KeyFlag
		name    string
		version KeyVersion
		value   string
	}{
		{"n=v", Std, "n", 0, "v"},
		{"(!)n=v", Modified, "n", 0, "v"},
		{"(x)n=v", Deleted, "n", 0, "v"},
		{"n{1}=v", Std, "n", 1, "v"},
		{"(!)n{1}=v", Modified, "n", 1, "v"},
	}

	for i, testCase := range tests {
		flag, name, version, value := parseKey(testCase.s)
		if flag != testCase.flag || name != testCase.name || version != testCase.version || value != testCase.value {
			t.Errorf("Test case %d failed", i)
		}
	}
}
