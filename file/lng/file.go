package lng

import "time"

// File описывате lng-файл
type File struct {
	UpdateTime   time.Time
	Changed      bool
	sectionsList []string
	sections     map[string]*Section
}

// NewFile создаёт lng-файл
func NewFile() *File {
	return &File{
		UpdateTime:   time.Now(),
		Changed:      false,
		sectionsList: make([]string, 0, 10),
		sections:     make(map[string]*Section),
	}
}

// NewSection добавляет новую секцию
func (f *File) NewSection(name string) (*Section, error) {
	s := newSection(f, name)
	f.changed()
	f.sectionsList = append(f.sectionsList, name)
	f.sections[name] = s
	return s, nil
}

// Sections возвращает список секций
func (f *File) Sections() []*Section {
	res := make([]*Section, 0, len(f.sectionsList))
	for _, n := range f.sectionsList {
		res = append(res, f.sections[n])
	}
	return res
}

// Section возвращает секцию по имени
func (f *File) Section(name string) *Section {
	return f.sections[name]
}

func (f *File) changed() {
	f.Changed = true
}
