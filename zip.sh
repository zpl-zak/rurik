#!/bin/sh

mkdir package 2> /dev/null
cp -Rf assets package/assets 
cp build/rurik.exe package/${1-game.exe} -f

cd package
7za a build.zip assets/ ${1-game.exe} -bso0 -bsp0 -tzip -mx7 -xr!*.aseprite
cd ..