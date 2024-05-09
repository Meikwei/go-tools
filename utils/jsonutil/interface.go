/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 20:31:46
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 20:33:21
 * @FilePath: \go-tools\utils\jsonutil\interface.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package jsonutil

// Json 定义了一个处理JSON数据的接口。
type Json interface {
	// Interface 返回底层数据。
	Interface() any

	// Encode 将其 marshaled 数据作为 `[]byte` 返回。
	Encode() ([]byte, error)

	// EncodePretty 将其 marshaled 数据作为带缩进的 `[]byte` 返回。
	EncodePretty() ([]byte, error)

	// MarshalJSON 实现了 json.Marshaler 接口。
	MarshalJSON() ([]byte, error)

	// Set 通过键值对修改 `Json` 映射。
	// 适用于轻松更改 `Json` 对象中的单个键/值。
	Set(key string, val any)

	// SetPath 递归检查/创建映射键，然后写入值到指定路径的 `Json` 中。
	SetPath(branch []string, val any)

	// Del 通过键删除 `Json` 映射中的条目，如果存在的话。
	Del(key string)

	// Get 返回一个新的 `Json` 对象的指针，
	// 用于其映射表示中的 `key`。
	//
	// 适用于链式操作（遍历嵌套的 JSON）：
	//    js.Get("top_level").Get("dict").Get("value").Int()
	Get(key string) Json

	// GetPath 通过指定路径查找项，
	// 无需深入使用 Get() 进行查找。
	//
	//   js.GetPath("top_level", "dict")
	GetPath(branch ...string) Json

	// CheckGet 返回一个新的 `Json` 对象的指针和
	// 一个标识成功或失败的 `bool` 值。
	//
	// 适用于当成功很重要时的链式操作：
	//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
	//        log.Println(data)
	//    }
	CheckGet(key string) (Json, bool)

	// Map 将类型断言为 `map`。
	Map() (map[string]any, error)

	// Array 将类型断言为一个 `array`。
	Array() ([]any, error)

	// Bool 将类型断言为 `bool`。
	Bool() (bool, error)

	// String 将类型断言为 `string`。
	String() (string, error)

	// Bytes 将类型断言为 `[]byte`。
	Bytes() ([]byte, error)

	// StringArray 将类型断言为一个 `string` 的数组。
	StringArray() ([]string, error)

	// UnmarshalJSON 实现了 json.Unmarshaler 接口。
	UnmarshalJSON(p []byte) error

	// Float64 强制转换为 float64。
	Float64() (float64, error)

	// Int 强制转换为 int。
	Int() (int, error)

	// Int64 强制转换为 int64。
	Int64() (int64, error)

	// Uint64 强制转换为 uint64。
	Uint64() (uint64, error)
}
