/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-04-29 21:52:30
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-01 15:59:25
 * @FilePath: \tools\mw\gin.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// Copyright © 2023 OpenIM. All rights reserved.
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

package mw

import (
	"net/http"

	"github.com/Meikwei/go-tools/apiresp"
	"github.com/Meikwei/go-tools/errs"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/constant"
)

// CorsHandler 是 Gin 框架的跨域配置处理函数。
// 此函数返回一个 gin.HandlerFunc，用于设置响应中的 CORS 头部，
// 实现广泛的跨域访问控制，以支持跨源请求。
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置标准 CORS 头部，允许任何来源、方法和头部。
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		// 指定预检请求中暴露的头部及最大存活时间。
		c.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar",
		)
		c.Header(
			"Access-Control-Max-Age",
			"172800",
		)

		// 设置是否支持凭证（如 cookies、授权头、TLS 客户端证书等）。
		c.Header(
			"Access-Control-Allow-Credentials",
			"false",
		)

		// 设置响应内容类型为 JSON。
		c.Header(
			"content-type",
			"application/json",
		)

		// 预检请求处理，返回 JSON 响应并中断处理链。
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}

		// 对于非预检请求，继续处理链中的下一个中间件。
		c.Next()
	}
}

// GinParseOperationID 用于解析并提取请求头中的 OperationID。
// 此函数返回一个 gin.HandlerFunc，检查是否为 POST 请求，
// 提取 OperationID 头部，并将其设置到 gin 上下文中。
// 如果缺少 OperationID，则中断请求链并返回错误响应。
func GinParseOperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断请求方法是否为 POST。
		if c.Request.Method == http.MethodPost {
			// 尝试从请求头中获取 OperationID。
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				// 如果 OperationID 缺失，返回错误响应并中断请求。
				err := errs.New("header must have operationID")
				apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
				c.Abort()
				return
			}
			// 将提取的 OperationID 设置到 gin 上下文中以供后续使用。
			c.Set(constant.OperationID, operationID)
		}
		// 继续执行处理链中的下一个中间件。
		c.Next()
	}
}