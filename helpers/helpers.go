package helpers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"mime"
)

// Minifier is a type used to minifiy the source.
type Minifier interface {
	Minify(mimetype string, w io.Writer, r io.Reader) error
}

// Turn this on to disregard the file existence check.
var ReloadAlways bool

/*
This file contains the helpers that come with the templates.
*/

var Helpers template.FuncMap

type FileError struct {
	Filename string
	Error    error
}

type BuildError []FileError

func (be BuildError) Error() string {
	errstr := "We had the following errors:{ "
	for _, fe := range be {
		errstr += fmt.Sprintf("%v:“%v”,", fe.Filename, fe.Error)
	}
	errstr = errstr[:len(errstr)-1]
	errstr += " }"
	return errstr
}

type KeyValueType struct {
	Key   string
	Value string
}

type OrderedMapType []KeyValueType

func (omap OrderedMapType) Exists(key string) bool {
	if len(omap) == 0 {
		return false
	}
	for _, item := range omap {
		if item.Key == key {
			return true
		}
	}
	return false
}

func (omap OrderedMapType) ExistsAt(key string) (at int, ok bool) {
	if len(omap) == 0 {
		return 0, false
	}
	for i, item := range omap {
		if item.Key == key {
			return i, true
		}
	}
	return 0, false
}

func (omap OrderedMapType) Set(key, value string) OrderedMapType {
	if at, ok := omap.ExistsAt(key); ok {
		omap[at].Value = value
		return omap
	}
	omap = append(omap, KeyValueType{key, value})
	return omap
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
// Gotten from: http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func moveFileContents(src, dst string) (err error) {
	closeInFile := true
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeInFile {
			in.Close()
		}
	}()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}

	in.Close()
	closeInFile = false
	os.Remove(src)
	return nil
}

// BuildFile will take the list of files concatenated and minify (if a minifier is provided) them.
// As the order of the files can be important we take the order of the file provided.
// If a file is listed more then once, only the first listing will be included.
func BuildFile(dist string, min Minifier, mimetype, oldname string, filenames ...string) (filename string, err error) {
	// First we have to check if the oldname file exists. If it does, then we don't do anything.
	if oldname != "" && !ReloadAlways {

		// If the file already exists on the File system do nothing.
		filename = filepath.Join(dist, oldname)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			return oldname, nil
		}
	}
	// This map will hold the files we have already seend and their sha1
	var filesum OrderedMapType
	prefix := "txtbuild"
	ext := ".txt"
	if mimetype != "" {
		exts, err := mime.ExtensionsByType(mimetype)
		if err != nil || len(exts) == 0 {
			goto DONE_PREFIX
		}
		// exts will be an array of exts with the leading period.
		ext = exts[0]
		prefix = ext[1:] + "build"
	}
DONE_PREFIX:

	// This is where we will do our work while we are building things out.
	tmpfile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up
	var be BuildError
	for _, filename := range filenames {
		if filesum.Exists(filename) {
			// Already saw the file; skip.
			continue
		}
		filesum = filesum.Set(filename, "")
		file, err := os.Open(filename)
		if err != nil {
			be = append(be, FileError{filename, err})
			continue
		}
		h := sha1.New()
		mw := io.MultiWriter(h, tmpfile)
		if min != nil {
			err = min.Minify(mimetype, mw, file)
			if err != nil {
				return filename, err
			}
		} else {
			io.Copy(mw, file)
		}

		file.Close()
		sha1sum := fmt.Sprintf("%x", h.Sum(nil))
		filesum = filesum.Set(filename, sha1sum)
	}
	if len(be) != 0 {
		return "", be
	}
	jsobj, err := json.Marshal(filesum)
	if err != nil {
		return "", err
	}

	fname := fmt.Sprintf("%v-%x%v", prefix, sha1.Sum(jsobj), ext)
	filename = filepath.Join(dist, fname)
	os.MkdirAll(dist, os.ModePerm)

	if err := moveFileContents(tmpfile.Name(), filename); err != nil {
		return "", err
	}
	return fname, nil
}


