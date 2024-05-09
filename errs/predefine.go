package errs

// 定义了一系列常见的错误类型，方便在应用程序中识别和处理特定的错误情况。
var (
	ErrArgs             = NewCodeError(ArgsError, "ArgsError")             // 参数错误，表明传入的参数不满足要求。
	ErrNoPermission     = NewCodeError(NoPermissionError, "NoPermissionError")     // 无权限错误，表示当前操作没有足够的权限。
	ErrInternalServer   = NewCodeError(ServerInternalError, "ServerInternalError")   // 内部服务器错误，通常指服务器在处理请求时发生了未预期的状况。
	ErrRecordNotFound   = NewCodeError(RecordNotFoundError, "RecordNotFoundError")   // 记录未找到错误，表示根据给定的条件无法找到匹配的记录。
	ErrDuplicateKey     = NewCodeError(DuplicateKeyError, "DuplicateKeyError")     // 键重复错误，表明尝试插入或更新的记录的键已存在于数据库中。
	ErrTokenMalformed   = NewCodeError(TokenMalformedError, "TokenMalformedError")   // Token格式错误，表示提供的Token格式不正确或缺失必要字段。
	ErrTokenNotValidYet = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") // Token尚未生效错误，表明提供的Token虽然有效，但其生效时间还未到达。
	ErrTokenUnknown     = NewCodeError(TokenUnknownError, "TokenUnknownError")     // Token未知错误，表示提供的Token无法被识别或已失效。
	ErrTokenExpired     = NewCodeError(TokenExpiredError, "TokenExpiredError")     // Token过期错误，表示提供的Token已超过其有效期。
)