package categories

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/text/unicode/norm"
)

type CategoriesProvider interface {
	GetAll() []Category
	Find(name string) (Category, bool)
}

type InMemoryCategoriesProvider struct {
	categories  map[string]Category
	updatesChan <-chan Category
	context     context.Context
	mutex       *sync.RWMutex
}

func NewProvider(seed []Category, updates <-chan Category, context context.Context) (*InMemoryCategoriesProvider, error) {
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

	provider := &InMemoryCategoriesProvider{
		categories:  categories,
		updatesChan: updates,
		context:     context,
		mutex:       &sync.RWMutex{},
	}
	go provider.runSync()
	return provider, nil
}

func (provider *InMemoryCategoriesProvider) GetAll() []Category {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	arr := []Category{}
	for _, v := range provider.categories {
		arr = append(arr, v)
	}
	return arr
}

func (provider *InMemoryCategoriesProvider) Find(name string) (Category, bool) {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()
	normalized := norm.NFC.String(name)
	val, exists := provider.categories[normalized]
	return val, exists
}

func (provider *InMemoryCategoriesProvider) runSync() {
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
