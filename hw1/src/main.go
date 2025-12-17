package main

import (
	"fmt"

	"github.com/cchrisris/go-course/hw1/book"
	"github.com/cchrisris/go-course/hw1/library"
	"github.com/cchrisris/go-course/hw1/storage"
)

func main() {
	fmt.Println("Library demo")

	initialBooks := []book.Book{
		{
			Title:         "Strange Case of Dr Jekyll and Mr Hyde",
			Author:        "Robert Louis Stevenson",
			Language:      "en",
			Publisher:     "Longmans, Green & Co.",
			PublishedYear: 1886,
			Pages:         144,
		},
		{
			Title:         "The Divine Comedy",
			Author:        "Dante Alighieri",
			Language:      "en",
			Publisher:     "None",
			PublishedYear: 1320,
			Pages:         798,
		},
		{
			Title:         "Cat's Cradle",
			Author:        "Kurt Vonnegut",
			Language:      "en",
			Publisher:     "Holt, Rinehart and Winston",
			PublishedYear: 1963,
			Pages:         304,
		},
		{
			Title:         "Concerning the Spiritual in Art",
			Author:        "Wassily Kandinsky",
			Language:      "en",
			Publisher:     "None",
			PublishedYear: 1911,
			Pages:         256,
		},
		{
			Title:         "Behave: The Biology of Humans at Our Best and Worst",
			Author:        "Robert Sapolsky",
			Language:      "en",
			Publisher:     "Penguin Press",
			PublishedYear: 2017,
			Pages:         800,
		},
	}

	libraryInstance := library.NewLibrary(storage.NewSliceStorage(), library.SequentialIDGenerator(1))
	panicIfError(libraryInstance.Load(initialBooks))

	fmt.Println("\nFind books")
	printLookup(libraryInstance, "Cat's Cradle")
	printLookup(libraryInstance, "The Divine Comedy")
	printLookup(libraryInstance, "Unknown Book")

	fmt.Println("\nSwap ID generator, add a book")
	panicIfError(libraryInstance.ReplaceIDGenerator(library.SequentialIDGenerator(1000)))
	addedBook, addError := libraryInstance.AddBook(book.Book{
		Title:         "The Prince",
		Author:        "Niccolo Machiavelli",
		Language:      "en",
		Publisher:     "None",
		PublishedYear: 1532,
		Pages:         164,
	})
	panicIfError(addError)
	fmt.Printf("Added: %s\n", addedBook.String())
	printLookup(libraryInstance, "The Prince")

	fmt.Println("\nSwap storage (slice -> map)")
	panicIfError(libraryInstance.RebuildStorage(storage.NewMapStorage()))
	printLookup(libraryInstance, "Strange Case of Dr Jekyll and Mr Hyde")
	printLookup(libraryInstance, "Behave: The Biology of Humans at Our Best and Worst")

	fmt.Println("\nLoad more books")
	panicIfError(libraryInstance.Load([]book.Book{
		{
			Title:         "The Art of Color",
			Author:        "Johannes Itten",
			Language:      "en",
			Publisher:     "John Wiley & Sons",
			PublishedYear: 1961,
			Pages:         196,
		},
		{
			Title:         "Faust",
			Author:        "Johann Wolfgang von Goethe",
			Language:      "en",
			Publisher:     "None",
			PublishedYear: 1808,
			Pages:         464,
		},
		{
			Title:         "The Long Hard Road Out of Hell",
			Author:        "Marilyn Manson",
			Language:      "en",
			Publisher:     "ReganBooks",
			PublishedYear: 1998,
			Pages:         288,
		},
		{
			Title:         "Leviathan",
			Author:        "Thomas Hobbes",
			Language:      "en",
			Publisher:     "Andrew Crooke",
			PublishedYear: 1651,
			Pages:         736,
		},
	}))

	printLookup(libraryInstance, "Faust")
	printLookup(libraryInstance, "The Art of Color")

	fmt.Println("\nDone.")
}

func printLookup(libraryInstance *library.Library, title string) {
	foundBook, isFound := libraryInstance.GetBook(title)
	fmt.Printf("- %s: ", title)
	if !isFound {
		fmt.Println("not found")
		return
	}
	fmt.Println(foundBook.String())
}

func panicIfError(operationError error) {
	if operationError != nil {
		panic(operationError)
	}
}
