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
    if [ ! -d binPack ]; then   
        echo  binPack" dir not exist"
        exit -1
    fi

    mkdir -p  binPack/licensemgr-v$BUILD_VERSION
    cp licensemgr binPack/licensemgr-v$BUILD_VERSION
    cp licensemgr.exe binPack/licensemgr-v$BUILD_VERSION
    cp licensemgr.macho binPack/licensemgr-v$BUILD_VERSION
    cd binPack
    # tar -zcvf licensemgr.tar.gz licensemgr
    zip -r licensemgr-v$BUILD_VERSION.zip licensemgr-v$BUILD_VERSION
    rm -rf licensemgr-v$BUILD_VERSION
    cd $ROOT_PATH
}

packCli() {
    if [ ! -d binPack ]; then   
        echo  binPack" dir not exist"
        exit -1
    fi
    mkdir -p binPack/license-v$BUILD_VERSION
    cp tools/linklib/libplugin.so binPack/license-v$BUILD_VERSION
    cp tools/linklib/libshared.so binPack/license-v$BUILD_VERSION
    cp tools/linklib/libshared.h binPack/license-v$BUILD_VERSION

    cp tools/register/register binPack/license-v$BUILD_VERSION
    cp tools/register/register.exe binPack/license-v$BUILD_VERSION
    cp tools/register/register.macho binPack/license-v$BUILD_VERSION

    cd binPack
    zip -r license-v$BUILD_VERSION.zip license-v$BUILD_VERSION
    rm -rf license-v$BUILD_VERSION
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
    if [ -d binPack ]; then   
        rm -rf binPack
    fi 
    mkdir -p binPack
    packLicenseMgr
    packCli
fi


