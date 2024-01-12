package build

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/evanw/esbuild/pkg/api"
)

func Build() {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"./static/app.jsx"},
		Bundle:      true,
		Outfile:     "./out.js",
		Platform:    api.PlatformBrowser,
		Loader: map[string]api.Loader{
			".jsx": api.LoaderJSX,
		},
		JSXFactory:  "h",
		JSXFragment: "Fragment",
		Write:       true,
	})

	if len(result.Errors) != 0 {
		spew.Dump(result.Errors)
		os.Exit(1)
	}
}
