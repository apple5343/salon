package generator

import (
	"errors"
	"math/rand/v2"
)

type Storage[T any] struct {
	available map[string]T
	pending   map[string]T
	created   map[string]T
}

func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		available: make(map[string]T),
		pending:   make(map[string]T),
		created:   make(map[string]T),
	}
}

func (s *Storage[T]) AddAvailable(id string, item T) {
	s.available[id] = item
}

func (s *Storage[T]) AddPending(id string, item T) {
	s.pending[id] = item
}

func (s *Storage[T]) GetAvailable(id string) (T, bool) {
	item, ok := s.available[id]
	return item, ok
}

func (s *Storage[T]) GetPending(id string) (T, bool) {
	item, ok := s.pending[id]
	return item, ok
}

func (s *Storage[T]) GetCreated(id string) (T, bool) {
	item, ok := s.created[id]
	return item, ok
}

func (s *Storage[T]) GetOneCreated() (T, error) {
	if len(s.created) == 0 {
		return *new(T), errors.New("no available items")
	}
	keys := make([]string, 0, len(s.created))
	for k := range s.created {
		keys = append(keys, k)
	}
	id := keys[rand.IntN(len(keys))]
	item := s.created[id]
	return item, nil
}

func (s *Storage[T]) PickOne() (T, error) {
	if len(s.available) == 0 {
		return *new(T), errors.New("no available items")
	}
	keys := make([]string, 0, len(s.available))
	for k := range s.available {
		keys = append(keys, k)
	}
	id := keys[rand.IntN(len(keys))]
	item := s.available[id]
	delete(s.available, id)
	s.pending[id] = item
	return item, nil
}

func (s *Storage[T]) MoveToCreated(id string, serviceID string, serviceItem T) error {
	if _, ok := s.pending[id]; !ok {
		return errors.New("item not pending")
	}
	s.created[serviceID] = serviceItem
	delete(s.pending, id)
	return nil
}

func (s *Storage[T]) MoveToAvailable(id string, serviceID string, serviceItem T) error {
	if _, ok := s.pending[id]; !ok {
		return errors.New("item not pending")
	}
	s.available[serviceID] = serviceItem
	delete(s.pending, id)
	return nil
}

func (s *Storage[T]) AvaliableCount() int {
	return len(s.available)
}

func (s *Storage[T]) CreatedCount() int {
	return len(s.created)
}
