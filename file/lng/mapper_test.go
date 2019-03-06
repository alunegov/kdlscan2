package lng

import (
	"testing"
	"time"
)

func TestParseKey(t *testing.T) {
	var tests = []struct {
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

func TestParseUpdateTime(t *testing.T) {
	var tests = []struct {
		EncTime string
		expTime time.Time
		expErr  bool
	}{
		{
			"2018-08-07T15:14:49+05:00",
			time.Date(2018, 8, 7, 15, 14, 49, 0, time.FixedZone("UTC+5", 5*60*60)),
			false,
		},
		{
			"26.07.2018 14:40:59",
			time.Date(2018, 7, 26, 14, 40, 59, 0, time.Local),
			false,
		},
		{
			"wrong time",
			time.Time{},
			true,
		},
	}

	for i, test := range tests {
		resTime, err := parseUpdateTime(test.EncTime)
		if test.expErr && err == nil {
			t.Errorf("Test case %d failed. Expected %v but got %v", i, test.expErr, err == nil)
		}
		if !test.expTime.Equal(resTime) {
			t.Errorf("Test case %d failed. Expected %s but got %s", i, test.expTime, resTime)
		}
	}
}
