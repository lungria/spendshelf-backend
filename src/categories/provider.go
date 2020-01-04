package categories

import (
	"context"
	"errors"
	"sync"
)

type Provider interface {
	GetAll() []Category
	Find(name string) (Category, bool)
}

type InMemoryProvider struct {
	categories  map[string]Category
	updatesChan <-chan Category
	context     context.Context
	mutex       *sync.RWMutex
}

func NewProvider(context context.Context, seed []Category, updates <-chan Category) (*InMemoryProvider, error) {
	if seed == nil {
		return nil, errors.New("seed map must not be nil")
	}
	if updates == nil {
		return nil, errors.New("updates channel must not be nil")
	}
	categories := make(map[string]Category)
	for _, v := range seed {
		categories[v.NormalizedName] = v
	}

	provider := &InMemoryProvider{
		categories:  categories,
		updatesChan: updates,
		context:     context,
		mutex:       &sync.RWMutex{},
	}
	go provider.runSync()
	return provider, nil
}

func (provider *InMemoryProvider) GetAll() []Category {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	arr := []Category{}
	for _, v := range provider.categories {
		arr = append(arr, v)
	}
	return arr
}

func (provider *InMemoryProvider) Find(normalizedName string) (Category, bool) {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	val, exists := provider.categories[normalizedName]
	return val, exists
}

func (provider *InMemoryProvider) runSync() {
	for {
		select {
		case <-provider.context.Done():
			{
				return
			}
		case val, ok := <-provider.updatesChan:
			{
				if !ok {
					return
				}
				provider.mutex.Lock()
				key := val.NormalizedName
				provider.categories[key] = val
				provider.mutex.Unlock()
			}
		}
	}
}
