package zen

// RouterGroup abstracts a group of routes sharing a prefix and middleware.
type RouterGroup interface {
	GET(path string, h ...Handler)
	POST(path string, h ...Handler)
	PUT(path string, h ...Handler)
	DELETE(path string, h ...Handler)
	PATCH(path string, h ...Handler)
	Use(mw ...Handler)
	Group(prefix string, mw ...Handler) RouterGroup
}
