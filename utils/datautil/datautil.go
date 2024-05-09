package datautil

import (
	"reflect"
	"sort"

	"github.com/Meikwei/go-tools/errs"
	"github.com/Meikwei/go-tools/utils/jsonutil"
	"github.com/jinzhu/copier"
)

// SliceSub 函数返回在切片 a 中但不在切片 b 中的元素（a - b）。
// 参数 a 为原始切片，参数 b 为要从中剔除元素的切片。
// 返回值为一个新的切片，包含切片 a 中不包含在切片 b 中的元素。
func SliceSub[E comparable](a, b []E) []E {
	// 如果切片 b 为空，则直接返回切片 a，因为不存在需要剔除的元素。
	if len(b) == 0 {
		return a
	}

	// 使用 map k 来记录切片 b 中的所有元素，以便后续查找。
	k := make(map[E]struct{})
	for i := 0; i < len(b); i++ {
		k[b[i]] = struct{}{}
	}

	// 使用 map t 来记录已经处理过的元素，避免重复添加。
	t := make(map[E]struct{})
	
	// 初始化结果切片 rs，预分配足够的空间以优化内存使用。
	rs := make([]E, 0, len(a))

	// 遍历切片 a，找出不在切片 b 中的元素。
	for i := 0; i < len(a); i++ {
		e := a[i]
		// 如果元素 e 已经被处理过，则跳过。
		if _, ok := t[e]; ok {
			continue
		}
		// 如果元素 e 在切片 b 中，则跳过。
		if _, ok := k[e]; ok {
			continue
		}
		// 将元素 e 添加到结果切片 rs，并标记为已处理。
		rs = append(rs, e)
		t[e] = struct{}{}
	}

	// 返回结果切片 rs。
	return rs
}

// SliceSubAny returns elements in slice a that are not present in slice b (a - b).
// SliceSubAny 是一个函数，它通过将 slice b 中的元素通过转换函数 fn 转换为与 slice a 中元素可比较的形式，
// 然后从 a 中去除这些转换后的元素，返回一个新的切片。
// 
// 参数：
// E - a 和转换后元素的类型，必须是可比较的。
// T - b 中元素的类型，可以是任意类型。
// a - 原始切片 a，其元素类型为 E。
// b - 要进行转换和去除的切片 b，其元素类型为 T。
// fn - 一个函数，用于将 b 中的每个元素转换为 E 类型。
// 
// 返回值：
// 返回一个新的切片，它是由 a 中不包含通过 fn 转换后的 b 中元素的部分组成的。
func SliceSubAny[E comparable, T any](a []E, b []T, fn func(t T) E) []E {
	// 使用 SliceSub 函数和转换后的 b 切片来从 a 中去除元素，然后返回结果。
	return SliceSub(a, Slice(b, fn))
}

// SliceAnySub returns elements in slice a that are not present in slice b (a - b).
// fn is a function that extracts a comparable value from elements of slice a.
func SliceAnySub[E any, T comparable](a, b []E, fn func(t E) T) []E {
	m := make(map[T]E)
	for i := 0; i < len(b); i++ {
		v := b[i]
		m[fn(v)] = v
	}
	var es []E
	for i := 0; i < len(a); i++ {
		v := a[i]
		if _, ok := m[fn(v)]; !ok {
			es = append(es, v)
		}
	}
	return es
}

// DistinctAny duplicate removal.
func DistinctAny[E any, K comparable](es []E, fn func(e E) K) []E {
	v := make([]E, 0, len(es))
	tmp := map[K]struct{}{}
	for i := 0; i < len(es); i++ {
		t := es[i]
		k := fn(t)
		if _, ok := tmp[k]; !ok {
			tmp[k] = struct{}{}
			v = append(v, t)
		}
	}
	return v
}

func DistinctAnyGetComparable[E any, K comparable](es []E, fn func(e E) K) []K {
	v := make([]K, 0, len(es))
	tmp := map[K]struct{}{}
	for i := 0; i < len(es); i++ {
		t := es[i]
		k := fn(t)
		if _, ok := tmp[k]; !ok {
			tmp[k] = struct{}{}
			v = append(v, k)
		}
	}
	return v
}

func Distinct[T comparable](ts []T) []T {
	if len(ts) < 2 {
		return ts
	} else if len(ts) == 2 {
		if ts[0] == ts[1] {
			return ts[:1]
		} else {
			return ts
		}
	}
	return DistinctAny(ts, func(t T) T {
		return t
	})
}

// Delete Delete slice elements, support negative number to delete the reciprocal number
func Delete[E any](es []E, index ...int) []E {
	switch len(index) {
	case 0:
		return es
	case 1:
		i := index[0]
		if i < 0 {
			i = len(es) + i
		}
		if len(es) <= i {
			return es
		}
		return append(es[:i], es[i+1:]...)
	default:
		tmp := make(map[int]struct{})
		for _, i := range index {
			if i < 0 {
				i = len(es) + i
			}
			tmp[i] = struct{}{}
		}
		v := make([]E, 0, len(es))
		for i := 0; i < len(es); i++ {
			if _, ok := tmp[i]; !ok {
				v = append(v, es[i])
			}
		}
		return v
	}
}

// DeleteAt Delete slice elements, support negative number to delete the reciprocal number
func DeleteAt[E any](es *[]E, index ...int) []E {
	v := Delete(*es, index...)
	*es = v
	return v
}

// IndexAny get the index of the element
func IndexAny[E any, K comparable](e E, es []E, fn func(e E) K) int {
	k := fn(e)
	for i := 0; i < len(es); i++ {
		if fn(es[i]) == k {
			return i
		}
	}
	return -1
}

// IndexOf get the index of the element
func IndexOf[E comparable](e E, es ...E) int {
	return IndexAny(e, es, func(t E) E {
		return t
	})
}

// Contain Whether to include
func Contain[E comparable](e E, es ...E) bool {
	return IndexOf(e, es...) >= 0
}

// DuplicateAny Whether there are duplicates
func DuplicateAny[E any, K comparable](es []E, fn func(e E) K) bool {
	t := make(map[K]struct{})
	for _, e := range es {
		k := fn(e)
		if _, ok := t[k]; ok {
			return true
		}
		t[k] = struct{}{}
	}
	return false
}

// Duplicate Whether there are duplicates
func Duplicate[E comparable](es []E) bool {
	return DuplicateAny(es, func(e E) E {
		return e
	})
}

// SliceToMapOkAny slice to map (Custom type, filter)
func SliceToMapOkAny[E any, K comparable, V any](es []E, fn func(e E) (K, V, bool)) map[K]V {
	kv := make(map[K]V)
	for i := 0; i < len(es); i++ {
		t := es[i]
		if k, v, ok := fn(t); ok {
			kv[k] = v
		}
	}
	return kv
}

// SliceToMapAny slice to map (Custom type)
func SliceToMapAny[E any, K comparable, V any](es []E, fn func(e E) (K, V)) map[K]V {
	return SliceToMapOkAny(es, func(e E) (K, V, bool) {
		k, v := fn(e)
		return k, v, true
	})
}

// SliceToMap slice to map
func SliceToMap[E any, K comparable](es []E, fn func(e E) K) map[K]E {
	return SliceToMapOkAny(es, func(e E) (K, E, bool) {
		k := fn(e)
		return k, e, true
	})
}

// SliceSetAny slice to map[K]struct{}
func SliceSetAny[E any, K comparable](es []E, fn func(e E) K) map[K]struct{} {
	return SliceToMapAny(es, func(e E) (K, struct{}) {
		return fn(e), struct{}{}
	})
}

func Filter[E, T any](es []E, fn func(e E) (T, bool)) []T {
	rs := make([]T, 0, len(es))
	for i := 0; i < len(es); i++ {
		e := es[i]
		if t, ok := fn(e); ok {
			rs = append(rs, t)
		}
	}
	return rs
}

// Slice Converts slice types in batches
func Slice[E any, T any](es []E, fn func(e E) T) []T {
	v := make([]T, len(es))
	for i := 0; i < len(es); i++ {
		v[i] = fn(es[i])
	}
	return v
}

// SliceSet slice to map[E]struct{}
func SliceSet[E comparable](es []E) map[E]struct{} {
	return SliceSetAny(es, func(e E) E {
		return e
	})
}

// HasKey get whether the map contains key
func HasKey[K comparable, V any](m map[K]V, k K) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

// Min get minimum value
func Min[E Ordered](e ...E) E {
	v := e[0]
	for _, t := range e[1:] {
		if v > t {
			v = t
		}
	}
	return v
}

// Max get maximum value
func Max[E Ordered](e ...E) E {
	v := e[0]
	for _, t := range e[1:] {
		if v < t {
			v = t
		}
	}
	return v
}

func Paginate[E any](es []E, pageNumber int, showNumber int) []E {
	if pageNumber <= 0 {
		return []E{}
	}
	if showNumber <= 0 {
		return []E{}
	}
	start := (pageNumber - 1) * showNumber
	end := start + showNumber
	if start >= len(es) {
		return []E{}
	}
	if end > len(es) {
		end = len(es)
	}
	return es[start:end]
}

// BothExistAny gets elements that are common in the slice (intersection)
func BothExistAny[E any, K comparable](es [][]E, fn func(e E) K) []E {
	if len(es) == 0 {
		return []E{}
	}
	var idx int
	ei := make([]map[K]E, len(es))
	for i := 0; i < len(ei); i++ {
		e := es[i]
		if len(e) == 0 {
			return []E{}
		}
		kv := make(map[K]E)
		for j := 0; j < len(e); j++ {
			t := e[j]
			k := fn(t)
			kv[k] = t
		}
		ei[i] = kv
		if len(kv) < len(ei[idx]) {
			idx = i
		}
	}
	v := make([]E, 0, len(ei[idx]))
	for k := range ei[idx] {
		all := true
		for i := 0; i < len(ei); i++ {
			if i == idx {
				continue
			}
			if _, ok := ei[i][k]; !ok {
				all = false
				break
			}
		}
		if !all {
			continue
		}
		v = append(v, ei[idx][k])
	}
	return v
}

// BothExist Gets the common elements in the slice (intersection)
func BothExist[E comparable](es ...[]E) []E {
	return BothExistAny(es, func(e E) E {
		return e
	})
}

//func CompleteAny[K comparable, E any](ks []K, es []E, fn func(e E) K) bool {
//	if len(ks) == 0 && len(es) == 0 {
//		return true
//	}
//	kn := make(map[K]uint8)
//	for _, e := range Distinct(ks) {
//		kn[e]++
//	}
//	for k := range SliceSetAny(es, fn) {
//		kn[k]++
//	}
//	for _, n := range kn {
//		if n != 2 {
//			return false
//		}
//	}
//	return true
//}

// Complete whether a and b are equal after deduplication (ignore order)
func Complete[E comparable](a []E, b []E) bool {
	return len(Single(a, b)) == 0
}

// Keys get map keys
func Keys[K comparable, V any](kv map[K]V) []K {
	ks := make([]K, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

// Values get map values
func Values[K comparable, V any](kv map[K]V) []V {
	vs := make([]V, 0, len(kv))
	for k := range kv {
		vs = append(vs, kv[k])
	}
	return vs
}

// Sort basic type sorting
func Sort[E Ordered](es []E, asc bool) []E {
	SortAny(es, func(a, b E) bool {
		if asc {
			return a < b
		} else {
			return a > b
		}
	})
	return es
}

// SortAny custom sort method
func SortAny[E any](es []E, fn func(a, b E) bool) {
	sort.Sort(&sortSlice[E]{
		ts: es,
		fn: fn,
	})
}

// If true -> a, false -> b
func If[T any](isa bool, a, b T) T {
	if isa {
		return a
	}
	return b
}

func ToPtr[T any](t T) *T {
	return &t
}

// Equal Compares slices to each other (including element order)
func Equal[E comparable](a []E, b []E) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Single exists in a and does not exist in b or exists in b and does not exist in a
func Single[E comparable](a, b []E) []E {
	kn := make(map[E]uint8)
	for _, e := range Distinct(a) {
		kn[e]++
	}
	for _, e := range Distinct(b) {
		kn[e]++
	}
	v := make([]E, 0, len(kn))
	for k, n := range kn {
		if n == 1 {
			v = append(v, k)
		}
	}
	return v
}

// Order sorts ts by es
func Order[E comparable, T any](es []E, ts []T, fn func(t T) E) []T {
	if len(es) == 0 || len(ts) == 0 {
		return ts
	}
	kv := make(map[E][]T)
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		k := fn(t)
		kv[k] = append(kv[k], t)
	}
	rs := make([]T, 0, len(ts))
	for _, e := range es {
		vs := kv[e]
		delete(kv, e)
		rs = append(rs, vs...)
	}
	for k := range kv {
		rs = append(rs, kv[k]...)
	}
	return rs
}

func OrderPtr[E comparable, T any](es []E, ts *[]T, fn func(t T) E) []T {
	*ts = Order(es, *ts, fn)
	return *ts
}

func UniqueJoin(s ...string) string {
	data, _ := jsonutil.JsonMarshal(s)
	return string(data)
}

type sortSlice[E any] struct {
	ts []E
	fn func(a, b E) bool
}

func (o *sortSlice[E]) Len() int {
	return len(o.ts)
}

func (o *sortSlice[E]) Less(i, j int) bool {
	return o.fn(o.ts[i], o.ts[j])
}

func (o *sortSlice[E]) Swap(i, j int) {
	o.ts[i], o.ts[j] = o.ts[j], o.ts[i]
}

// Ordered types that can be sorted
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// NotNilReplace sets old to new_ when new_ is not null
func NotNilReplace[T any](old, new_ *T) {
	if new_ == nil {
		return
	}
	*old = *new_
}

func StructFieldNotNilReplace(dest, src any) {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		srcField := srcVal.Field(i)

		// Check if the source field is valid
		if srcField.IsValid() {
			// Check if the target field can be set
			if destField.CanSet() {
				// Handling fields of slice type
				if destField.Kind() == reflect.Slice && srcField.Kind() == reflect.Slice {
					elemType := destField.Type().Elem()
					// Check if a slice element is a pointer to a structure
					if elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
						// Create a new slice to store the copied elements
						newSlice := reflect.MakeSlice(destField.Type(), srcField.Len(), srcField.Cap())
						for j := 0; j < srcField.Len(); j++ {
							newElem := reflect.New(elemType.Elem())
							// Recursive update, retaining non-zero values
							StructFieldNotNilReplace(newElem.Interface(), srcField.Index(j).Interface())
							// Checks if the field of the new element is zero-valued, and if so, preserves the value at the corresponding position in the original slice
							for k := 0; k < newElem.Elem().NumField(); k++ {
								if newElem.Elem().Field(k).IsZero() {
									newElem.Elem().Field(k).Set(destField.Index(j).Elem().Field(k))
								}
							}
							newSlice.Index(j).Set(newElem)
						}
						destField.Set(newSlice)
					} else {
						destField.Set(srcField)
					}
				} else {
					// For non-sliced fields, update the source field if it is non-zero, otherwise keep the original value
					if !srcField.IsZero() {
						destField.Set(srcField)
					}
				}
			}
		}
	}
}

func Batch[T any, V any](fn func(T) V, ts []T) []V {
	if ts == nil {
		return nil
	}
	res := make([]V, 0, len(ts))
	for i := range ts {
		res = append(res, fn(ts[i]))
	}
	return res
}

func InitSlice[T any](val *[]T) {
	if val != nil && *val == nil {
		*val = []T{}
	}
}

func InitMap[K comparable, V any](val *map[K]V) {
	if val != nil && *val == nil {
		*val = map[K]V{}
	}
}

func GetSwitchFromOptions(Options map[string]bool, key string) (result bool) {
	if Options == nil {
		return true
	}
	if flag, ok := Options[key]; !ok || flag {
		return true
	}
	return false
}

func SetSwitchFromOptions(options map[string]bool, key string, value bool) {
	if options == nil {
		options = make(map[string]bool, 5)
	}
	options[key] = value
}

// copy a by b  b->a
func CopyStructFields(a any, b any, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func GetElemByIndex(array []int, index int) (int, error) {
	if index < 0 || index >= len(array) {
		return 0, errs.New("index out of range", "index", index, "array", array).Wrap()
	}

	return array[index], nil
}
