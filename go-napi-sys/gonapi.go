package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo CXXFLAGS:  -I./include/
#cgo CFLAGS: -I./include/ -DNAPI_EXPERIMENTAL=1
#cgo LDFLAGS: -L./lib/ -lnode_api
#include <stdlib.h>
#include "gonapi.h"
#include <node_api.h>
*/
import "C"
import (
	"bytes"
	"unsafe"
)

// ccstring function transforms a Go string into a C string (array of characters)
// and returns the pointer to the first element.
func cstring(s string) unsafe.Pointer {
	p := make([]byte, len(s)+1)
	copy(p, s)
	return unsafe.Pointer(&p[0])
}

// Aliases for JavaScript types
// Basic N-API Data Types
// N-API exposes the following fundamental datatypes as abstractions that are
// consumed by the various APIs. These APIs should be treated as opaque,
// introspectable only with other N-API calls.

// NapiEnv is used to represent a context that the underlying N-API
// implementation can use to persist VM-specific state. This structure is passed
// to native functions when they're invoked, and it must be passed back when
// making N-API calls. Specifically, the same napi_env that was passed in when
// the initial native function was called must be passed to any subsequent nested
// N-API calls. Caching the napi_env for the purpose of general reuse is not
// allowed.
type NapiEnv C.napi_env

// NapiValue is an opaque pointer that is used to represent a JavaScript value.
type NapiValue C.napi_value

// NapiRef is an abstraction to use to reference a NapiValue. This allows for
// users to manage the lifetimes of JavaScript values, including defining their
// minimum lifetimes explicitly.
type NapiRef C.napi_ref

// NapiHandleScope is an abstraction used to control and modify the lifetime of
// objects created within a particular scope. In general, N-API values are
// created within the context of a handle scope. When a native method is called
// from JavaScript, a default handle scope will exist. If the user does not
// explicitly create a new handle scope, N-API values will be created in the
// default handle scope. For any invocations of code outside the execution of a
// native method (for instance, during a libuv callback invocation), the module
// is required to create a scope before invoking any functions that can result
// in the creation of JavaScript values.
// Handle scopes are created using NapiOnpenHandleScope and are destroyed using
// NapiCloseHandleScope. Closing the scope can indicate to the GC that all
// NapiValues created during the lifetime of the handle scope are no longer
// referenced from the current stack frame.
type NapiHandleScope C.napi_handle_scope

// NapiEscapableHandleScope represents a special type of handle scope to return
// values created within a particular handle scope to a parent scope.
type NapiEscapableHandleScope C.napi_escapable_handle_scope

// NapiCallbackInfo is an opaque datatype that is passed to a callback function.
// It can be used for getting additional information about the context in which
// the callback was invoked.
type NapiCallbackInfo C.napi_callback_info

// NapiDeferred known sa "deferred" object is created and returned alongside a
// Promise. The deferred object is bound to the created Promise and is the only
// means to resolve or reject the Promise. The  deferred object will be
// automatically freed on rejection or on resolving the Promise.
type NapiDeferred C.napi_deferred

// This is a struct used as container for types of property atrributes.
type propertyAttributes struct {
	NapiDefault      int
	NapiWritable     int
	NapiEnumerable   int
	NapiConfigurable int
	// Used with napi_define_class to distinguish static properties
	// from instance properties. Ignored by napi_define_properties.
	NapiStatic int
}

// PropertyAttributes contains the flags to control the  behavior of properties
// set on a JavaScript object. They can be one or more of the following bitflags:
// - NapiDefault - Used to indicate that no explicit attributes are set on the
// given property. By default, a property is read only, not enumerable and not
// configurable.
// - NapiWritable - Used to indicate that a given property is writable.
// - NapiEnumerable - Used to indicate that a given property is enumerable.
// - NapiConfigurable -  Used to indicate that a given property is configurable,
// as defined in Section 6.1.7.1 of the ECMAScript Language Specification.
// - NapiStatic - Used to indicate that the property will be defined as a static
// property on a class as opposed to an instance property, which is the default.
// This is used only by NapiDefineClass. It is ignored by NapiDfineProperties.
var PropertyAttributes = &propertyAttributes{
	NapiDefault:      C.napi_default,
	NapiWritable:     C.napi_writable,
	NapiEnumerable:   C.napi_enumerable,
	NapiConfigurable: C.napi_configurable,
	NapiStatic:       C.napi_static,
}

// NapiPropertyAttributes represents the flags used to control the behavior of
// properties set on a JavaScript object.
// Other than napi_static they correspond to the attributes listed in
// Section 6.1.7.1 of the ECMAScript Language Specification.
// Currently they can be one or more of the following bitflags:
// napi_default - Used to indicate that no explicit attributes are set on the
// given property. By default, a property is read only, not enumerable and not
// configurable.
// napi_writable - Used to indicate that a given property is writable.
// napi_enumerable - Used to indicate that a given property is enumerable.
// napi_configurable - Used to indicate that a given property is configurable,
// as defined in Section 6.1.7.1 of the ECMAScript Language Specification.
// napi_static - Used to indicate that the property will be defined as a static
// property on a class as opposed to an instance property, which is the default.
// This is used only by NapiDefineClass. It is ignored by NapiDefineProperties.
type NapiPropertyAttributes C.napi_property_attributes

// This is a struct used as container for types of NapiValue.
type valueType struct {
	// ES6 types (corresponds to typeof)
	NapiUndefined int
	NapiNull      int
	NapiBoolean   int
	NapiNumber    int
	NapiString    int
	NapiSymbol    int
	NapiObject    int
	NapiFunction  int
	NapiExternal  int
	NapiBigint    int
}

// ValueType contains the type of a NapiValue. This generally corresponds to the
// types described in Section 6.1 of the ECMAScript Language Specification. In
// addition to types in that section, ValueType can also represent Functions and
// Objects with external data. A JavaScript value of type NapiExternal appears in
// JavaScript as a plain object such that no properties can be set on it, and no
//prototype.
var ValueType = &valueType{
	NapiUndefined: C.napi_undefined,
	NapiNull:      C.napi_null,
	NapiBoolean:   C.napi_boolean,
	NapiNumber:    C.napi_number,
	NapiString:    C.napi_string,
	NapiSymbol:    C.napi_symbol,
	NapiObject:    C.napi_object,
	NapiFunction:  C.napi_function,
	NapiExternal:  C.napi_external,
	NapiBigint:    C.napi_bigint,
}

// NapiValueType describes the type of NapiValue. This generally corresponds to
// the types described in Section 6.1 of the ECMAScript Language Specification.
// In addition to types in that section, NapiValueType can also represent
// Functions and Objects with external data.
// A JavaScript value of type napi_external appears in JavaScript as a plain
// object such that no properties can be set on it, and no prototype.
// Currently the following types are supported:
//  napi_undefined,
//  napi_null,
//  napi_boolean,
//  napi_number,
//  napi_string,
//  napi_symbol,
//  napi_object,
//  napi_function,
//  napi_external,
//  napi_bigint,
type NapiValueType C.napi_valuetype

// This is a struct used as container for types used in TypedArray.
type typedArrayType struct {
	NapiInt8Array         int
	NapiUInt8Array        int
	NapiUInt8ClampedArray int
	NapiInt16Array        int
	NapiUInt16Array       int
	NapiInt32Array        int
	NapiUInt32Array       int
	NapiFloat32Array      int
	NapiFloat64Array      int
	NapiBigInt64Array     int
	NapiBigUInt64Array    int
}

// TypedArrayType contains the underlying binary scalar datatype of the
// TypedArray defined in sectiontion 22.2 of the ECMAScript Language
// Specification.
var TypedArrayType = &typedArrayType{
	NapiInt8Array:         C.napi_int8_array,
	NapiUInt8Array:        C.napi_uint8_array,
	NapiUInt8ClampedArray: C.napi_uint8_clamped_array,
	NapiInt16Array:        C.napi_int16_array,
	NapiUInt16Array:       C.napi_uint16_array,
	NapiInt32Array:        C.napi_int32_array,
	NapiUInt32Array:       C.napi_uint32_array,
	NapiFloat32Array:      C.napi_float32_array,
	NapiFloat64Array:      C.napi_float64_array,
	NapiBigInt64Array:     C.napi_bigint64_array,
	NapiBigUInt64Array:    C.napi_biguint64_array,
}

// NapiTypedArrayType represents the underlying binary scalar datatype of the
// TypedArray defined in sectiontion 22.2 of the ECMAScript Language
// Specification.
type NapiTypedArrayType C.napi_typedarray_type

// This is a struct used as container for N-API status.
type status struct {
	NapiOK                    int
	NapiInvalidArg            int
	NapiObjectExpected        int
	NapiStringExpected        int
	NapiNameExpected          int
	NapiFunctionExpected      int
	NapiNumberExpected        int
	NapiBooleanExpected       int
	NapiArrayExpected         int
	NapiGenericFailure        int
	NapiPendingException      int
	NapiCancelled             int
	NapiEscapeCalledTwice     int
	NapiHandleScopeMismatch   int
	NapiCallbackScopeMismatch int
	NapiQueueFull             int
	NapiClosing               int
	NapiBigintExpected        int
	NapiDateExpected          int
}

// Status contains the status code indicating the success or failure of
// a N-API call. Currently, the following status codes are supported:
//  napi_ok
//  napi_invalid_arg
//  napi_object_expected
//  napi_string_expected
//  napi_name_expected
//  napi_function_expected
//  napi_number_expected
//  napi_boolean_expected
//  napi_array_expected
//  napi_generic_failure
//  napi_pending_exception
//  napi_cancelled
//  napi_escape_called_twice
//  napi_handle_scope_mismatch
//  napi_callback_scope_mismatch
//  napi_queue_full
//  napi_closing
//  napi_bigint_expected
//  napi_date_expected
// If additional information is required upon an API returning a failed status,
// it can be obtained by calling NapiGetLastErrorInfo.
var Status = &status{
	NapiOK:                    C.napi_ok,
	NapiInvalidArg:            C.napi_invalid_arg,
	NapiObjectExpected:        C.napi_object_expected,
	NapiStringExpected:        C.napi_string_expected,
	NapiNameExpected:          C.napi_name_expected,
	NapiFunctionExpected:      C.napi_function_expected,
	NapiNumberExpected:        C.napi_number_expected,
	NapiBooleanExpected:       C.napi_boolean_expected,
	NapiArrayExpected:         C.napi_array_expected,
	NapiGenericFailure:        C.napi_generic_failure,
	NapiPendingException:      C.napi_pending_exception,
	NapiCancelled:             C.napi_cancelled,
	NapiEscapeCalledTwice:     C.napi_escape_called_twice,
	NapiHandleScopeMismatch:   C.napi_handle_scope_mismatch,
	NapiCallbackScopeMismatch: C.napi_callback_scope_mismatch,
	NapiQueueFull:             C.napi_queue_full,
	NapiClosing:               C.napi_closing,
	NapiBigintExpected:        C.napi_bigint_expected,
	NapiDateExpected:          C.napi_date_expected,
}

// NapiStatus represent the status code indicating the success or failure of
// a N-API call. Currently, the following status codes are supported:
//  napi_ok
//  napi_invalid_arg
//  napi_object_expected
//  napi_string_expected
//  napi_name_expected
//  napi_function_expected
//  napi_number_expected
//  napi_boolean_expected
//  napi_array_expected
//  napi_generic_failure
//  napi_pending_exception
//  napi_cancelled
//  napi_escape_called_twice
//  napi_handle_scope_mismatch
//  napi_callback_scope_mismatch
//  napi_queue_full
//  napi_closing
//  napi_bigint_expected
// If additional information is required upon an API returning a failed status,
// it can be obtained by calling NapiGetLastErrorInfo.
type NapiStatus C.napi_status

// NapiCallback represents a function pointer type for user-provided native
// functions which are to be exposed to JavaScript via N-API. Callback functions
// should satisfy the following signature:
// typedef napi_value (*napi_callback)(napi_env, napi_callback_info);
type NapiCallback C.napi_callback

// NapiFinalize represents a function pointer type for add-on provided functions
// that allow the user to be notified when externally-owned data is ready to be
// cleaned up because the object with which it was associated with, has been
// garbage-collected. The user must provide a function satisfying the following
// signature which would get called upon the object's collection. Currently,
// `napi_finalize` can be used for finding out when objects that have external
// data are collected. Finalize functions hould satisfy the following signature:
// typedef void (*napi_finalize)(napi_env env,
//								 void* finalize_data,
//								 void* finalize_hint);
type NapiFinalize C.napi_finalize

// NapiPropertyDescriptor is a data structure that used to define the properties
// of a JavaScript object.
type NapiPropertyDescriptor C.napi_property_descriptor

// NapiExtendedErrorInfo contains additional information about a failed status
// happened on an N-API call.
// The NapiStatus return value provides a VM-independent representation of the
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

// NapiCallbackScope represents
type NapiCallbackScope C.napi_callback_scope

// NapiAyncContext represents the context for the async operation that is
// invoking a callback. This should normally be a value previously obtained from
// `napi_async_init`. However `NULL` is also allowed, which indicates the current
// async context (if any) is to be used for the callback.
type NapiAyncContext C.napi_async_context

// NapiAsyncWork represents the handle for the newly created asynchronous work
// and it is used to execute logic asynchronously.
type NapiAsyncWork C.napi_async_work

// NapiThreadsafeFunction is an opaque pointer that represents a JavaScript
// function which can be called asynchronously from multiple threads.
type NapiThreadsafeFunction C.napi_threadsafe_function

// This is a struct used as container for modes to release a
// NapiThreadSafeFunction.
type tsfnReleaseMode struct {
	NapiTsfnRelease int
	NapiTsfnAbort   int
}

// TsfnReleaseMode contains values to be given to NapiReleaseThreadsafeFunction()
// to indicate whether the thread-safe function is to be closed immediately
// (NapiTsfnAbort) or merely released (NapiTsfnRelease) and thus available for
// subsequent use via NapiAcquireThreadsafeFunction() and
// NapiCallThreadsafeFunction().
var TsfnReleaseMode = &tsfnReleaseMode{
	NapiTsfnRelease: C.napi_tsfn_release,
	NapiTsfnAbort:   C.napi_tsfn_abort,
}

// NapiTheradsafeFunctionReleaseMode represents a value to be given to
// NapiReleaseThreadsafeFunction() to indicate whether the thread-safe function
// is to be closed immediately (NapiTsfnAbort) or merely released
// (NapiTsfnRelease) and thus available for subsequent use via
// NapiAcquireThreadsafeFunction() and NapiCallThreadsafeFunction().
type NapiTheradsafeFunctionReleaseMode C.napi_threadsafe_function_release_mode

// This is a struct used as container for types used to call a
// NapiThreadSafeFunction.
type tsfnCallMode struct {
	NapiTsfnNonBlocking int
	NapiTsfnBlocking    int
}

// TsfnCallMode contains values to be given to NapiCallThreadsafeFunction() to
// indicate whether the call should block whenever the queue associated with the
// thread-safe function is full.
var TsfnCallMode = &tsfnCallMode{
	NapiTsfnNonBlocking: C.napi_tsfn_nonblocking,
	NapiTsfnBlocking:    C.napi_tsfn_blocking,
}

// NapiThreadsafeFunctionCallMode contains values used to indicate whether the
// call should block whenever the queue associated with the thread-safe function
// is full.
type NapiThreadsafeFunctionCallMode C.napi_threadsafe_function_call_mode

// NapiAsyncExecuteCallback is a function pointer used with functions that
// support asynchronous operations. Callback functions must statisfy the
// following signature:
// typedef void (*napi_async_execute_callback)(napi_env env, void* data);
// Implementations of this type of function should avoid making any N-API calls
// that could result in the execution of JavaScript or interaction with
// JavaScript objects.
type NapiAsyncExecuteCallback C.napi_async_execute_callback

// NapiAsyncCompleteCallback is a function pointer used with functions that
// support asynchronous operations. Callback functions must statisfy the
// following signature:
// typedef void (*napi_async_complete_callback)(napi_env env,
//												napi_status status,
//												void* data);
type NapiAsyncCompleteCallback C.napi_async_complete_callback

// NapiThreadsafeFunctionCallJS is a function pointer used with asynchronous
// thread-safe function calls. The callback will be called on the main thread.
// Its purpose is to use a data item arriving via the queue from one of the
// secondary threads to construct the parameters necessary for a call into
// JavaScript.
// Callback functions must satisfy the following signature:
// typedef void (*napi_threadsafe_function_call_js)(napi_env env,
//													napi_value js_callback,
//													void* context,
//													void* data);
type NapiThreadsafeFunctionCallJS C.napi_threadsafe_function_call_js

// NapiNodeVersion is a structure that contains informations about the version
// of Node.js instance.
// Currently, the following fields are exposed:
//  major
//  minor
//  patch
//  release
type NapiNodeVersion *C.napi_node_version

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

// NapiGetLastErrorInfo function returns the information for the last N-API call
// that was made.
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

// Object Lifetime management
// As N-API calls are made, handles to objects in the heap for the underlying VM
// may be returned as napi_values. These handles must hold the objects 'live'
// until they are no longer required by the native code, otherwise the objects
// could be collected before the native code was finished using them. As object
// handles are returned they are associated with a 'scope'. The lifespan for the
// default scope is tied to the lifespan of the native method call. The result is
// that, by default, handles remain valid and the objects associated with these
// handles will be held live for the lifespan of the native method call. In many
// cases, however, it is necessary that the handles remain valid for either a
// shorter or longer lifespan than that of the native method.
// N-API only supports a single nested hierarchy of scopes. There is only one
// active scope at any time, and all new handles will be associated with that
// scope while it is active. Scopes must be closed in the reverse order from
// which they are opened. In addition, all scopes created within a native method
// must be closed before returning from that method.
// When nesting scopes, there are cases where a handle from an inner scope needs
// to live beyond the lifespan of that scope. N-API supports an 'escapable scope'
// in order to support this case. An escapable scope allows one handle to be
// 'promoted' so that it 'escapes' the current scope and the lifespan of the
// handle changes from the current scope to that of the outer scope.

// NapiOnpenHandleScope function opens a new scope.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func NapiOnpenHandleScope(env NapiEnv) (NapiHandleScope, NapiStatus) {
	var res C.napi_handle_scope
	var status = C.napi_open_handle_scope(env, &res)
	return NapiHandleScope(res), NapiStatus(status)
}

// NapiCloseHandleScope function closes the scope passed in. Scopes must be
// closed in the reverse order from which they were created.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the scope to be closed.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func NapiCloseHandleScope(env NapiEnv, scope NapiHandleScope) NapiStatus {
	return NapiStatus(C.napi_close_handle_scope(env, scope))
}

// NapiOnpenEscapableHandleScope function opens a new scope from which one object
// can be promoted to the outer scope.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func NapiOnpenEscapableHandleScope(env NapiEnv) (NapiEscapableHandleScope, NapiStatus) {
	var res C.napi_escapable_handle_scope
	var status = C.napi_open_escapable_handle_scope(env, &res)
	return NapiEscapableHandleScope(res), NapiStatus(status)
}

// NapiCloseEscapableHandleScope function closes the scope passed in. Scopes must
// be closed in the reverse order from which they were created.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the scope to be closed.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func NapiCloseEscapableHandleScope(env NapiEnv, scope NapiEscapableHandleScope) NapiStatus {
	return NapiStatus(C.napi_close_escapable_handle_scope(env, scope))
}

// NapiEscapeHandle function promotes the handle to the JavaScript object so that
// it is valid for the lifetime of the outer scope. It can only be called once
// per scope. If it is called more than once an error will be returned.
// [in] env: The environment that the API is invoked under.
// [in] scope: napi_value representing the current scope.
// [in] escapee: napi_value representing the JavaScript Object to be escaped.
// This function can be called even if there is a pending JavaScript exception.
// N-API version: 1
func NapiEscapeHandle(env NapiEnv, scope NapiEscapableHandleScope, escapee NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_escape_handle(env, scope, escapee, &res)
	return NapiValue(res), NapiStatus(status)
}

// References to objects with a lifespan longer than that of the native method
// In some cases an addon will need to be able to create and reference objects
// with a lifespan longer than that of a single native method invocation.
// For example, to create a constructor and later use that constructor in a
// request to creates instances, it must be possible to reference the constructor
// object across many different instance creation requests. This would not be
// possible with a normal handle returned as a NapiValue as described in the
// earlier section. The lifespan of a normal handle is managed by scopes and all
// scopes must be closed before the end of a native method.

// N-API provides methods to create persistent references to an object. Each
// persistent reference has an associated count with a value of 0 or higher. The
// count determines if the reference will keep the corresponding object live.
// References with a count of 0 do not prevent the object from being collected
// and are often called 'weak' references. Any count greater than 0 will prevent
// the object from being collected.

// References must be deleted once they are no longer required by the addon.
// When a reference is deleted it will no longer prevent the corresponding object
// from being collected. Failure to delete a persistent reference will result in
// a 'memory leak' with both the native memory for the persistent reference and
// the corresponding object on the heap being retained forever.

// There can be multiple persistent references created which refer to the same
// object, each of which will either keep the object live or not based on its
// individual count.

// NapiCreateReference function creates a new reference with the specified
// reference count to the Object passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing the Object to which we want a reference.
// [in] initial_refcount: Initial reference count for the new reference.
// N-API version: 1
func NapiCreateReference(env NapiEnv, value NapiValue, refCount uint) (NapiRef, NapiStatus) {
	var res C.napi_ref
	var status = C.napi_create_reference(env, value, C.uint(refCount), &res)
	return NapiRef(res), NapiStatus(status)
}

// NapiDeleteReference function deletes the reference passed in.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref to be deleted.
// This API can be called even if there is a pending JavaScript exception.
// N-API version: 1
func NapiDeleteReference(env NapiEnv, ref NapiRef) NapiStatus {
	var status = C.napi_delete_reference(env, ref)
	return NapiStatus(status)
}

// NapiReferenceRef function  increments the reference count for the reference
// passed in and returns the resulting reference count.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which the reference count will be incremented.
// N-API version: 1
func NapiReferenceRef(env NapiEnv, ref NapiRef) (uint, NapiStatus) {
	var res C.uint
	var status = C.napi_reference_ref(env, ref, &res)
	return uint(res), NapiStatus(status)
}

// NapiReferenceUnref function ecrements the reference count for the reference
// passed in and returns the resulting reference count.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which the reference count will be decremented.
// N-API version: 1
func NapiReferenceUnref(env NapiEnv, ref NapiRef) (uint, NapiStatus) {
	var res C.uint
	var status = C.napi_reference_unref(env, ref, &res)
	return uint(res), NapiStatus(status)
}

// NapiGetReferenceValue function returns the NapiValue representing the
// JavaScript Object associated with the NapiRef. Otherwise, result will be
// NULL.
// [in] env: The environment that the API is invoked under.
// [in] ref: napi_ref for which we requesting the corresponding Object.
// N-API version: 1
func NapiGetReferenceValue(env NapiEnv, ref NapiRef) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_reference_value(env, ref, &res)
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

// NapiCreateArray function returns an N-API value corresponding to a JavaScript
// Array type. JavaScript arrays are described in Section 22.1 of the ECMAScript
// Language Specification.
// [in] env: The environment that the N-API call is invoked under.
// N-API version: 1
func NapiCreateArray(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_array(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateArrayWithLength function returns an N-API value corresponding to a
// JavaScript Array type. The Array's length property is set to the passed-in
// length parameter. However, the underlying buffer is not guaranteed to be
// pre-allocated by the VM when the array is created - that behavior is left to
// the underlying VM implementation.
// avaScript arrays are described in Section 22.1 of the ECMAScript Language
// Specification.
// [in] env: The environment that the API is invoked under.
// [in] length: The initial length of the Array.
// N-API version: 1
func NapiCreateArrayWithLength(env NapiEnv, length uint) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_array_with_length(env, C.size_t(length), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateArrayBuffer function returns N-API value corresponding to a
// JavaScript `ArrayBuffer`. ArrayBuffer is a data stucture used to represent
// fixed-length binary data buffers. They are normally used as backing-buffer for
// `TypedArray` objects. The ArrayBuffer allocated will have an underlying byte
// buffer whose size is determined by the length parameter that's passed in. The
// underlying buffer is optionally returned back to the caller in case the caller
// wants to directly manipulate the buffer.
// This buffer can only be written to directly from native code.
// To write to this buffer from JavaScript, a typed array or DataView object
// would need to be created.
// JavaScript ArrayBuffer objects are described in Section 24.1 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] length: The length in bytes of the array buffer to create.
// [out] data: Pointer to the underlying byte buffer of the ArrayBuffer.
// [out] result: A napi_value representing a JavaScript ArrayBuffer.
// N-API version: 1
func NapiCreateArrayBuffer(env NapiEnv, length uint) (NapiValue, unsafe.Pointer, NapiStatus) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_arraybuffer(env, C.size_t(length), &data, &res)
	return NapiValue(res), data, NapiStatus(status)
}

// NapiCreateBuffer function returns N-API value that allocates a node::Buffer
// object. While this is still a fully-supported data structure, in most cases
// musing a TypedArray will suffice.
// [in] env: The environment that the API is invoked under.
// [in] size: Size in bytes of the underlying buffer.
// [out] data: Raw pointer to the underlying buffer.
// [out] result: A napi_value representing a node::Buffer.
// N-API version: 1
func NapiCreateBuffer(env NapiEnv, length uint) (NapiValue, unsafe.Pointer, NapiStatus) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_buffer(env, C.size_t(length), &data, &res)
	return NapiValue(res), data, NapiStatus(status)
}

// NapiCreateBufferCopy function  allocates a node::Buffer object and initializes
// it with data copied from the passed-in buffer. While this is still a
// fully-supported data structure, in most cases using a TypedArray will suffice.
// [in] env: The environment that the API is invoked under.
// [in] length: Size in bytes of the input buffer (should be the same as the size
// of the new buffer).
// [in] data: Raw pointer to the underlying buffer to copy from.
// [out] result_data: Pointer to the new Buffer's underlying data buffer.
// [out] result: A napi_value representing a node::Buffer.
// N-API version: 1
func NapiCreateBufferCopy(env NapiEnv, length uint, raw unsafe.Pointer) (NapiValue, unsafe.Pointer, NapiStatus) {
	var res C.napi_value
	var data unsafe.Pointer
	var status = C.napi_create_buffer_copy(env, C.size_t(length), raw, &data, &res)
	return NapiValue(res), data, NapiStatus(status)
}

// NapiCreateExternal function allocates a JavaScript value with external data
// attached to it. This is used to pass external data through JavaScript code, so
// it can be retrieved later by native code. The API allows the caller to pass in
// a finalize callback, in case the underlying native resource needs to be
// cleaned up when the external JavaScript value gets collected.
// [in] env: The environment that the API is invoked under.
// [in] data: Raw pointer to the external data.
// [in] finalize_cb: Optional callback to call when the external value is being
// collected.
// [in] finalize_hint: Optional hint to pass to the finalize callback during
// collection.
// [out] result: A napi_value representing an external value.
// The created value is not an object, and therefore does not support additional
// properties. It is considered a distinct value type `napi_external`.
// N-API version: 1
func NapiCreateExternal(env NapiEnv, raw unsafe.Pointer) (NapiValue, NapiStatus) {
	var res C.napi_value
	// Remember to handle napi_finalize finalize_cb and void* finalize_hint
	var status = C.napi_create_external(env, raw, nil, nil, &res)
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

// NapiCreateObject function allocates a default JavaScript Object. It is the
// equivalent of doing new Object() in JavaScript.
// The JavaScript Object type is described in Section 6.1.7 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [out] result: A napi_value representing a JavaScript Object.
// Returns napi_ok if the API succeeded.
// N-API version: 1
func NapiCreateObject(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_object(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateSymbol function creates a JavaScript Symbol object from a
// UTF8-encoded C string.
// The JavaScript Symbol type is described in Section 19.4 of the ECMAScript
// Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] description: Optional napi_value which refers to a JavaScript String to
// be set as the description for the symbol.
// [out] result: A napi_value representing a JavaScript Symbol.
// Returns napi_ok if the API succeeded.
// N-API version: 1
func NapiCreateSymbol(env NapiEnv, value NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_symbol(env, value, &res)
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

// NapiCreateInt32 function creates JavaScript Number from the C int32_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateInt32(env NapiEnv, value int32) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_int32(env, C.int(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateUInt32 function creates JavaScript Number from the C uint32_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateUInt32(env NapiEnv, value uint32) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_uint32(env, C.uint(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateInt64 function creates JavaScript Number from the C int64_t type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateInt64(env NapiEnv, value int64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_int64(env, C.int64_t(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateDouble function creates JavaScript Number from the C double type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// The JavaScript Number type is described in Section 6.1.6 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateDouble(env NapiEnv, value float64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_double(env, C.double(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintInt64 function creates JavaScript BigInt from the C int64_t
// type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// N-API version: -
func NapiCreateBigintInt64(env NapiEnv, value int64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_bigint_int64(env, C.int64_t(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintUInt64 function creates JavaScript BigInt from the C uint64_t
// type.
// [in] env: The environment that the API is invoked under.
// [in] value: Integer value to be represented in JavaScript.
// N-API version: -
func NapiCreateBigintUInt64(env NapiEnv, value uint64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_bigint_uint64(env, C.uint64_t(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateBigintWords function converts an array of unsigned 64-bit words into
// a single BigInt value.
// [in] env: The environment that the API is invoked under.
// [in] sign_bit: Determines if the resulting BigInt will be positive or
// negative.
// [in] word_count: The length of the words array.
// [in] words: An array of uint64_t little-endian 64-bit words.
// [out] result: A napi_value representing a JavaScript BigInt.
// Returns napi_ok if the API succeeded.
// N-API version: -
func NapiCreateBigintWords(env NapiEnv, sign int, words []uint64) (NapiValue, NapiStatus) {
	var res C.napi_value
	var raw = (unsafe.Pointer(&words[0]))
	defer C.free(raw)
	var status = C.napi_create_bigint_words(env, C.int(sign), C.size_t(len(words)), (*C.uint64_t)(raw), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringLatin1 function creates a JavaScript String object from an
// ISO-8859-1-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing an ISO-8859-1-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateStringLatin1(env NapiEnv, str string) (NapiValue, NapiStatus) {
	var res C.napi_value
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_latin1(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringUtf16 function creates a JavaScript String object from a
// UTF16-LE-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing a UTF16-LE-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateStringUtf16(env NapiEnv, str string) (NapiValue, NapiStatus) {
	var res C.napi_value
	cstr := (*C.ushort)(cstring(str))
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_utf16(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateStringUtf8 function creates a JavaScript String object from a
// UTF8-encoded C string. The native string is copied.
// [in] env: The environment that the API is invoked under.
// [in] str: Character buffer representing a UTF8-encoded string.
// The JavaScript String type is described in Section 6.1.4 of the ECMAScript
// Language Specification.
// N-API version: 1
func NapiCreateStringUtf8(env NapiEnv, str string) (NapiValue, NapiStatus) {
	var res C.napi_value
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_utf8(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiGetArrayLength function returns the length of an array.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing the JavaScript Array whose length is
// being queried.
// [out] result: uint32 representing length of the array.
// Returns napi_ok if the API succeeded.
// Array length is described in Section 22.1.4.1 of the ECMAScript Language
// Specification.
// N-API version: 1
func NapiGetArrayLength(env NapiEnv, value NapiValue) (uint32, NapiStatus) {
	var res C.uint32_t
	var status = C.napi_get_array_length(env, value, &res)
	return uint32(res), NapiStatus(status)
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

// NapiGetValueBool function returns the C boolean primitive equivalent of the
// given JavaScript Boolean.
// [in] env: The environment that the API is invoked under.
// [in] value: NapiValue representing JavaScript Boolean.
// Returns napi_ok if the API succeeded. If a non-boolean NapiValue is passed
// in it returns napi_boolean_expected.
// N-API version: 1
func NapiGetValueBool(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_get_value_bool(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiGetValueDouble function returns the C double primitive equivalent of the
// given JavaScript Number.
// [in] env: The environment that the API is invoked under.
// [in] value: NapiValue representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// N-API version: 1
func NapiGetValueDouble(env NapiEnv, value NapiValue) (float64, NapiStatus) {
	var res C.double
	var status = C.napi_get_value_double(env, value, &res)
	return float64(res), NapiStatus(status)
}

// NapiGetValueBigintInt64 function returns the C int64_t primitive equivalent of
// the given JavaScript BigInt. If needed it will truncate the value, setting
// lossless to false.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript BigInt.
// Returns napi_ok if the API succeeded. If a non-BigInt is passed in it returns
// napi_bigint_expected.
// N-API version: -
func NapiGetValueBigintInt64(env NapiEnv, value NapiValue) (int64, bool, NapiStatus) {
	var res C.int64_t
	var lossless C.bool
	var status = C.napi_get_value_bigint_int64(env, value, &res, &lossless)
	return int64(res), bool(lossless), NapiStatus(status)
}

// NapiGetValueBigintUInt64 function returns the C uint64_t primitive equivalent
// of the given JavaScript BigInt. If needed it will truncate the value, setting
// lossless to false.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript BigInt.
// Returns napi_ok if the API succeeded. If a non-BigInt is passed in it returns
// napi_bigint_expected.
// N-API version: -
func NapiGetValueBigintUInt64(env NapiEnv, value NapiValue) (uint64, bool, NapiStatus) {
	var res C.uint64_t
	var lossless C.bool
	var status = C.napi_get_value_bigint_uint64(env, value, &res, &lossless)
	return uint64(res), bool(lossless), NapiStatus(status)
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

// NapiGetValueInt32 function returns the C int32 primitive equivalent of the
// given JavaScript Number.
// If the number exceeds the range of the 32 bit integer, then the result is
// truncated to the equivalent of the bottom 32 bits. This can result in a large
// positive number becoming a negative number if the value is > 2^31 -1.
// Non-finite number values (NaN, +Infinity, or -Infinity) set the result to
// zero.
// N-API version: 1
func NapiGetValueInt32(env NapiEnv, value NapiValue) (int32, NapiStatus) {
	var res C.int32_t
	var status = C.napi_get_value_int32(env, value, &res)
	return int32(res), NapiStatus(status)
}

// NapiGetValueInt64 function returns the C int64 primitive equivalent of the
// given JavaScript Number.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// Number values outside the range of Number.MIN_SAFE_INTEGER -(2^53 - 1) - Number.MAX_SAFE_INTEGER (2^53 - 1)
// will lose precision.
// Non-finite number values (NaN, +Infinity, or -Infinity) set the result to
// zero.
// N-API version: 1
func NapiGetValueInt64(env NapiEnv, value NapiValue) (int64, NapiStatus) {
	var res C.int64_t
	var status = C.napi_get_value_int64(env, value, &res)
	return int64(res), NapiStatus(status)
}

// NapiGetValueStringLatin1 function returns the ISO-8859-1-encoded string
// corresponding the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the ISO-8859-1-encoded string into. If NULL is
// passed in, the length of the string (in bytes) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of bytes copied into the buffer, excluding the null
// terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func NapiGetValueStringLatin1(env NapiEnv, value NapiValue, len uint) (string, NapiStatus) {
	var buf (*C.char)
	var res C.size_t
	var status = C.napi_get_value_string_latin1(env, value, buf, C.size_t(len), &res)
	return string(C.GoStringN(buf, C.int(res))), NapiStatus(status)
}

// NapiGetValueStringUtf8 function returns the UTF16-encoded string corresponding
// the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the UTF8-encoded string into. If NULL is passed in,
// the length of the string (in bytes) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of bytes copied into the buffer, excluding the null
// terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func NapiGetValueStringUtf8(env NapiEnv, value NapiValue, len uint) (string, NapiStatus) {
	var buf (*C.char)
	var res C.size_t
	var status = C.napi_get_value_string_utf8(env, value, buf, C.size_t(len), &res)
	return string(C.GoStringN(buf, C.int(res))), NapiStatus(status)
}

// NapiGetValueStringUtf16 function returns the UTF16-encoded string
// corresponding the value passed in.
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript string.
// [in] buf: Buffer to write the UTF16-LE-encoded string into. If NULL is passed
// in, the length of the string (in 2-byte code units) is returned.
// [in] bufsize: Size of the destination buffer. When this value is insufficient,
// the returned string will be truncated.
// [out] result: Number of 2-byte code units copied into the buffer, excluding
// the null terminator.
// Returns napi_ok if the API succeeded. If a non-String napi_value is passed in
// it returns napi_string_expected.
// N-API version: 1
func NapiGetValueStringUtf16(env NapiEnv, value NapiValue, len uint) (string, NapiStatus) {
	var buf (*C.ushort)
	var res C.size_t
	var status = C.napi_get_value_string_utf16(env, value, buf, C.size_t(len), &res)
	var str = bytes.NewBuffer(C.GoBytes(unsafe.Pointer(buf), C.int(res))).String()
	return str, NapiStatus(status)
}

// NapiGetValueUint32 function returns the C primitive equivalent of the
// given NapiValue as a uint32_t
// [in] env: The environment that the API is invoked under.
// [in] value: napi_value representing JavaScript Number.
// Returns napi_ok if the API succeeded. If a non-number NapiValue is passed in
// it returns napi_number_expected.
// N-API version: 1
func NapiGetValueUint32(env NapiEnv, value NapiValue) (uint32, NapiStatus) {
	var res C.uint32_t
	var status = C.napi_get_value_uint32(env, value, &res)
	return uint32(res), NapiStatus(status)
}

// NapiGetBoolean function returns the JavaScript singleton object that is used
// to represent the given boolean value.
// [in] env: The environment that the API is invoked under.
// [in] value: The value of the boolean to retrieve.
// N-API version: 1
func NapiGetBoolean(env NapiEnv, value bool) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_boolean(env, C.bool(value), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiGetGlobal function returns JavaScript global object.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func NapiGetGlobal(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_global(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNull function returns JavaScript null object.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func NapiGetNull(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_null(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiGetUndefined function returns JavaScript Undefined value.
// [in] env: The environment that the API is invoked under.
// N-API version: 1
func NapiGetUndefined(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_undefined(env, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToBool function implements the abstract operation ToBoolean() as
// defined in Section 7.1.2 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func NapiCoerceToBool(env NapiEnv, value NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_coerce_to_bool(env, value, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToNumber function implements the abstract operation ToNumber() as
// defined in Section 7.1.3 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func NapiCoerceToNumber(env NapiEnv, value NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_coerce_to_number(env, value, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToObject function implements the abstract operation ToObject() as
// defined in Section 7.1.13 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func NapiCoerceToObject(env NapiEnv, value NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_coerce_to_object(env, value, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiCoerceToString function mplements the abstractoperation ToString() as
// defined in Section 7.1.13 of the ECMAScript Language Specification.
// This function can be re-entrant if getters are defined on the passed-in
// Object.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to coerce.
// N-API version: 1
func NapiCoerceToString(env NapiEnv, value NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_coerce_to_string(env, value, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiTypeOf function is similar to invoke the typeof Operator on the object as
// defined in Section 12.5.5 of the ECMAScript Language Specification.
// However, it has support for detecting an External value. If value has a type
// that is invalid, an error is returned.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value whose type to query.
// If the type of value is not a known ECMAScript type and value is not an
// External value napi_invalid_arg will be returned.
// N-API version: 1
func NapiTypeOf(env NapiEnv, value NapiValue) (NapiValueType, NapiStatus) {
	var res C.napi_valuetype
	var status = C.napi_typeof(env, value, &res)
	return NapiValueType(res), NapiStatus(status)
}

// NapiInstanceOf function is similar to invoke the instanceof Operator on the
// object as defined in Section 12.10.4 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] object: The JavaScript value to check.
// [in] constructor: The JavaScript function object of the constructor function
// to check against.
// N-API version: 1
func NapiInstanceOf(env NapiEnv, object NapiValue, constructor NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_instanceof(env, object, constructor, &res)
	return bool(res), NapiStatus(status)
}

// NapiIsArray function is similar to invoke the IsArray operation on the object
// as defined in Section 7.2.2 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func NapiIsArray(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_array(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiIsArrayBuffer function checks if the Object passed in is an array buffer.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func NapiIsArrayBuffer(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_arraybuffer(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiIsBuffer function  checks if the Object passed in is a buffer.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func NapiIsBuffer(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_buffer(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiIsTypedArray function checks if the Object passed in is a typed array.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func NapiIsTypedArray(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_typedarray(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiIsDataview function checks if the Object passed in is a DataView.
// [in] env: The environment that the API is invoked under.
// [in] value: The JavaScript value to check.
// N-API version: 1
func NapiIsDataview(env NapiEnv, value NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_is_dataview(env, value, &res)
	return bool(res), NapiStatus(status)
}

// NapiStrictEquals function is simnilar to invoke the Strict Equality algorithm
// as defined in Section 7.2.14 of the ECMAScript Language Specification.
// [in] env: The environment that the API is invoked under.
// [in] lhs: The JavaScript value to check.
// [in] rhs: The JavaScript value to check against.
// N-API version: 1
func NapiStrictEquals(env NapiEnv, lhs NapiValue, rhs NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_strict_equals(env, lhs, rhs, &res)
	return bool(res), NapiStatus(status)
}

// NapiGetPropertyNames function returns the names of the enumerable properties
// of object as an array of strings. The properties of object whose key is a
// symbol will not be included.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the properties.
// N-API version: 1
func NapiGetPropertyNames(env NapiEnv, object NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_property_names(env, object, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiSetProperty function set a property on the Object passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object on which to set the property.
// [in] key: The name of the property to set.
// [in] value: The property value.
// N-API version: 1
func NapiSetProperty(env NapiEnv, object NapiValue, key NapiValue, value NapiValue) NapiStatus {
	var status = C.napi_set_property(env, object, key, value)
	return NapiStatus(status)
}

// NapiGetProperty function gets the requested property from the Object
// passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] key: The name of the property to retrieve.
// N-API version: 1
func NapiGetProperty(env NapiEnv, object NapiValue, key NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_property(env, object, key, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiHasProperty function checks if the Object passed in has the named
// property.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the property whose existence to check.
// N-API version: 1
func NapiHasProperty(env NapiEnv, object NapiValue, key NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_has_property(env, object, key, &res)
	return bool(res), NapiStatus(status)
}

// NapiDeleteProperty function attempts to delete the key own property from
// object.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the property to delete.
// N-API version: 1
func NapiDeleteProperty(env NapiEnv, object NapiValue, key NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_delete_property(env, object, key, &res)
	return bool(res), NapiStatus(status)
}

// NapiHasOwnProperty function checks if the Object passed in has the named own
// property. key must be a string or a Symbol, or an error will be thrown. N-API
// will not perform any conversion between data types.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] key: The name of the own property whose existence to check.
// N-API version: 1
func NapiHasOwnProperty(env NapiEnv, object NapiValue, key NapiValue) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_has_own_property(env, object, key, &res)
	return bool(res), NapiStatus(status)
}

// NapiSetNamedProperty function set a property on the Object passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object on which to set the property.
// [in] utf8Name: The name of the property to set.
// [in] value: The property value.
// N-API version: 1
func NapiSetNamedProperty(env NapiEnv, object NapiValue, key string, value NapiValue) NapiStatus {
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_set_named_property(env, object, ckey, value)
	return NapiStatus(status)
}

// NapiGetNamedProperty function gets the requested property from the Object
// passed in.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] utf8Name: The name of the property to get.
// N-API version: 1
func NapiGetNamedProperty(env NapiEnv, object NapiValue, key string) (NapiValue, NapiStatus) {
	var res C.napi_value
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_get_named_property(env, object, ckey, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiHasNamedProperty function checks if the Object passed in has the named
// property.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] utf8Name: The name of the property whose existence to check.
// N-API version: 1
func NapiHasNamedProperty(env NapiEnv, object NapiValue, key string) (bool, NapiStatus) {
	var res C.bool
	var ckey = C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var status = C.napi_has_named_property(env, object, ckey, &res)
	return bool(res), NapiStatus(status)
}

// NapiSetElement function sets and element on the Object passed in.
// [in] object: The object from which to set the properties.
// [in] index: The index of the property to set.
// [in] value: The property value.
// N-API version: 1
func NapiSetElement(env NapiEnv, object NapiValue, index uint, value NapiValue) NapiStatus {
	var status = C.napi_set_element(env, object, C.uint32_t(index), value)
	return NapiStatus(status)
}

// NapiGetElement function gets the element at the requested index.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the property.
// [in] index: The index of the property to get.
// N-API version: 1
func NapiGetElement(env NapiEnv, object NapiValue, index uint) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_element(env, object, C.uint32_t(index), &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiHasElement function returns if the Object passed in has an element at the
// requested index.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] index: The index of the property whose existence to check.
// N-API version: 1
func NapiHasElement(env NapiEnv, object NapiValue, index uint) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_has_element(env, object, C.uint32_t(index), &res)
	return bool(res), NapiStatus(status)
}

// NapiDeleteElement function attempts to delete the specified index from object.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object to query.
// [in] index: The index of the property to delete.
// N-API version: 1
func NapiDeleteElement(env NapiEnv, object NapiValue, index uint) (bool, NapiStatus) {
	var res C.bool
	var status = C.napi_delete_element(env, object, C.uint32_t(index), &res)
	return bool(res), NapiStatus(status)
}

// NapiDefineProperties function allows the efficient definition of multiple
// properties on a given object. The properties are defined using property
// descriptors.
// Given an array of such property descriptors, this function will set the
// properties on the object one at a time, as defined by DefineOwnProperty()
// described in Section 9.1.6 of the ECMA262 specification.
// [in] env: The environment that the N-API call is invoked under.
// [in] object: The object from which to retrieve the properties.
// [in] property_count: The number of elements in the properties array.
// [in] properties: The array of property descriptors.
// N-API version: 1
func NapiDefineProperties(env NapiEnv) NapiStatus {
	// TODO  napi_define_properties(napi_env env, napi_value object, size_t property_count, const napi_property_descriptor* properties);
	var status = C.napi_ok
	return NapiStatus(status)
}

// Working with JavaScript Functions
// N-API provides a set of APIs that allow JavaScript code to call back into
// native code.  N-API APIs that support calling back into native code take in a
// callback functions represented by the NapiCallback type.
// When the JavaScript VM calls back to native code, the NapiCallback function
// provided is invoked.
// Additionally, N-API provides a set of functions which allow calling JavaScript
// functions from native code. One can either call a function like a regular
// JavaScript function call, or as a constructor function.

// Any non-NULL data which is passed to this API via the data field of the
// NapiPropertyDescriptor items can be associated with object and freed whenever
// object is garbage-collected by passing both object and the data to
// NapiAddFinalizer.

// NapiCallFunction function allows a JavaScript function object to be called
// from a native add-on. This is the primary mechanism of calling back from the
// add-on's native code into JavaScript.
// For the special case of calling into JavaScript after an async operation,
// see NapiMakeCallback.
// [in] env: The environment that the API is invoked under.
// [in] recv: The this object passed to the called function.
// [in] func: napi_value representing the JavaScript function to be invoked.
// [in] argc: The count of elements in the argv array.
// [in] argv: Array of napi_values representing JavaScript values passed in as
// arguments to the function.
// N-API version: 1
func NapiCallFunction(env NapiEnv, receiver NapiValue, function NapiValue, args []NapiValue) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO  napi_call_function (napi_env env, napi_value recv, napi_value func, int argc, const napi_value* argv, napi_value* result)
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiCreateFunction function allows an add-on author to create a function
// object in native code. This is the primary mechanism to allow calling into
// the add-on's native code from JavaScript.
// [in] env: The environment that the API is invoked under.
// [in] utf8Name: The name of the function encoded as UTF8. This is visible
// within JavaScript as the new function object's name property.
// [in] length: The length of the utf8name in bytes, or NAPI_AUTO_LENGTH if it
// is null-terminated.
// [in] cb: The native function which should be called when this function object
// is invoked.
// [in] data: User-provided data context. This will be passed back into the
// function when invoked later.
// N-API version: 1
func NapiCreateFunction(env NapiEnv, name string, cb NapiCallback) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO create_function(napi_env env, const char* utf8name, size_t length, napi_callback cb, void* data, napi_value* result);
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetCbInfo function is used within a callback function to retrieve details
// about the call like the arguments and the this pointer from a given callback
// info.
// [in] env: The environment that the API is invoked under.
// [in] cbinfo: The callback info passed into the callback function.
// [in-out] argc: Specifies the size of the provided argv array and receives the
// actual count of arguments.
// [out] argv: Buffer to which the napi_value representing the arguments are copied. If there are more arguments than the provided count, only the requested number of arguments are copied. If there are fewer arguments provided than claimed, the rest of argv is filled with napi_value values that represent undefined.
// [out] this: Receives the JavaScript this argument for the call.
// [out] data: Receives the data pointer for the callback.
// N-API version: 1
func NapiGetCbInfo(env NapiEnv, cbinfo NapiCallbackInfo) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_get_cb_info(napi_env env, napi_callback_info cbinfo, size_t* argc, napi_value* argv, napi_value* thisArg, void** data)
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiGetNewTarget function returns the new.target of the constructor call. If
// the current callback is not a constructor call, the result is NULL.
// [in] env: The environment that the API is invoked under.
// [in] cbinfo: The callback info passed into the callback function.
// N-API version: 1
func NapiGetNewTarget(env NapiEnv, cbinfo NapiCallbackInfo) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_get_new_target(env, cbinfo, &res)
	return NapiValue(res), NapiStatus(status)
}

// NapiNewInstance function  is used to instantiate a new JavaScript value using
// a given NapiValue that represents the constructor for the object.
// [in] env: The environment that the API is invoked under.
// [in] cons: napi_value representing the JavaScript function to be invoked as a
// constructor.
// [in] argc: The count of elements in the argv array.
// [in] argv: Array of JavaScript values as napi_value representing the
// arguments to the constructor.
// [out] result: napi_value representing the JavaScript object returned, which in
// this case is the constructed object.
// N-API version: 1
func NapiNewInstance(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_new_instance(napi_env env, napi_value cons, size_t argc, napi_value* argv, napi_value* result)
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

//Object Wrap
// N-API offers a way to "wrap" C++ classes and instances so that the class
// constructor and methods can be called from JavaScript.
// The NapiDefineClass function defines a JavaScript class with constructor, s
// tatic properties and methods, and instance properties and methods that
// correspond to the C++ class.
// When JavaScript code invokes the constructor, the constructor callback uses
// NapiWrap to wrap a new C++ instance in a JavaScript object, then returns the
// wrapper object.
// When JavaScript code invokes a method or property accessor on the class, the
// corresponding NapiCallback C++ function is invoked.
// For wrapped objects it may be difficult to distinguish between a function
// called on a class prototype and a function called on an instance of a class.
// A common pattern used to address this problem is to save a persistent
// reference to the class constructor for later instanceof checks.

// NapiDefineClass function defines a JavaScript class that corresponds to
// a C++ class.
// The C++ constructor callback should be a static method on the class that calls
// the actual class constructor, then wraps the new C++ instance in a JavaScript
// object, and returns the wrapper object.
// The JavaScript constructor function returned from napi_define_class is often
// saved and used later, to construct new instances of the class from native
// code, and/or check whether provided values are instances of the class. In that
// case, to prevent the function value from being garbage-collected, create a
// persistent reference to it using NapiCreateReference and ensure the
// reference count is kept >= 1.
// [in] env: The environment that the API is invoked under.
// [in] utf8name: Name of the JavaScript constructor function; this is not
// required to be the same as the C++ class name, though it is recommended for
// clarity.
// [in] length: The length of the utf8name in bytes, or NAPI_AUTO_LENGTH if it
// is null-terminated.
// [in] constructor: Callback function that handles constructing instances of
// the class. (This should be a static method on the class, not an actual C++
// constructor function.)
// [in] data: Optional data to be passed to the constructor callback as the data
// property of the callback info.
// [in] property_count: Number of items in the properties array argument.
// [in] properties: Array of property descriptors describing static and instance
// data properties, accessors, and methods on the class.
// See documentation for NapiPropertyDescriptor function.
// [out] result: A napi_value representing the constructor function for the
// class.
// Any non-NULL data which is passed to this API via the data parameter or via
// the data field of the NapiPropertyDescriptor array items can be associated
// with the resulting JavaScript constructor (which is returned in the result
// parameter) and freed whenever the class is garbage-collected by passing both
// the JavaScript function and the data to NapiAddFinalizer.
// N-API version: 1
func NapiDefineClass(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_define_class(napi_env env, const char* utf8name, size_t length, napi_callback constructor, void* data, size_t property_count, const napi_property_descriptor* properties, napi_value* result);
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiWrap function wraps a native instance in a JavaScript object. The native
// instance can be retrieved later using NapiUnwrap().
// [in] env: The environment that the API is invoked under.
// [in] js_object: The JavaScript object that will be the wrapper for the native
// object.
// [in] native_object: The native instance that will be wrapped in the
// JavaScript object.
// [in] finalize_cb: Optional native callback that can be used to free the native
// instance when the JavaScript object is ready for garbage-collection.
// [in] finalize_hint: Optional contextual hint that is passed to the finalize
// callback.
// [out] result: Optional reference to the wrapped object.
// When JavaScript code invokes a constructor for a class that was defined using
// NapiDefineClass(), the NapiCallback for the constructor is invoked. After
// constructing an instance of the native class, the callback must then call
// NapiWrap() to wrap the newly constructed instance in the already-created
// JavaScript object that is the this argument to the constructor callback. That
// this object was created from the constructor function's prototype, so it
// already has definitions of all the instance properties and methods.
// Typically when wrapping a class instance, a finalize callback should be
// provided that simply deletes the native instance that is received as the data
// argument to the finalize callback.
// The optional returned reference is initially a weak reference, meaning it has
// a reference count of 0. Typically this reference count would be incremented
// temporarily during async operations that require the instance to remain valid.
// N-API version: 1
func NapiWrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_wrap(napi_env env, napi_value js_object, void* native_object, napi_finalize finalize_cb, void* finalize_hint, napi_ref* result);
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiUnwrap function retrieves a native instance that was previously wrapped
// in a JavaScript object using NapiWrap().
// [in] env: The environment that the API is invoked under.
// [in] js_object: The object associated with the native instance.
// [out] result: Pointer to the wrapped native instance.
// When JavaScript code invokes a method or property accessor on the class, the
// corresponding NapiCallback is invoked. If the callback is for an instance
// method or accessor, then the this argument to the callback is the wrapper
// object; the wrapped C++ instance that is the target of the call can be
// obtained then by calling NapiUnwrap() on the wrapper object.
// N-API version: 1
func NapiUnwrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// napi_remove_wrap(napi_env env, napi_value js_object, void** result)
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiRemoveWrap function retrieves a native instance that was previously
// wrapped in the JavaScript object using NapiWrap() and removes the wrapping.
// If a finalize callback was associated with the wrapping, it will no longer be
// called when the JavaScript object becomes garbage-collected.
// [in] env: The environment that the API is invoked under.
// [in] js_object: The object associated with the native instance.
// [out] result: Pointer to the wrapped native instance.
// N-API version: 1
func NapiRemoveWrap(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_remove_wrap(napi_env env, napi_value js_object, void** result)s
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
}

// NapiAddFinalizer function adds a NapiFinalize callback which will be called
// when the JavaScript object is ready for garbage collection.
// [in] env: The environment that the API is invoked under.
// [in] js_object: The JavaScript object to which the native data will be
// attached.
// [in] native_object: The native data that will be attached to the JavaScript
// object.
// [in] finalize_cb: Native callback that will be used to free the native data
// when the JavaScript object is ready for garbage-collection.
// [in] finalize_hint: Optional contextual hint that is passed to the finalize
// callback.
// [out] result: Optional reference to the JavaScript object.
// This API is similar to NapiWrap() except that:
//  - the native data cannot be retrieved later using Napinwrap(),
//  - nor can it be removed later using NapiRemoveWrap(),
//  - the API can be called multiple times with different data items in order to
//    attach each of them to the JavaScript object.
// Caution: The optional returned reference (if obtained) should be deleted via
// NapiDeleteReference ONLY in response to the finalize callback invocation. If
// it is deleted before, then the finalize callback may never be invoked.
// Therefore, when obtaining a reference a finalize callback is also required in
// order to enable correct disposal of the reference.
// N-API version: 1
func NapiAddFinalizer(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	// TODO napi_add_finalizer(napi_env env, napi_value js_object, void* native_object, napi_finalize finalize_cb, void* finalize_hint, napi_ref* result)
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

// NapiGetNodeVersion function fills the version struct with the major, minor,
// and patch version of Node.js that is currently running, and the release field
// with the value of process.release.name.
// [in] env: The environment that the API is invoked under.
// The returned buffer is statically allocated and does not need to be freed.
// N-API version: 1
func NapiGetNodeVersion(env NapiEnv) (NapiNodeVersion, NapiStatus) {
	var res *C.napi_node_version
	var status = C.napi_get_node_version(env, &res)
	return NapiNodeVersion(res), NapiStatus(status)
}

// NapiGetVersion function returns the highest version of N-API supported.
// [in] env: The environment that the API is invoked under.
// This function returns the highest N-API version supported by the Node.js
// runtime. N-API is planned to be additive such that newer releases of Node.js
// may support additional API functions. In order to allow an addon to use a
// newer function when running with versions of Node.js that support it, while
// providing fallback behavior when running with Node.js versions that don't
// support it.
// N-API version: 1
func NapiGetVersion(env NapiEnv) (uint32, NapiStatus) {
	var res C.uint32_t
	var status = C.napi_get_version(env, &res)
	return uint32(res), NapiStatus(status)
}

// NapiAdjustExternalMemory function gives V8 an indication of the amount of
// externally allocated memory that is kept alive by JavaScript objects
// (i.e. a JavaScript object that points to its own memory allocated by a native
// module). Registering externally allocated memory will trigger global garbage
// collections more often than it would otherwise.
// [in] env: The environment that the API is invoked under.
// [in] change_in_bytes: The change in externally allocated memory that is kept
// alive by JavaScript objects.
// N-API version: 1
func NapiAdjustExternalMemory(env NapiEnv, changeInBytes int64) (int64, NapiStatus) {
	var res C.int64_t
	var status = C.napi_adjust_external_memory(env, C.int64_t(changeInBytes), &res)
	return int64(res), NapiStatus(status)
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

// CCallback  ...
type CCallback func(NapiEnv, NapiCallbackInfo) NapiValue

// Caller contains a callback to call
type Caller struct {
	cb CCallback
}

/*func (s *Caller) cb(env C.napi_env, info C.napi_callback_info) C.napi_value {
	value, _ := NapiCreateInt32(NapiEnv(env), 7)
	return C.napi_value(value)
}*/

func createInt32(env NapiEnv, info NapiCallbackInfo) NapiValue {
	value, _ := NapiCreateInt32(NapiEnv(env), 7)
	return value
}

//export ExecuteCallback
func ExecuteCallback(data unsafe.Pointer, env C.napi_env, info C.napi_callback_info) C.napi_value {
	caller := (*Caller)(data)
	return (C.napi_value)(caller.cb(NapiEnv(env), NapiCallbackInfo(info)))
}

//export Initialize
func Initialize(env NapiEnv, exports NapiValue) C.napi_value {
	name := C.CString("createInt32")
	defer C.free(unsafe.Pointer(name))
	caller := &Caller{
		cb: createInt32,
	}
	desc := NapiPropertyDescriptor{
		utf8name:   name,
		name:       nil,
		method:     (C.napi_callback)(C.CallbackMethod(unsafe.Pointer(caller))), //nil,
		getter:     nil,
		setter:     nil,
		value:      nil,
		attributes: C.napi_default,
		data:       nil,
	}
	C.napi_define_properties(env, exports, 1, (*C.napi_property_descriptor)(&desc))
	return (C.napi_value)(exports)
}

func main() {}
