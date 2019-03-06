package kdlscan2

import (
	"errors"
	"log"

	"github.com/alunegov/kdlscan2/file/lng"
)

// Sync синхронизирует перевод ресурсов между двумя lng-файлами
func Sync(targetFileName string, refFileName string) error {
	if err := createBackup(targetFileName); err != nil {
		return err
	}

	target, err := lng.Load(targetFileName + backupFileExt)
	if err != nil {
		return err
	}
	ref, err := lng.Load(refFileName)
	if err != nil {
		return err
	}

	if err := syncResourceStrings(target, ref); err != nil {
		return err
	}

	if err := lng.Save(target, targetFileName); err != nil {
		return err
	}

	return nil
}

func syncResourceStrings(targetFile *lng.File, refFile *lng.File) error {
	targetSection := targetFile.Section(resourceStringSection)
	if targetSection == nil {
		return errors.New("ResourceStrings section in target is absent")
	}
	refSection := refFile.Section(resourceStringSection)
	if refSection == nil {
		return errors.New("ResourceStrings section in ref is absent")
	}

	refSection.FilterKey(func(k *lng.Key) bool {
		return k.Flag() == lng.Std
	})
	// если это файл после kdl или чистый lng-файл, удаляем коды
	refSection.ForEachKey(func(k *lng.Key) {
		k.StripResID()
	})

	targetSection.ForEachKey(func(k *lng.Key) {
		if k.Flag() == lng.Deleted {
			return
		}

		_, name, _ := k.DecodeName() // если это файл после kdl или чистый lng-файл, берём только имя, без кода
		refKey, _ := refSection.Key(name)
		if refKey != nil && (k.Value() != refKey.Value() || k.Flag() == lng.Modified) {
			log.Printf("%s:\n", k.Name())
			log.Printf("    '%s'   ->   '%s'\n", k.Value(), refKey.Value())

			k.SetFlag(lng.Std)
			k.SetValue(refKey.Value())
		}
	})

	return nil
}
