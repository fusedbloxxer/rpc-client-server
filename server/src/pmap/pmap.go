package pmap

import (
	"sync"
)

type ConcurrentMap interface {
	Exists(string) bool
	Keys() []string
	Delete(string)
	Add(string)
	Size() int
}

type PMap struct {
	Map	map[string]bool
	Mutex sync.Mutex
}

func (pm *PMap) Exists(key string) bool {
	pm.Mutex.Lock()
	_, exists := pm.Map[key]
	pm.Mutex.Unlock()
	return exists
}

func (pm *PMap) Add(key string) {
	pm.Mutex.Lock()
	pm.Map[key] = true
	pm.Mutex.Unlock()
}

func (pm *PMap) Delete(key string) {
	pm.Mutex.Lock()
	delete(pm.Map, key)
	pm.Mutex.Unlock()
}

func (pm *PMap) Size() (sz int) {
	pm.Mutex.Lock()
	sz = len((*pm).Map)
	pm.Mutex.Unlock()
	return
}

func (pm *PMap) Keys() (keys []string) {
	pm.Mutex.Lock()
	keys = make([]string, 0, len((*pm).Map))
	for key := range pm.Map {
		keys = append(keys, key)
	}
	pm.Mutex.Unlock()
	return
}
