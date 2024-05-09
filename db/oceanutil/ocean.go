/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 10:39:57
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-07 00:05:39
 * @FilePath: \go-tools\db\oceanutil\ocean.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oceanutil

import (
	"context"

	"github.com/Meikwei/go-tools/db/tx"
	"github.com/Meikwei/go-tools/errs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 是一个用于存储配置信息的结构体。
// 它包含了与服务器连接相关的所有必要信息。
type Config struct {
	Dns 				string // 连接服务器的URI
	Host        string // 服务器主机地址
	Username    string // 连接服务器的用户名
	TenantName  string // 租户名称
	ClusterName string // 集群名称
	Password    string // 连接服务器的密码
	Port        int    // 服务器端口
	SchemaName  string // 使用的数据库模式名称
	DataType    string // 数据类型，mysql、oracle
	Charset     string // 字符集名称
	MaxRetry    int    // 最大重试次数
}
type Client struct {
	tx tx.Tx
	db *gorm.DB
}

func (c *Client) GetDB() *gorm.DB {
	return c.db
}

func (c *Client) GetTx() tx.Tx {
	return c.tx
}

func NewOceanbase(ctx context.Context,config *Config) (*Client, error) {
	if err := config.ValidateAndSetDefaults(); err != nil {
		return nil, err
	}
	db, err := initializeDBWithRetry(ctx, config)
	if err !=nil{
		return nil, errs.WrapMsg(err, "failed to connect to oceanbase", "Dns", config.Dns)
	}
	otx,err:=NewOceanTx(ctx,db)
	if err !=nil{
		return nil,err
	}
	return &Client{
		tx: otx,
		db: db,
	},nil
}

func connectOcean(config *Config) (*gorm.DB, error){
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.Dns,
		// DefaultStringSize: 256, // string 类型字段的默认长度
		// DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		// DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		// DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err == nil {
		if result := db.Raw("SELECT 1 FROM DUAL"); result.Error != nil {
			// 数据库连接成功
			return db, nil
		}
		sql,_:=db.DB();
		defer sql.Close();
	}
	return db, err
}

