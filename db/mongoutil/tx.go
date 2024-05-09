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

package mongoutil

import (
	"context"

	"github.com/Meikwei/go-tools/db/tx"
	"github.com/Meikwei/go-tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewMongoTx 创建一个MongoDB事务对象。
//
// ctx: 上下文，用于控制请求的生命周期。
// client: MongoDB客户端，用于执行数据库操作。
//
// 返回一个tx.Tx接口实现和可能发生的错误。
func NewMongoTx(ctx context.Context, client *mongo.Client) (tx.MongoTx, error) {
    // 初始化mongoTx结构体
	mtx := mongoTx{
		client: client,
	}
    // 尝试初始化事务功能
	if err := mtx.init(ctx); err != nil {
		return nil, err
	}
	return &mtx, nil
}

// NewMongo 创建一个不带事务的MongoDB操作对象。
//
// client: MongoDB客户端，用于执行数据库操作。
//
// 返回一个tx.Tx接口实现。
func NewMongo(client *mongo.Client) tx.MongoTx {
	return &mongoTx{
		client: client,
	}
}

// mongoTx 是对tx.Tx接口的实现，封装了MongoDB的事务处理。
type mongoTx struct {
	client *mongo.Client
	tx     func(context.Context, func(ctx context.Context) error) error
}

// init 初始化mongoTx的事务功能。
//
// ctx: 上下文，用于控制初始化过程。
//
// 返回可能发生的错误。
func (m *mongoTx) init(ctx context.Context) error {
    // 检查MongoDB是否部署在集群中
	var res map[string]any
	if err := m.client.Database("admin").RunCommand(ctx, bson.M{"isMaster": 1}).Decode(&res); err != nil {
		return errs.WrapMsg(err, "check whether mongo is deployed in a cluster")
	}
    // 如果不是集群，不支持事务
	if _, allowTx := res["setName"]; !allowTx {
		return nil // non-clustered transactions are not supported
	}
    // 设置事务执行函数
	m.tx = func(fnctx context.Context, fn func(ctx context.Context) error) error {
		sess, err := m.client.StartSession()
		if err != nil {
			return errs.WrapMsg(err, "mongodb start session failed")
		}
		defer sess.EndSession(fnctx)
        // 使用MongoDB session执行事务
		_, err = sess.WithTransaction(fnctx, func(sessCtx mongo.SessionContext) (any, error) {
			return nil, fn(sessCtx)
		})
		return errs.WrapMsg(err, "mongodb transaction failed")
	}
	return nil
}

// Transaction 执行事务或直接执行函数（取决于事务是否可用）。
//
// ctx: 上下文，用于控制操作的生命周期。
// fn: 需要在事务中执行的函数。
//
// 返回可能发生的错误。
func (m *mongoTx) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
    // 如果事务功能可用，则使用事务执行；否则直接执行函数
	if m.tx == nil {
		return fn(ctx)
	}
	return m.tx(ctx, fn)
}
