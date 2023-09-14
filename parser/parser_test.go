package parser

import (
	"reflect"
	"testing"

	"github.com/kaiquecaires/go-import-movies/models"
)

func TestParseLine(t *testing.T) {
	t.Run("Should return the movie with the correct year", func(t *testing.T) {
		input := []string{"123", "Foo bar (2023)", "foo|bar"}
		expected := models.Movie{
			Id:     123,
			Title:  "Foo bar",
			Year:   2023,
			Genres: []string{"foo", "bar"},
		}
		result, _ := ParseLine(input[:])

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v but got %+v", expected, result)
		}
	})

	t.Run("Should return the movie with the year empty", func(t *testing.T) {
		input := []string{"123", "Foo bar", "foo|bar"}
		expected := models.Movie{
			Id:     123,
			Title:  "Foo bar",
			Year:   0,
			Genres: []string{"foo", "bar"},
		}
		result, _ := ParseLine(input[:])
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v but got %+v", expected, result)
		}
	})

	t.Run("Should return an error if the id is not a string", func(t *testing.T) {
		input := []string{"invalid id", "Foo bar", "foo|bar"}
		expected := ""
		_, err := ParseLine(input[:])
		if err == nil {
			t.Errorf("Expected %+v but got %+v", expected, err)
		}
	})

	t.Run("Should ignore incorrect year", func(t *testing.T) {
		input := []string{"123", "Foo bar (202a)", "foo|bar"}
		expected := models.Movie{
			Id:     123,
			Title:  "Foo bar (202a)",
			Year:   0,
			Genres: []string{"foo", "bar"},
		}
		result, _ := ParseLine(input[:])

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v but got %+v", expected, result)
		}
	})
}
