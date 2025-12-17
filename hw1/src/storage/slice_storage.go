package storage

import (
	"fmt"

	"github.com/cchrisris/go-course/hw1/book"
)

type SliceStorage struct {
	items []book.IdentifiedBook
}

func NewSliceStorage() *SliceStorage {
	return &SliceStorage{items: make([]book.IdentifiedBook, 0)}
}

func (sliceStorage *SliceStorage) AddBook(identifiedBook book.IdentifiedBook) error {
	if err := identifiedBook.Validate(); err != nil {
		return fmt.Errorf("invalid book: %w", err)
	}

	for index := range sliceStorage.items {
		if sliceStorage.items[index].ID == identifiedBook.ID {
			sliceStorage.items[index] = identifiedBook
			return nil
		}
	}
	sliceStorage.items = append(sliceStorage.items, identifiedBook)
	return nil
}

func (sliceStorage *SliceStorage) GetBook(id uint64) (book.IdentifiedBook, bool) {
	for _, item := range sliceStorage.items {
		if item.ID == id {
			return item, true
		}
	}
	return book.IdentifiedBook{}, false
}

func (sliceStorage *SliceStorage) DeleteBook(id uint64) bool {
	for index := range sliceStorage.items {
		if sliceStorage.items[index].ID == id {
			sliceStorage.items = append(sliceStorage.items[:index], sliceStorage.items[index+1:]...)
			return true
		}
	}
	return false
}

func (sliceStorage *SliceStorage) Clear() {
	sliceStorage.items = make([]book.IdentifiedBook, 0)
}

func (sliceStorage *SliceStorage) GetAllBooks() []book.IdentifiedBook {
	return append([]book.IdentifiedBook(nil), sliceStorage.items...)
}
