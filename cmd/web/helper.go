package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form"
	"github.com/justinas/nosurf"
)

// serverError helper writer an error message and stack trace to the errorLog.
// then sends a generic 500 internal server error respose to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// client helper sends a specific status code and corresponding description
// responser like 400 "Bad Request" when a problem occurs
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper. This is simply a convenience wrapper around clientError
// which sends a 404 Not Found response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exists", page)
		app.serverError(w, err)
		return
	}
	// initialzie a new buffer
	buf := new(bytes.Buffer)

	//writing template to the buffer if there is an error then we will now of it
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//status code like 200 ok , 400 bad request
	w.WriteHeader(status)

	/* err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	} */

	//write the contents of the buffer to the http.ResponseWritter..
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		//if we try to use an invalid target destination, the decode() method will return an error with
		//the type *form.invalidDecideError.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		//for all other errors, we use them as normal
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAUthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAUthenticated
}
