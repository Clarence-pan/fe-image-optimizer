#ifndef _LOAD_RESOURCE_WINDOWS_H_
#define _LOAD_RESOURCE_WINDOWS_H_

typedef unsigned char byte;
typedef long long int64;

void* loadAppResourceById(int64 resId, byte** startAddr, int64* size);

int getLastError();
const char* getLastErrorAsString();

#endif // define _LOAD_RESOURCE_WINDOWS_H_