package zen

import "net/http"

// RouterGroup abstracts a group of routes sharing a prefix and middleware.
// All standard gin routing methods are exposed for full compatibility.
type RouterGroup interface {
	GET(path string, h ...Handler)
	POST(path string, h ...Handler)
	PUT(path string, h ...Handler)
	DELETE(path string, h ...Handler)
	PATCH(path string, h ...Handler)
	HEAD(path string, h ...Handler)
	OPTIONS(path string, h ...Handler)
	Any(path string, h ...Handler)
	Handle(method, path string, h ...Handler)
	Use(mw ...Handler)
	Group(prefix string, mw ...Handler) RouterGroup
	StaticFile(relativePath, filepath string)
	Static(relativePath, root string)
	StaticFS(relativePath string, fs http.FileSystem)
}
