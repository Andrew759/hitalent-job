package base

import (
	"hitalent/cmd/service"
	"net/http"
)

type DIContainer struct {
	DBDecorator service.DBDecorator
}

type Controller struct {
	ServeMux     *http.ServeMux
	Dependencies DIContainer
}

type RequestHandler interface {
	HandleRequest()
}
