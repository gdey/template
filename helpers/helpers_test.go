package helpers

import (
	"testing"

	"os"

	"github.com/gdey/tbltest"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

func TestBuildJS(t *testing.T) {
	type testcase struct {
		dist         string
		min          Minifier
		mimetype     string
		oldfilename  string
		expected     string
		jsfilenames  []string
		expectedErr  error
		expectedSha1 string
		dontCleanup  bool
	}

	min := minify.New()
	min.AddFunc("text/javascript", js.Minify)

	tests := tbltest.Cases(
		testcase{
			jsfilenames: []string{
				"assets/js/1.js",
				"assets/js/2.js",
			},
			mimetype: "text/javascript",
			dist:     "assets/dist",
			expected: "jsbuild-6a141653d4cbe0f00f3b485746777a0927da0d5d.js",
		},
		testcase{
			jsfilenames: []string{
				"assets/js/1.js",
				"assets/js/2.js",
				"assets/js/4.js",
			},
			min:      min,
			mimetype: "text/javascript",
			dist:     "assets/dist",
			expected: "jsbuild-4bb902a19c32ea4ad7b8114351ef023d39efe456.js",
		},
		testcase{
			jsfilenames: []string{
				"assets/js/1.js",
				"assets/js/2.js",
				"assets/js/3.js", // this would error if we were not short-circuiting the logic.
			},
			mimetype:    "text/javascript",
			oldfilename: "jsbuild-6e756c6cda39a3ee5e6b4b0d3255bfef95601890afd80709.js",
			dist:        "assets/distexists",
			expected:    "jsbuild-6e756c6cda39a3ee5e6b4b0d3255bfef95601890afd80709.js",
			dontCleanup: true,
		},
	)

	tests.Run(func(idx int, test testcase) {
		gotFilename, err := BuildFile(test.dist, test.min, test.mimetype, test.oldfilename, test.jsfilenames...)
		if err != test.expectedErr {
			t.Errorf("Test %v: Expected error %v but got %v", idx, test.expectedErr, err)
			return
		}
		if test.expectedErr != nil {
			return
		}
		if gotFilename != test.expected {
			t.Errorf("Test: %v: Expected filesname: %v but got %v", idx, test.expected, gotFilename)
			return
		}
		// Check the sha1 contents to see if they match.
		//clean up file.
		if !test.dontCleanup {
			os.Remove(gotFilename)
		}

	})

}
