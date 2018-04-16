#include <windows.h>
#include "load-resource_windows.h"

void* loadAppResourceById(int64 resId, byte** startAddr, int64* size) {
    HMODULE hIns;
    HRSRC hRes;
    HGLOBAL hResLoad;

    hIns = GetModuleHandle(NULL);
    hRes = FindResource(hIns, MAKEINTRESOURCE(resId), MAKEINTRESOURCE(300));
    if (hRes == NULL) {
        return NULL;
    }

    hResLoad = LoadResource(NULL, hRes);
    if (hResLoad == NULL){
        return NULL;
    }

    *startAddr = LockResource(hResLoad);
    if (*startAddr == NULL) {
        return NULL;
    }

    *size = SizeofResource(hIns, hRes);

    return hRes;
}

int getLastError(){
    return GetLastError();
}

const char* getLastErrorAsString(){
    LPSTR messageBuffer = NULL;
    DWORD errorMessageID = GetLastError();
    if (errorMessageID == 0) {
        return messageBuffer;
    }

    FormatMessageA(FORMAT_MESSAGE_ALLOCATE_BUFFER | FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS,
                NULL, errorMessageID, MAKELANGID(LANG_ENGLISH, SUBLANG_ENGLISH_US), (LPSTR)&messageBuffer, 0, NULL);

    return messageBuffer;
}
