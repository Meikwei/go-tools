/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 22:06:00
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 22:08:49
 * @FilePath: \go-tools\db\tx\mongoTx.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tx

import "context"

type MongoTx interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}