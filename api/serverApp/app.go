package serverApp

import (
	"net/http"
)

// routes: is a map of route to method to HandlerFunc
type App struct {
	Opts     *AppOpts
	handlers map[string]map[string]appHandler
}

type AppOpts struct {
	PrefixRoute  string
	ErrorHandler IErrorHandler
}

// Constructor of app
func NewApp(opts *AppOpts) *App {
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = &DefaultErrorHandler{}
	}
	app := &App{
		Opts: opts,
	}
	initOpts(app)
	// avoid nil map when initialize the app
	app.handlers = make(map[string]map[string]appHandler)
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
func (a *App) findHandler(path string, method string) (appHandler, bool, bool) {
	_, pathExists := a.handlers[path]
	handler, methodExists := a.handlers[path][method]
	return handler, pathExists, methodExists
}

// Set all the handlers of the application
func (a *App) SetHandler(path string, method string, handlers ...IHandler) {
	_, exists := a.handlers[path]
	if !exists {
		a.handlers[path] = make(map[string]appHandler)
	}
	appHandlers := handlersToAppHandlers(handlers...)
	handler := chainHandlers(appHandlers...)
	a.handlers[path][method] = handler
}

// Method to simulate an inheritance in go, to include the Handle method defined for the user
// And chained for the framework
func handlersToAppHandlers(handlers ...IHandler) []appHandler {
	appHandlers := []appHandler{}
	for _, handler := range handlers {
		appHandlerInstance := appHandler{
			Handler: handler,
		}
		appHandlers = append(appHandlers, appHandlerInstance)
	}
	return appHandlers
}

// Auxiliar function to connect all the handlers in a chain
func chainHandlers(handlers ...appHandler) appHandler {
	handler := handlers[0]
	for idx := 1; idx < len(handlers); idx++ {
		handlers[idx-1].SetNext(&handlers[idx])
	}
	return handler
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
	response, errResponse := handler.Handle(res, req)
	if errResponse != nil {
		a.Opts.ErrorHandler.Handle(errResponse, res, req)
		return
	}
	res.WriteHeader(response.StatusCode)
	jsonResponse(res, response.Body)
}
