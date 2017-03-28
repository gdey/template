package helpers

func BuildJSFile(dest, oldfilename string, filenames ...string) (string, error) {
	/*
		min := minify.New()
		min.AddFunc("text/javascript", js.Minify)
	*/
	return BuildFile(dest, nil, "text/javascript", oldfilename, filenames...)
}

func BuildCSSFile(dest, oldfilename string, filenames ...string) (string, error) {
	/*
		min := minify.New()
		min.AddFunc("text/css", css.Minify)
	*/
	return BuildFile(dest, nil, "text/css", oldfilename, filenames...)
}
