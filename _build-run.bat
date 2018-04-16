@echo off
REM Requires: windows x64
REM Outputs:
REM    - fe-image-optimizer.exe  - the console executable
REM    - fe-image-optimizerw.exe - the GUI executable

cd /d %~dp0 ^
    && .\_build.bat ^
    && echo running... ^
    && .\fe-image-optimizer.exe ^
    && echo done.

