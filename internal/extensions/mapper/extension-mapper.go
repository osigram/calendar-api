package mapper

import (
	"calendar-api/internal/extensions"
	"errors"
)

type ExtensionMapper struct {
	extensions map[uint]extensions.Extension
}

func NewExtensionMapper() *ExtensionMapper {
	return &ExtensionMapper{make(map[uint]extensions.Extension)}
}

func (em *ExtensionMapper) RegisterExtension(id uint, extension extensions.Extension) {
	em.extensions[id] = extension
}

func (em *ExtensionMapper) Get(id uint) (extensions.Extension, error) {
	if extension, ok := em.extensions[id]; ok {
		return extension, nil
	}

	return nil, errors.New("extension not found")
}
