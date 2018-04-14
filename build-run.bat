cd /d %~dp0 ^
    && rsrc -manifest app.manifest -o rsrc.syso ^
    && go build -ldflags="-H windowsgui" -o test.exe ^
    && .\test.exe

