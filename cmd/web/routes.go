package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//mux := http.NewServeMux()
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// This line of code is setting up a route for serving static files. It is using the `router.Handler`
	// method to handle HTTP GET requests to the path "/static/*filepath". The `http.StripPrefix` function
	// is used to strip the "/static" prefix from the request URL path before serving the static files
	// from the specified directory using `http.FileServer`. This allows the application to serve static
	// files like CSS, JavaScript, images, etc., located in the "./ui/static" directory.
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	/* mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate) */

	//creating a new middleware chain containing the middleware specific to our dynamic application route
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))

	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
