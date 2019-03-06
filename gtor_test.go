package kdlscan2

import (
	"testing"

	"github.com/alunegov/kdlscan2/file/lng"
)

func TestCleanFlags(t *testing.T) {
	var tests = []struct {
		name string
		flag lng.KeyFlag
	}{
		{"0", lng.Std},
		{"1", lng.Modified},
		{"2", lng.Deleted},
	}

	file := lng.NewFile()

	section, _ := file.NewSection("test")
	for _, testCase := range tests {
		section.NewKey(testCase.flag, testCase.name, 0, "")
	}

	if err := cleanFlags(file); err != nil {
		t.Error(err)
	}

	if len(section.Keys()) != 2 {
		t.Error("len")
	}
	for _, key := range section.Keys() {
		if key.Flag() != lng.Std {
			t.Errorf("%s is not Std", key.Name())
		}
	}
}
