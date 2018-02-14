package util

import (
	"path"
	"runtime"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

type Set struct {
	buff map[string]struct{}
}

func NewSet() *Set {
	return &Set{buff: map[string]struct{}{}}
}

func (s *Set) Append(i string) {
	s.buff[i] = struct{}{}
}

func (s *Set) Include(i string) bool {
	return Include(s.ToSlice(), i)
}

func (s *Set) ToSlice() []string {
	keys := make([]string, 0, len(s.buff))
	for k := range s.buff {
		keys = append(keys, k)
	}
	return keys
}

func GetCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}