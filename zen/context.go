package zen

import (
	"mime/multipart"
	"net/http"
)

// Context abstracts the HTTP request/response context.
// It wraps the underlying engine's context but exposes a framework-neutral API.
// All gin.Context methods are available, ensuring full compatibility.
type Context interface {
	// --- Path / Query params ---
	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	QueryMap(key string) map[string]string

	// --- Request headers ---
	Header(key string) string
	SetHeader(key, value string)
	ContentType() string
	GetRawData() ([]byte, error)
	IsWebsocket() bool

	// --- Bind (writes HTTP 400 on error) ---
	Bind(v any) error
	BindJSON(v any) error
	BindQuery(v any) error
	BindUri(v any) error
	BindHeader(v any) error

	// --- ShouldBind (returns error, no auto response) ---
	ShouldBind(v any) error
	ShouldBindJSON(v any) error
	ShouldBindQuery(v any) error
	ShouldBindUri(v any) error
	ShouldBindHeader(v any) error

	// --- Form ---
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	PostFormMap(key string) map[string]string

	// --- Multipart / File upload ---
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error

	// --- Cookie ---
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetSameSite(samesite http.SameSite)

	// --- Response: JSON variants ---
	JSON(code int, v any)
	SecureJSON(code int, v any)
	JSONP(code int, obj any)
	IndentedJSON(code int, obj any)
	AsciiJSON(code int, obj any)
	PureJSON(code int, obj any)

	// --- Response: other formats ---
	XML(code int, obj any)
	YAML(code int, obj any)
	String(code int, format string, values ...any)
	Data(code int, contentType string, data []byte)
	Status(code int)
	Redirect(code int, location string)

	// --- Underlying request / response ---
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	ClientIP() string
	FullPath() string

	// --- Context store ---
	Set(key string, value any)
	Get(key string) (any, bool)
	MustGet(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	GetStringMapString(key string) map[string]string

	// --- Flow control ---
	Next()
	Abort()
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, v any)
	IsAborted() bool
}

// Handler is the function signature for HTTP handlers.
type Handler func(Context)

// Middleware is a semantic alias for Handler.
type Middleware = Handler
