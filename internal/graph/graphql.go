package graph

import (
	"github.com/calvarado2004/go-movies-backend/internal/models"
	"github.com/graphql-go/graphql"
	"strings"
)

type Graph struct {
	Movies      []*models.Movie
	QueryString string
	Config      graphql.SchemaConfig
	fields      graphql.Fields
	movieType   *graphql.Object
}

func NewGraph(movies []*models.Movie) *Graph {

	var movieType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Movie",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"description": &graphql.Field{
					Type: graphql.String,
				},
				"release_date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"runtime": &graphql.Field{
					Type: graphql.Int,
				},
				"mpaa_rating": &graphql.Field{
					Type: graphql.String,
				},
				"created_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"updated_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"image": &graphql.Field{
					Type: graphql.String,
				},
				"genres": &graphql.Field{
					Type: graphql.NewList(graphql.String),
				},
			},
		},
	)

	var fields = graphql.Fields{
		"list": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Get all movies",
			Resolve: func(params graphql.ResolveParams) (any, error) {
				return movies, nil
			},
		},

		"searching": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Search movies by title",
			Args: graphql.FieldConfigArgument{
				"titleContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (any, error) {
				var filteredMovies []*models.Movie
				titleContains, ok := params.Args["titleContains"].(string)
				if ok {
					for _, movie := range movies {
						if strings.Contains(strings.ToLower(movie.Title), strings.ToLower(titleContains)) {
							filteredMovies = append(filteredMovies, movie)
						}
					}
				}

				return filteredMovies, nil
			},
		},

		"get": &graphql.Field{
			Type:        movieType,
			Description: "Get movie by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (any, error) {
				id, ok := params.Args["id"].(int)
				if ok {
					for _, movie := range movies {
						if movie.ID == id {
							return movie, nil
						}
					}
				}

				return nil, nil
			},
		},
	}

	return &Graph{
		Movies:    movies,
		fields:    fields,
		movieType: movieType,
	}
}
