package server

import (
	"FileServer/ui"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *Server) RegisterRoutes() http.Handler {
	fileSystem := http.FS(ui.Files)
	t := &Template{
		templates: template.Must(template.ParseFS(ui.Files, "html/*.html")),
	}
	assetHandler := http.FileServer(fileSystem)
	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Recover())

	e.Use(middleware.Static("."))
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/*", s.fileExploreHandler)
	e.POST("/*", s.uploadHandler)
	e.GET("/static/*", echo.WrapHandler(assetHandler))
	e.DELETE("/delete/*", s.deleteHandler)
	e.PATCH("/update/*", s.renameHandler)
	e.GET("/update/*", s.renameGetHandler)

	// Used for when i needed to refresh the page after sending an htmx request
	e.GET("/refresh", echo.HandlerFunc(func(c echo.Context) error {
		c.Response().Header().Add("HX-Refresh", "true")
		return nil
	}))

	return e
}
