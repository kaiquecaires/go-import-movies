package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiquecaires/go-import-movies/models"
)

type Writer struct {
	Pool *pgxpool.Pool
}

func (repo *Writer) InsertMovie(movie models.Movie) error {
	_, err := repo.Pool.Exec(
		context.Background(),
		"INSERT INTO movies (id, title, year, Genres) values ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET created_at = NOW();",
		movie.Id,
		movie.Title,
		movie.Year,
		movie.Genres,
	)

	return err
}
