package oceanutil

import (
	"context"
	"database/sql"

	"github.com/Meikwei/go-tools/db/tx"
	"github.com/Meikwei/go-tools/errs"
	"gorm.io/gorm"
)
type oceanTx struct {
	client *gorm.DB
	tx     func(context.Context, func(tx *gorm.DB) error,...*sql.TxOptions) error
}
// WithTransaction 是 GORM 的事务包装函数
func (o oceanTx) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	tx := o.client.Begin(opts...)
	defer func() {
		if r := recover(); r != nil {
			o.client.Rollback()
		}
	}()
	o.client.WithContext(ctx)
	if err := o.client.Error; err != nil {
		return err
	}
	// 执行事务中的函数
	err := fn(tx)
	if err != nil {
		tx.Rollback() // 如果有错误，回滚事务
		return err
	}
	return tx.Commit().Error // 提交事务
}
func (o *oceanTx) init() error {
	// 持续会话模式禁用事务
	o.client.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	})
	o.tx=func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error{
		err:=o.WithTransaction(ctx,func(tx *gorm.DB)(error){
			return fn(tx)
		},opts...)
		return errs.WrapMsg(err,"oceanbase transaction failed")
	}
	return nil;
}

func (o *oceanTx) Transaction(ctx context.Context,fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error{
	if o.tx==nil{
		o.client.WithContext(ctx)
		if err := o.client.Error; err != nil {
			return err
		}
		return fn(o.client)
		
	}
	return o.tx(ctx,fn,opts...)
}

func NewOceanTx(ctx context.Context,client *gorm.DB) (tx.Tx,error){
	otx:=oceanTx{client:client}
	if err:=otx.init(); err != nil{
		return nil,err
	}
	return &otx,nil
}

func NewOcean(client *gorm.DB) tx.Tx{
	return &oceanTx{client:client}
}
