// +build debug

package template_test

import (
	"testing"

	"github.com/gdey/template"
)

func TestTemplateDebug(t *testing.T) {

	fixture := FileList{
		BaseDir: "tpl",
		Files: []FileType{
			{"parsefile.template", `This is a template! {{.}}`},
		},
	}
	defer fixture.CreateFilesOFail(t).RemoveAll()

	tpl := template.Must(
		template.Must(
			template.New("parsefile.template",
				template.ParseFile("tpl/parsefile.template"),
				template.DistRoot("tpl/dist"),
			)).ParseFiles())

	ExecuteTemplateOrFail(t, tpl, "hello", "This is a template! hello")
	fixture.SetFile("parsefile.template", `This is a template after! {{.}}`).CreateFileOrFail(t, "parsefile.template")

	ExecuteTemplateOrFail(t, tpl, "hello", "This is a template after! hello")
}

func TestTemplateBuildLinkJS1Debug(t *testing.T) {

	fixture := FileList{
		BaseDir: "tpl",
		Files: []FileType{
			{"views/1.js", `alert(1);`},
			{"parsefile.template", "{{buildLinkToJSFiles `tpl/views/1.js`}}"},
		},
	}
	defer fixture.CreateFilesOFail(t).RemoveAll()

	tpl := template.Must(
		template.Must(
			template.New("parsefile.template",
				template.ParseFile("tpl/parsefile.template"),
				template.DistRoot("tpl/dist"),
			)).ParseFiles())

	ExecuteTemplateOrFail(t, tpl, "hello", `<script type="text/javascript" src="jsbuild-cbe88841a0d1c699592e2a61e3ffa7c33b61a4f7.js"></script>`)

	fixture.SetFile("views/1.js", "alert(2);").CreateFileOrFail(t, "views/1.js")
	ExecuteTemplateOrFail(t, tpl, "hello", `<script type="text/javascript" src="jsbuild-5cf0442672da09dfce402e2f3adbe5bf0139d0ec.js"></script>`)
}
