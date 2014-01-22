package task

import (
	"strconv"
	"sync"
)

type Env struct {
	Mutex sync.RWMutex
	Vars  map[string]string
}

func NewEnvironment() *Env {
	return &Env{
		Vars: make(map[string]string),
	}
}

func (e *Env) Get(key string) string {
	e.Mutex.RLock()
	defer e.Mutex.RUnlock()
	return e.Vars[key]
}

func (e *Env) Set(key string, val string) {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	e.Vars[key] = val
}

func (e *Env) GetString(key string) string {
	return e.Get(key)
}

func (e *Env) SetString(key string, val string) {
	e.Set(key, val)
}

func (e *Env) GetInt(key string) int {
	val, _ := strconv.ParseInt(e.Get(key), 10, 0)
	return int(val)
}

func (e *Env) SetInt(key string, val int) {
	e.Set(key, strconv.FormatInt(int64(val), 10))
}

func (e *Env) GetBool(key string) bool {
	val, _ := strconv.ParseBool(e.Get(key))
	return val
}

func (e *Env) SetBool(key string, val bool) {
	e.Set(key, strconv.FormatBool(val))
}

func (e *Env) Exists(key string) bool {
	return e.Get(key) != ""
}
