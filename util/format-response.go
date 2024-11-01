package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jramirezm00/preview-task-acquia/models"
	"github.com/jramirezm00/preview-task-acquia/services"
)

func CreateResponse() error {
	// get all the films
	filmsRs, filmsRsErr := services.ListStarWarsFilms()
	if filmsRsErr != nil {
		return filmsRsErr
	}

	// get all the planets
	planets, planetsErr := services.HandlePlanetsPagination()
	if planetsErr != nil {
		return planetsErr
	}

	//build film-planets response
	response := buildResponse(filmsRs.Results, planets)

	fmt.Println(response)

	return nil
}

// func for creating the final response
func buildResponse(films []services.Film, planets []services.Planet) string {
	planetsMap := make(map[string]services.Planet, 0)
	var finalResponse models.FinalResponse

	//create a map so it's easier to find the necessary planets
	for _, planet := range planets {
		planetsMap[planet.Url] = planet
	}

	// no need to sort the films, already in asc order from the swapi
	for _, film := range films {
		var filmPlanets []string

		//add the planet name to the respective film
		for _, planet := range film.Planets {
			filmPlanets = append(filmPlanets, planetsMap[planet].Name)
		}

		//sort planets alphabetically
		sort.Strings(filmPlanets)

		finalResponse = append(finalResponse, models.FinalResponseItem{
			Title:   film.Title,
			Planets: filmPlanets,
		})
	}

	finalResponseBytes, finalResponseBytesErr := finalResponse.MarshalJSON()
	if finalResponseBytesErr != nil {
		panic(finalResponseBytesErr)
	}

	dst := &bytes.Buffer{}
	if err := json.Indent(dst, finalResponseBytes, "", "  "); err != nil {
		panic(err)
	}

	return dst.String()
}
