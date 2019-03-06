package lng

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// KeyFlag описывает флаг ключа
type KeyFlag int

const (
	// Std обозначает Стд
	Std KeyFlag = iota
	// Modified обозначает Изменено
	Modified
	// Deleted обозначает Удалено
	Deleted
)

func (f KeyFlag) String() string {
	switch f {
	case Modified:
		return "(!)"
	case Deleted:
		return "(x)"
	default:
		return ""
	}
}

// KeyVersion описывает версию ключа
type KeyVersion int

func (v KeyVersion) String() string {
	if v > 0 {
		return fmt.Sprintf("{%d}", v)
	}
	return ""
}

// Key описывает ключ секции lng-файла
type Key struct {
	file    *File
	flag    KeyFlag
	name    string
	version KeyVersion
	value   string
}

// newKey создаёт ключ
func newKey(file *File, flag KeyFlag, name string, version KeyVersion, value string) *Key {
	return &Key{
		file:    file,
		flag:    flag,
		name:    name,
		version: version,
		value:   value,
	}
}

// Flag возвращает флаг
func (k *Key) Flag() KeyFlag { return k.flag }

// SetFlag задаёт флаг
func (k *Key) SetFlag(flag KeyFlag) {
	if k.flag != flag {
		k.changed()
		k.flag = flag
	}
}

// Name возвращает имя
func (k *Key) Name() string { return k.name }

// SetName задаёт имя
func (k *Key) SetName(name string) {
	if k.name != name {
		k.changed()
		k.name = name
	}
}

// StripResID удаляет код ресурса из имени
func (k *Key) StripResID() error {
	_, resName, err := k.DecodeName()
	if err != nil {
		return err
	}
	k.changed()
	k.name = resName
	return nil
}

// RestoreResID добавляет код ресурса к имени
func (k *Key) RestoreResID(resID int) error {
	if _, _, err := k.DecodeName(); err == nil {
		return errors.New("already with resID")
	}
	k.changed()
	k.name = strconv.Itoa(resID) + "_" + k.name
	return nil
}

// DecodeName декодирует имя, как имя ресурса
func (k *Key) DecodeName() (resID int, resName string, err error) {
	resName = k.name // на случай, если это не имя ресурса

	delim := strings.Index(k.name, "_")
	if delim == -1 {
		err = errors.New("no delim")
		return
	}
	var err2 error
	if resID, err2 = strconv.Atoi(k.name[:delim]); err2 != nil {
		err = errors.New("no resID")
		return
	}
	resName = k.name[delim+1:]

	return
}

// Version возвращает версию
func (k *Key) Version() KeyVersion { return k.version }

// SetVersion задаёт версию
func (k *Key) SetVersion(version KeyVersion) {
	if k.version != version {
		k.changed()
		k.version = version
	}
}

// Value возвращает значение
func (k *Key) Value() string { return k.value }

// SetValue задаёт значение
func (k *Key) SetValue(value string) {
	if k.value != value {
		k.changed()
		k.value = value
	}
}

func (k *Key) String() string {
	return fmt.Sprintf("%s%s%s=%s", k.flag, k.name, k.version, k.value)
}

// changed выставляет флаг Изменено у lng-файла
func (k *Key) changed() {
	k.file.Changed = true
}
