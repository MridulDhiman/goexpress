package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"strings"
)

// type RouterConfig struct {
// 	method string
// 	path string
// 	routeGrp string
// }

type RouteGroup struct {
router *Router
prefix string
}

type Router struct {
	routes   map[string]http.HandlerFunc
}

type App struct {
	newAddr string
	router  *Router
}

type fnSign func()

// type apiFunc func (w http.ResponseWriter, r* http.Request) error

func (a *App) NewApp() *App {
	router:= a.Router()
	return &App{
		router: router,
	}
}


func (a *App) Router() *Router {
	return &Router{
		routes:   make(map[string]http.HandlerFunc),
	}
}

// http.HandleFunc interface has ServeHTTP method
// So, we can reverse engineer in a way like, I want to forward my ResponseWriter and request to Handler Func
// So, I can use req. methods and paths to map to appropriate handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find the handler from the routes
	// handler (or we can call it, controller in express.js context)
	if handler := r.findHandler(req); handler != nil {
		handler(w, req)
	} else {
		w.WriteHeader(404)
		temp := []string{req.URL.Path, req.Method, "Not Found"}
		fmt.Fprintln(w, strings.Join(temp, " "))
	}
}

func (r *Router) findHandler(req *http.Request) http.HandlerFunc {
	mapKey := []string{req.Method, req.URL.Path}
	route := strings.Join(mapKey, ":")
	fmt.Println("Route key", route, "findHandler()")
	if handler, ok := r.routes[route]; ok {
		return handler
	}
	fmt.Println("could not find handler: findHandler()")
	return nil
}

// forward all the HTTP Request to router.Handle Method
func (r* Router) Handle(method string, route string, handler http.HandlerFunc) {
	// it will http request methods and path against the handlers
	mapKey := []string{method, route}
	r.routes[strings.Join(mapKey, ":")] = handler
}

func (r *Router) Get(route string, handler http.HandlerFunc) {
	r.Handle("GET", path.Join(r.routeGrp, route), handler)
}

func (r *Router) Post(route string, handler http.HandlerFunc) {
	r.Handle("POST", path.Join(r.routeGrp, route), handler)
}

func (r *Router) Put(route string, handler http.HandlerFunc) {
	r.Handle("PUT", path.Join(r.routeGrp, route), handler)
}

func (r *Router) Delete(route string, handler http.HandlerFunc) {
	r.Handle("DELETE", path.Join(r.routeGrp, route), handler)
}

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.router.Get(path, handler)
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.router.Post(path, handler)
}

func (a *App) Put(path string, handler http.HandlerFunc) {
	a.router.Put(path, handler)
}

func (a *App) Delete(path string, handler http.HandlerFunc) {
	a.router.Delete(path, handler)
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

func (a *App) Group(pattern string, router *Router) *Router {
	return &Router {

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

	app.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// sending favicon.ico
		w.WriteHeader(200)
	})

	helloRouter := app.Group("/user")
	helloRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from /user route")
	});

	app.Listen(PORT, func() { fmt.Printf("Server is listening on port %s\n", PORT) })
}
