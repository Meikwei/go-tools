/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-09 22:58:17
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-10 20:35:09
 * @FilePath: \go-tools\db\oceanutil\uril.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oceanutil

import (
	"github.com/Meikwei/go-tools/db/pagination"
	"github.com/Meikwei/go-tools/errs"
	"gorm.io/gorm"
)

// anes 函数是一个泛型函数，用于将一个类型为 T 的切片转换为一个类型为 any 的切片。
// 参数:
//   - ts []T: 一个类型为 T 的切片，T 可以是任意类型。
// 返回值:
//   - []any: 一个类型为 any 的切片，包含了输入切片 ts 的所有元素。
func anes[T any](ts []T) []any {
    // 创建一个类型为 any、长度与输入切片 ts 相同的切片。
    val := make([]any, len(ts))
    // 遍历输入切片 ts，并将每个元素转换为 any 类型后，添加到 val 切片中。
    for i := range ts {
        val[i] = ts[i]
    }
    // 返回转换后的切片 val。
    return val
}
func InsertMany[T any](coll *gorm.DB,val []T) (*gorm.DB,error) {
	result:=coll.Create(anes(val))
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase insert many")
	}
	return result,nil
}

func UpdateMany[T any](coll *gorm.DB,val T,where any,args ...any) (*gorm.DB,error) {
	result:=coll.Where(where,args...).Updates(val)
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase update many")
	}
	return result,nil
}

func Find[T any](coll *gorm.DB,val any,where any,args ...any)(*gorm.DB,error){
	result:=coll.Where(where,args...).Find(val);
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase find")
	}
	return result,nil
}

func FindOne[T any](coll *gorm.DB,val any,where any,args ...any)(*gorm.DB,error){
	result:=coll.Where(where,args...).First(val);
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase find one")
	}
	return result,nil
}

func FindOneAndUpdate[T any](coll *gorm.DB,val T,where any,args ...any)(*gorm.DB,error){
	result:=coll.Where(where,args...).FirstOrCreate(val);
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase find one and update")
	}
	return result,nil
}

func FindPage[T any](coll *gorm.DB,val any,pagination pagination.Pagination,where any,args ...any)(*gorm.DB,error){

	result:=coll.Where(where,args...).Offset(int(pagination.GetPageNumber())-1).Limit(int(pagination.GetShowNumber())).Find(val);
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase find page")
	}
	return result,nil
}

func deleteMany[T any](coll *gorm.DB,val any,where any,args ...any)(*gorm.DB,error){
	result:=coll.Where(where,args...).Delete(val);
	if result.Error != nil {
		return nil,errs.WrapMsg(result.Error,"oceanbase delete many")
	}
	return result,nil
}

func Aggregate[T any](coll *gorm.DB,val any,key string, value any)error{
	result:=coll.Set(key,value).AutoMigrate(val);
	if result != nil {
		return errs.WrapMsg(result,"oceanbase aggregate")
	}
	return nil
}








