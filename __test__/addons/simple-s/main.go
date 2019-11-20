package main

import (
	"C"
	"go-napi-sys/__test__/addons/simple-s/napisys"
	"unsafe"
)

func hello(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	value, _ := napisys.CreateStringUtf8(napisys.Env(env), "world")
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: hello,
	}
	desc := napisys.Property{
		Name:   "hello",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
