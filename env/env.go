/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 20:19:51
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 20:24:07
 * @FilePath: \go-tools\env\env.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package env

import (
	"os"
	"strconv"

	"github.com/Meikwei/go-tools/errs"
)

// GetString 根据给定的键返回环境变量的值。
// 如果该键未设置，则返回提供的默认值。
func GetString(key, defaultValue string) string {
	// 查找与键关联的环境变量
	v, ok := os.LookupEnv(key)
	if ok {
		// 如果找到键，则返回其值
		return v
	}
	// 如果键未找到，则返回默认值
	return defaultValue
}

// GetInt 返回环境变量解析为整数的值，或在未设置时返回默认值。
// 它将与键关联的值解析为整数。
func GetInt(key string, defaultValue int) (int, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		// 尝试将环境变量值转换为整数
		value, err := strconv.Atoi(v)
		if err != nil {
			// 如果转换失败，使用 errs 包封装错误并返回默认值
			return defaultValue, errs.WrapMsg(err, "Atoi failed", "value", v)
		}
		return value, nil
	}
	// 键未设置，直接返回默认值和nil错误
	return defaultValue, nil
}

// GetFloat64 返回环境变量解析为浮点数的值，或在未设置时返回默认值。
// 它将与键关联的值解析为64位浮点数。
func GetFloat64(key string, defaultValue float64) (float64, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		// 尝试将环境变量值转换为浮点数
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			// 如果转换失败，使用 errs 包封装错误并返回默认值
			return defaultValue, errs.WrapMsg(err, "ParseFloat failed", "value", v)
		}
		return value, nil
	}
	// 键未设置，直接返回默认值和nil错误
	return defaultValue, nil
}

// GetBool 返回环境变量解析为布尔值的值，或在未设置时返回默认值。
// 它将与键关联的值解析为布尔值。
func GetBool(key string, defaultValue bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if ok {
		// 尝试将环境变量值转换为布尔值
		value, err := strconv.ParseBool(v)
		if err != nil {
			// 如果转换失败，使用 errs 包封装错误并返回默认值
			return defaultValue, errs.WrapMsg(err, "ParseBool failed", "value", v)
		}
		return value, nil
	}
	// 键未设置，直接返回默认值和nil错误
	return defaultValue, nil
}
