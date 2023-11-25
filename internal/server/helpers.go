package server

import (
	"FileServer/internal/logger"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

// Matches (). with any digit in the middle including the .
// filename(̲3)̲.̲png including the .
var rx = regexp.MustCompile(`\(\d\).`)

// Returns a unique name acounting for duplicate files
// If you try to upload a file that already exists on disk
// folder/File.txt  --->  folder/File(2).txt
// And if 3 copies exist on disk
// folder/File.txt  --->  folder/File(4).txt
// Will remove old duplication marks and update them
// folder/File(3).txt  --->  folder/File(4).txt
func getUniqueName(path string) (string, error) {
	// rawFileName := filepath.Base(path)
	// dirPath := filepath.Dir(path)
	dirPath, rawFileName := filepath.Split(path)
	// remove copy marks from for aesthetics
	// duplicateFile(2).txt  --->  duplicateFile.txt
	// du(1)ff(2).txt  --->   du(1)ff.txt
	// We then check duplicates ourselves and give it the proper number
	// Note that this only works on that particural style of marks
	// Which is fine because this doesn't actually need to happen
	// duplicate( copy 2).txt  --->  duplicate( copy 2)(2).txt
	// That is ugly but it works, and very unlikely to stumble on
	// this style.
	nameNoRegEx := string(rx.ReplaceAll([]byte(rawFileName), []byte(".")))

	fileName := filepath.Join(dirPath, nameNoRegEx)
	fileName = fmt.Sprintf(".%s", fileName)
	ext := filepath.Ext(fileName)
	// We remove the extention so we can add it later with
	// a * right before it   file.txt  -->  file*.txt
	// So we can search for duplicates with different numbers
	fileNoExt := strings.TrimSuffix(fileName, ext)
	globPattern := fmt.Sprintf("%s*%s", fileNoExt, ext)
	files, err := filepath.Glob(globPattern)
	if err != nil {
		logger.Log.Infof("path: %s", path)
		logger.Log.Infof("nameNoRegex: %s", nameNoRegEx)
		log.Infof("fileName: %s", fileName)
		logger.Log.Error(err)
		return "", err
	}

	var matches int

	for _, file := range files {
		tmp := string(rx.ReplaceAll([]byte(file), []byte(".")))
		if filepath.Base(tmp) == nameNoRegEx {
			matches++
		}

	}
	verNum := ""

	newName := fmt.Sprintf("%s%s%s", fileNoExt, verNum, ext)
	return newName, nil
}
