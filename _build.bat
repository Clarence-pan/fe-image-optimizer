@echo off
REM Requires: windows x64
REM Outputs:
REM    - fe-image-optimizer.exe  - the console executable
REM    - fe-image-optimizerw.exe - the GUI executable
REM You can also use the following command to compile resources:
REM    rsrc -arch=amd64 -manifest app.manifest -ico app.ico -o rsrc.syso

cd /d %~dp0 ^
    && echo compiling resources... ^
    && windres -o rsrc.syso -F pe-x86-64 app.rc ^
    && echo building... ^
    && go build -ldflags="-H windowsgui" -o fe-image-optimizer.exe ^
    && echo built.

