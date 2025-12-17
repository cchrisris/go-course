package storage

import (
	"fmt"

	"github.com/cchrisris/go-course/hw1/book"
)

type MapStorage struct {
	byID map[uint64]book.IdentifiedBook
}

func NewMapStorage() *MapStorage {
	return &MapStorage{byID: make(map[uint64]book.IdentifiedBook)}
}

func (mapStorage *MapStorage) AddBook(identifiedBook book.IdentifiedBook) error {
	if err := identifiedBook.Validate(); err != nil {
		return fmt.Errorf("invalid book: %w", err)
	}
	if mapStorage.byID == nil {
		mapStorage.byID = make(map[uint64]book.IdentifiedBook)
	}
	mapStorage.byID[identifiedBook.ID] = identifiedBook
	return nil
}

func (mapStorage *MapStorage) GetBook(id uint64) (book.IdentifiedBook, bool) {
	found, ok := mapStorage.byID[id]
	return found, ok
}

func (mapStorage *MapStorage) DeleteBook(id uint64) bool {
	if _, ok := mapStorage.byID[id]; !ok {
		return false
	}
	delete(mapStorage.byID, id)
	return true
}

func (mapStorage *MapStorage) Clear() {
	mapStorage.byID = make(map[uint64]book.IdentifiedBook)
}

func (mapStorage *MapStorage) GetAllBooks() []book.IdentifiedBook {
	result := make([]book.IdentifiedBook, 0, len(mapStorage.byID))
	for _, value := range mapStorage.byID {
		result = append(result, value)
	}
	return result
}
