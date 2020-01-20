#include <stdio.h>
#include <dlfcn.h>
#include <string.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

#include "liblicense.h"

const char *licenseDirPath  = "../../register/";
const char *productNameString =  "switch-directory-chain";
const char *libshardPath = "../../../linklib/plugin/liblicense.so";
const char *licenseLogPath = "./license.log";

static const char *ErrList[] = {
	  "uninitialized object",
		"unknown error",
		"dir does not exist",
		"new watcher object failed",
		"watcher add dir failed",
		"failed to load public key",
		"reading authorization file failed",
		"decode authorization file failed",
		"failed to verify signature",
		"unmarshal license object failed",
		"failed to get machine code",
		"product name does not match",
		"license is expired",
		"license used before issued",
		"machine id does not match",
     NULL
};

/*
export LD_LIBRARY_PATH=../../../linklib/shared/
gcc -g -o main main.c -ldl -llicense -L ../../../linklib/shared  -I ../../../linklib/shared/
./main  "../../register"  "switch-directory-chain" "../../../linklib/shared/liblicense.so"
*/

char* (*NewLicenseFunc)(GoString p0, GoString p1, GoString p2);
char* (*FreeLicenseFunc)();
int (*VerifyLicenseFunc)(GoString p0, GoString p1);
char* (*ReadLicenseFunc)(GoString p0, GoString p1);
long long int (*GetExpireSecFunc)(GoString p0, GoString p1);

typedef struct
{
	void *handle;
    GoString licensePath;
    GoString productName;
}MulArg;


void *VerifyLicenseFc(void *arg){
    MulArg *ma = (MulArg*)arg;

    //验证license
    VerifyLicenseFunc = dlsym(ma->handle, "VerifyLicense");
    if(VerifyLicenseFunc == NULL) {
        printf("VerifyLicenseFunc is null\n");
        pthread_exit((void *)-1);
        return NULL;
    }

    while(1){
        //sleep(1);    
        int errCode = VerifyLicenseFunc(ma->licensePath, ma->productName);
        if(errCode <= 0) 
        {
            printf("验证失败, %d, %s\n", errCode, ErrList[-errCode]);
        } 
        else 
        {
            printf("验证成功\n");
        }
    }
    pthread_exit((void *)-1);
    return NULL;
}

void *ReadLicenseFc(void *arg){
    MulArg *ma = (MulArg*)arg;

    //读取license配置文件
    ReadLicenseFunc = dlsym(ma->handle, "ReadLicense");
    if(ReadLicenseFunc == NULL) {
        printf("ReadLicenseFunc is null\n");
        pthread_exit((void *)-1);
        return NULL;
    }

    while(1){
        //sleep(1);    
        char *resp = ReadLicenseFunc(ma->licensePath, ma->productName);
        printf("ReadLicenseFunc %s\n", resp);
        free(resp);
    }
    pthread_exit((void *)-1);
    return NULL;
}

void *GetExpireSecFc(void *arg){
    MulArg *ma = (MulArg*)arg;

    //检测过期时间
    GetExpireSecFunc = dlsym(ma->handle, "GetExpireSec");
    if(GetExpireSecFunc == NULL) {
        printf("GetExpireSecFunc is null\n");
        pthread_exit((void *)-1);
        return NULL;
    }

    while(1){
        //sleep(1);  
        long long int seconds = GetExpireSecFunc(ma->licensePath, ma->productName);
        printf("GetExpireSecFunc, 剩余秒数:%lld\n", seconds);
    }
    pthread_exit((void *)-1);
    return NULL;
}


int main(int argc,char *argv[])
{
    int count;
    char *licenseFilePtr;
    char *productNamePtr;
    char *libshardFilePtr;

    pthread_t tid1;
    pthread_t tid2;
    pthread_t tid3;
    void *retval1;
    void *retval2;
    void *retval3;

    if(argc > 1) {
    // printf("The command line has %d arguments :\n",argc-1);
    for (count = 1; count < argc; ++count) {
        // printf("%d: %s\n",count,argv[count]);
        if(count == 1) {
            licenseFilePtr = argv[count];
            licenseDirPath = licenseFilePtr;
        } else if(count == 2) {
            productNamePtr = argv[count];
            productNameString = productNamePtr;
        } else if(count == 3) {
            libshardFilePtr = argv[count];
            libshardPath = libshardFilePtr;
        } 
    }
    }


    GoString licensePath ={
        p: licenseDirPath,
        n: strlen(licenseDirPath)
    };

    GoString productName ={
        p: productNameString,
        n: strlen(productNameString)
    };

    GoString licenseLog ={
        p: licenseLogPath,
        n: strlen(licenseLogPath)
    };

    void* handle = dlopen(libshardPath, RTLD_LAZY);
    if(handle == NULL) {
        printf("%s\n","handle is null");
        return 0;
    }

    //创建license对象
    {
    NewLicenseFunc = dlsym(handle, "NewLicense");
    char *resp = NewLicenseFunc(licensePath, productName, licenseLog);
    if(resp != NULL && resp[0] == '\0') {
        printf("%s\n", "创建License对象成功");
    } else {
        printf("创建License对象失败 %s\n", resp);
    }

    free(resp);
    }

    MulArg mulArgs  = {
        handle = handle,
        licensePath=licensePath,
        productName=productName,
    };

    if(pthread_create(&tid1, NULL, VerifyLicenseFc, (void*)&mulArgs) ) 
    { 
        perror("pthread_create VerifyLicense error "); 
        exit(EXIT_FAILURE); 
    } 

    if(pthread_create(&tid1, NULL, ReadLicenseFc, (void*)&mulArgs) ) 
    { 
        perror("pthread_create ReadLicense error "); 
        exit(EXIT_FAILURE); 
    } 

    if(pthread_create(&tid1, NULL, GetExpireSecFc, (void*)&mulArgs) ) 
    { 
        perror("pthread_create GetExpireSec error "); 
        exit(EXIT_FAILURE); 
    } 


    pthread_join(tid1, &retval1);
    pthread_join(tid1, &retval1);
    pthread_join(tid1, &retval1);

    //销毁license对象
    {
        FreeLicenseFunc = dlsym(handle, "FreeLicense");
        char *resp = FreeLicenseFunc();
        if(resp != NULL && resp[0] == '\0') {
            printf("%s\n", "销毁License对象成功");
        } else {
            printf("销毁License对象失败 %s\n", resp);
        }

        free(resp);
    }

    dlclose(handle);
    return 0;
}

