package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo CXXFLAGS:  -I./include/
#cgo CFLAGS: -I./include/
#cgo LDFLAGS: -L./lib/ -lnode_api
#include <stdlib.h>
#include "gonapi.h"
#include <node_api.h>
*/
import "C"
import "unsafe"

// Aliases for JavaScript types

// NapiEnv represents ...
type NapiEnv C.napi_env

// NapiValue represents ...
type NapiValue C.napi_value

// NapiRef represents ...
type NapiRef C.napi_ref

// NapiHandleScope represents ...
type NapiHandleScope C.napi_handle_scope

// NapiEscapableHandleScope represents ...
type NapiEscapableHandleScope C.napi_escapable_handle_scope

// NapiCallbackInfo represents ...
type NapiCallbackInfo C.napi_callback_info

// NapiDeferred represents ...
type NapiDeferred C.napi_deferred

// NapiPropertyAttributes represents ...
type NapiPropertyAttributes C.napi_property_attributes

// NapiValueType represents ...
type NapiValueType C.napi_valuetype

// NapiTypedArrayType represents ...
type NapiTypedArrayType C.napi_typedarray_type

// NapiStatus represents ...
type NapiStatus C.napi_status

// NapiCallback represents ...
type NapiCallback C.napi_callback

// NapiFinalize represents ...
type NapiFinalize C.napi_finalize

// NapiPropertyDescriptor represents ...
type NapiPropertyDescriptor C.napi_property_descriptor

// NapiExtendedErrorInfo represents ...
type NapiExtendedErrorInfo C.napi_extended_error_info

// Aliases for types strickly connected with the runtime

// NapiCallbackScope represents ...
type NapiCallbackScope C.napi_callback_scope

// NapiAyncContext represents ...
type NapiAyncContext C.napi_async_context

// NapiAsyncWork represents ...
type NapiAsyncWork C.napi_async_work

// NapiThreadsafeFunction represents ...
type NapiThreadsafeFunction C.napi_threadsafe_function

// NapiTheradsafeFunctionReleaseMode represents ...
type NapiTheradsafeFunctionReleaseMode C.napi_threadsafe_function_release_mode

// NapiThreadsafeFunctionCallMode represents ...
type NapiThreadsafeFunctionCallMode C.napi_threadsafe_function_call_mode

// NapiAsyncExecuteCallback represents ...
type NapiAsyncExecuteCallback C.napi_async_execute_callback

// NapiAsyncCompleteCallback represents ...
type NapiAsyncCompleteCallback C.napi_async_complete_callback

// NapiThreadsafeFunctionCallJS represents ...
type NapiThreadsafeFunctionCallJS C.napi_threadsafe_function_call_js

// NapiNodeVersion represents ...
type NapiNodeVersion C.napi_node_version

// NapiGetLastErrorInfo function ...
func NapiGetLastErrorInfo() {}

// NapiThrow function ...
func NapiThrow() {}

// NapiThrowError function ...
func NapiThrowError() {}

// NapiThrowTypeError function ...
func NapiThrowTypeError() {}

// NapiThrowRangError function ...
func NapiThrowRangError() {}

// NapiIsError function ...
func NapiIsError() {}

// NapiCreateError function ...
func NapiCreateError() {}

// NapiCreateTypeError function ...
func NapiCreateTypeError() {}

// NapiCreateRangeError function ...
func NapiCreateRangeError() {}

// NapiGetAndClearLastException function ...
func NapiGetAndClearLastException() {}

// NapiIsExceptionPending function ...
func NapiIsExceptionPending() {}

// NapiFatalException function ...
func NapiFatalException() {}

// NapiFatalError function ...
func NapiFatalError() {}

// NapiOnpenHandleScope function ...
func NapiOnpenHandleScope() {}

// NapiClosesHandleScope function ...
func NapiClosesHandleScope() {}

// NapiOnpenEscapableHandleScope function ...
func NapiOnpenEscapableHandleScope() {}

// NapiClosesEscapableHandleScope function ...
func NapiClosesEscapableHandleScope() {}

// NapiEscapeHandle function ...
func NapiEscapeHandle() {}

// NapiCreateReference function ...
func NapiCreateReference() {}

// NapiDeleteReference function ...
func NapiDeleteReference() {}

// NapiReferenceRef function ...
func NapiReferenceRef() {}

// NapiReferenceUnref function ...
func NapiReferenceUnref() {}

// NapiGetReferenceValue function ...
func NapiGetReferenceValue() {}

// NapiAddEnvCleanupHook function ...
func NapiAddEnvCleanupHook() {}

// NapiRemoveCleaupHook function ...
func NapiRemoveCleaupHook() {}

// NapiCreateArray function ...
func NapiCreateArray() {}

// NapiCreateArrayWithLength function ...
func NapiCreateArrayWithLength() {}

// NapiCreateArrayBuffer function ...
func NapiCreateArrayBuffer() {}

// NapiCreateBuffer function ...
func NapiCreateBuffer() {}

// NapiCreateBufferCopy function ...
func NapiCreateBufferCopy() {}

// NapiCreateExternal function ...
func NapiCreateExternal() {}

// NapiCreateExternalArrayBuffer function ...
func NapiCreateExternalArrayBuffer() {}

// NapiCreateExternalBuffer function ...
func NapiCreateExternalBuffer() {}

// NapiCreateObject function ...
func NapiCreateObject() {}

// NapiCreateSymbol function ...
func NapiCreateSymbol() {}

// NapiCreateTypedArray function ...
func NapiCreateTypedArray() {}

// NapiCreateDataview function ...
func NapiCreateDataview() {}

// NapiCreateInt32 function ...
func NapiCreateInt32() {}

// NapiCreateUInt32 function ...
func NapiCreateUInt32() {}

// NapiCreateInt64 function ...
func NapiCreateInt64() {}

// NapiCreateDouble function ...
func NapiCreateDouble() {}

// NapiCreateBigintInt64 function ...
func NapiCreateBigintInt64() {}

// NapiCreateBigintUInt64 function ...
func NapiCreateBigintUInt64() {}

// NapiCreateBigintWords function ...
func NapiCreateBigintWords() {}

// NapiCreateStringLatin1 function ...
func NapiCreateStringLatin1() {}

// NapiCreateStringUtf16 function ...
func NapiCreateStringUtf16() {}

// NapiCreateStringUtf8 function ...
func NapiCreateStringUtf8() {}

// NapiGetArrayLength function ...
func NapiGetArrayLength() {}

// NapiGetArrayBufferInfo function ...
func NapiGetArrayBufferInfo() {}

// NapiGetPrototype function ...
func NapiGetPrototype() {}

// NapiGetTypedArrayInfo function ...
func NapiGetTypedArrayInfo() {}

// NapiGetDataviewInfo function ...
func NapiGetDataviewInfo() {}

// NapiGetValueBool function ...
func NapiGetValueBool() {}

// NapiGetValueDouble function ...
func NapiGetValueDouble() {}

// NapiGetValueBigintInt64 function ...
func NapiGetValueBigintInt64() {}

// NapiGetValueBigintUInt64 function ...
func NapiGetValueBigintUInt64() {}

// NapiGetValueBigintWords function ...
func NapiGetValueBigintWords() {}

//NapiGetValueExternal function ...
func NapiGetValueExternal() {}

// NapiGetValueInt32 function ...
func NapiGetValueInt32() {}

// NapiGetValueInt64 function ...
func NapiGetValueInt64() {}

// NapiGetValueStringLatin1 function ...
func NapiGetValueStringLatin1() {}

// NapiGetValueStringUtf8 function ...
func NapiGetValueStringUtf8() {}

// NapiGetValueStringUtf16 function ...
func NapiGetValueStringUtf16() {}

// NapiGetValueUint32 function ...
func NapiGetValueUint32() {}

// NapiGetBoolean function ...
func NapiGetBoolean() {}

// NapiGetGlobal function ...
func NapiGetGlobal() {}

// NapiGetNull function ...
func NapiGetNull() {}

// NapiGetUndefined function ...
func NapiGetUndefined() {}

// NapiCoerceToBool function ...
func NapiCoerceToBool() {}

// NapiCoerceToNumber function ...
func NapiCoerceToNumber() {}

// NapiCoerceToObject function ...
func NapiCoerceToObject() {}

// NapiCoerceToString function ...
func NapiCoerceToString() {}

// NapiTypeOf function ...
func NapiTypeOf() {}

// NapiInstanceOf function ...
func NapiInstanceOf() {}

// NapiIsArray function ...
func NapiIsArray() {}

// NapiIsArrayBuffer function ...
func NapiIsArrayBuffer() {}

// NapiIsBuffer function ...
func NapiIsBuffer() {}

// NapiIsTypedArray function ...
func NapiIsTypedArray() {}

// NapiIsDataview function ...
func NapiIsDataview() {}

// NapiStrictEquals function ...
func NapiStrictEquals() {}

// NapiGetPropertyNames function ...
func NapiGetPropertyNames() {}

// NapiSetProperty function ...
func NapiSetProperty() {}

// NapiGetProperty function ...
func NapiGetProperty() {}

// NapiHasProperty function ...
func NapiHasProperty() {}

// NapiDeleteProperty function ...
func NapiDeleteProperty() {}

// NapiHasOwnProperty function ...
func NapiHasOwnProperty() {}

// NapiSetNamedProperty function ...
func NapiSetNamedProperty() {}

// NapiGetNamedProperty function ...
func NapiGetNamedProperty() {}

// NapiHasNamedProperty function ...
func NapiHasNamedProperty() {}

// NapiSetElement function ...
func NapiSetElement() {}

// NapiGetElement function ...
func NapiGetElement() {}

// NapiHasElement function ...
func NapiHasElement() {}

// NapiDeleteElement function ...
func NapiDeleteElement() {}

// NapiDefineProperties function ...s
func NapiDefineProperties() {}

// NapiCallFunction function ...s
func NapiCallFunction() {}

// NapiCreateFunction function ...
func NapiCreateFunction() {}

// NapiGetCbInfo function ...
func NapiGetCbInfo() {}

// NapiGetNewTarget function ...
func NapiGetNewTarget() {}

// NapiNewInstance function ...
func NapiNewInstance() {}

// NapiDefineClass function ...
func NapiDefineClass() {}

// NapiWrap function ...
func NapiWrap() {}

// NapiUnwrap function ...
func NapiUnwrap() {}

// NapiRemoveWrap function ...
func NapiRemoveWrap() {}

// NapiAddFinalizer function ...
func NapiAddFinalizer() {}

// NapiCreateAsyncWork function ...
func NapiCreateAsyncWork() {}

// NapiDeleteAsyncWork function ...
func NapiDeleteAsyncWork() {}

// NapiQueueAsyncWork function ...
func NapiQueueAsyncWork() {}

// NapiCancelAsyncWork function ...
func NapiCancelAsyncWork() {}

// NapiAsyncInit function ...
func NapiAsyncInit() {}

// NapiAsyncDestroy function ...
func NapiAsyncDestroy() {}

// NapiMakeCallback function ...
func NapiMakeCallback() {}

// NapiOpenCallbackScope function ...
func NapiOpenCallbackScope() {}

// NapiCloseCallbackScope function ...
func NapiCloseCallbackScope() {}

// NapiGetNodeVersion function ...
func NapiGetNodeVersion() {}

// NapiGetVersion function ...
func NapiGetVersion() {}

// NapiAdjustExternalMemory unction ...f
func NapiAdjustExternalMemory() {}

// NapiCreatePromise function ...
func NapiCreatePromise() {}

// NapiResolveDeferred function ...
func NapiResolveDeferred() {}

// NapiRejectDeferred function ...
func NapiRejectDeferred() {}

// NapiIsPromise function ...
func NapiIsPromise() {}

// NapiRunScript function ...
func NapiRunScript() {}

// NapiGetUvEventLoop function ...
func NapiGetUvEventLoop() {}

// NapiCreateThreadsafeFunction function ...
func NapiCreateThreadsafeFunction() {}

// NapiGetThreadsafeFunctionContext function ...
func NapiGetThreadsafeFunctionContext() {}

// NapiCallThreadsafeFunction function ...
func NapiCallThreadsafeFunction() {}

// NapiAcquireThreadsafeFunction function ...
func NapiAcquireThreadsafeFunction() {}

// NapiReleaseThreadsafeFunction function ...
func NapiReleaseThreadsafeFunction() {}

// NapiRefThreadsafeFunction function ...
func NapiRefThreadsafeFunction() {}

// NapiUnrefThreadsafeFunction function ...
func NapiUnrefThreadsafeFunction() {}

// Caller contains a callback to call
type Caller struct{}

func (s *Caller) cb(env C.napi_env, info C.napi_callback_info) C.napi_value {
	var res C.napi_value
	C.napi_create_int32(env, C.int(5), &res)
	return res
}

//export ExecuteCallback
func ExecuteCallback(data unsafe.Pointer, env C.napi_env, info C.napi_callback_info) C.napi_value {
	caller := (*Caller)(data)
	return caller.cb(env, info)
}

//export Initialize
func Initialize(env C.napi_env, exports C.napi_value) C.napi_value {
	name := C.CString("createInt32")
	defer C.free(unsafe.Pointer(name))
	caller := &Caller{}
	desc := C.napi_property_descriptor{
		utf8name:   name,
		name:       nil,
		method:     (C.napi_callback)(C.CallbackMethod(unsafe.Pointer(&caller))), //nil,
		getter:     nil,
		setter:     nil,
		value:      nil,
		attributes: C.napi_default,
		data:       nil,
	}
	C.napi_define_properties(env, exports, 1, &desc)
	return exports
}

func main() {}
