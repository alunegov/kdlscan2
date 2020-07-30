package drc

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

type resource struct {
	id  int
	str string
}

// File описывает DRC(Delphi Resource String)-файл
type File struct {
	// Строковые ресурсы в виде имя->{id, текст}
	stringsMap map[string]*resource
}

// Load загружает (и разбирает) drc-файл с диска
// encoding - кодировка drc-файла ("", либо имя кодировки из http://www.w3.org/TR/encoding)
func Load(name string, encoding string) (*File, error) {
	res := &File{
		stringsMap: make(map[string]*resource),
	}
	if err := res.parse(name, encoding); err != nil {
		return nil, err
	}
	return res, nil
}

// GetKey возвращает текст ресурса по id
func (f *File) GetKey(id int) (string, bool) {
	for k, v := range f.stringsMap {
		if v.id == id {
			return k, true
		}
	}
	return "", false
}

// GetID возвращает id ресурса по имени
func (f *File) GetID(key string) (int, bool) {
	if r, ok := f.stringsMap[key]; ok {
		return r.id, true
	}
	return 0, false
}

// GetStr возвращает текст ресурса по имени
func (f *File) GetStr(key string) (string, bool) {
	if r, ok := f.stringsMap[key]; ok {
		return r.str, true
	}
	return "", false
}

const (
	parseComment         = iota // обработка примечания
	parseDefines                // обработка defines
	parseTextStringTable        // обработка текста STRINGTABLE\nBEGIN
	parseStrings                // обработка строк
	parseTextEnd                // обработка текста END
)

// parse разбирает drc-файл
func (f *File) parse(fileName string, encoding string) error {
	ff, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(_f *os.File) {
		_ = _f.Close()
	}(ff)

	var t io.Reader
	if encoding == "" {
		t = ff
	} else {
		enc, err := htmlindex.Get(encoding)
		if err != nil {
			return err
		}

		t = transform.NewReader(ff, enc.NewDecoder())
	}

	state := parseComment

	s := bufio.NewScanner(t)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		switch state {
		case parseComment:
			if strings.HasPrefix(line, "*/") {
				state = parseDefines
			}
		case parseDefines:
			if line == "STRINGTABLE" {
				state = parseTextStringTable
			} else {
				key, id := parseDefine(line)
				f.stringsMap[key] = &resource{id: id}
			}
		case parseTextStringTable:
			if line == "BEGIN" {
				state = parseStrings
			}
		case parseStrings:
			if line == "END" {
				state = parseTextEnd
			} else {
				key, str := parseString(line)
				f.stringsMap[key].str = str
			}
		case parseTextEnd:
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

// parseDefine разбирает строку с id ресурса
// Формат строки - `#define RbDiagnosisReportService_SWaitDiag 62944`
func parseDefine(s string) (key string, id int) {
	f := strings.Fields(s)
	// TODO: bounds check
	key = f[1]
	id, _ = strconv.Atoi(f[2])
	return
}

// parseString разбирает строку с тестом ресурса
// Формат строки - `	RbDiagnosisReportService_SWaitDiag,	"Ждите, идет диагностика..."` в D7 или
// `	RbDiagnosisReportService_SWaitDiag,	L"\x0416\x0434\x0438\x0442\x0435, \x0438\x0434\x0435\x0442 \x0434\x0438\x0430\x0433\x043d\x043e\x0441\x0442\x0438\x043a\x0430..."` в DXE
func parseString(s string) (key string, str string) {
	i := strings.Index(s, ",")
	// TODO: bounds check
	key = strings.TrimSpace(s[:i])
	if s[i+2] == 'L' {
		str = s[i+3:]
		// in Delphi '\x' means HEX (encoding are UTF-16); in Go '\x' means UTF-8, '\u' means UTF-16
		str = strings.Replace(str, `\x`, `\u`, -1)
		str, _ = strconv.Unquote(str)
	} else {
		str = s[i+3 : len(s)-1]
	}
	return
}
