package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Books struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CheckedIn int32  `json:"checked"`
}

/*func startUp() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "Ba$h2202",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "alexandria",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}*/

func checkOut(c *gin.Context) {
	var bk Books
	id := c.Param("id")

	row := db.QueryRow("SELECT * FROM library WHERE id = ?", id)
	if err := row.Scan(&bk.ID, &bk.Title, &bk.Author, &bk.CheckedIn); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Failed to query db %v\n", err)
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "Title not found, maybe try adding it?"})
			return
		}
		log.Printf("Failed to query db %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to contact db"})
		return
	}
	bk.CheckedIn = 0
	result, err := db.Exec("UPDATE library SET checked = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("Failed to update db %v", err)
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "Failed to update db"})
		return
	}
	change, err := result.RowsAffected()
	if err != nil || change < 1 || change > 2 {
		log.Printf("Failed to update db %v", err)
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "Failed to update db"})
		return
	}

	c.IndentedJSON(http.StatusOK, bk)
	return
}

func checkIn(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("UPDATE library SET checked = 1 WHERE id = ?", id)
	if err != nil {
		log.Printf("Failed to update db %v", err)
		c.IndentedJSON(http.StatusPreconditionFailed, gin.H{"message": "Failed to update db"})
	}
	res, err := result.RowsAffected()
	if err != nil || res < 1 {
		log.Printf("Failed to update db %v", err)
		c.IndentedJSON(http.StatusPreconditionFailed, gin.H{"message": "Failed to update db"})
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "Succesfully returned"})
	return
}

func index(c *gin.Context) {
	var book []Books

	rows, err := db.Query("SELECT * FROM library")
	if err != nil {
		log.Printf("failed to query * from db %v", err)
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Failed to contact db"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var bk Books
		if err := rows.Scan(&bk.ID, &bk.Title, &bk.Author, &bk.CheckedIn); err != nil {
			log.Printf("Failed to process data from db %v", err)
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "Failed to contact db"})
			return
		}
		book = append(book, bk)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to contact db %v", err)
		c.IndentedJSON(http.StatusPreconditionFailed, gin.H{"message": "Failed to contact db"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
	return
}

func aboutUs(c *gin.Context) {
	c.File("exec/aboutUs.html")
}

func donate(c *gin.Context) {
	var newBook Books
	if err := c.BindJSON(&newBook); err != nil {
		log.Printf("Failed to add new book %v", err)
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "improperly formated application"})
		return
	}
	results, err := db.Exec("INSERT INTO library (title, author, checked) VALUES (?, ?, ?)", newBook.Title, newBook.Author, newBook.CheckedIn)
	if err != nil {
		log.Printf("Failed to add new book %v", err)
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "Failed to update db"})
	}
	_, writeErr := results.LastInsertId()
	if writeErr != nil {
		log.Printf("Failed to update id num %v", err)
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "Failed to update db num"})
	}

	c.IndentedJSON(http.StatusCreated, newBook)
	return
}

func main() {
	// db := startUp()
	cfg := mysql.Config{
		User:   "root",
		Passwd: "Ba$h2202",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "alexandria",
	}

	cfg.AllowNativePasswords = true

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	router := gin.Default()

	router.GET("/", index)
	router.GET("/about", aboutUs)
	router.POST("/checkIn/:id", checkIn)
	router.GET("/checkOut/:id", checkOut)
	router.POST("/donate", donate)

	router.Run("localhost:42069")
}
