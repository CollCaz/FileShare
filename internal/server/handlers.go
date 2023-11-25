package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Bios-Marcel/wastebasket"
	"github.com/labstack/echo/v4"
)

func (s *Server) uploadHandler(c echo.Context) error {
	c.Response().Header().Add("HX-Refresh", "true")
	file, err := c.FormFile("file")
	if err != nil {
		s.logger.Error(err)
		return err
	}

	safeURL, err := url.PathUnescape(fmt.Sprint(c.Request().URL))
	if err != nil {
		s.logger.Error(err)
		return err
	}
	fileName := file.Filename
	path := filepath.Join(safeURL, fileName)
	// path = fmt.Sprintf(".%s", path)
	newName, err := getUniqueName(path)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(newName)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s *Server) deleteHandler(c echo.Context) error {
	c.Response().Header().Add("HX-Refresh", "true")
	// The static middleware for some god damn reason
	// gets applied to ALL routes, not just ones
	// declared after it, so sending a delete request
	// to /aa.png will serve aa.png instead of running
	// this handler, so i prefix delete requests with
	// /delete to work around that issue, then remove it.
	file := c.ParamValues()[0]
	file = strings.TrimPrefix(file, "delete/")
	file = strings.TrimSuffix(file, "-delete")
	file = fmt.Sprintf(".%s", file)
	err := wastebasket.Trash(file)
	if err != nil {
		s.logger.Info(c.Request().URL)
		s.logger.Error(err)
		return err
	}
	return nil
}

// Responsible for showing the user the contents of the current directory
func (s *Server) fileExploreHandler(c echo.Context) error {
	params := c.ParamValues()

	// Finds the "current diectory" by checking the url Params
	// If the first param (the url path) is empty, show the
	// root directory "."
	// fileshare.com/  --> root directory
	// fileshare.com/folder1/subfolder2  --> ./folder1/subfolder2
	var path string
	var err error
	if len(params[0]) > 0 {
		// Unescape the url, Arabic letters and other
		// non standard characters will break everything
		// if unescaped
		path, err = url.PathUnescape(params[0])
		if err != nil {
			s.logger.Info(c.Request().URL)
			s.logger.Error(err)
			return err
		}
	} else {
		path = "."
	}

	var pathForCumbs string
	if path != "." {
		pathForCumbs = path
	}

	crumbsData := make(map[int][]string)
	crumbName := strings.Split(pathForCumbs, "/")
	if path == "." {
		crumbName[0] = "/"
	}
	for i, crumb := range crumbName {
		crumbPath := strings.Join(crumbName[0:i+1], "/")
		crumbsData[i] = []string{crumb, crumbPath}
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		s.logger.Info(c.Request().URL)
		s.logger.Error("PANIC: ", err)
		panic("could not read current directory")
	}

	// index[directory_path, directory_name]
	dirMap := make(map[int][]string)
	// index[file_path, file_name, ID]
	fileMap := make(map[int][]string)

	for i, item := range dirs {
		filesToExclude := []string{"/.log.txt", "/favicon.ico"}
		// dir/file.ext will take you to a url
		// /dir/file.ext will download the file
		// Becaue we do params[0] + "/"
		// if we have no params it will become a leading
		// slash, and that is bad
		// href="/dir" will redirect to http://dir
		// href="dir" will do what we want
		pre := ""
		if len(params[0]) > 0 {
			pre = "/"
		}

		// This feels more readable than Sprintf(%d%d%d%d)
		fileName := pre + params[0] + "/" + item.Name()
		if slices.Contains(filesToExclude, fileName) {
			continue
		}
		if item.IsDir() {
			dirMap[len(dirMap)] = []string{fileName, item.Name()}
		} else {
			fileMap[len(fileMap)] = []string{fileName, item.Name(), fmt.Sprintf("id%d", i)}
		}
	}

	Data := struct {
		Files  map[int][]string
		Dirs   map[int][]string
		Crumbs map[int][]string
		Pls    map[int][]string
	}{
		Files:  fileMap,
		Dirs:   dirMap,
		Crumbs: crumbsData,
	}

	err = c.Render(http.StatusOK, "base", &Data)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s *Server) renameHandler(c echo.Context) error {
	c.Response().Header().Add("HX-Refresh", "true")
	// The static middleware for some god damn reason
	// gets applied to ALL routes, not just ones
	// declared after it, so sending a patch request
	// to /aa.png will serve aa.png instead of running
	// this handler, so i prefix delete requests with
	// /update AND postfix with -update to work around that issue, then remove it.
	// hxUrl := c.Request().Header.Get("HX-Current-URL")
	safeURL, err := url.PathUnescape(fmt.Sprint(c.Request().URL))
	s.logger.Info(safeURL)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	path := fmt.Sprint(safeURL)
	path = strings.TrimPrefix(path, "/update")
	path = strings.TrimSuffix(path, "-update")
	s.logger.Info(path)
	ext := filepath.Ext(path)

	newName := c.FormValue("new-name")
	newName = strings.ReplaceAll(path, filepath.Base(path), newName)
	newName = newName + ext
	newName, err = getUniqueName(newName)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	// Make the path relative
	path = fmt.Sprintf(".%s", path)
	err = os.Rename(path, newName)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

// Responsible for serving the html input form when renaming a file
func (s *Server) renameGetHandler(c echo.Context) error {
	file := c.ParamValues()[0]
	filePath, err := url.PathUnescape(file)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	filePath = strings.TrimSuffix(filePath, "-update")
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	fileNoExt := strings.TrimSuffix(fileName, ext)
	data := struct {
		Path string
		Name string
		Ext  string
	}{
		Path: filePath,
		Name: fileNoExt,
		Ext:  ext,
	}
	err = c.Render(http.StatusSeeOther, "nameForm", data)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}
