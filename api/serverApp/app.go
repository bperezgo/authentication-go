package serverApp

import (
	"context"
	"net/http"
)

// routes: is a map of route to method to HandlerFunc
type App struct {
	Opts     *AppOpts
	handlers map[string]map[string]IHandler
}

type AppOpts struct {
	PrefixRoute string
}

// Constructor of app
func NewApp(opts *AppOpts) *App {
	app := &App{
		Opts: opts,
	}
	initOpts(app)
	// avoid nil map when initialize the app
	app.handlers = make(map[string]map[string]IHandler)
	return app
}

// Initialize all the options of the app
func initOpts(app *App) {
	if app.Opts == nil {
		return
	}
	// TODO: Verify why does not work the prefix
	http.Handle(app.Opts.PrefixRoute, app)
}

// Return The HandlerFunc, if path exists and if method exists
func (a *App) findHandler(path string, method string) (IHandler, bool, bool) {
	_, pathExists := a.handlers[path]
	handler, methodExists := a.handlers[path][method]
	return handler, pathExists, methodExists
}

// Set all the handlers of the application
func (a *App) SetHandler(path string, method string, handler IHandler) {
	_, exists := a.handlers[path]
	if !exists {
		a.handlers[path] = make(map[string]IHandler)
	}
	a.handlers[path][method] = handler
}

// Method to implement http.Handler
func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler, pathExists, methodExists := a.findHandler(req.URL.Path, req.Method)
	if !pathExists {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if !methodExists {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := context.Background()
	response, errResponse := handler.Handle(ctx, req)
	if errResponse != nil {
		res.WriteHeader(errResponse.StatusCode)
		jsonResponse(res, errResponse.Body)
		return
	}
	res.WriteHeader(response.StatusCode)
	jsonResponse(res, response.Body)
}
