/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 20:33:49
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 20:35:15
 * @FilePath: \go-tools\utils\jsonutil\json_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package jsonutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJsonMarshal 测试JsonMarshal函数。
// 该函数用于检验JsonMarshal函数在对不同数据类型进行JSON编码时的正确性。
// 参数:
//   - t *testing.T: 测试环境的句柄，用于报告测试失败和日志记录。
func TestJsonMarshal(t *testing.T) {
    // 尝试对结构体数据进行JSON编码
    structData := struct{ Name string }{"John"}
    structBytes, err := JsonMarshal(structData)
    assert.NoError(t, err) // 确保编码过程中没有错误发生
    assert.JSONEq(t, `{"Name":"John"}`, string(structBytes)) // 确保编码结果符合预期

    // 尝试对json.RawMessage类型数据进行JSON编码
    marshalerData := json.RawMessage(`{"type":"raw"}`)
    marshalerBytes, err := JsonMarshal(marshalerData)
    assert.NoError(t, err) // 确保编码过程中没有错误发生
    assert.Equal(t, `{"type":"raw"}`, string(marshalerBytes)) // 确保编码结果符合预期
}

// TestJsonUnmarshal 测试JsonUnmarshal函数。
// 该函数用于检验JsonUnmarshal函数在对不同格式的JSON数据进行解码时的正确性。
// 参数:
//   - t *testing.T: 测试环境的句柄，用于报告测试失败和日志记录。
func TestJsonUnmarshal(t *testing.T) {
    // 尝试对JSON字符串进行解码到结构体
    structBytes := []byte(`{"Name":"Jane"}`)
    var structData struct{ Name string }
    err := JsonUnmarshal(structBytes, &structData)
    assert.NoError(t, err) // 确保解码过程中没有错误发生
    assert.Equal(t, "Jane", structData.Name) // 确保解码结果符合预期

    // 尝试对JSON字符串进行解码到json.RawMessage
    marshalerBytes := []byte(`{"type":"unmarshal"}`)
    var marshalerData json.RawMessage
    err = JsonUnmarshal(marshalerBytes, &marshalerData)
    assert.NoError(t, err) // 确保解码过程中没有错误发生
    assert.Equal(t, `{"type":"unmarshal"}`, string(marshalerData)) // 确保解码结果符合预期
}
