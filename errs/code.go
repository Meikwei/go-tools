/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 10:49:29
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-05 10:58:07
 * @FilePath: \go-tools\errors\code.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package errs

// 定义一系列的错误码常量，用于标识不同的错误情况

const (
	// 通用错误码
	ServerInternalError = 500  // 服务器内部错误
	ArgsError           = 1001 // 输入参数错误
	NoPermissionError   = 1002 // 权限不足
	DuplicateKeyError   = 1003 // 键重复错误
	RecordNotFoundError = 1004 // 记录不存在错误

	// 与令牌相关的错误码
	TokenExpiredError     = 1501 // 令牌过期错误
	TokenMalformedError   = 1503 // 令牌格式错误
	TokenNotValidYetError = 1504 // 令牌尚未生效错误
	TokenUnknownError     = 1505 // 未知令牌错误
)