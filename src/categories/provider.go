package categories

import (
	"context"
	"errors"
	"sync"
)

type provider interface {
	GetAll() []Category
	Find(name string) (Category, bool)
}

type inMemoryProvider struct {
	categories  map[string]Category
	updatesChan <-chan Category
	context     context.Context
	mutex       *sync.RWMutex
}

func newProvider(context context.Context, seed []Category, updates <-chan Category) (*inMemoryProvider, error) {
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

	provider := &inMemoryProvider{
		categories:  categories,
		updatesChan: updates,
		context:     context,
		mutex:       &sync.RWMutex{},
	}
	go provider.runSync()
	return provider, nil
}

// GetAll return all categories from in memory cache
func (provider *inMemoryProvider) GetAll() []Category {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	arr := []Category{}
	for _, v := range provider.categories {
		arr = append(arr, v)
	}
	return arr
}

// Find specific record from in memory cache
func (provider *inMemoryProvider) Find(normalizedName string) (Category, bool) {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	val, exists := provider.categories[normalizedName]
	return val, exists
}

func (provider *inMemoryProvider) runSync() {
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
