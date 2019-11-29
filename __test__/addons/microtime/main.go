package main

import (
	"C"
	"go-napi-sys/__test__/addons/microtime/napisys"
	"unsafe"
	"time"
)

func unix(env napisys.Env, info napisys.CallbackInfo) napisys.Value {
	now := time.Now()
	value, _ := napisys.CreateInt64(env, now.UnixNano())
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	caller := &napisys.Caller{
		Cb: unix,
	}
	desc := napisys.Property{
		Name:   "unix",
		Method: caller,
	}
	props := []napisys.Property{desc}
	napisys.DefineProperties((napisys.Env)(env), (napisys.Value)(exports), props)
	return exports
}

func main() {}
