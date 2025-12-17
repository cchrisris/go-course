package book

import (
	"fmt"
	"strings"
)

type Book struct {
	Title         string
	Author        string
	Language      string
	Publisher     string
	PublishedYear int
	Pages         int
}

func (book Book) Validate() error {
	if strings.TrimSpace(book.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(book.Author) == "" {
		return fmt.Errorf("author is required")
	}
	if book.Pages < 0 {
		return fmt.Errorf("pages must be non-negative")
	}
	if book.PublishedYear < 0 {
		return fmt.Errorf("published year must be non-negative")
	}
	return nil
}

func (book Book) String() string {
	year := ""
	if book.PublishedYear > 0 {
		year = fmt.Sprintf("%d", book.PublishedYear)
	}

	parts := []string{
		strings.TrimSpace(book.Title),
		strings.TrimSpace(book.Author),
	}
	if year != "" {
		parts = append(parts, year)
	}
	return strings.Join(parts, " — ")
}
