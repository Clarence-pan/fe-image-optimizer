@echo off

cd /d %~dp0 ^
    && .\_build.bat ^
    && echo packing... ^
    && copy /y .\fe-image-optimizer.exe    .\ͼƬ�Ż�����\ͼƬ�Ż�����.exe ^
    && copy /y .\config.json   .\ͼƬ�Ż�����\ ^
    && copy /y .\lib\*.* .\ͼƬ�Ż�����\lib ^
    && echo done.

