@echo off


cd /d "%~dp0" && pngquant --force --verbose --speed 1 --quality 50-90 --strip --output test-png.optimized.png test-png.png



