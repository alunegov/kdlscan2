package kdlscan2

import (
	"strconv"
	"testing"

	"github.com/alunegov/kdlscan2/file/lng"
)

func TestApplyOldTranslation(t *testing.T) {
	refLngFile := lng.NewFile()
	oldLngFile := lng.NewFile()

	refSection, _ := refLngFile.NewSection("[TIDSMainForm]")
	oldSection, _ := oldLngFile.NewSection("[TIDSMainForm]")

	_, _ = refSection.NewKey(lng.Modified, "1", 1, "ref")
	_, _ = oldSection.NewKey(lng.Modified, "1", 1, "ref")
	_, _ = oldSection.NewKey(lng.Deleted, "1", 0, "old")

	for i := 2; i < 6; i++ {
		s := strconv.Itoa(i)
		_, _ = refSection.NewKey(lng.Std, s, 0, "ref"+s)
		_, _ = oldSection.NewKey(lng.Std, s, 0, "ref"+s)
	}

	_, _ = oldSection.NewKey(lng.Std, "8", 0, "old")

	if err := applyOldTranslation(refLngFile, oldLngFile, true, true); err != nil {
		t.Error(err)
	}

	if len(refSection.Keys()) != 7 {
		t.Errorf("exp len %d, actual len %d", 7, len(refSection.Keys()))
	}
}

func TestApplyOldTranslation_WillPreserveUpdateTimeWhenNoChanges(t *testing.T) {
	refLngFile := lng.NewFile()
	oldLngFile := lng.NewFile()

	refSection, _ := refLngFile.NewSection("[TIDSMainForm]")
	oldSection, _ := oldLngFile.NewSection("[TIDSMainForm]")

	_, _ = refSection.NewKey(lng.Modified, "1", 0, "ref")
	_, _ = oldSection.NewKey(lng.Modified, "1", 0, "old")

	refLngFile.Changed = false

	if err := applyOldTranslation(refLngFile, oldLngFile, true, true); err != nil {
		t.Error(err)
	}

	if refLngFile.Changed {
		t.Error("err")
	}
}

func TestApplyOldTranslation_ShouldPreserveVersion(t *testing.T) {
	refLngFile := lng.NewFile()
	oldLngFile := lng.NewFile()

	refSection, _ := refLngFile.NewSection("[TIDSMainForm]")
	oldSection, _ := oldLngFile.NewSection("[TIDSMainForm]")

	_, _ = refSection.NewKey(lng.Modified, "1", 0, "ref")
	_, _ = oldSection.NewKey(lng.Modified, "1", 1, "old")
	_, _ = oldSection.NewKey(lng.Deleted, "1", 0, "old_pre")

	if err := applyOldTranslation(refLngFile, oldLngFile, true, true); err != nil {
		t.Error(err)
	}

	if refSection.Keys()[0].Version() != 1 {
		t.Error("err")
	}
}

func TestApplyMarkConf(t *testing.T) {
	var tests = []struct {
		flag         lng.KeyFlag
		markModified bool
		markDeleted  bool
		expFlag      lng.KeyFlag
	}{
		{lng.Std, true, true, lng.Std},
		{lng.Modified, true, true, lng.Modified},
		{lng.Modified, false, true, lng.Std},
		{lng.Deleted, true, true, lng.Deleted},
		{lng.Deleted, true, false, lng.Std},
	}

	for i, testCase := range tests {
		flag := applyMarkConf(testCase.flag, testCase.markModified, testCase.markDeleted)
		if flag != testCase.expFlag {
			t.Errorf("Test case %d failed", i)
		}
	}
}
