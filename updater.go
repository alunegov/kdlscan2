package kdlscan2

import (
	"time"

	"github.com/alunegov/kdlscan2/file/lng"
)

// Update обновляет lng-файл на основе референсного lng-файла
func Update(targetFileName string, refFileName string, markModified bool, markDeleted bool) error {
	if err := createBackup(targetFileName); err != nil {
		return err
	}

	refLngFile, err := lng.Load(refFileName)
	if err != nil {
		return err
	}
	oldLngFile, err := lng.Load(targetFileName + backupFileExt)
	if err != nil {
		return err
	}

	// если открыли lng-файл после kdl, то нужно удалить коды ресурсов, строки не сортируем (нам они нужны временно)
	if err := stripResourceID(oldLngFile, false); err != nil {
		return err
	}
	if err := applyOldTranslation(refLngFile, oldLngFile, markModified, markDeleted); err != nil {
		return err
	}
	if refLngFile.Changed {
		refLngFile.UpdateTime = time.Now()
	} else {
		refLngFile.UpdateTime = oldLngFile.UpdateTime
	}

	if err := lng.Save(refLngFile, targetFileName); err != nil {
		return err
	}

	return nil
}

// applyOldTranslation применяет перевод из старого файла
func applyOldTranslation(refLngFile *lng.File, oldLngFile *lng.File, markModified bool, markDeleted bool) error {
	// refLngFile по алгоритму "загрязнится" (k.SetFlag, k.SetValue, s.NewKeyAt и s.NewKey), поэтому следим за
	// изменением сами (fileChanged, sectionChanged и keyChanged). За начало берём текущее состояние файла.
	fileChanged := refLngFile.Changed

	for _, section := range refLngFile.Sections() {
		sectionOld := oldLngFile.Section(section.Name())
		sectionChanged := false

		// восстанавливаем переведённые строки
		for _, key := range section.Keys() {
			// выставляем флаг Изменено, далее он может измениться на Стд
			key.SetFlag(applyMarkConf(lng.Modified, markModified, markDeleted))
			keyChanged := true

			if sectionOld != nil {
				keyOld, keyOldDeleted := sectionOld.Key(key.Name())
				if keyOld != nil {
					// в старом файле есть такая секция и строка/ключ
					if keyOld.Version() < key.Version() {
						// новая версия строки останется с флагом Изменено, старую добавим с флагом Удалено
						keyOldDeleted = keyOld
					} else if keyOld.Flag() != lng.Deleted {
						// восстанавлием старый флаг, версию и перевод
						key.SetFlag(applyMarkConf(keyOld.Flag(), markModified, markDeleted))
						key.SetVersion(keyOld.Version())
						key.SetValue(keyOld.Value())
						keyChanged = false
					}
					// else старый ключ той-же версии, но с флагом Удалено
					// TODO: сброс keyChanged?

					// Добавляем старую строку с флагом Удалено - это может быть как Удалённая строка, так и строка
					// предыдущей версии. Добавляемые ключи не попадут в текущий range.
					// sectionChanged не выставляем - строки с флагом Удалено
					if markDeleted && keyOldDeleted != nil {
						_, _ = section.NewKeyAt(key.Name(), lng.Deleted, keyOldDeleted.Name(), keyOldDeleted.Version(),
							keyOldDeleted.Value())
					}
				}
				// else в старой версии нет такого ключа - ключ остаётся с флагом Изменено
			}
			// else в старой версии нет такой секции - ключ остаётся с флагом Изменено

			if keyChanged {
				sectionChanged = true
			}
		}

		// добавляем строки из старого файла, которых нет в новом, с флагом Удалено
		// sectionChanged не выставляем - строки с флагом Удалено
		if markDeleted && sectionOld != nil {
			for _, key := range sectionOld.Keys() {
				if kk, _ := section.Key(key.Name()); kk == nil {
					_, _ = section.NewKey(lng.Deleted, key.Name(), key.Version(), key.Value())
				}
			}
		}

		if sectionChanged {
			fileChanged = true
		}
	}

	refLngFile.Changed = fileChanged

	return nil
}

// applyMarkConf учитывает конфигурацию при задании флага
func applyMarkConf(f lng.KeyFlag, markModified bool, markDeleted bool) lng.KeyFlag {
	if f == lng.Modified && markModified {
		return f
	} else if f == lng.Deleted && markDeleted {
		return f
	} else {
		return lng.Std
	}
}
