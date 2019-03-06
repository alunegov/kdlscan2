package kdlscan2

import (
	"errors"

	"github.com/alunegov/kdlscan2/file/drc"
	"github.com/alunegov/kdlscan2/file/lng"
)

// Generate генерирует lng-файл для kdl - с кодами ресурсов и без флагов Изменено и Удалено
func Generate(targetFileName string, lngFileName string, drcFileName string, drcFileEncoding string) error {
	if err := createBackup(targetFileName); err != nil {
		return nil
	}

	lngFile, err := lng.Load(lngFileName)
	if err != nil {
		return err
	}
	drcFile, err := drc.Load(drcFileName, drcFileEncoding)
	if err != nil {
		return err
	}

	if err := cleanFlags(lngFile); err != nil {
		return err
	}
	if err := restoreResourceID(lngFile, drcFile); err != nil {
		return err
	}

	if err := lng.Save(lngFile, targetFileName); err != nil {
		return err
	}

	return nil
}

// cleanFlags удаляет строки/ключи с флагом Удалено и заменяет флаг Изменено на Стд
func cleanFlags(lngFile *lng.File) error {
	for _, s := range lngFile.Sections() {
		s.FilterKey(func(k *lng.Key) bool {
			return k.Flag() != lng.Deleted
		})

		s.ForEachKey(func(k *lng.Key) {
			k.SetFlag(lng.Std)
		})
	}

	return nil
}

// restoreResourceID восстанавливает коды ресурсов на основе drc-файла
func restoreResourceID(lngFile *lng.File, drcFile *drc.File) error {
	if s := lngFile.Section(resourceStringSection); s != nil {
		for _, k := range s.Keys() {
			id, ok := drcFile.GetID(k.Name())
			if !ok {
				return errors.New("no res with name " + k.Name())
			}

			_ = k.RestoreResID(id)
		}

		// kdl ожидает отсортированные по коду ключи
		s.SortKeys()
	}

	return nil
}
