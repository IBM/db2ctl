package route

import (
	"net/http"

	"github.com/IBM/db2ctl/pkg/rest/handler"
	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

//NewRouter creates a new mux router for application
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	subrouter.Use(loggingMiddleware)
	subrouter.Use(panicHandlerMiddleware)
	for _, route := range routesForApp {
		subrouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return subrouter
}

var routesForApp = routes{
	route{
		"Index Page",
		"GET",
		"/",
		handler.Index,
	},
	route{
		"Web socket",
		"GET",
		"/ws",
		handler.WebSocket,
	},
	route{
		"Configuration Toggle",
		"POST",
		"/config",
		handler.Configuration,
	},
	route{
		"Install module",
		"POST",
		"/install",
		handler.InstallModule,
	},
	route{
		"Get state",
		"GET",
		"/state",
		handler.GetState,
	},
	route{
		"Cancel execution",
		"PUT",
		"/cancel",
		handler.Cancel,
	},
}
