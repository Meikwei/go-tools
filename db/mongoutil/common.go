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

package mongoutil

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const (
	defaultMaxPoolSize = 100
	defaultMaxRetry    = 3
)

// buildMongoURI 构造 MongoDB 的 URI 从提供的配置。
// 参数:
//   config *Config: 包含 MongoDB 连接所需配置信息的结构体指针。
// 返回值:
//   string: 格式化的 MongoDB URI 字符串。
func buildMongoURI(config *Config) string {
	credentials := ""
	if config.Username != "" && config.Password != "" {
		// 如果配置了用户名和密码，则构造认证信息
		credentials = fmt.Sprintf("%s:%s@", config.Username, config.Password)
	}
	// 使用配置信息构造 MongoDB URI
	return fmt.Sprintf("mongodb://%s%s/%s?maxPoolSize=%d", credentials, strings.Join(config.Address, ","), config.Database, config.MaxPoolSize)
}

// shouldRetry 判断一个错误是否应该触发重试。
// 参数:
//   ctx context.Context: 上下文，用于控制函数的生命周期。
//   err error: 执行操作过程中发生的错误。
// 返回值:
//   bool: 如果错误应该触发重试则返回 true，否则返回 false。
func shouldRetry(ctx context.Context, err error) bool {
	select {
	case <-ctx.Done():
		// 如果上下文被取消或超时，则不重试
		return false
	default:
		// 检查错误是否为 mongo.CommandError 类型
		if cmdErr, ok := err.(mongo.CommandError); ok {
			// 不重试特定的错误代码
			return cmdErr.Code != 13 && cmdErr.Code != 18
		}
		// 对于非 mongo.CommandError 类型的错误，默认重试
		return true
	}
}