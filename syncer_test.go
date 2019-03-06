package kdlscan2

import (
	"testing"

	"github.com/alunegov/kdlscan2/file/lng"
)

func TestSyncResourceStrings(t *testing.T) {
	var tests = []struct {
		name        string
		targetFlag  lng.KeyFlag
		targetValue string
		refFlag     lng.KeyFlag
		refValue    string
		expFlag     lng.KeyFlag
		expValue    string
	}{
		{"0", lng.Std, "t", lng.Std, "r", lng.Std, "r"},
		{"1", lng.Std, "t", lng.Modified, "r", lng.Std, "t"},
		{"2", lng.Std, "t", lng.Deleted, "r", lng.Std, "t"},

		{"3", lng.Modified, "t", lng.Std, "r", lng.Std, "r"},
		{"4", lng.Modified, "t", lng.Modified, "r", lng.Modified, "t"},
		{"5", lng.Modified, "t", lng.Deleted, "r", lng.Modified, "t"},

		{"6", lng.Deleted, "t", lng.Std, "r", lng.Deleted, "t"},
		{"7", lng.Deleted, "t", lng.Modified, "r", lng.Deleted, "t"},
		{"8", lng.Deleted, "t", lng.Deleted, "r", lng.Deleted, "t"},

		{"9", lng.Std, "t", lng.Std, "t", lng.Std, "t"},
		{"10", lng.Std, "t", lng.Modified, "t", lng.Std, "t"},
		{"11", lng.Std, "t", lng.Deleted, "t", lng.Std, "t"},

		{"12", lng.Modified, "t", lng.Std, "t", lng.Std, "t"},
		{"13", lng.Modified, "t", lng.Modified, "t", lng.Modified, "t"},
		{"14", lng.Modified, "t", lng.Deleted, "t", lng.Modified, "t"},

		{"15", lng.Deleted, "t", lng.Std, "t", lng.Deleted, "t"},
		{"16", lng.Deleted, "t", lng.Modified, "t", lng.Deleted, "t"},
		{"17", lng.Deleted, "t", lng.Deleted, "t", lng.Deleted, "t"},
	}

	targetFile := lng.NewFile()
	refFile := lng.NewFile()

	targetSection, _ := targetFile.NewSection(resourceStringSection)
	refSection, _ := refFile.NewSection(resourceStringSection)

	for _, testCase := range tests {
		_, _ = targetSection.NewKey(testCase.targetFlag, testCase.name, 0, testCase.targetValue)
		_, _ = refSection.NewKey(testCase.refFlag, testCase.name, 0, testCase.refValue)
	}

	if err := syncResourceStrings(targetFile, refFile); err != nil {
		t.Error(err)
	}

	for i, testCase := range tests {
		key, _ := targetSection.Key(testCase.name)
		if key.Flag() != testCase.expFlag || key.Value() != testCase.expValue {
			t.Errorf("Test case %d failed", i)
		}
	}
}
