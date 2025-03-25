package mapper

import (
	"calendar-api/pkg/extensions"
	"context"
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

func (em *ExtensionMapper) Get(ctx context.Context, id uint) (extensions.Extension, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if extension, ok := em.extensions[id]; ok {
			return extension, nil
		}

		return nil, errors.New("extension not found")
	}

}
