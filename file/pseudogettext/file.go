package pseudogettext

import (
	"bufio"
	"os"
	"strings"
)

const sectionMarker = ";"

// File описывает псевдо-gettext-файл
type File struct {
	lines []string
}

// Load загружает (и рабирает) псевдо-gettext-файл с диска
func Load(fileName string) (*File, error) {
	res := &File{}
	var err error

	if res.lines, err = loadLines(fileName); err != nil {
		return nil, err
	}

	return res, nil
}

func loadLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(_f *os.File) {
		_ = _f.Close()
	}(f)

	res := make([]string, 0)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		// берём все кроме пустых строк и названий секций
		if len(line) > 0 && !strings.HasPrefix(line, sectionMarker) {
			res = append(res, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Lines возвращает список строк
func (f *File) Lines() []string {
	res := make([]string, len(f.lines))
	if i := copy(res, f.lines); i != len(f.lines) {
		return []string{}
	}
	return res
}
