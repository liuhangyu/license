#include <stdio.h>
#include <dlfcn.h>
#include <string.h>
#include <stdlib.h>
#include "libshared.h"

const char *licenseDirPath  = "../cli";
const char *productNameString =  "switch-directory-chain";
const char *libshardPath = "../linklib/libshared.so";
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




int main(int argc,char *argv[])
{
  int count;
  char *licenseFilePtr;
  char *productNamePtr;
  char *libshardFilePtr;
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
  char* (*NewLicenseFunc)(GoString p0, GoString p1, GoString p2);
  char* (*FreeLicenseFunc)();

  int (*VerifyLicenseFunc)(GoString p0, GoString p1);
  char* (*ReadLicneseFunc)(GoString p0, GoString p1);
  long long int (*GetExpireSecFunc)(GoString p0, GoString p1);

  if(handle == NULL) {
    printf("%s","handle is null");
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

  //验证license
  {
    VerifyLicenseFunc = dlsym(handle, "VerifyLicense");
    int errCode = VerifyLicenseFunc(licensePath, productName);
    if(errCode <= 0) 
    {
			printf("验证失败, %d, %s\n", errCode, ErrList[-errCode]);
		} 
    else 
    {
      printf("验证成功\n");
    }
  }

  //读取license配置文件
  {
    ReadLicneseFunc = dlsym(handle, "ReadLicnese");
    char *resp = ReadLicneseFunc(licensePath, productName);
    printf("ReadLicneseFunc %s\n", resp);
    free(resp);
  }

  //检测过期时间
  {
    GetExpireSecFunc = dlsym(handle, "GetExpireSec");
    long long int seconds = GetExpireSecFunc(licensePath, productName);
    printf("GetExpireSecFunc, 剩余秒数:%lld\n", seconds);
  }


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

