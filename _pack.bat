@echo off

cd /d %~dp0 ^
    && .\_build.bat ^
    && echo packing... ^
    && copy /y .\fe-image-optimizer.exe    .\图片优化工具\图片优化工具.exe ^
    && copy /y .\config.json   .\图片优化工具\ ^
    && copy /y .\lib\*.* .\图片优化工具\lib ^
    && echo done.

