echo Cleaning previous build ... && \
rm -rf libgoaddon.a && \
rm -rf libgoaddon.h && \
rm -rf ./build && \
rm -rf ./napisys && \
cp -r ../../napisys ./napisys #&& \



#echo Start prebuild process ... && \
#echo Adding libraries ... && \
#cp ../go-napi-sys/libgoaddon.a libgoaddon.a && \
#cp ../go-napi-sys/libgoaddon.h libgoaddon.h && \
#echo Start building ... && \
# Remember for Node.js version less than 12 the MACOSX_DEPLOYMENT_TARGET need to 
# be set to 10.7
#export MACOSX_DEPLOYMENT_TARGET=10.10 && \
#npm install && \
#echo Build finished. && \
#echo Test ...
#npm test && \
#echo Test executed with success.