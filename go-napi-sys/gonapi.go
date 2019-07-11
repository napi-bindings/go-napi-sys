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
// The napi_status return value provides a VM-independent representation of the
// error which occurred. In some cases it is useful to be able to get more
// detailed information, including a string representing the error as well as
// VM (engine)-specific information.
// error_message: UTF8-encoded string containing a VM-neutral description of the
// error.
// engine_reserved: Reserved for VM-specific error details. This is currently
// not implemented for any VM.
// engine_error_code: VM-specific error code. This is currently not implemented
// for any VM.
// error_code: The N-API status code that originated with the last error.
// Do not rely on the content or format of any of the extended information as it
// is not subject to SemVer and may change at any time. It is intended only for
// logging purposes.
type NapiExtendedErrorInfo *C.napi_extended_error_info

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

// Error Handling
// N-API uses both return values and JavaScript exceptions for error handling.
// The following sections explain the approach for each case.
// All of the N-API functions share the same error handling pattern. The return
// type of all API functions is napi_status.
// The return value will be napi_ok if the request was successful and no uncaught
// JavaScript exception was thrown. If an error occurred AND an exception was
// thrown, the napi_status value for the error will be returned. If an exception
// was thrown, and no error occurred, napi_pending_exception will be returned.

// In cases where a return value other than napi_ok or napi_pending_exception is
// returned, napi_is_exception_pending must be called to check if an exception is
// pending. See the section on exceptions for more details.

// The napi_status return value provides a VM-independent representation of the
// error which occurred. In some cases it is useful to be able to get more
// detailed information, including a string representing the error as well as
// VM (engine)-specific information.

// Any N-API function call may result in a pending JavaScript exception. This is
// obviously the case for any function that may cause the execution of
// JavaScript, but N-API specifies that an exception may be pending on return
// from any of the API functions. If the napi_status returned by a function is
// napi_ok then no exception is pending and no additional action is required. If
// the napi_status returned is anything other than napi_ok or
// napi_pending_exception, in order to try to recover and continue instead of
// simply returning immediately, napi_is_exception_pending must be called in
// order to determine if an exception is pending or not. In many cases when an
// N-API function is called and an exception is already pending, the function
// will return immediately with a napi_status of napi_pending_exception.
// However, this is not the case for all functions. N-API allows a subset of the
// functions to be called to allow for some minimal cleanup before returning to
// JavaScript. In that case, napi_status will reflect the status for the
// function. It will not reflect previous pending exceptions. To avoid confusion,
// check the error status after every function call.

// When an exception is pending one of two approaches can be employed.:
// The first approach is to do any appropriate cleanup and then return so that
// execution will return to JavaScript. As part of the transition back to
// JavaScript the exception will be thrown at the point in the JavaScript code
// where the native method was invoked. The behavior of most N-API calls is
// unspecified while an exception is pending, and many will simply return
// napi_pending_exception, so it is important to do as little as possible and
// then return to JavaScript where the exception can be handled.
// The second approach is to try to handle the exception. There will be cases
// where the native code can catch the exception, take the appropriate action,
// and then continue. This is only recommended in specific cases where it is
// known that the exception can be safely handled.

// The Node.js project is adding error codes to all of the errors generated
// internally. The goal is for applications to use these error codes for all
// error checking. The associated error messages will remain, but will only be
// meant to be used for logging and display with the expectation that the message
// can change without SemVer applying. In order to support this model with N-API,
// both in internal functionality and for module specific functionality
// (as its good practice), the throw_ and create_ functions take an optional code
// parameter which is the string for the code to be added to the error object. If
// the optional parameter is NULL then no code will be associated with the error.

// NapiGetLastErrorInfo function ...
// [in] env: The environment that the API is invoked under.
// This API retrieves a napi_extended_error_info structure with information about
// the last error that occurred.
// The content of the napi_extended_error_info returned is only valid up until an
// n-api function is called on the same env.
// Do not rely on the content or format of any of the extended information as it
// is not subject to SemVer and may change at any time. It is intended only for
// logging purposes.
// The function can be called even if there is a pending JavaScript exception.
func NapiGetLastErrorInfo(env NapiEnv) (NapiExtendedErrorInfo, NapiStatus) {
	var res *C.napi_extended_error_info
	var status = C.napi_get_last_error_info(env, &res)
	return NapiExtendedErrorInfo(res), NapiStatus(status)
}

// NapiThrow function throws the JavaScript value provided.
// [in] env: The environment that the API is invoked under.
// [in] error: The JavaScript value to be thrown.
// N-API version: 1
func NapiThrow(env NapiEnv, error NapiValue) NapiStatus {
	return NapiStatus(C.napi_throw(env, error))
}

// NapiThrowError function throws a JavaScript Error with the text provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func NapiThrowError(env NapiEnv, msg string, code string) NapiStatus {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return NapiStatus(C.napi_throw_error(env, ccode, cmsg))
}

// NapiThrowTypeError function  throws a JavaScript TypeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func NapiThrowTypeError(env NapiEnv, msg string, code string) NapiStatus {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return NapiStatus(C.napi_throw_type_error(env, ccode, cmsg))
}

// NapiThrowRangError function throws a JavaScript RangeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional error code to be set on the error.
// [in] msg: C string representing the text to be associated with the error.
// N-API version: 1
func NapiThrowRangError(env NapiEnv, msg string, code string) NapiStatus {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	var ccode = C.CString(code)
	defer C.free(unsafe.Pointer(ccode))
	return NapiStatus(C.napi_throw_range_error(env, ccode, cmsg))
}

// NapiIsError function queries a napi_value to check if it represents an error
// object.
// [in] env: The environment that the API is invoked under.
// [in] value: The napi_value to be checked.
// Boolean value that is set to true if napi_value represents an error, false
// otherwise.
// N-API version: 1
func NapiIsError(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_error(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiCreateError function creates a JavaScript Error with the text provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func NapiCreateError(env NapiEnv, msg NapiValue, code NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_error(env, code, msg, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateTypeError function creates a JavaScript TypeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func NapiCreateTypeError(env NapiEnv, code NapiValue, msg NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_type_error(env, code, msg, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateRangeError function creates a JavaScript RangeError with the text
// provided.
// [in] env: The environment that the API is invoked under.
// [in] code: Optional napi_value with the string for the error code to be
// associated with the error.
// [in] msg: napi_value that references a JavaScript String to be used as the
// message for the Error.
// N-API version: 1
func NapiCreateRangeError(env NapiEnv, code NapiValue, msg NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_range_error(env, code, msg, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiGetAndClearLastException function returns true if an exception is pending.
// This function can be called even if there is a pending JavaScript exception.
// [in] env: The environment that the API is invoked under.
// The function returns the exception if one is pending, NULL otherwise.
// N-API version: 1
func NapiGetAndClearLastException(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_and_clear_last_exception(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiIsExceptionPending function ...
// [in] env: The environment that the API is invoked under.
// Boolean value that is set to true if an exception is pending.
// N-API version: 1
func NapiIsExceptionPending(env NapiEnv) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_exception_pending(env, &res)
	return bool(res), NapiStatus(status)
}

// NapiFatalException function triggers an 'uncaughtException' in JavaScript.
// Useful if an async callback throws an exception with no way to recover.
// [in] env: The environment that the API is invoked under.
// [in] err: The error that is passed to 'uncaughtException'.
// N-API version: 3
func NapiFatalException(env NapiEnv, error NapiValue) NapiStatus {
	return NapiStatus(C.napi_fatal_exception(env, error))
}

// NapiFatalError function thrown a fatal error o immediately terminate the
// process.
// [in] location: Optional location at which the error occurred.
// [in] location_len: The length of the location in bytes, or NAPI_AUTO_LENGTH
// if it is null-terminated.
// [in] message: The message associated with the error.
// [in] message_len: The length of the message in bytes, or NAPI_AUTO_LENGTH if
// it is null-terminated.
// This function can be called even if there is a pending JavaScript exception.
// The function call does not return, the process will be terminated.
// N-API version: 1
func NapiFatalError(location string, msg string) {
	clocation := C.CString(location)
	defer C.free(unsafe.Pointer(clocation))
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	C.napi_fatal_error(clocation, C.NAPI_AUTO_LENGTH, cmsg, C.NAPI_AUTO_LENGTH)
	return
}

// NapiOnpenHandleScope function ...
func NapiOnpenHandleScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiClosesHandleScope function ...
func NapiClosesHandleScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiOnpenEscapableHandleScope function ...
func NapiOnpenEscapableHandleScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiClosesEscapableHandleScope function ...
func NapiClosesEscapableHandleScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiEscapeHandle function ...
func NapiEscapeHandle(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateReference function ...
func NapiCreateReference(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDeleteReference function ...
func NapiDeleteReference(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiReferenceRef function ...
func NapiReferenceRef(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiReferenceUnref function ...
func NapiReferenceUnref(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetReferenceValue function ...
func NapiGetReferenceValue(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAddEnvCleanupHook function ...
func NapiAddEnvCleanupHook(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRemoveCleaupHook function ...
func NapiRemoveCleaupHook(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateArray function ...
func NapiCreateArray(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateArrayWithLength function ...
func NapiCreateArrayWithLength(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateArrayBuffer function ...
func NapiCreateArrayBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBuffer function ...
func NapiCreateBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBufferCopy function ...
func NapiCreateBufferCopy(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateExternal function ...
func NapiCreateExternal(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateExternalArrayBuffer function ...
func NapiCreateExternalArrayBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateExternalBuffer function ...
func NapiCreateExternalBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateObject function ...
func NapiCreateObject(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateSymbol function ...
func NapiCreateSymbol(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateTypedArray function ...
func NapiCreateTypedArray(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateDataview function ...
func NapiCreateDataview(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateInt32 function ...
func NapiCreateInt32(env NapiEnv, value int32) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_int32(env, C.int(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateUInt32 function ...
func NapiCreateUInt32(env NapiEnv, value uint32) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_uint32(env, C.uint(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateInt64 function ...
func NapiCreateInt64(env NapiEnv, value int64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_int64(env, C.int64_t(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateDouble function ...
func NapiCreateDouble(env NapiEnv, value float64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_double(env, C.double(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintInt64 function ...
func NapiCreateBigintInt64(env NapiEnv, value int64) (NapiValue, NapiStatus) {
	/*var res C.napi_value
	var status = C.napi_create_bigint_int64(env, C.int64_t(value), &res)
	return NapiValue(res), NapiStatus(status)*/
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintUInt64 function ...
func NapiCreateBigintUInt64(env NapiEnv, value uint64) (NapiValue, NapiStatus) {
	/*var res C.napi_value
	var status = C.napi_create_bigint_uint64(env, C.uint64_t(value), &res)
	return NapiValue(res), NapiStatus(status)*/
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintWords function ...
func NapiCreateBigintWords() (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringLatin1 function ...
func NapiCreateStringLatin1(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringUtf16 function ...
func NapiCreateStringUtf16(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringUtf8 function ...
func NapiCreateStringUtf8(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetArrayLength function ...
func NapiGetArrayLength(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetArrayBufferInfo function ...
func NapiGetArrayBufferInfo(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetPrototype function ...
func NapiGetPrototype(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetTypedArrayInfo function ...
func NapiGetTypedArrayInfo(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetDataviewInfo function ...
func NapiGetDataviewInfo(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueBool function ...
func NapiGetValueBool(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueDouble function ...
func NapiGetValueDouble(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueBigintInt64 function ...
func NapiGetValueBigintInt64(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueBigintUInt64 function ...
func NapiGetValueBigintUInt64(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueBigintWords function ...
func NapiGetValueBigintWords(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

//NapiGetValueExternal function ...
func NapiGetValueExternal(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueInt32 function ...
func NapiGetValueInt32(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueInt64 function ...
func NapiGetValueInt64(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueStringLatin1 function ...
func NapiGetValueStringLatin1(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueStringUtf8 function ...
func NapiGetValueStringUtf8(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueStringUtf16 function ...
func NapiGetValueStringUtf16(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetValueUint32 function ...
func NapiGetValueUint32(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetBoolean function ...
func NapiGetBoolean(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetGlobal function ...
func NapiGetGlobal(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNull function ...
func NapiGetNull(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetUndefined function ...
func NapiGetUndefined(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToBool function ...
func NapiCoerceToBool(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToNumber function ...
func NapiCoerceToNumber(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToObject function ...
func NapiCoerceToObject(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToString function ...
func NapiCoerceToString(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiTypeOf function ...
func NapiTypeOf(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiInstanceOf function ...
func NapiInstanceOf(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsArray function ...
func NapiIsArray(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsArrayBuffer function ...
func NapiIsArrayBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsBuffer function ...
func NapiIsBuffer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsTypedArray function ...
func NapiIsTypedArray(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsDataview function ...
func NapiIsDataview(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiStrictEquals function ...
func NapiStrictEquals(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetPropertyNames function ...
func NapiGetPropertyNames(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiSetProperty function ...
func NapiSetProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetProperty function ...
func NapiGetProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiHasProperty function ...
func NapiHasProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDeleteProperty function ...
func NapiDeleteProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiHasOwnProperty function ...
func NapiHasOwnProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiSetNamedProperty function ...
func NapiSetNamedProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNamedProperty function ...
func NapiGetNamedProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiHasNamedProperty function ...
func NapiHasNamedProperty(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiSetElement function ...
func NapiSetElement(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetElement function ...
func NapiGetElement(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiHasElement function ...
func NapiHasElement(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDeleteElement function ...
func NapiDeleteElement(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDefineProperties function ...s
func NapiDefineProperties(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCallFunction function ...s
func NapiCallFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateFunction function ...
func NapiCreateFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetCbInfo function ...
func NapiGetCbInfo(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNewTarget function ...
func NapiGetNewTarget(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiNewInstance function ...
func NapiNewInstance(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDefineClass function ...
func NapiDefineClass(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiWrap function ...
func NapiWrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiUnwrap function ...
func NapiUnwrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRemoveWrap function ...
func NapiRemoveWrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAddFinalizer function ...
func NapiAddFinalizer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateAsyncWork function ...
func NapiCreateAsyncWork(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiDeleteAsyncWork function ...
func NapiDeleteAsyncWork(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiQueueAsyncWork function ...
func NapiQueueAsyncWork(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCancelAsyncWork function ...
func NapiCancelAsyncWork(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAsyncInit function ...
func NapiAsyncInit(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAsyncDestroy function ...
func NapiAsyncDestroy(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiMakeCallback function ...
func NapiMakeCallback(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiOpenCallbackScope function ...
func NapiOpenCallbackScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCloseCallbackScope function ...
func NapiCloseCallbackScope(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNodeVersion function ...
func NapiGetNodeVersion(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetVersion function ...
func NapiGetVersion(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAdjustExternalMemory unction ...f
func NapiAdjustExternalMemory(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreatePromise function ...
func NapiCreatePromise(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiResolveDeferred function ...
func NapiResolveDeferred(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRejectDeferred function ...
func NapiRejectDeferred(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiIsPromise function ...
func NapiIsPromise(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRunScript function ...
func NapiRunScript(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetUvEventLoop function ...
func NapiGetUvEventLoop(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateThreadsafeFunction function ...
func NapiCreateThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetThreadsafeFunctionContext function ...
func NapiGetThreadsafeFunctionContext(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCallThreadsafeFunction function ...
func NapiCallThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAcquireThreadsafeFunction function ...
func NapiAcquireThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiReleaseThreadsafeFunction function ...
func NapiReleaseThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRefThreadsafeFunction function ...
func NapiRefThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiUnrefThreadsafeFunction function ...
func NapiUnrefThreadsafeFunction(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// Caller contains a callback to call
type Caller struct{}

func (s *Caller) cb(env C.napi_env, info C.napi_callback_info) C.napi_value {
	value, _ := NapiCreateInt32(NapiEnv(env), 7)
	return C.napi_value(value)
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
