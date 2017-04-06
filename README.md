# Template

This is a simple wrapper around [html/template](http://godoc.org/html/template). The main
benefit is that it allows reloading of assets without having to rebuild the 
binary in debug mode.

This library uses build tags to control debug v.s. nondebug mode.
To build with with debug mode turned on use:
```sh
$ go build -tags="debug" 
```

In addition these additional helper functions have been added by
default to the system:

-----------------------------------

 Command                 | Description
-------------------- | ------------
 buildJSFiles              | Concatenates the given file list, the list is expected to be in a comma separated string, into a file and returns the new file's name. 
 buildLinkToJSFiles    | Same as the buildJSFiles but will return a script tag contain the appropriate URL. 
 buildCSSFiles            | Concatenates the given file list, the list is expected to be in a comma separated string, into a file and returns the new file's name. 
 buildLinkToCSSFiles | Same as the buildCSSFiles but will return a link tag contain the appropriate URL.
 
 ---------------------------------

 Look at `examples/parsefilemin` for an example of how to use the package.
 
 ```go
package main

import (
	"os"

	"github.com/gdey/template"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

func main() {
    min := minify.New()
    min.AddFunc(helpers.JSMimeType, js.Minify)
    min.AddFunc(helpers.CSSMimeType, css.Minify)
	tpl := template.Must(
		template.Must(template.New("main.template",
			template.ParseFileList("tpl/parsefile.txt"),
			template.URLBase("static"),
			template.DistRoot("examples/tmp"),
			template.Minifier(helpers.JSMimeType, min),
            template.Minifier(helpers.CSSMimeType, min),
		)).ParseFiles())
	if err := tpl.Execute(os.Stdout, "No Data"); err != nil {
		panic(err)
	}
}

/*
output:
<html>
<head>


	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">


<script type="text/javascript" src="static/jsbuild-0b96c8ce82fd02cfea234be7618cd86358b7119e.js"></script>

</head>
<body>
</body>

*/


```
