package main

import (
	"C"
	"go-napi-sys/__test__/addons/pass-function/napisys"
	"unsafe"
)

func runCallback(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	
}

func main() {}
