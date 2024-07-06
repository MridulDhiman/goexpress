package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type ContextKey string
type fnSign func()

const (
	QueryParamsContext ContextKey = "query-params"
	PathParamsContext  ContextKey = "path-params"
)

type QueryParams struct {
	key   string
	value string
}

type Params struct {
	key   string
	value string
}

type Router struct {
	routes      map[string]http.HandlerFunc
	queryParams map[string][]*QueryParams
	patterns    map[string][]string
	params      map[string]*Params
	middlewares map[string][]*http.HandlerFunc
}

type RouteGroup struct {
	router *Router
	prefix string
}

type App struct {
	newAddr string
	router  *Router
}

func (a *App) NewApp() *App {
	router := a.Router()
	return &App{
		router: router,
	}
}

func (a *App) Router() *Router {
	return &Router{
		routes:      make(map[string]http.HandlerFunc),
		queryParams: make(map[string][]*QueryParams),
		params:      make(map[string]*Params),
		patterns:    make(map[string][]string),
		middlewares: make(map[string][]*http.HandlerFunc),
	}
}

// http.HandleFunc interface has ServeHTTP method
// So, we can reverse engineer in a way like, I want to forward my ResponseWriter and request to Handler Func
// So, I can use req. methods and paths to map to appropriate handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find the handler from the routes
	// handler (or we can call it, controller in express.js context)

	params, paramsOk := r.FindParams(req)

	if handler, queryParams := r.FindHandlerAndParams(req); handler != nil {
		// add context here

		if queryParams != nil || paramsOk {
			ctx := req.Context()
			ctx = r.SetContext(req, QueryParamsContext, queryParams)
			ctx = r.SetContext(req, PathParamsContext, params)
			handler(w, req.WithContext(ctx))
		} else {
			handler(w, req)
		}
	} else {
		w.WriteHeader(404)
		temp := []string{req.URL.Path, req.Method, "Not Found"}
		fmt.Fprintln(w, strings.Join(temp, " "))
	}
}

func (r *Router) FindParams(req *http.Request) ([]*Params, bool) {
	isDynamicRoute := r.isDynamicRoute(req.URL.Path)
	kvStore := make([]*Params, 1)
	if isDynamicRoute {
		pattern := r.MakePattern(req.URL.Path)
		for _, oldPattern := range r.patterns {
			params, ok := r.MatchPattern(pattern, oldPattern)

			if ok {
				return params, true
			} else {
				return nil, false
			}
		}
	}

	return kvStore, false
}

func (r *Router) FindHandlerAndParams(req *http.Request) (http.HandlerFunc, []*QueryParams) {

	mapKey := []string{req.Method, req.URL.Path}
	route := strings.Join(mapKey, ":")
	fmt.Println("Route key", route, "findHandler()")
	handler, ok := r.routes[route]
	queryParams, queryOk := r.queryParams[route]
	if ok || queryOk {
		if ok && !queryOk {
			return handler, nil
		}
		if !ok && queryOk {
			return nil, queryParams
		}
		if ok && queryOk {
			return handler, queryParams
		}

	}
	fmt.Println("could not find handler: findHandler()")

	return nil, nil
}

func (r *Router) GetContext(req *http.Request, key any) any {
	value := req.Context().Value(key)
	return value
}

func (r *Router) SetContext(req *http.Request, key any, value any) context.Context {
	return context.WithValue(req.Context(), key, value)
}

// forward all the HTTP Request to router.Handle Method
func (a *App) Handle(method string, route string, handlers ...http.HandlerFunc) {
	// it will http request methods and path against the handlers
	a.router.Handle(method, route, handlers...)
}

func (r *RouteGroup) Get(route string, handler ...http.HandlerFunc) {
	r.router.Handle("GET", path.Join(r.prefix, route), handler...)
}

func (r *RouteGroup) Post(route string, handler ...http.HandlerFunc) {
	r.router.Handle("POST", path.Join(r.prefix, route), handler...)
}

func (r *RouteGroup) Put(route string, handler ...http.HandlerFunc) {
	r.router.Handle("PUT", path.Join(r.prefix, route), handler...)
}
func (r *RouteGroup) Delete(route string, handler ...http.HandlerFunc) {
	r.router.Handle("DELETE", path.Join(r.prefix, route), handler...)
}

func (a *App) NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		router: a.router,
		prefix: prefix,
	}
}

func (r *Router) SetQueryParams(searchKey string, queryParams url.Values) {
	r.queryParams[searchKey] = make([]*QueryParams, 1)
	for key, value := range queryParams {
		r.queryParams[searchKey] = append(r.queryParams[searchKey], &QueryParams{
			key:   key,
			value: strings.Join(value, ""),
		})
	}
}

func (r *Router) MatchPattern(newPattern []string, oldPattern []string) ([]*Params, bool) {
	patternLength := len(newPattern)
	dynamicRouteIndices := []int{}
	for i := 0; i < patternLength; i++ {
		if r.isDynamicRoute(oldPattern[i]) {
			dynamicRouteIndices = append(dynamicRouteIndices, i)
		} else if oldPattern[i] != newPattern[i] {
			// diff. route
			return nil, false
		}
	}
	kvStore := make([]*Params, 1)

	for _, ind := range dynamicRouteIndices {
		key := strings.TrimPrefix(oldPattern[ind], ":")
		value := newPattern[ind]
		kvStore = append(kvStore, &Params{
			key:   key,
			value: value,
		})
	}
	return kvStore, true
}

func (r *Router) MakePattern(path string) []string {
	if !strings.HasSuffix(path, "/") {
		path = strings.Join([]string{path, "/"}, "")
	}
	pattern := `/[a-zA-Z0-9]+/`
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Could not compile the string pattern to *regexp.RegExp", "Router.MakePattern() method")
	}
	patterns := re.FindAllString(path, -1)
	fmt.Println("Patterns: ", patterns)

	finalPatterns:= make([]string, 1)
	for _, pattern := range patterns {
		finalPatterns = append(finalPatterns, strings.Trim(pattern, "/"))
	}
	return finalPatterns
}

func (r *Router) isDynamicRoute(path string) bool {
	pattern := `:[a-zA-Z0=9]+`
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Could not compile the string pattern to *regexp.RegExp", "Router.isDynamicRoute() method")
	}
	return re.MatchString(path)
}

func (r *Router) Handle(method string, route string, handler ...http.HandlerFunc) {
	mapKey := []string{method, route}
	x, err := url.Parse(route)
	if err != nil {
		fmt.Println("error in parsing url: ", err)
	}
	r.SetQueryParams(strings.Join(mapKey, ":"), x.Query())
	if r.isDynamicRoute(route) {
		r.patterns[route] = r.MakePattern(route)
	}
	r.routes[strings.Join(mapKey, ":")] = handler[len(handler)-1]
}

func (a *App) Get(path string, handler ...http.HandlerFunc) {
	a.Handle("GET", path, handler...)
}

func (a *App) Post(path string, handler ...http.HandlerFunc) {
	a.Handle("POST", path, handler...)
}

func (a *App) Put(path string, handler ...http.HandlerFunc) {
	a.Handle("PUT", path, handler...)
}

func (a *App) Delete(path string, handler ...http.HandlerFunc) {
	a.Handle("DELETE", path, handler...)
}

func (a *App) Listen(port string, Callback fnSign) {
	newAddr := strings.Join([]string{":", port}, "")
	a.newAddr = newAddr
	Callback()
	// ListenAndServe will be blocking the main thread, and will only release in case of error, and server shut down.
	err := http.ListenAndServe(a.newAddr, a.router)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	express := new(App)
	app := express.NewApp()
	PORT := "3000"

	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// writing to response
		fmt.Fprintf(w, "Hello")
	})

	app.Get("/:id", func(w http.ResponseWriter, r *http.Request) {

	})

	app.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// sending favicon.ico
		w.WriteHeader(200)
	})

	helloRouter := app.NewRouteGroup("/user")
	helloRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from /user route")
	})

	app.Listen(PORT, func() { fmt.Printf("Server is listening on port %s\n", PORT) })
}
