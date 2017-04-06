package main

import (
	"os"

	"github.com/gdey/template"
)

// Example is the example.
func Example() {
	tpl := template.Must(
		template.Must(template.New("main.template",
			template.ParseFileList("tpl/parsefile.txt"),
			template.URLBase("static"),
			template.DistRoot("examples/tmp"),
		)).ParseFiles())
	if err := tpl.Execute(os.Stdout, "No Data"); err != nil {
		panic(err)
	}

	// Output:
	// <html>
	// <head>
	//
	//
	//	<meta charset="utf-8">
	//	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	//
	//
	// <script type="text/javascript" src="static/jsbuild-0b96c8ce82fd02cfea234be7618cd86358b7119e.js"></script>
	//
	// </head>
	// <body>
	// </body>
}

func main() { Example() }
