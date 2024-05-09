/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-06 21:34:13
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-06 21:35:43
 * @FilePath: \go-tools\db\tx\tx.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tx

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type Tx interface {
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}