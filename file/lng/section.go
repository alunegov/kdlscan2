package lng

import (
	"errors"
	"sort"
)

// Section описывает секцию lng-файла
type Section struct {
	name string
	keys []*Key
}

// newSection создаёт секцию
func newSection(name string) *Section {
	return &Section{
		name: name,
		keys: make([]*Key, 0),
	}
}

// Name возвращает имя
func (s *Section) Name() string { return s.name }

// Keys возвращает список ключей
func (s *Section) Keys() []*Key {
	res := make([]*Key, len(s.keys))
	copy(res, s.keys)
	return res
}

// NewKey добавляет новый ключ
func (s *Section) NewKey(flag KeyFlag, name string, version KeyVersion, value string) (*Key, error) {
	k := newKey(flag, name, version, value)
	s.keys = append(s.keys, k)
	return k, nil
}

// NewKeyAt добавляет новый ключ после ключа с указанным именем
func (s *Section) NewKeyAt(atName string, flag KeyFlag, name string, version KeyVersion, value string) (*Key, error) {
	for i, k := range s.keys {
		if k.Name() == atName {
			kk := newKey(flag, name, version, value)
			// insert kk at i+1
			s.keys = append(s.keys, nil)
			copy(s.keys[i+2:], s.keys[i+1:])
			s.keys[i+1] = kk
			return kk, nil
		}
	}
	return nil, errors.New("no item with name " + atName)
}

// Key возвращает ключ по имени, и одноименный ключ с другим флагом
// TODO: ускорить получение ключа по имени (map[string]*Key)
func (s *Section) Key(name string) (*Key, *Key) {
	for i, k := range s.keys {
		if k.Name() == name {
			var deleted *Key
			if (i+1) < len(s.keys) && s.keys[i+1].Name() == name {
				deleted = s.keys[i+1]
			}
			return k, deleted
		}
	}
	return nil, nil
}

// KeyByValue возвращает ключ по значению
func (s *Section) KeyByValue(value string) *Key {
	for _, k := range s.keys {
		if k.Value() == value {
			return k
		}
	}
	return nil
}

// FilterKey фильтрует ключи (удаляет ключи, не соответствующие условию)
func (s *Section) FilterKey(op func(k *Key) bool) {
	filtered := make([]*Key, 0, len(s.keys))
	for _, k := range s.keys {
		if op(k) {
			filtered = append(filtered, k)
		}
	}

	s.keys = filtered
}

// ForEachKey обходит ключи, выполняя указанную операцию
// Во время обхода добавлять новые ключи нельзя
func (s *Section) ForEachKey(op func(k *Key)) {
	for _, k := range s.keys {
		op(k)
	}
}

// SortKeys сортирует ключи по имени
func (s *Section) SortKeys() {
	sort.Slice(s.keys, func(i, j int) bool {
		return s.keys[i].Name() < s.keys[j].Name()
	})
}
