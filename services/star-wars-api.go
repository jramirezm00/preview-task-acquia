package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type FilmsResponse struct {
	Results []Film `json:"results"`
}

type Film struct {
	Title       string   `json:"title"`
	ReleaseDate string   `json:"release_date"`
	Planets     []string `json:"planets"`
}

type PlanetsResponse struct {
	Next     *string  `json:"next"`
	Previous string   `json:"previous"`
	Results  []Planet `json:"results"`
}

type Planet struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const (
	listFilmsURL = "https://swapi.dev/api/films/"
)

func ListStarWarsFilms() (*FilmsResponse, error) {
	resp, err := http.Get(listFilmsURL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("request to fetch films failed with code %d and status %s", resp.StatusCode, resp.Status)
		panic(err)
	}

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	var response FilmsResponse

	unmarshalErr := json.Unmarshal(body, &response)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	return &response, nil
}

func ListStarWarsPlanets(url string, ch chan<- *PlanetsResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("request to fetch planet failed with code %d and status %s", resp.StatusCode, resp.Status)
		panic(err)
	}

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}

	var response PlanetsResponse

	unmarshalErr := json.Unmarshal(body, &response)
	if unmarshalErr != nil {
		panic(unmarshalErr)
	}

	ch <- &response
}

func HandlePlanetsPagination() ([]Planet, error) {
	var planets []Planet
	listPlanetsURLs := []string{"https://swapi.dev/api/planets/?page=1", "https://swapi.dev/api/planets/?page=2", "https://swapi.dev/api/planets/?page=3",
		"https://swapi.dev/api/planets/?page=4", "https://swapi.dev/api/planets/?page=5", "https://swapi.dev/api/planets/?page=6"}

	ch := make(chan *PlanetsResponse)
	var wg sync.WaitGroup

	for _, url := range listPlanetsURLs {
		wg.Add(1)
		go ListStarWarsPlanets(url, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for planetRs := range ch {
		planets = append(planets, planetRs.Results...)
	}

	return planets, nil
}
