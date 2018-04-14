cd /d %~dp0 ^
    && rsrc -arch=amd64 -manifest app.manifest -o rsrc.syso ^
    && go build -o fe-image-optimizer.exe ^
    && go build -ldflags="-H windowsgui" -o fe-image-optimizerw.exe ^
    && .\fe-image-optimizerw.exe

