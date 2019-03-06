package kdlscan2

import (
	"os"
	"sort"

	"github.com/alunegov/kdlscan2/file/lng"
	"github.com/alunegov/kdlscan2/file/pseudogettext"
)

// Scan генерирует референсный lng-файл - сканирование исполняемого файла, добавление строк из псевдо-gettext файлов
// Пока исполняемый файл не сканируется, вместо этого используется lng-файл после kdlscan.
func Scan(targetFileName string, lngFileName string, extraFileNames []string) error {
	if err := createBackup(targetFileName); err != nil {
		return err
	}

	lngFile, err := lng.Load(lngFileName)
	if err != nil {
		return err
	}

	if err := stripResourceID(lngFile, true); err != nil {
		return err
	}
	if err := addPseudoGettextFiles(lngFile, extraFileNames); err != nil {
		return err
	}

	if err := lng.Save(lngFile, targetFileName); err != nil {
		return err
	}

	return nil
}

func createBackup(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil
	}

	oldFileName := fileName + backupFileExt
	if err := os.Rename(fileName, oldFileName); err != nil {
		return err
	}

	return nil
}

// stripResourceID удаляет коды ресурсов
func stripResourceID(lngFile *lng.File, needSort bool) error {
	if s := lngFile.Section(resourceStringSection); s != nil {
		s.ForEachKey(func(k *lng.Key) {
			_ = k.StripResID()
		})

		// ранее ключи были отсортированы по коду (не по имени модуля/ресурса). сортируем по имени, чтобы было меньше
		// изменений в VCS.
		if needSort {
			s.SortKeys()
		}
	}

	return nil
}

// addPseudoGettextFiles добавляет в lng-файл строки из псевдо-gettext файлов
// Дубли строк отбрасываются, строки сортируются и добавляются в виде строка=строка
func addPseudoGettextFiles(file *lng.File, fileNames []string) error {
	lines := make([]string, 0)
	linesMap := make(map[string]struct{}, 0)

	for _, fileName := range fileNames {
		pgFile, err := pseudogettext.Load(fileName)
		if err != nil {
			return err
		}

		for _, l := range pgFile.Lines() {
			if _, ok := linesMap[l]; ok {
				// у нас уже есть такая строка
				continue
			}
			lines = append(lines, l)
			linesMap[l] = struct{}{}
		}
	}

	if len(lines) == 0 {
		return nil
	}

	// kdl-PseudoGettext ожидает отсортированные по имени ключи. на самом деле нет, там используется поиск перебором,
	// потому что UTF8 преобразуется в используемую кодировку CP и имена ключей коверкаются.
	sort.Strings(lines)

	s, err := file.NewSection(pseudoGettextSection)
	if err != nil {
		return err
	}
	for _, l := range lines {
		if _, err := s.NewKey(lng.Std, l, 0, l); err != nil {
			return err
		}
	}

	return nil
}
