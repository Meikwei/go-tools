/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 10:43:12
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-05 11:00:51
 * @FilePath: \go-tools\errs\error.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 10:43:12
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-05 11:00:40
 * @FilePath: \go-tools\errors\error.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package errs

import (
	"bytes"
	"fmt"
)

// Error 定义了一个错误接口，增加了错误的可扩展性。
type Error interface {
	// Is 检查当前错误是否与另一个错误相等。
	Is(err error) bool
	// Wrap 返回当前错误的一个封装，允许错误链的构建。
	Wrap() error
	// WrapMsg 返回当前错误的一个封装，并添加额外的消息。
	WrapMsg(msg string, kv ...any) error
	error
}

// New 创建并返回一个新的Error实例，允许错误消息和键值对的传递。
func New(s string, kv ...any) Error {
	return &errorString{
		s: toString(s, kv),
	}
}

// errorString 是一个实现了Error接口的错误类型。
type errorString struct {
	s string
}

// Is 检查当前错误是否与另一个错误相等。
func (e *errorString) Is(err error) bool {
	if err == nil {
		return false
	}
	t, ok := err.(*errorString)
	return ok && e.s == t.s
}

// Error 返回当前错误的字符串表示。
func (e *errorString) Error() string {
	return e.s
}

// Wrap 返回当前错误的一个封装，允许错误链的构建。
func (e *errorString) Wrap() error {
	return Wrap(e)
}

// WrapMsg 返回当前错误的一个封装，并添加额外的消息。
func (e *errorString) WrapMsg(msg string, kv ...any) error {
	return WrapMsg(e, msg, kv...)
}


// toString 用于构造错误消息，支持在错误消息中添加键值对信息。
func toString(s string, kv []any) string {
	if len(kv) == 0 {
		return s
	} else {
		var buf bytes.Buffer
		buf.WriteString(s)

		for i := 0; i < len(kv); i += 2 {
			if buf.Len() > 0 {
				buf.WriteString(", ")
			}

			key := fmt.Sprintf("%v", kv[i])
			buf.WriteString(key)
			buf.WriteString("=")

			if i+1 < len(kv) {
				value := fmt.Sprintf("%v", kv[i+1])
				buf.WriteString(value)
			} else {
				buf.WriteString("MISSING")
			}
		}
		return buf.String()
	}
}