package zen

import (
	"context"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"time"
)

// Context abstracts the HTTP request/response context.
// It embeds context.Context so it can be passed directly to context-aware functions.
//
// For gin-specific features that require gin types (Negotiate, Render,
// ShouldBindWith, MustBindWith, etc.) use zgin.GinContext(c) to obtain the
// underlying *gin.Context directly.
type Context interface {
	context.Context

	// --- URL params ---
	Param(key string) string
	AddParam(key, value string)

	// --- Query string ---
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) []string
	GetQueryArray(key string) ([]string, bool)
	QueryMap(key string) map[string]string
	GetQueryMap(key string) (map[string]string, bool)

	// --- Request info ---
	Header(key string) string
	SetHeader(key, value string)
	ContentType() string
	GetRawData() ([]byte, error)
	IsWebsocket() bool
	RemoteIP() string
	ClientIP() string
	FullPath() string
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	// --- Handler introspection ---
	HandlerName() string
	HandlerNames() []string

	// --- Bind (writes HTTP 400 on error) ---
	Bind(v any) error
	BindJSON(v any) error
	BindXML(v any) error
	BindQuery(v any) error
	BindYAML(v any) error
	BindTOML(v any) error
	BindPlain(v any) error
	BindUri(v any) error
	BindHeader(v any) error

	// --- ShouldBind (returns error, no auto response) ---
	ShouldBind(v any) error
	ShouldBindJSON(v any) error
	ShouldBindXML(v any) error
	ShouldBindQuery(v any) error
	ShouldBindYAML(v any) error
	ShouldBindTOML(v any) error
	ShouldBindPlain(v any) error
	ShouldBindUri(v any) error
	ShouldBindHeader(v any) error

	// ShouldBindBodyWith* re-reads the body from an internal cache — safe for
	// multiple consecutive binds in the same request.
	ShouldBindBodyWithJSON(v any) error
	ShouldBindBodyWithXML(v any) error
	ShouldBindBodyWithYAML(v any) error
	ShouldBindBodyWithTOML(v any) error
	ShouldBindBodyWithPlain(v any) error

	// --- Form ---
	PostForm(key string) string
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) []string
	GetPostFormArray(key string) ([]string, bool)
	PostFormMap(key string) map[string]string
	GetPostFormMap(key string) (map[string]string, bool)

	// --- Multipart / file upload ---
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string, perm ...fs.FileMode) error

	// --- Cookie ---
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetCookieData(cookie *http.Cookie)
	SetSameSite(samesite http.SameSite)

	// --- Response: JSON variants ---
	JSON(code int, obj any)
	SecureJSON(code int, obj any)
	JSONP(code int, obj any)
	IndentedJSON(code int, obj any)
	AsciiJSON(code int, obj any)
	PureJSON(code int, obj any)

	// --- Response: other formats ---
	HTML(code int, name string, obj any)
	XML(code int, obj any)
	YAML(code int, obj any)
	TOML(code int, obj any)
	ProtoBuf(code int, obj any)
	BSON(code int, obj any)
	String(code int, format string, values ...any)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	Status(code int)
	Redirect(code int, location string)

	// --- File responses ---
	File(filepath string)
	FileFromFS(filepath string, fs http.FileSystem)
	FileAttachment(filepath, filename string)

	// --- Streaming ---
	SSEvent(name string, message any)
	Stream(step func(w io.Writer) bool) bool

	// --- Content negotiation ---
	NegotiateFormat(offered ...string) string
	SetAccepted(formats ...string)

	// --- Context store ---
	Set(key string, value any)
	Get(key string) (any, bool)
	MustGet(key string) any
	Delete(key string)
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt8(key string) int8
	GetInt16(key string) int16
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint8(key string) uint8
	GetUint16(key string) uint16
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat32(key string) float32
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string

	// --- Error tracking ---
	// AddError appends err to the request's error list.
	// Errors are available to downstream middleware (e.g., logging, recovery).
	AddError(err error)

	// --- Flow control ---
	Next()
	Abort()
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, obj any)
	AbortWithStatusPureJSON(code int, obj any)
	AbortWithError(code int, err error)
	IsAborted() bool
}

// Handler is the function signature for HTTP handlers.
type Handler func(Context)

// Middleware is a semantic alias for Handler.
type Middleware = Handler
