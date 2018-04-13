@echo off


cd /d "%~dp0" && jpegoptim --force --verbose --max=70 --stdout test-jpg.jpg > test-jpg.optimized.jpg



