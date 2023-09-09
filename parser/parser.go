package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kaiquecaires/go-import-movies/models"
)

func ParseLine(line []string) (models.Movie, error) {
	model := models.Movie{}

	id, err := strconv.Atoi(line[0])

	if err != nil {
		return model, err
	}

	model.Id = id

	regex := regexp.MustCompile(`^(.*?)\s\((\d{4})\)$`)
	matches := regex.FindStringSubmatch(line[1])

	if len(matches) == 3 {
		model.Title = matches[1]

		year, err := strconv.Atoi(matches[2])

		if err == nil {
			model.Year = year
		}
	} else {
		model.Title = line[1]
	}

	if len(line) == 3 {
		model.Genres = strings.Split(line[2], "|")
	}

	return model, err
}
