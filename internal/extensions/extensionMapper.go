package extensions

import "errors"

type ExtensionMapper struct {
	extensions map[uint]Extension
}

func NewExtensionMapper() *ExtensionMapper {
	return &ExtensionMapper{make(map[uint]Extension)}
}

func (em *ExtensionMapper) RegisterExtension(id uint, extension Extension) {
	em.extensions[id] = extension
}

func (em *ExtensionMapper) Get(id uint) (Extension, error) {
	if extension, ok := em.extensions[id]; ok {
		return extension, nil
	}

	return nil, errors.New("extension not found")
}
