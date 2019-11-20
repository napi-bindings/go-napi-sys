package main

import (
	"C"
	"go-napi-sys/__test__/addons/pass-arguments/napisys"
	"unsafe"
)

func add(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	params, _, _, _ := napisys.GetCbInfo(env, info)
	//v1, _ := napisys.GetValueInt32(env, params[0])
	//v2, _ := napisys.GetValueInt32(env, params[1])
	value, _ := napisys.CreateInt32(env, int32(len(params)) )

	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: add,
	}
	desc := napisys.Property{
		Name:   "add",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
