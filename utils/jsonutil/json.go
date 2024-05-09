package jsonutil

import (
	"encoding/json"

	"github.com/Meikwei/go-tools/errs"
)

// JsonMarshal 将 Go 语言数据结构 v 转换为 JSON 字节切片。
// 参数:
//   v any - 需要被转换为 JSON 格式的任意类型数据。
// 返回值:
//   []byte - 转换后的 JSON 字节切片。
//   error - 在转换过程中遇到错误时返回的错误信息。
func JsonMarshal(v any) ([]byte, error) {
	m, err := json.Marshal(v) // 使用 Go 标准库的 json.Marshal 进行 JSON 序列化
	return m, errs.Wrap(err)  // 包装并返回可能的错误
}

// JsonUnmarshal 将 JSON 字节切片 b 反序列化为 Go 语言数据结构 v。
// 参数:
//   b []byte - 需要被反序列化的 JSON 字节切片。
//   v any - 用于存储反序列化结果的任意类型数据。
// 返回值:
//   error - 在反序列化过程中遇到错误时返回的错误信息。
func JsonUnmarshal(b []byte, v any) error {
	return errs.Wrap(json.Unmarshal(b, v)) // 使用 Go 标准库的 json.Unmarshal 进行 JSON 反序列化
}

// StructToJsonString 将 Go 语言数据结构转换为 JSON 字符串。
// 参数:
//   param any - 需要被转换为 JSON 字符串的任意类型数据。
// 返回值:
//   string - 转换后的 JSON 字符串。
func StructToJsonString(param any) string {
	dataType, _ := JsonMarshal(param)   // 将数据序列化为 JSON 字节切片
	dataString := string(dataType)       // 将字节切片转换为字符串
	return dataString                    // 返回 JSON 字符串
}

// JsonStringToStruct 将 JSON 字符串转换为 Go 语言数据结构。
// 参数:
//   s string - 需要被转换为 Go 数据结构的 JSON 字符串。
//   args any - 用于存储转换结果的任意类型数据。
// 返回值:
//   error - 在转换过程中遇到错误时返回的错误信息。
func JsonStringToStruct(s string, args any) error {
	err := json.Unmarshal([]byte(s), args) // 使用 Go 标准库的 json.Unmarshal 进行 JSON 反序列化
	return err                              // 返回可能的错误
}