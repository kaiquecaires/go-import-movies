package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kaiquecaires/go-import-movies/db"
	"github.com/kaiquecaires/go-import-movies/parser"
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

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(reader, &wg, &writer)
	}

	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	fmt.Println("Done. Elapsed Milliseconds: ", elapsedTime.Milliseconds())
}

var mu sync.Mutex

func worker(reader *csv.Reader, wg *sync.WaitGroup, writer *db.Writer) {
	defer wg.Done()
	for {
		mu.Lock()
		record, err := reader.Read()
		mu.Unlock()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading CSV record:", err)
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
}
