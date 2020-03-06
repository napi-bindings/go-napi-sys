echo Cleaning previous build ... && \
rm -rf *.a && \
rm -rf libgoaddon.h && \
echo Start prebuild process ... && \
echo Cleaning libraries ... && \
rm -rf lib && \
echo Cleaning includes && \
cd include && \
rm -rf *.*
cd .. && \
echo Adding libraries ... && \
mkdir lib && \
#cp ../../../napi-stub/libnode_api.a lib/libnode_api.a && \
echo Adding includes ... && \
cp ../../../napi-stub/js_native_api_types.h include/js_native_api_types.h && \
cp ../../../napi-stub/js_native_api.h include/js_native_api.h && \
cp ../../../napi-stub/node_api_types.h include/node_api_types.h && \
cp ../../../napi-stub/node_api.h include/node_api.h && \
echo Start building ... && \
# Remember for Node.js version less than 12 the MACOSX_DEPLOYMENT_TARGET need to 
# be set to 10.7
export MACOSX_DEPLOYMENT_TARGET=10.10 && \
# export GOPATH=$(pwd) && \
go build -a -x -o libgoaddon.a -buildmode=c-archive . && \
echo Build finished.

