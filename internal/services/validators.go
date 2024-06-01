package services

import (
	Error "github.com/mayank12gt/free-audiobooks-backend/internal/errors"
)

func (q *Query) Validate() error {

	err := Error.NewError()
	if q.PageSize > 50 || q.PageSize < 1 {
		err = err.Set("page_size", "max value is 50 and min value is 1")
	}

	if q.Page < 1 {
		err = err.Set("page", "min value is 1")
	}

	if q.TotalTimeRange.TotalTimeMin == 0 && q.TotalTimeRange.TotalTimeMax != 0 || q.TotalTimeRange.TotalTimeMin != 0 && q.TotalTimeRange.TotalTimeMax == 0 {
		err.Set("length", "Both lengthMin and lengthMax should be set")
	}
	if q.TotalTimeRange.TotalTimeMax < 0 || q.TotalTimeRange.TotalTimeMin < 0 {
		err.Set("length", "lengthMin and lengthMax should be positive")
	}

	if q.TotalTimeRange.TotalTimeMax != 0 && q.TotalTimeRange.TotalTimeMin != 0 {
		if q.TotalTimeRange.TotalTimeMax <= q.TotalTimeRange.TotalTimeMin {
			err.Set("length", "lengthMax must be > lengthMin")
		}
	}

	if len(err.E) == 0 {
		return nil
	}

	return err
}
