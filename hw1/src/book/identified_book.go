package book

import "fmt"

type IdentifiedBook struct {
	ID uint64
	Book
}

func (identifiedBook IdentifiedBook) Validate() error {
	if identifiedBook.ID == 0 {
		return fmt.Errorf("id must be non-zero")
	}
	return identifiedBook.Book.Validate()
}

func (identifiedBook IdentifiedBook) String() string {
	return fmt.Sprintf("#%d: %s", identifiedBook.ID, identifiedBook.Book.String())
}
