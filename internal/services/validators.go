package services

import (
	Error "github.com/mayank12gt/free-audiobooks-backend/internal/errors"
)

func (q *Query) Validate() error {

	err := Error.NewError()
	if q.PageSize > 50 {
		err = err.Set("page_size", "max value is 50")
	}
	if q.Page < 1 {
		err = err.Set("page", "min value is 1")

	}

	if len(err.E) == 0 {
		return nil
	}

	return err
}
