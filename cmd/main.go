package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaiquecaires/go-import-movies/db"
	"github.com/kaiquecaires/go-import-movies/parser"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	file, err := os.Open("./assets/movie.csv")

	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)

	pgpool, err := db.CreatePostgresPool()

	if err != nil {
		fmt.Println("Error opening pool with postgres:", err)
	}

	writer := db.Writer{
		Pool: pgpool,
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading CSV record:", err)
			continue
		}

		movie, err := parser.ParseLine(record)

		if err != nil {
			continue
		}

		err = writer.InsertMovie(movie)

		if err != nil {
			fmt.Println("Error inserting value:", err)
			continue
		}
	}
}
