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

    mkdir -p  binPack/licensemgr-$BUILD_VERSION
    cp licensemgr binPack/licensemgr-$BUILD_VERSION
    cp licensemgr.exe binPack/licensemgr-$BUILD_VERSION
    cp licensemgr.macho binPack/licensemgr-$BUILD_VERSION
    cd binPack
    # tar -zcvf licensemgr.tar.gz licensemgr
    zip -r licensemgr-$BUILD_VERSION.zip licensemgr-$BUILD_VERSION
    rm -rf licensemgr-$BUILD_VERSION
    cd $ROOT_PATH
}

packCli() {
    if [ ! -d binPack ]; then   
        echo  binPack" dir not exist"
        exit -1
    fi
    mkdir -p binPack/license
    cp tools/linklib/libplugin.so binPack/license
    cp tools/linklib/libshared.so binPack/license
    cp tools/linklib/libshared.h binPack/license

    cp tools/register/register binPack/license
    cp tools/register/register.exe binPack/license
    cp tools/register/register.macho binPack/license

    cd binPack
    zip -r license.zip license
    rm -rf license
    cd $ROOT_PATH
}

backUpMgr() {
    if [ ! -d licenseMgrVersion ]; then   
        mkdir -p licenseMgrVersion
    fi
    cp -f binPack/licensemgr*  licenseMgrVersion
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
    backUpMgr
fi


