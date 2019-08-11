package set

import (
	"sort"
	"sync"
)

type Set struct {
	sync.RWMutex
	m map[string]bool
}

// 新建集合对象
func New(items []string) *Set {
	s := &Set{
		m: make(map[string]bool, len(items)),
	}
	s.Add(items)
	return s
}

// 添加元素
func (s *Set) Add(items []string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// 无序列表
func (s *Set) List() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// 排序列表
func (s *Set) SortList() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}

// 差集
func (s *Set) Minus(sets ...*Set) *Set {
	r := New(s.List())
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}
