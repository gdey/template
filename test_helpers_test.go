package template_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"fmt"
	"runtime"

	"io/ioutil"

	"github.com/gdey/template"
)

type FileType struct {
	Filename string
	Content  string
}

func (fl *FileType) CreateFile(base string) error {

	filename := filepath.Join(base, fl.Filename)
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(fl.Content); err != nil {
		return err
	}
	return nil
}

type FileList struct {
	BaseDir string
	Files   []FileType
}

func (fl *FileList) CreateFilesOFail(t *testing.T) *FileList {
	os.MkdirAll(fl.BaseDir, os.ModePerm)
	for _, f := range fl.Files {
		if err := f.CreateFile(fl.BaseDir); err != nil {
			fl.RemoveAll()
			t.Fatalf("Failed to create file (%v) : %v ", f.Filename, err)
		}
	}
	return fl
}

func (fl *FileList) CreateFileOrFail(t *testing.T, filename string) *FileList {
	for _, f := range fl.Files {
		if f.Filename == filename {
			if err := f.CreateFile(fl.BaseDir); err != nil {
				t.Fatalf("Failed to create file (%v) : %v ", f.Filename, err)
			}
			return fl
		}
	}
	t.Fatalf("Failed to create file (%v) : was not able to find it my list of files. ", filename)
	return fl
}

func (fl *FileList) RemoveAll() {
	for _, f := range fl.Files {
		os.RemoveAll(filepath.Join(fl.BaseDir, f.Filename))
	}
}

func (fl *FileList) SetFile(filename, content string) *FileList {
	for i, f := range fl.Files {
		if f.Filename == filename {
			fl.Files[i].Content = content
			return fl
		}
	}
	// We did not find the file. So, add it.
	fl.Files = append(fl.Files, FileType{filename, content})
	return fl
}

func (fl *FileList) CatFile(filename string) string {
	fname := filepath.Join(fl.BaseDir, filename)
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// MyCallerFileLine returns the FileLine of the caller of the function that called it :)
func MyCallerFileLine() string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	// return its name
	filename, line := fun.FileLine(fpcs[0] - 1)
	filename = filepath.Base(filename)
	return fmt.Sprintf("%v:%v", filename, line)
}

func ExecuteTemplateOrFail(t *testing.T, tpl *template.Template, data interface{}, expected string) string {
	caller := MyCallerFileLine()
	b := bytes.NewBufferString("")
	if err := tpl.Execute(b, data); err != nil {
		t.Fatal(caller, ":Got error Executing template:", err)
	}
	got := b.String()
	if got != expected {

		t.Fatalf("%v:expected: “%v”  got:“%v”", caller, expected, got)
	}
	return got
}
