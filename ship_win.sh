#!/bin/sh

rm -rf package

# Ship the latest Windows build
echo Shipping for Windows...
make
./zip.sh game.exe
mv package/build.zip package/demo-win64-x86_64.zip
butler push package/demo-win64-x86_64.zip zaklaus/rurik-framework:demo-windows

echo Done!