package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	Error "github.com/mayank12gt/free-audiobooks-backend/internal/errors"
	"github.com/mayank12gt/free-audiobooks-backend/internal/repos"
	"github.com/mayank12gt/free-audiobooks-backend/internal/services"
)

type Response struct {
	Metadata   repos.Metadata     `json:"metadata"`
	Audiobooks []*repos.Audiobook `json:"audiobooks"`
}
type GenresResponse struct {
	Metadata repos.Metadata    `json:"metadata"`
	Genres   []*repos.GenreDTO `json:"genres"`
}

type ApiError struct {
	Error map[string]error `json:"error"`
}

func (app *app) listHandler() func(c echo.Context) error {
	return func(c echo.Context) error {

		search := c.QueryParam("search")
		language := c.QueryParam("language")
		sortBy := c.QueryParam("sort_by")

		var totalTimeMin, totalTimeMax int
		var err error
		if c.QueryParam("lengthMin") != "" {
			totalTimeMin, err = strconv.Atoi(c.QueryParam("lengthMin"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("length", "Must be an integer"))
			}
		} else {
			totalTimeMin = 0
		}
		if c.QueryParam("lengthMax") != "" {
			totalTimeMax, err = strconv.Atoi(c.QueryParam("lengthMax"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("length", "Must be an integer"))
			}
		} else {
			totalTimeMax = 0
		}

		query := services.Query{
			Search:   search,
			Language: language,
			TotalTimeRange: services.TimeRange{
				TotalTimeMin: int64(totalTimeMin),
				TotalTimeMax: int64(totalTimeMax),
			},
			Sort: sortBy,
		}

		if c.QueryParam("genres") != "" {
			query.Genres = strings.Split(c.QueryParam("genres"), ",")
		} else {
			query.Genres = []string{}
		}

		if c.QueryParam("page") != "" {
			query.Page, err = strconv.Atoi(c.QueryParam("page"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("page", "Must be an integer"))
			}
		} else {
			query.Page = 1
		}

		if c.QueryParam("page_size") != "" {
			query.PageSize, err = strconv.Atoi(c.QueryParam("page_size"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("page_size", "Must be an integer"))
			}
		} else {
			query.PageSize = 20
		}

		error := query.Validate()
		if error != nil {
			log.Print(error)
			return c.JSON(400, error)
		}

		audiobooks, meta, err := app.services.AudiobooksService.List(query)
		if err != nil {
			log.Print(err)
			if len(audiobooks) == 0 {
				return c.JSON(404, err)
			}
			return c.JSON(500, err)
		}
		return c.JSON(200, Response{
			Metadata:   meta,
			Audiobooks: audiobooks,
		})

	}
}

func (app *app) GetHandler() func(c echo.Context) error {
	return func(c echo.Context) error {

		id := (c.Param("id"))
		// if err != nil {
		// 	return c.JSON(400, err.Error())
		// }

		audiobook, err := app.services.AudiobooksService.Get(id)
		if err != nil {
			return c.JSON(400, err)
		}

		return c.JSON(200, audiobook)

	}
}

func (app *app) ListGenresHandler() func(c echo.Context) error {
	return func(c echo.Context) error {

		var page, page_size int
		var err error

		if c.QueryParam("page_size") != "" {
			page_size, err = strconv.Atoi(c.QueryParam("page_size"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("page_size", "Must be an integer"))
			}
		} else {
			page_size = 20
		}

		if c.QueryParam("page") != "" {
			page, err = strconv.Atoi(c.QueryParam("page"))
			if err != nil {
				return c.JSON(http.StatusBadRequest, Error.NewError().Set("page", "Must be an integer"))
			}
		} else {
			page = 1
		}
		genres, meta, err := app.services.AudiobooksService.GetGenres(page, page_size)

		if err != nil {
			return c.JSON(500, "Error")
		}

		return c.JSON(200, GenresResponse{
			Metadata: meta,
			Genres:   genres,
		})
	}
}
