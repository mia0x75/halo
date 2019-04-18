package tools

import "reflect"

// Contains 判断一个元素是否存在与某一个数组中
func Contains(array interface{}, val interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}
