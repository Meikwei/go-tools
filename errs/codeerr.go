package errs

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)
const initialCapacity = 3
const minimumCodesLength = 2
var DefaultCodeRelation = newCodeRelation()
// CodeError 接口定义了一个带有错误码、错误信息和详细信息的错误类型接口。
type CodeError interface {
    Code() int    // 返回错误码
    Msg() string  // 返回错误信息
    Detail() string  // 返回错误的详细信息
    WithDetail(detail string) CodeError  // 添加详细信息到错误对象，返回新的CodeError实例
    Error
}

// NewCodeError 创建并返回一个新的CodeError实例。
func NewCodeError(code int, msg string) CodeError {
    return &codeError{
        code: code,
        msg: msg,
    }
}

// codeError 是对CodeError接口的实现。
type codeError struct {
    code   int
    msg    string
    detail string
}

// Code 实现CodeError接口的Code方法。
func (e *codeError) Code() int {
    return e.code
}

// Msg 实现CodeError接口的Msg方法。
func (e *codeError) Msg() string {
    return e.msg
}

// Detail 实现CodeError接口的Detail方法。
func (e *codeError) Detail() string {
    return e.detail
}

// WithDetail 实现CodeError接口的WithDetail方法，用于添加错误的详细信息。
func (e *codeError) WithDetail(detail string) CodeError {
    var d string
    if e.detail == "" {
        d = detail
    } else {
        d = e.detail + ", " + detail
    }
    return &codeError{
        code:   e.code,
        msg:    e.msg,
        detail: d,
    }
}

// Wrap 方法将codeError转换为标准错误类型。
func (e *codeError) Wrap() error {
    return Wrap(e)
}

// WrapMsg 方法为错误添加额外的消息，并转换为标准错误类型。
func (e *codeError) WrapMsg(msg string, kv ...any) error {
    return WrapMsg(e, msg, kv...)
}

// Is 方法用于检查当前错误是否与另一个错误匹配。
func (e *codeError) Is(err error) bool {
    codeErr, ok := Unwrap(err).(CodeError)
    if !ok {
        if err == nil && e == nil {
            return true
        }
        return false
    }
    if e == nil {
        return false
    }
    code := codeErr.Code()
    if e.code == code {
        return true
    }
    return DefaultCodeRelation.Is(e.code, code)
}

// Error 方法返回错误的字符串表示。
func (e *codeError) Error() string {
    v := make([]string, 0, initialCapacity)
    v = append(v, strconv.Itoa(e.code), e.msg)

    if e.detail != "" {
        v = append(v, e.detail)
    }

    return strings.Join(v, " ")
}

// Unwrap 方法用于解开嵌套的错误，直到找到非nil的普通错误类型。
func Unwrap(err error) error {
    for err != nil {
        unwrap, ok := err.(interface {
            Unwrap() error
        })
        if !ok {
            break
        }
        err = unwrap.Unwrap()
    }
    return err
}

// Wrap 方法为错误添加堆栈信息。
func Wrap(err error) error {
    return errors.WithStack(err)
}

// WrapMsg 方法为错误添加额外的消息和堆栈信息。
func WrapMsg(err error, msg string, kv ...any) error {
    if err == nil {
        return nil
    }
    withMessage := errors.WithMessage(err, toString(msg, kv))
    return errors.WithStack(withMessage)
}

// CodeRelation 接口定义了错误码之间的关系。
type CodeRelation interface {
    Add(codes ...int) error              // 添加错误码之间的关系
    Is(parent, child int) bool          // 检查错误码之间是否存在关系
}

// newCodeRelation 创建并返回一个CodeRelation的实现实例。
func newCodeRelation() CodeRelation {
    return &codeRelation{m: make(map[int]map[int]struct{})}
}

// codeRelation 是对CodeRelation接口的实现，用于管理错误码之间的关系。
type codeRelation struct {
    m map[int]map[int]struct{}
}

// Add 方法用于建立错误码之间的父子关系。
func (r *codeRelation) Add(codes ...int) error {
    if len(codes) < minimumCodesLength {
        return New("codes length must be greater than 2", "codes", codes).Wrap()
    }
    for i := 1; i < len(codes); i++ {
        parent := codes[i-1]
        s, ok := r.m[parent]
        if !ok {
            s = make(map[int]struct{})
            r.m[parent] = s
        }
        for _, code := range codes[i:] {
            s[code] = struct{}{}
        }
    }
    return nil
}

// Is 方法用于检查两个错误码之间是否存在关系。
func (r *codeRelation) Is(parent, child int) bool {
    if parent == child {
        return true
    }
    s, ok := r.m[parent]
    if !ok {
        return false
    }
    _, ok = s[child]
    return ok
}
