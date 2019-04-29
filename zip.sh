#!/bin/sh

mkdir package 2> /dev/null
cp -Rf data package/data 
cp build/rurik.exe package/${1-game.exe} -f

cd package
7za a build.zip data/ ${1-game.exe} -bso0 -bsp0 -tzip -mx7 -xr!*.aseprite
cd ..