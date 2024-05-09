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
	"context"
	"fmt"
	"strings"

	"github.com/Meikwei/go-tools/errs"
	"github.com/Meikwei/go-tools/log"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/errinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GrpcClient 创建并返回一个gRPC的DialOption，配置了客户端的链式拦截器。
func GrpcClient() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(RpcClientInterceptor)
}

// RpcClientInterceptor 是一个gRPC的链式拦截器函数，用于在客户端发起的RPC调用前后添加额外的处理逻辑。
// ctx: 上下文，用于传递请求的元数据和控制请求的生命周期。
// method: 要调用的RPC方法名。
// req: 请求的消息体。
// resp: 响应的消息体。
// cc: gRPC的客户端连接实例。
// invoker: gRPC的调用者，用于执行实际的RPC调用。
// opts: gRPC调用的选项。
// 返回值: 执行过程中可能出现的错误。
func RpcClientInterceptor(ctx context.Context, method string, req, resp any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {

	// 检查上下文是否为nil
	if ctx == nil {
		return errs.ErrInternalServer.WrapMsg("call rpc request context is nil")
	}

	// 通过method生成和丰富上下文
	ctx, err = getRpcContext(ctx, method)
	if err != nil {
		return err
	}

	// 记录RPC客户端请求日志
	log.ZDebug(ctx, fmt.Sprintf("RPC Client Request - %s", extractFunctionName(method)), "funcName", method, "req", req, "conn target", cc.Target())
	
	// 执行实际的RPC调用
	err = invoker(ctx, method, req, resp, cc, opts...)
	if err == nil {
		// 记录RPC客户端响应成功日志
		log.ZInfo(ctx, fmt.Sprintf("RPC Client Response Success - %s", extractFunctionName(method)), "funcName", method, "resp", rpcString(resp))
		return nil
	}

	// 记录RPC客户端响应错误日志
	log.ZError(ctx, fmt.Sprintf("RPC Client Response Error - %s", extractFunctionName(method)), err, "funcName", method)

	// 处理gRPC错误
	rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
	if !ok {
		return errs.ErrInternalServer.WrapMsg(err.Error())
	}
	sta := rpcErr.GRPCStatus()
	if sta.Code() == 0 {
		return errs.NewCodeError(errs.ServerInternalError, err.Error()).Wrap()
	}

	// 处理错误详情，如果有
	if details := sta.Details(); len(details) > 0 {
		errInfo, ok := details[0].(*errinfo.ErrorInfo)
		if ok {
			s := strings.Join(errInfo.Warp, "->") + errInfo.Cause
			return errs.NewCodeError(int(sta.Code()), sta.Message()).WithDetail(s).Wrap()
		}
	}
	return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
}

// getRpcContext 生成或更新上下文，添加自定义的头部信息和一些关键的上下文变量。
// ctx: 输入的上下文。
// method: 当前的RPC方法名。
// 返回值: 更新后的上下文和可能遇到的错误。
func getRpcContext(ctx context.Context, method string) (context.Context, error) {
	md := metadata.Pairs()
	
	// 从ctx中提取自定义头部并添加到metadata中
	if keys, _ := ctx.Value(constant.RpcCustomHeader).([]string); len(keys) > 0 {
		for _, key := range keys {
			val, ok := ctx.Value(key).([]string)
			if !ok {
				return nil, errs.ErrInternalServer.WrapMsg("ctx missing key", "key", key)
			}
			if len(val) == 0 {
				return nil, errs.ErrInternalServer.WrapMsg("ctx key value is empty", "key", key)
			}
			md.Set(key, val...)
		}
		md.Set(constant.RpcCustomHeader, keys...)
	}

	// 提取并添加operationID到metadata
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		log.ZWarn(ctx, "ctx missing operationID", errs.New("ctx missing operationID"), "funcName", method)
		return nil, errs.ErrArgs.WrapMsg("ctx missing operationID")
	}
	md.Set(constant.OperationID, operationID)

	// 提取并添加opUserID到metadata
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok {
		md.Set(constant.OpUserID, opUserID)
	}

	// 提取并添加opUserIDPlatformID到metadata
	opUserIDPlatformID, ok := ctx.Value(constant.OpUserPlatform).(string)
	if ok {
		md.Set(constant.OpUserPlatform, opUserIDPlatformID)
	}

	// 提取并添加connID到metadata
	connID, ok := ctx.Value(constant.ConnID).(string)
	if ok {
		md.Set(constant.ConnID, connID)
	}

	// 返回带有更新后metadata的上下文
	return metadata.NewOutgoingContext(ctx, md), nil
}

// extractFunctionName 从方法名中提取函数名。
// funcName: 要处理的方法名。
// 返回值: 提取的函数名。
func extractFunctionName(funcName string) string {
	parts := strings.Split(funcName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
