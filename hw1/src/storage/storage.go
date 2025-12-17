package storage

import "github.com/cchrisris/go-course/hw1/book"

type Storage interface {
	AddBook(book book.IdentifiedBook) error
	GetBook(id uint64) (book.IdentifiedBook, bool)
	DeleteBook(id uint64) bool
	Clear()
	GetAllBooks() []book.IdentifiedBook
}
