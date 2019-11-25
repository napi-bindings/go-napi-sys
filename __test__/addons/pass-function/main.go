package main

import (
	"C"
	"go-napi-sys/__test__/addons/pass-function/napisys"
	"unsafe"
)

func runCallback(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	params, _, _, _ := napisys.GetCbInfo(env, info)
	global, _ := napisys.GetGlobal(env)
	cb := params[0]
	str, _ := napisys.CreateStringUtf8(env, "hello world")
	arguments := []napisys.Value{str}
	res, _ := napisys.CallFunction(env, global, cb, arguments)
	return res
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	newExports, _ := napisys.CreateFunction((napisys.Env)(env), "", runCallback)
	return unsafe.Pointer(newExports)
}

func main() {}
