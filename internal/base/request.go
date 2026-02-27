package base

import (
	"fmt"
	"net/http"
	"strconv"
)

type Request struct {
	*http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
	}
}

func (r *Request) HTTPId() (id int, err error) {
	id, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid URL")
	}
	return
}
