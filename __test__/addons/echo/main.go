package main

import (
	"C"
	"go-napi-sys/__test__/addons/echo/napisys"
	"unsafe"
)

func echo(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	params, _, _, _ := napisys.GetCbInfo(env, info)
	
	value, _ := napisys.CreateInt32(napisys.Env(env), int32(len(params)))
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: echo,
	}
	desc := napisys.Property{
		Name:   "echo",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
