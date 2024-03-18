package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []Book{
	{
		ID:       "978-0544003415",
		Name:     "The Hobbit",
		Author:   "J.R.R. Tolkien",
		Quantity: 5,
	},
	{
		ID:       "978-0261103573",
		Name:     "The Lord of the Rings",
		Author:   "J.R.R. Tolkien",
		Quantity: 4,
	},
	{
		ID:       "978-0061120084",
		Name:     "To Kill a Mockingbird",
		Author:   "Harper Lee",
		Quantity: 2,
	},
}

func getbooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createbooks(c *gin.Context) {
	var newbook Book
	if err := c.BindJSON(&newbook); err != nil {
		return
	}
	books = append(books, newbook)
	c.IndentedJSON(http.StatusCreated, books)
}

func bookbyID(c *gin.Context) {
	id := c.Param("id")
	book, err := getbookbyID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Bad Request"})
		return
	}
	c.IndentedJSON(http.StatusFound, book)
}

func getbookbyID(id string) (*Book, error) {
	for i := 0; i < len(books); i++ {
		if books[i].ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func checkoutbook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}
	book, err := getbookbyID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity -= 1
	response := gin.H{
		"message":      "Book checked out successfully",
		"book_details": book,
	}
	c.IndentedJSON(http.StatusOK, response)

}

func returnbook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query patrameter id"})
		return
	}
	book, err := getbookbyID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity += 1
	response := gin.H{
		"message":      "Book returned successfully",
		"book_details": book,
	}
	c.IndentedJSON(http.StatusOK, response)
}

func main() {

	router := gin.Default()
	router.GET("/books", getbooks)
	router.POST("/add_book", createbooks)
	router.GET("/getbook/:id", bookbyID)
	router.PATCH("/checkout", checkoutbook)
	router.PATCH("/returnbook", returnbook)
	router.Run("localhost:8081")

}
