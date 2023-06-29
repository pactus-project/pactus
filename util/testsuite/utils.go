package testsuite

import "reflect"

func DeepCopy[T any](src *T) *T {
	srcVal := reflect.ValueOf(*src)
	dstVal := reflect.New(srcVal.Type()).Elem()
	reflect.Copy(dstVal, srcVal)

	return dstVal.Interface().(*T)
}
