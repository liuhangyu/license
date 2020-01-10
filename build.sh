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
    if [ ! -d lincese.v$BUILD_VERSION ]; then   
        echo  lincese.v$BUILD_VERSION" dir not exist"
        exit -1
    fi

    mkdir -p  lincese.v$BUILD_VERSION/licensemgr
    cp licensemgr lincese.v$BUILD_VERSION/licensemgr
    cp licensemgr.exe lincese.v$BUILD_VERSION/licensemgr
    cp licensemgr.macho lincese.v$BUILD_VERSION/licensemgr
    cd lincese.v$BUILD_VERSION
    # tar -zcvf licensemgr.tar.gz licensemgr
    zip -r licensemgr.zip licensemgr
    rm -rf licensemgr
    cd $ROOT_PATH
}

packCli() {
    if [ ! -d lincese.v$BUILD_VERSION ]; then   
        echo  lincese.v$BUILD_VERSION" dir not exist"
        exit -1
    fi
    mkdir -p lincese.v$BUILD_VERSION/cli
    cp tools/linklib/libplugin.so lincese.v$BUILD_VERSION/cli
    cp tools/linklib/libshared.so lincese.v$BUILD_VERSION/cli
    cp tools/linklib/libshared.h lincese.v$BUILD_VERSION/cli

    cp tools/register/register lincese.v$BUILD_VERSION/cli
    cp tools/register/register.exe lincese.v$BUILD_VERSION/cli
    cp tools/register/register.macho lincese.v$BUILD_VERSION/cli

    cd lincese.v$BUILD_VERSION
    zip -r cli.zip cli
    rm -rf cli
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
    if [ -d lincese.v$BUILD_VERSION ]; then   
        rm -rf lincese.v$BUILD_VERSION
    fi 
    mkdir -p lincese.v$BUILD_VERSION
    packLicenseMgr
    packCli
fi


