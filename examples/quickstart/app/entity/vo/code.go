package vo

const (
	INVALID_REQUEST_PARAMETERS = "INVALID_REQUEST_PARAMETERS" // 请求参数无效
	SUCCESS                    = "SUCCESS"                    // 操作成功
	FAILED                     = "FAILED"                     // 操作失败
	CREATION_FAILED            = "CREATION_FAILED"            // 创建失败
	CREATION_SUCCESS           = "CREATION_SUCCESS"           // 创建成功
	UPDATE_FAILED              = "UPDATE_FAILED"              // 更新失败
	UPDATE_SUCCESS             = "UPDATE_SUCCESS"             // 更新成功
	DELETE_FAILED              = "DELETE_FAILED"              // 删除失败
	DELETE_SUCCESS             = "DELETE_SUCCESS"             // 删除成功
	LOGIN_ERROR                = "LOGIN_ERROR"                // 登录失败
	LOGIN_SUCCESS              = "LOGIN_SUCCESS"              // 登录成功
	CAPTCHA_ERROR              = "CAPTCHA_ERROR"              // 验证码错误

	AUTH_ERROR            = "AUTH_ERROR"            // 认证失败
	AUTH_DEVICE_ERROR     = "AUTH_DEVICE_ERROR"     // 设备不存在
	AUTH_USER_ERROR       = "AUTH_USER_ERROR"       // 用户不存在
	PASSWORD_ERROR        = "PASSWORD_ERROR"        // 密码错误
	USERNAME_EXISTS       = "USERNAME_EXISTS"       // 用户名已存在
	CHANNEL_CODE_EXISTS   = "CHANNEL_CODE_EXISTS"   // 渠道代码已存在
	PERMISSIONS_ERROR     = "PERMISSIONS_ERROR"     //权限不足
	CATEGORY_HAS_PRODUCTS = "CATEGORY_HAS_PRODUCTS" // 分类下存在产品，无法删除
)

var Messages = map[string]string{
	"INVALID_REQUEST_PARAMETERS": "请求参数无效，请检查输入。",
	"SUCCESS":                    "操作成功！",
	"FAILED":                     "操作失败，请稍后重试。",
	"CREATION_FAILED":            "创建失败，请稍后重试。",
	"CREATION_SUCCESS":           "创建成功！",
	"UPDATE_FAILED":              "更新失败，请稍后重试。",
	"UPDATE_SUCCESS":             "更新成功！",
	"DELETE_FAILED":              "删除失败，请稍后重试。",
	"DELETE_SUCCESS":             "删除成功！",
	"LOGIN_ERROR":                "登录失败，请检查用户名或密码。",
	"LOGIN_SUCCESS":              "登录成功，欢迎回来！",
	"CAPTCHA_ERROR":              "验证码错误，请重新输入。",

	"AUTH_ERROR":            "认证失败，请重新登录。",
	"AUTH_DEVICE_ERROR":     "认证失败，设备不存在。",
	"AUTH_USER_ERROR":       "认证失败，用户不存在。",
	"PASSWORD_ERROR":        "密码错误，请重新输入。",
	"USERNAME_EXISTS":       "用户名已存在，请更换其他用户名。",
	"CHANNEL_CODE_EXISTS":   "渠道唯一标识已存在，请更换其他渠道代码。",
	"PERMISSIONS_ERROR":     "权限不足，无法执行此操作。",
	"CATEGORY_HAS_PRODUCTS": "该分类下存在产品，无法删除",
}

func GetMessage(code string) string {
	if msg, exists := Messages[code]; exists {
		return msg
	}
	if code != "" {

		return code
	}
	return "操作未成功，请稍后再试。若多次失败，请联系客服。"
}
