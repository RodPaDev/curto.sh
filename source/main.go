// so this logic is supposed to be on lambdas but haven't figured out how to do it yet
// so for now while I prototype locally I'll use echo

package main

import (
	"fmt"
	"html/template"
	"io"
	"regexp"

	"url-shortener/source/db"
	"url-shortener/source/nanoid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var dbInstance db.DB = *db.NewDB()

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type ShortenedUrl struct {
	Short string
	Long  string
	Title string
}

func newShortenedUrl() ShortenedUrl {
	return ShortenedUrl{
		Short: "",
		Long:  "",
	}
}

type State struct {
	Form     FormData
	Data     ShortenedUrl
	Title    string
	NotFound bool
}

func newState() State {
	return State{
		Form:     newFormData(),
		Data:     newShortenedUrl(),
		Title:    "curto.sh",
		NotFound: false,
	}
}

func main() {
	// new echo instance
	e := echo.New()
	// use logger middleware
	e.Use(middleware.Logger())
	// register template renderer
	e.Renderer = newTemplate()
	// serve static files
	e.Static("/css", "css")
	e.Static("/assets", "assets")
	e.Static("/fonts", "fonts")

	// state
	s := newState()

	// routes
	e.GET("/", func(c echo.Context) error {
		s = newState()
		return c.Render(200, "index", s)
	})

	e.GET("/another-one", func(c echo.Context) error {
		s = newState()
		return c.Render(200, "url-shortener", s.Form)
	})

	e.POST("/shorten", func(c echo.Context) error {
		url := c.FormValue("url")

		// validate url
		regexp := regexp.MustCompile(`https?:\/\/(?:www\.)?[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.[^\s]{2,}`)
		if !regexp.MatchString(url) {
			s.Form.Errors["url"] = "Please enter a valid URL"
			c.Response().Header().Set("HX-Retarget", "#error")
			return c.Render(422, "error", s.Form)
		}
		// formData := newFormData()
		short := nanoid.Gen(6)
		dbInstance.Set(short, url)

		// if  c.Request().Host is local than use http else use https
		shortUrl := c.Request().Host
		if shortUrl == "localhost:3033" || shortUrl == "127.0.0.1:3033" {
			shortUrl = fmt.Sprintf("http://%s", shortUrl)
		} else {
			shortUrl = fmt.Sprintf("https://%s", shortUrl)
		}

		s.Data.Short = fmt.Sprintf("%s/%s", shortUrl, short)
		s.Data.Long = url

		return c.Render(200, "shortened-url", s.Data)
	})

	e.GET("/:url", func(c echo.Context) error {
		short := c.Param("url")
		fmt.Println(short)
		url, ok := dbInstance.Get(short)
		if !ok {
			s.NotFound = true
			return c.Render(200, "index", s)
		}
		return c.Redirect(302, url)
	})

	e.GET("/surprise", func(c echo.Context) error {
		return c.Redirect(302, "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
