package lng

// File описывате lng-файл
type File struct {
	sectionsList []string
	sections     map[string]*Section
}

// NewFile создаёт lng-файл
func NewFile() *File {
	return &File{
		sectionsList: make([]string, 0, 10),
		sections:     make(map[string]*Section),
	}
}

// NewSection добавляет новую секцию
func (f *File) NewSection(name string) (*Section, error) {
	s := newSection(name)
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
