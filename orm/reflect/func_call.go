package reflect

import "reflect"

func IterateFunc(entity any, args ...any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()

	res := make(map[string]FuncInfo)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func
		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn)
		inputValue := make([]reflect.Value, 0, numIn)
		inputValue = append(inputValue, reflect.ValueOf(entity))
		for j := 0; j < numIn; j++ {
			in := fn.Type().In(j)
			input = append(input, in)
			if j != 0 {
				inputValue = append(inputValue, reflect.Zero(in))
			}
		}

		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))
		}

		resValue := fn.Call(inputValue)
		result := make([]any, 0, len(resValue))

		for _, value := range resValue {
			result = append(result, value.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}

	return res, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
