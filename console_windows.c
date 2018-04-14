#include <windows.h>
#include "console_windows.h"

int freeConsole(){
    return (int)FreeConsole();
}

