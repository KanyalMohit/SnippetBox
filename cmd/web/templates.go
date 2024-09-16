package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.mohit.net/internal/models"
)

// templateData struct to hold data passed to templates
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// creating a function which will return a nicely formatted string
func humanDate(t time.Time) string {

	if t.IsZero() {
		return ""
	}
	
	return t.Format("02 Jan 2006 at 15:04")
}

/*
Initialize a template.FuncMap object abd storing it in a global variable. This is
essentially a string-keyed map which points to our custom template function
*/
var function = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	//initialize new map to act as cache
	cache := map[string]*template.Template{}

	/*using filePath.Glob() function to get a slice of all filepath that
	  matches the pattern
	*/
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extravting the file name like home.html and assigning it to a variable
		name := filepath.Base(page)

		/* files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		} */
		//parsing the files into a template set..

		/*we using template.New() to create an empty template set,using the funcs() method to
		register the template.FuncMap and then parse the file as normal
		*/
		ts, err := template.New(name).Funcs(function).ParseFiles("./ui/html/base.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
