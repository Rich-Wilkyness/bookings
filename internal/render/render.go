package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/Rich-Wilkyness/bookings/internal/config"
	"github.com/Rich-Wilkyness/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// ---------------------------------------------------- Advanced Render -------------------------------------------
// Advantage - no longer have to keep track of how many files are in the template folder
// Advantage - how many of those files are using page.tmpl vs layout.tmpl

var functions = template.FuncMap{}

// creates a global variable to access our config / cache
var app *config.AppConfig
var pathToTemplates = "./templates"

// sets config for the render package to have access
func NewTemplates(a *config.AppConfig) {
	app = a
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// this function will allow us to get data we want on every single page (things like a session)
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash") // PopString will add something to our session until a new page is hit and then it will be taken out of our session automatically
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// td is new - it is any data we are going to send to our template - see TemplateData
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// create template cache - instead to improve our caching - where we do not create a new cache everytime
	// we are going make our cache on main and call the rendering from that cache here
	// we made a "global" cache in our config package

	var tc map[string]*template.Template

	// determine if we are going to use our cache or create a new cache
	if app.UseCache {
		// need to now get our cache from our AppConfig
		// create new func NewTemplates()
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl] // t = the template, ok will be true if t exists
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer) // this is for finer error handling

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td) // this tells us that the error comes from the value stored in the map (helps us isolate where the problem is)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// we could create a map this way, but we'll do it another way
	// myCache := make(map[string]*template.Template)

	myCache := map[string]*template.Template{} // this is the same functionality as the make above, we just made an empty map
	// need to cache everything, when rendering the first thing you need to parse is the template(s), then the layout(s)

	// get all of the files named *.page.tmpl from ./templates
	// Glob is used to return a pattern of files that match the pattern
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// the code snippet is iterating over a collection of page filenames (pages), extracting the base filename for each page,
	// creating a new template with that filename, and then parsing the contents of each page file to associate them with their respective templates.
	// range through all files ending with *.page.tmpl
	for _, page := range pages { // remember the first variable in a for loop is the index which we leave blank "_"
		// page = the full path to the file
		filename := filepath.Base(page) // base returns the last element of the path. Here that will be the name of the file ending in ".page.tmpl"
		// if page is "path/to/file/page.tmpl", filename will be "page.tmpl".
		// ts is template set
		ts, err := template.New(filename).ParseFiles(page) // first we parse the page or get the content of the page, and then associate it with the new template
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			// remember parse means to associate the content. So we are associating the layout(s) to our template set(s) (ts) in a for loop
			// so each page is associated with all of the layout(s) if any
			if err != nil {
				return myCache, err
			}
		}

		myCache[filename] = ts
	}
	return myCache, nil
}
