package template

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gdey/template/helpers"
)

// The Source of the parsefile
type sourceType uint

const (
	// SrcFileList is a file that contains a list of ParseFiles.
	SrcFileList = sourceType(iota)
	// SrcParseFile is just a ParseFile
	SrcParseFile
	// SrcGlobFile is a glob that may or may not result in ParseFiles.
	SrcGlobFile
)

type parseFileSources struct {
	Type sourceType
	File string
}

// DefaultBase is the default base directory for finding resources.
var DefaultBase string

func execDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "/"
	}
	return dir
}

func init() {

	var err error
	if DefaultBase, err = os.Getwd(); err != nil {
		log.Println("Unable to get current working directory, defaulting to exec dir for base.")
		DefaultBase = execDir()
	}

}

// Template is the main template object.
type Template struct {
	*template.Template

	// Initial name of the Template.
	name string
	// The directory where build assets are placed into.
	dist string
	// The base directory to find assets.
	base string

	// The base url
	root string

	// helpers are the Helpers that the user is adding
	helpers template.FuncMap

	// Lock for parseFileSouces and parseFiles
	parseLock sync.Mutex
	// This is the list of source for files to parse.
	parseFilesSources []parseFileSources

	// This is a list file to parse.
	parseFiles []string

	buildLock sync.Mutex
	// Build File cache; this holds a key and the old filename for the buildfile. This way,
	// Only the first called for a set of files generates the build file in a template.
	buildFileOldFilenameCaché map[string]string

	// minifiers are the list of minifiers that can be used to minify files; indexed by mimetype.
	minifiers map[string]helpers.Minifier
}

type anOption func(t *Template) error

type BConfig []anOption

// BaseConfig creates a base config that can be used to create other templates from it without have to respecify the
// options.
func BaseConfig(options ...anOption) BConfig {
	return BConfig(options)
}

// NewTemplate returns a template object based on the options set in the base config.
func (bc BConfig) NewTemplate(name string, options ...anOption) (*Template, error) {
	var opts []anOption
	opts = append(opts, bc...)
	opts = append(opts, options...)
	return New(name, opts...)
}

// DistRoot is the directory that the build files are going to get written to.
func DistRoot(dist string) anOption {
	return func(t *Template) error {
		if filepath.IsAbs(dist) {
			t.dist = dist
		} else {
			t.dist = filepath.Join(DefaultBase, dist)
		}
		return nil
	}
}

// Helpers allow you add helper methods to the template. If the value of the map is not a function that can be accepted or the name
// of the function is not something that can be a function name, this method will panic.
func Helpers(helpers ...template.FuncMap) anOption {
	return func(t *Template) error {
		for _, helper := range helpers {
			for k, v := range helper {
				t.helpers[k] = v
			}
		}
		return nil
	}
}

// ResourceRoot sets the base directory to use when resolving any resource.
func ResourceRoot(base string) anOption {
	return func(t *Template) error {
		// If the base is abs, we
		if filepath.IsAbs(base) {
			t.base = base
		} else {
			t.base = filepath.Join(DefaultBase, base)
		}
		return nil
	}
}

// URLBase set the base of the url that is generated by the LinkTo* helper functions.
func URLBase(root string) anOption {
	return func(t *Template) error {
		t.root = root
		return nil
	}
}

// Minifier to use for the given mimetype. Only one minifier is allowed per mimetype.
func Minifier(mimetype string, minifier helpers.Minifier) anOption {
	return func(t *Template) error {
		if _, ok := t.minifiers[mimetype]; ok {
			return fmt.Errorf("Minifier for “%v” already provided.", mimetype)
		}
		t.minifiers[mimetype] = minifier
		return nil
	}
}

func parsePossibleGlob(base, glob string) ([]string, error) {
	if base == "" {
		base = DefaultBase
	}
	glob = filepath.Join(base, glob)
	matches, err := filepath.Glob(glob)
	if err == nil && len(matches) == 0 {
		log.Printf("WARNING: Glob(%v) did not match any files.", glob)
	}
	return matches, err
}

// parsePossilbleGlob will take a string that could possibly be a glob and a base. First it makes sure the glob is relative to
// the base, then converts the glob to a set of files names.
func (t *Template) parsePossibleGlob(glob string) ([]string, error) {
	return parsePossibleGlob(t.base, glob)
}
func (t *Template) fullLock() {
	t.buildLock.Lock()
	t.parseLock.Lock()
}

func (t *Template) fullUnlock() {
	t.parseLock.Unlock()
	t.buildLock.Unlock()
}

func parseFileList(t *Template, filename string) error {

	base := filepath.Dir(filename)
	// Now we need open up the file, each line of the file will be a file path comment or empty.
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// We need to read each line of the file, and add the line to parse file if it does not start with
	// # or is empty.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if len(txt) == 0 {
			continue
		}
		if txt[0] == '#' {
			continue
		}

		switch txt[0] {
		case '.':
			// no opt

		case '/':

			base, _ = os.Getwd()
		default:

			base = t.base
		}

		matches, err := parsePossibleGlob(base, txt)
		switch err {
		default:
			return err
		case filepath.ErrBadPattern:
			t.parseFiles = append(t.parseFiles, filepath.Join(base, txt))
		case nil:

			t.parseFiles = append(t.parseFiles, matches...)
		}
	}
	return nil
}

// ParseFileList will add the files from one or more file lists to the set of files to parse for the template. File
// are reparsed in debug more for each execute statement.
func ParseFileList(ffile string, files ...string) anOption {
	return func(t *Template) error {
		addSourceFile(t, SrcFileList, ffile)
		if err := parseFileList(t, ffile); err != nil {
			return err
		}

		for _, file := range files {
			addSourceFile(t, SrcFileList, file)
			if err := parseFileList(t, file); err != nil {
				return err
			}
		}
		return nil
	}
}

// ParseFile will add the file provided to the set of files to parse for the template.
func ParseFile(files ...string) anOption {
	return func(t *Template) error {
		for _, file := range files {
			addSourceFile(t, SrcParseFile, file)
			t.parseFiles = append(t.parseFiles, file)
		}
		return nil
	}
}

func parseGlob(t *Template, glob string) error {
	matches, err := t.parsePossibleGlob(glob)
	switch err {
	default:
		return err
	case filepath.ErrBadPattern:
		t.parseFiles = append(t.parseFiles, glob)
	case nil:
		t.parseFiles = append(t.parseFiles, matches...)
	}
	return nil
}

// ParseGlob will add files it finds from the provided globs to the list of files to parse for the template.
func ParseGlob(globs ...string) anOption {
	return func(t *Template) error {
		for _, glob := range globs {
			addSourceFile(t, SrcGlobFile, glob)
			if err := parseGlob(t, glob); err != nil {
				return err
			}
		}
		return nil
	}
}

// New creates a new template. Name of the template and a set of options
func New(name string, options ...anOption) (*Template, error) {
	t := Template{
		name:                      name,
		Template:                  template.New(name),
		minifiers:                 make(map[string]helpers.Minifier),
		buildFileOldFilenameCaché: make(map[string]string),
	}

	// New we need to install all our Helpers. We first install our Helpers, then
	// We install the users handlers, this does mean that the user can overwrite our
	// Helpers
	t.helpers = template.FuncMap{
		"buildJSFiles":        t.BuildJSFile,
		"buildLinkToJSFiles":  t.LinkToAndBuildJSFile,
		"buildCSSFiles":       t.BuildCSSFile,
		"buildLinkToCSSFiles": t.LinkToAndBuildCSSFile,
	}

	for _, opt := range options {
		if err := opt(&t); err != nil {
			return &t, err
		}
	}
	t.Template.Funcs(t.helpers)
	return &t, nil
}

// Must will panic if there is an error other returns the template.
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
