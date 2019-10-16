package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo CXXFLAGS:  -I./napisys/include/
#cgo CFLAGS: -I./napisys/include/ -DNAPI_EXPERIMENTAL=1
#cgo LDFLAGS: -L./napisys/lib/ -lnode_api
#include <stdlib.h>
#include "napisys/gonapi.h"
#include <node_api.h>
*/
import "C"
import (
	"go-napi-sys/__test__/addon-example/napisys"
	"unsafe"
)

func createInt32(env napisys.NapiEnv, info napisys.NapiCallbackInfo) napisys.NapiValue {
	value, _ := napisys.NapiCreateInt32(napisys.NapiEnv(env), 7)
	return value
}

//export Initialize
func Initialize(env unsafe.Pointer, exports unsafe.Pointer) unsafe.Pointer {
	name := C.CString("createInt32")
	defer C.free(unsafe.Pointer(name))
	/*caller := &Caller{
		cb: createInt32,
	}*/
	desc := napisys.NapiPropertyDescriptor{
		utf8name:   name,
		name:       nil,
		method:     nil, //(C.napi_callback)(C.CallbackMethod(unsafe.Pointer(caller))), //nil,
		getter:     nil,
		setter:     nil,
		value:      nil,
		attributes: napisys.PropertyAttributes.NapiDefault,
		data:       nil,
	}
	//C.napi_define_properties(env, exports, 1, (*C.napi_property_descriptor)(&desc))
	props := []napisys.NapiPropertyDescriptor{desc}
	napisys.NapiDefineProperties((napisys.NapiEnv)(env), (napisys.NapiValue)(exports), props)

	return exports
}

func main() {}
