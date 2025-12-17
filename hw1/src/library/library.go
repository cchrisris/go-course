package library

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cchrisris/go-course/hw1/book"
	"github.com/cchrisris/go-course/hw1/storage"
)

var (
	ErrMissingTitle       = errors.New("title is required")
	ErrTitleAlreadyExists = errors.New("book with this title already exists")
	ErrMissingIDGenerator = errors.New("id generator is nil")
	ErrMissingStorage     = errors.New("storage is nil")
)

type Library struct {
	storage     storage.Storage
	idGenerator func() uint64
	titleIndex  map[string]uint64
}

func NewLibrary(storage storage.Storage, idGenerator func() uint64) *Library {
	return &Library{
		storage:     storage,
		idGenerator: idGenerator,
		titleIndex:  make(map[string]uint64),
	}
}

func (library *Library) ReplaceIDGenerator(idGenerator func() uint64) error {
	if idGenerator == nil {
		return ErrMissingIDGenerator
	}
	library.idGenerator = idGenerator
	return nil
}

func (library *Library) RebuildStorage(newStorage storage.Storage) error {
	if newStorage == nil {
		return ErrMissingStorage
	}
	if library.storage == nil {
		library.storage = newStorage
		library.rebuildIndex()
		return nil
	}

	for _, identifiedBook := range library.storage.GetAllBooks() {
		if err := newStorage.AddBook(identifiedBook); err != nil {
			return fmt.Errorf("migrate book id=%d: %w", identifiedBook.ID, err)
		}
	}

	library.storage.Clear()
	library.storage = newStorage
	library.rebuildIndex()
	return nil
}

func (library *Library) AddBook(bookToAdd book.Book) (book.IdentifiedBook, error) {
	if library.storage == nil {
		return book.IdentifiedBook{}, ErrMissingStorage
	}
	if library.idGenerator == nil {
		return book.IdentifiedBook{}, ErrMissingIDGenerator
	}

	title := strings.TrimSpace(bookToAdd.Title)
	if title == "" {
		return book.IdentifiedBook{}, ErrMissingTitle
	}
	if _, exists := library.titleIndex[title]; exists {
		return book.IdentifiedBook{}, ErrTitleAlreadyExists
	}
	if err := bookToAdd.Validate(); err != nil {
		return book.IdentifiedBook{}, err
	}

	assignedID, generateError := library.generateUniqueID()
	if generateError != nil {
		return book.IdentifiedBook{}, generateError
	}

	identified := book.IdentifiedBook{ID: assignedID, Book: bookToAdd}
	if err := library.storage.AddBook(identified); err != nil {
		return book.IdentifiedBook{}, err
	}
	library.titleIndex[title] = assignedID
	return identified, nil
}

func (library *Library) Load(books []book.Book) error {
	for _, bookToAdd := range books {
		if _, err := library.AddBook(bookToAdd); err != nil {
			return err
		}
	}
	return nil
}

func (library *Library) GetBook(title string) (book.Book, bool) {
	if library.storage == nil {
		return book.Book{}, false
	}
	assignedID, hasTitle := library.titleIndex[strings.TrimSpace(title)]
	if !hasTitle {
		return book.Book{}, false
	}
	identifiedBook, found := library.storage.GetBook(assignedID)
	return identifiedBook.Book, found
}

func (library *Library) DeleteBook(title string) bool {
	if library.storage == nil {
		return false
	}
	normalizedTitle := strings.TrimSpace(title)
	assignedID, hasTitle := library.titleIndex[normalizedTitle]
	if !hasTitle {
		return false
	}
	if !library.storage.DeleteBook(assignedID) {
		return false
	}
	delete(library.titleIndex, normalizedTitle)
	return true
}

func (library *Library) rebuildIndex() {
	library.titleIndex = make(map[string]uint64)
	for _, identifiedBook := range library.storage.GetAllBooks() {
		normalizedTitle := strings.TrimSpace(identifiedBook.Title)
		if normalizedTitle == "" {
			continue
		}
		library.titleIndex[normalizedTitle] = identifiedBook.ID
	}
}

func (library *Library) generateUniqueID() (uint64, error) {
	const maxAttempts = 10000

	for attempt := 0; attempt < maxAttempts; attempt++ {
		assignedID := library.idGenerator()
		if assignedID == 0 {
			continue
		}
		if _, exists := library.storage.GetBook(assignedID); exists {
			continue
		}
		return assignedID, nil
	}
	return 0, fmt.Errorf("failed to generate unique id")
}
