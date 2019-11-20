package main

import (
	"C"
	"go-napi-sys/__test__/addons/simple-n/napisys"
	"unsafe"
)

func createInt32(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	value, _ := napisys.CreateInt32(napisys.Env(env), 7)
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: createInt32,
	}
	desc := napisys.Property{
		Name:   "createInt32",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
