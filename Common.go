package AnimaKit

import "github.com/robertkrimen/otto"

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func toValueOrPanic(any interface{}) otto.Value {
	ans, err := VM.ToValue(any)
	panicOnError(err)
	return ans
}
