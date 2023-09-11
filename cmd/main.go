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

	const numWorkers = 20
	dataChan := make(chan []string, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(dataChan, &wg, writer)
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

		dataChan <- record
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	fmt.Println("Done. Elapsed Milliseconds: ", elapsedTime.Milliseconds())
}

func worker(dataChan <-chan []string, wg *sync.WaitGroup, writer db.Writer) {
	defer wg.Done()

	for data := range dataChan {
		movie, err := parser.ParseLine(data)

		if err != nil {
			return
		}

		err = writer.InsertMovie(movie)

		if err != nil {
			fmt.Println("Error inserting value:", err)
		}
	}
}
