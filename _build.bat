@echo off
REM Requires: windows x64
REM Outputs:
REM    - fe-image-optimizer.exe  - the console executable
REM    - fe-image-optimizerw.exe - the GUI executable

cd /d %~dp0 ^
    && echo compiling resources... ^
    && rsrc -arch=amd64 -manifest app.manifest -o rsrc.syso ^
    && echo building cli program... ^
    && go build -o fe-image-optimizer.exe ^
    && echo building gui program... ^
    && go build -ldflags="-H windowsgui" -o fe-image-optimizerw.exe ^
    && echo built.

