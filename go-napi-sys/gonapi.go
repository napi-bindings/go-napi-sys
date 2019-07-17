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
import "unsafe"

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

// NapiCallbackInfo represents ...
type NapiCallbackInfo C.napi_callback_info

// NapiDeferred represents ...
type NapiDeferred C.napi_deferred

// NapiPropertyAttributes represents ...
type NapiPropertyAttributes C.napi_property_attributes

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

// NapiTypedArrayType represents ...
type NapiTypedArrayType C.napi_typedarray_type

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

// NapiCallback represents ...
type NapiCallback C.napi_callback

// NapiFinalize represents ...
type NapiFinalize C.napi_finalize

// NapiPropertyDescriptor represents ...
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
// // JavaScript arrays are described in Section 22.1 of the ECMAScript Language
// Specification.
// [in] env: The environment that the API is invoked under.
// [in] length: The initial length of the Array.
// N-API version: 1
func NapiCreateArrayWithLength(env NapiEnv, length uint) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_create_array_with_length(env, C.size_t(length), &res)
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

// NapiCreateBigintWords function ...
func NapiCreateBigintWords() (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
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
/*func NapiCreateStringUtf16(env NapiEnv, str string) (NapiValue, NapiStatus) {
	var res C.napi_value
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	var status = C.napi_create_string_utf16(env, cstr, C.NAPI_AUTO_LENGTH, &res)
	return NapiValue(res), NapiStatus(status)
}*/

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

// NapiGetValueStringUtf16 function ...s
func NapiGetValueStringUtf16(env NapiEnv) (NapiValue, NapiStatus) {
	var res C.napi_value
	var status = C.napi_ok
	return NapiValue(res), NapiStatus(status)
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

// NapiDeleteElement function ...
// N-API version: 1
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
