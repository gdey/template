// +build !debug

package template

import "io"

// addSourceFile will add the source of the file.
func addSourceFile(t *Template, typ sourceType, file string) {}

// genParseFileList will go through the data-structure and generate the file list to parse.
func genParseFileList(_ *Template) error { return nil }

// This will parse the files that have been build up
func (t *Template) ParseFiles() (*Template, error) {
	_, err := t.Template.ParseFiles(t.parseFiles...)
	return t, err
}

// Execute will reparse all the template, then execute the template with the given data.
func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.Template.Execute(w, data)
}
