package main

import (
	"C"
	"go-napi-sys/__test__/addons/nanotime/napisys"
	"unsafe"
	"time"
)

func unixNano(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	now := time.Now()
	value, _ := napisys.CreateInt64(env, now.UnixNano())
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: unixNano,
	}
	desc := napisys.Property{
		Name:   "unixNano",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
