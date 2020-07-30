package drc

import "testing"

func TestParseDefine(t *testing.T) {
	key, id := parseDefine("#define RbDiagnosisReportService_SWaitDiag 62944")
	if key != "RbDiagnosisReportService_SWaitDiag" {
		t.Error("key")
	}
	if id != 62944 {
		t.Error("id")
	}
}

func TestParseString(t *testing.T) {
	testCases := []struct {
		str     string
		expKey  string
		expText string
	}{
		{
			`	RbDiagnosisReportService_SWaitDiag,	"Ждите, идет диагностика..."`,
			"RbDiagnosisReportService_SWaitDiag",
			"Ждите, идет диагностика...",
		},
		{
			`	RbDiagnosisReportService_SWaitDiag,	L"\x0416\x0434\x0438\x0442\x0435, \x0438\x0434\x0435\x0442 \x0434\x0438\x0430\x0433\x043d\x043e\x0441\x0442\x0438\x043a\x0430..."`,
			"RbDiagnosisReportService_SWaitDiag",
			"Ждите, идет диагностика...",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			key, text := parseString(tc.str)
			if key != tc.expKey {
				t.Error("key " + key)
			}
			if text != tc.expText {
				t.Error("text " + text)
			}
		})
	}
}
