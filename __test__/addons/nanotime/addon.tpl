#include <node_api.h>

#include "libgoaddon.h"

napi_value Init(napi_env env, napi_value  exports) {
  return (napi_value) Initialize((void*) env, (void*) exports);
}

NAPI_MODULE(NODE_GYP_MODULE_NAME, Init)
