package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/satori/uuid"
	"math/rand"
	"os"
	"time"
)

const (
	addNote         = "INSERT INTO notes(id, data, create_time, update_time, owner_id) VALUES ($1, $2::json, $3, $4, $5);"
	getIdByUsername = "SELECT id FROM users WHERE username = $1 "
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var username string
	fmt.Print("Введите username: ")
	fmt.Scan(&username)

	var count int
	fmt.Print("Введите количество заметок: ")
	fmt.Scan(&count)

	var userID uuid.UUID
	err = db.QueryRow(context.Background(), getIdByUsername, username).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < count; i++ {
		noteTitle := gofakeit.Sentence(rand.Int()%8 + 3)
		noteContent := gofakeit.Paragraph(1, rand.Int()%10+1, rand.Int()%8+3, "\n")
		noteData := `{"title":"` + noteTitle + `","content":"` + noteContent + `"}`

		_, err = db.Exec(context.Background(), addNote, uuid.NewV4(), noteData, time.Now().UTC(), time.Now().UTC(), userID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("SUCCESS")
}
