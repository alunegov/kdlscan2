package kdlscan2

import (
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

	if err := lng.Save(refLngFile, targetFileName); err != nil {
		return err
	}

	return nil
}

// applyOldTranslation применяет перевод из старого файла
func applyOldTranslation(refLngFile *lng.File, oldLngFile *lng.File, markModified bool, markDeleted bool) error {
	for _, section := range refLngFile.Sections() {
		sectionOld := oldLngFile.Section(section.Name())

		// восстанавливаем переведённые строки
		for _, key := range section.Keys() {
			// выставляем флаг Изменено, далее он может измениться на Стд
			key.SetFlag(applyMarkConf(lng.Modified, markModified, markDeleted))

			if sectionOld != nil {
				keyOld, keyOldDeleted := sectionOld.Key(key.Name())
				if keyOld != nil {
					// в старом файле есть такая секция и строка/ключ
					if keyOld.Version() < key.Version() {
						// новая версия строки останется с флагом Изменено, старую добавим с флагом Удалено
						keyOldDeleted = keyOld
					} else {
						if keyOld.Flag() != lng.Deleted {
							// восстанавлием старый флаг и перевод
							key.SetFlag(applyMarkConf(keyOld.Flag(), markModified, markDeleted))
							key.SetValue(keyOld.Value())
						}
					}

					// добавляем старую строку с флагом Удалено. добавляемые ключи не попадут в текущий range
					if markDeleted && keyOldDeleted != nil {
						section.NewKeyAt(key.Name(), lng.Deleted, keyOldDeleted.Name(), keyOldDeleted.Version(),
							keyOldDeleted.Value())
					}
				}
			}
		}

		// добавляем строки из старого файла, которых нет в новом, с флагом Удалено
		if markDeleted && sectionOld != nil {
			for _, key := range sectionOld.Keys() {
				if kk, _ := section.Key(key.Name()); kk == nil {
					section.NewKey(lng.Deleted, key.Name(), key.Version(), key.Value())
				}
			}
		}
	}

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
