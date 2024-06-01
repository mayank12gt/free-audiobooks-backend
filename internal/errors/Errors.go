package Error

import (
	"log"
)

type Err struct {
	status_code int               `json:"-"`
	E           map[string]string `json:"errors"`
}

func (e Err) Error() string {
	var combinedErr string
	for _, value := range e.E {
		combinedErr = combinedErr + value
	}
	return combinedErr
}

func (e *Err) Set(key string, err string) *Err {
	if e == nil {
		e = NewError()
	}
	log.Print(e.Error())
	e.E[key] = e.E[key] + " " + err
	log.Print(e.Error())
	return e
}

func (e *Err) SetCode(code int) *Err {
	if e == nil {
		e = NewError()
	}
	e.status_code = code
	return e
}

func NewError() *Err {
	//log.Print(msg)
	return &Err{
		E: make(map[string]string),
	}
}
