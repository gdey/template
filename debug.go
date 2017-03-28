// +build debug

package template

import (
	"html/template"
	"io"
	"log"

	"github.com/gdey/template/helpers"
)

func init() {
	log.Println("Template running in debug mode. Templates will be reloaded.")
	helpers.ReloadAlways = true
}

// addSourceFile will add the source of the file.
func addSourceFile(t *Template, typ sourceType, file string) {
	t.parseFilesSources = append(t.parseFilesSources, parseFileSources{
		Type: typ,
		File: file,
	})
}

// genParseFileList will go through the data-structure and generate the file list to parse.
func genParseFileList(t *Template) error {
	t.fullLock()
	defer t.fullUnlock()
	// Clear whatever is in parsefiles.
	t.parseFiles = nil
	for _, filesrc := range t.parseFilesSources {
		switch filesrc.Type {
		case SrcGlobFile:
			if err := parseGlob(t, filesrc.File); err != nil {
				return err
			}
		case SrcFileList:
			if err := parseFileList(t, filesrc.File); err != nil {
				return err
			}
		case SrcParseFile:
			t.parseFiles = append(t.parseFiles, filesrc.File)
		}
	}
	return nil
}

// This will parse the files that have been build up
func (t *Template) ParseFiles() (*Template, error) {

	if err := genParseFileList(t); err != nil {
		return t, err
	}
	_, err := t.Template.ParseFiles(t.parseFiles...)

	return t, err
}

// Execute will reparse all the template, then execute the template with the given data.
func (t *Template) Execute(w io.Writer, data interface{}) error {

	t.fullLock()
	t.Template = template.New(t.name)
	t.Template.Funcs(t.helpers)
	t.fullUnlock()

	if _, err := t.ParseFiles(); err != nil {
		return err
	}
	err := t.Template.Execute(w, data)

	return err
}
