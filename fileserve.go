package zephyr

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

type defaultFS struct {
	root          string
	folderAliases map[string]string
	fs.FS
}

func DefaultFsHandler(fileRoot string) *defaultFS {
	dfs := os.DirFS(fileRoot)
	return &defaultFS{
		folderAliases: map[string]string{
			"@": "node_modules",
		},
		root: fileRoot,
		FS:   dfs,
	}

}

func (dfs *defaultFS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileType := path.Ext(r.RequestURI)
	switch fileType {
	case ".js", ".ts", ".jsx", ".tsx":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "unknown")
		w.Write([]byte("got" + fileType + ": unhandled filetype or empty"))
		return
	}

	_, remain, _ := strings.Cut(strings.TrimLeft(r.RequestURI, "/"), "/")

	for replaceStr, alias := range dfs.folderAliases {
		if !strings.Contains(remain, replaceStr) {
			continue
		}

		remain = strings.Replace(remain, replaceStr, alias, 1)
	}

	file, err := dfs.Open(remain)
	if err != nil {
		w.Write([]byte("error opening file:" + err.Error()))
		return
	}
	defer file.Close()

	out, err := io.ReadAll(file)
	w.Write(out)
}

func setHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}
