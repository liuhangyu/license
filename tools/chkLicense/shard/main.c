#include <stdio.h>
#include <dlfcn.h>
#include <string.h>
#include <stdlib.h>
#include "libshared.h"

const char *licenseFilePath  = "../../cli/license.dat";
const char *productNameString =  "switch-directory-chain";
const char *libshardPath = "../../linklib/libshared.so";

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
          licenseFilePath = licenseFilePtr;
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
    p: licenseFilePath,
    n: strlen(licenseFilePath)
  };

  GoString productName ={
    p: productNameString,
    n: strlen(productNameString)
  };

  void* handle = dlopen(libshardPath, RTLD_LAZY);
  char* (*VerifyLicenseFunc)(GoString p0, GoString p1);
  char* (*ReadLicneseFunc)(GoString p0);
  long long int (*GetExpireSecFunc)(GoString p0);

  //验证license
  {
    VerifyLicenseFunc = dlsym(handle, "VerifyLicense");
    char *resp = VerifyLicenseFunc(licensePath, productName);
    printf("VerifyLicenseFunc %s\n", resp);
    free(resp);
  }


  //读取license配置文件
  {
    ReadLicneseFunc = dlsym(handle, "ReadLicnese");
    char *resp = ReadLicneseFunc(licensePath);
    printf("ReadLicneseFunc %s\n", resp);
    free(resp);
  }

  //检测过期时间
  {
    GetExpireSecFunc = dlsym(handle, "GetExpireSec");
    long long int seconds = GetExpireSecFunc(licensePath);
    printf("GetExpireSecFunc %lld\n", seconds);
  }

  dlclose(handle);
  return 0;
}

/*
gcc -o main main.c -ldl -lshared -L ../../linklib  -I ../../linklib/
./main "../../cli/license.dat"  "switch-directory-chain" "../../linklib/libshared.so"
*/