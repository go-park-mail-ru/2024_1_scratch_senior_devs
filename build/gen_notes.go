package main

import (
	"context"
	"encoding/json"
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

func randomDate() time.Time {
	yearAgo := time.Now().AddDate(-1, 0, 0).Unix()
	now := time.Now().Unix()
	randomUnix := rand.Int63n(now-yearAgo) + yearAgo //nolint:all

	randomTime := time.Unix(randomUnix, 0)
	return randomTime.UTC()
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
		noteTitle := gofakeit.Sentence(rand.Int()%8 + 3)                                          //nolint:all
		noteContent := gofakeit.Paragraph(rand.Int()%10+1, rand.Int()%7+1, rand.Int()%10+3, "\n") //nolint:all

		noteData := map[string]string{
			"title":   noteTitle,
			"content": noteContent,
		}

		jsonData, err := json.Marshal(noteData)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = db.Exec(context.Background(), addNote, uuid.NewV4(), jsonData, time.Now().AddDate(-1, 0, 0).UTC(), randomDate(), userID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("SUCCESS")
}
