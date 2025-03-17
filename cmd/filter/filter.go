package filter

import (
	"log"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type RequestFilter interface {
	IsEnabled() bool
	ShouldRequestBeDenied() bool
}

type DefaultFilter struct {
	AssociatedEnvVar string
}

type TMDBFilter struct {
	Client tmdb.Client
	DefaultFilter
}

func (d *DefaultFilter) IsEnabled() bool {
	value, ok := os.LookupEnv(d.AssociatedEnvVar)
	if ok {
		log.Printf("Filter for '%s' is enabled and loaded with the value '%s'", d.AssociatedEnvVar, value)
	} else {
		log.Printf("Filter for '%s' is enabled", d.AssociatedEnvVar)
	}
	return ok
}

type MinRatingFilter struct {
	TMDBFilter
}

func (m *MinRatingFilter) ShouldRequestBeDenied() bool {
	return false
}

func LoadAllFilters() []RequestFilter {
	filters := []RequestFilter{
		&MinRatingFilter{},
	}

	return filters
}
