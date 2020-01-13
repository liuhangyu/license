#!/bin/bash
ROOT_PATH=$PWD
BUILD_VERSION=$(echo `git describe --abbrev=0 --tags`)

buildLicenseMgr() {
    echo "begin build licensemgr...."
    if [ $1 == "clean" ];then
        make clean
    else
        make
    fi
    echo "end build licensemgr....."
}

buildRegister() {
    echo "begin build register....."
    cd tools/register
    if [ $1 == "clean" ];then
        make clean
    else
        make
    fi
    cd $ROOT_PATH
    echo "end build register....."
}

buildDll() {
    echo "begin build dll....."
    cd tools/linklib
    if [ $1 == "clean" ];then
        make clean
    else
        make
    fi
    cd $ROOT_PATH
    echo "end build dll....."
}

#############打包#############
packLicenseMgr() {
    if [ ! -d lincese ]; then   
        echo  lincese" dir not exist"
        exit -1
    fi

    mkdir -p  lincese/licensemgr-v$BUILD_VERSION
    cp licensemgr lincese/licensemgr-v$BUILD_VERSION
    cp licensemgr.exe lincese/licensemgr-v$BUILD_VERSION
    cp licensemgr.macho lincese/licensemgr-v$BUILD_VERSION
    cd lincese
    # tar -zcvf licensemgr.tar.gz licensemgr
    zip -r licensemgr-v$BUILD_VERSION.zip licensemgr-v$BUILD_VERSION
    rm -rf licensemgr-v$BUILD_VERSION
    cd $ROOT_PATH
}

packCli() {
    if [ ! -d lincese ]; then   
        echo  lincese" dir not exist"
        exit -1
    fi
    mkdir -p lincese/cli-v$BUILD_VERSION
    cp tools/linklib/libplugin.so lincese/cli-v$BUILD_VERSION
    cp tools/linklib/libshared.so lincese/cli-v$BUILD_VERSION
    cp tools/linklib/libshared.h lincese/cli-v$BUILD_VERSION

    cp tools/register/register lincese/cli-v$BUILD_VERSION
    cp tools/register/register.exe lincese/cli-v$BUILD_VERSION
    cp tools/register/register.macho lincese/cli-v$BUILD_VERSION

    cd lincese
    zip -r cli-v$BUILD_VERSION.zip cli-v$BUILD_VERSION
    rm -rf cli-v$BUILD_VERSION
    cd $ROOT_PATH
}


MODE=$1
if [ "${MODE}" == "clean" ]; then
    buildLicenseMgr clean
    buildRegister clean
    buildDll clean
elif [ "${MODE}" == "build" ]; then
    buildLicenseMgr build
    buildRegister build
    buildDll build
elif [ "${MODE}" == "pack" ]; then
    echo "pack binary package..."
    echo $BUILD_VERSION
    if [ -d lincese ]; then   
        rm -rf lincese
    fi 
    mkdir -p lincese
    packLicenseMgr
    packCli
fi


