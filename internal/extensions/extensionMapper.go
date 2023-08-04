package extensions

import "errors"

type ExtensionMapper struct {
	extensions map[int64]Extension
}

func NewExtensionMapper() *ExtensionMapper {
	return &ExtensionMapper{make(map[int64]Extension)}
}

func (em *ExtensionMapper) RegisterExtension(id int64, extension Extension) {
	em.extensions[id] = extension
}

func (em *ExtensionMapper) Get(id int64) (Extension, error) {
	if extension, ok := em.extensions[id]; ok {
		return extension, nil
	}

	return nil, errors.New("extension not found")
}
