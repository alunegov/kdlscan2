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
	key, text := parseString("	RbDiagnosisReportService_SWaitDiag,	\"Ждите, идет диагностика...\"")
	if key != "RbDiagnosisReportService_SWaitDiag" {
		t.Error("key")
	}
	if text != "Ждите, идет диагностика..." {
		t.Error("text")
	}
}
