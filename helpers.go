package template

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdey/template/helpers"
)

func makeKey(filenames []string) (key string) {
	for _, f := range filenames {
		key += f
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(key)))
}

func filepatternToFilenames(base string, patterns []string) (filenames []string, err error) {

	for _, pat := range patterns {
		glob := filepath.Join(base, pat)

		files, err := filepath.Glob(glob)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				log.Printf("unable to get working directory(%v)", err)
				// We are setting wd for the following message; it's not used any where else.
				wd = "UNKNOWN"
			}
			log.Printf("WARNING: Glob(%v) did not match any files. Base Directory: (%v)", glob, wd)
		}
		filenames = append(filenames, files...)
	}
	return filenames, nil
}

// BuildMimeTypeFile is a helper function that takes a MimeType and a set of filenames and generated a combined (minimizied if a minimizer is provided)
// file.
func (t *Template) BuildMimeTypeFile(mimetype string, fnames string) (filename string, err error) {

	var patterns []string
	for _, fname := range strings.Split(fnames, ",") {
		fn := strings.TrimSpace(fname)
		if fn == "" {
			continue
		}
		patterns = append(patterns, fn)
	}
	filenames, err := filepatternToFilenames(t.base, patterns)
	if err != nil {
		return "", err
	}
	key := makeKey(filenames)
	t.buildLock.Lock()
	oldFilename := t.buildFileOldFilenameCaché[key]
	dest := t.dist
	t.buildLock.Unlock()

	if filename, err = helpers.BuildFile(dest, t.minifiers[mimetype], mimetype, oldFilename, filenames...); err != nil {
		return filename, err
	}
	if oldFilename != filename {
		t.buildLock.Lock()
		t.buildFileOldFilenameCaché[key] = filename
		t.buildLock.Unlock()
	}
	return filename, err
}

// BuildJSFile is a helper function that takes a set of filename and generated a combined (minimizied if a minimizer is provided)
// Javascript file.
func (t *Template) BuildJSFile(fnames string) (filename string, err error) {
	return t.BuildMimeTypeFile(helpers.JSMimeType, fnames)
}

// LinkToAndBuildJSFile is the same as the buildJSFiles but will return a script tag contain the appropriate URL.
func (t *Template) LinkToAndBuildJSFile(fnames string) (template.HTML, error) {
	filename, err := t.BuildJSFile(fnames)
	if err != nil {
		return "", err
	}

	if t.root != "" {
		filename = strings.Join([]string{t.root, filename}, "/")
	}
	return template.HTML(fmt.Sprintf(`<script type="text/javascript" src="%v"></script>`, filename)), nil
}

// BuildCSSFile is a helper function that takes a set of filename and generated a combined (minimizied if a minimizer is provided)
// Javascript file.
func (t *Template) BuildCSSFile(fnames string) (filename string, err error) {
	return t.BuildMimeTypeFile(helpers.CSSMimeType, fnames)
}

// LinkToAndBuildCSSFile is the same as the buildCSSFiles but will return a link tag contain the appropriate URL.
func (t *Template) LinkToAndBuildCSSFile(fnames string) (template.HTML, error) {
	filename, err := t.BuildCSSFile(fnames)
	if err != nil {
		return "", err
	}

	if t.root != "" {
		filename = strings.Join([]string{t.root, filename}, "/")
	}
	return template.HTML(fmt.Sprintf(`<link rel = "stylesheet" type="text/css" href="%v" />`, filename)), nil
}
