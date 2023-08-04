package reflect

import "reflect"

func IterateArray(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	l := val.Len()

	res := make([]any, 0, l)
	for i := 0; i < l; i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}

	return res, nil
}

func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	resKeys := make([]any, 0, val.Len())
	resValues := make([]any, 0, val.Len())
	keys := val.MapKeys()
	for _, key := range keys {
		v := val.MapIndex(key)
		resKeys = append(resKeys, key)
		resValues = append(resValues, v.Interface())
	}
	return resKeys, resValues, nil
}
