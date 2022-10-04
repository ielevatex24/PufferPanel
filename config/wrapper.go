package config

import (
	"github.com/spf13/viper"
	"sync"
)

type entry[T ValueType] struct {
	key   string
	value T
}

type StringEntry struct {
	entry[string]
}
type BoolEntry struct {
	entry[bool]
}
type IntEntry struct {
	entry[int]
}
type Int64Entry struct {
	entry[int64]
}

type ValueType interface {
	int | int64 | bool | string
}

var cache = make(map[string]interface{})

func (se StringEntry) Value() string {
	return viper.GetString(se.Key())
}
func (se BoolEntry) Value() bool {
	return viper.GetBool(se.Key())
}
func (se IntEntry) Value() int {
	return viper.GetInt(se.Key())
}
func (se Int64Entry) Value() int64 {
	return viper.GetInt64(se.Key())
}

func (e entry[T]) Key() string {
	return e.key
}

func (e entry[T]) Set(value T, save bool) error {
	e.value = value

	if save {
		return viper.WriteConfig()
	}
	return nil
}

var locker sync.Mutex

func asString(key string, def string) StringEntry {
	return StringEntry{entry: as[string](key, def)}
}
func asBool(key string, def bool) BoolEntry {
	return BoolEntry{entry: as[bool](key, def)}
}
func asInt(key string, def int) IntEntry {
	return IntEntry{entry: as[int](key, def)}
}
func asInt64(key string, def int64) Int64Entry {
	return Int64Entry{entry: as[int64](key, def)}
}

func as[T ValueType](key string, def T) entry[T] {
	locker.Lock()
	defer locker.Unlock()

	val, exists := cache[key]
	if exists {
		return val.(entry[T])
	}

	e := entry[T]{key: key, value: def}
	cache[key] = e
	viper.SetDefault(key, def)
	return e
}
