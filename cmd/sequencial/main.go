package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kaiquecaires/go-import-movies/parser"

	"github.com/joho/godotenv"
	"github.com/kaiquecaires/go-import-movies/db"
)

func main() {
	startTime := time.Now()

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

	i := -1
	line := 0

	for {
		line++
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(line-1, "Error reading CSV record:", err)
			continue
		}

		if i == -1 {
			i++
			continue
		}

		movie, err := parser.ParseLine(record)

		if err != nil {
			fmt.Println(record, err)
			continue
		}

		err = writer.InsertMovie(movie)

		if err != nil {
			fmt.Println("Error inserting value:", err)
		}
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	fmt.Println("Done. Elapsed Milliseconds: ", elapsedTime.Milliseconds())
}
