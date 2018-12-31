#!/bin/sh

rm -rf package

# Ship the latest Linux build
echo Shipping for Linux...
make
./zip.sh game
mv package/build.zip package/demo-linux-x86_64.zip
butler push package/demo-linux-x86_64.zip zaklaus/rurik-framework:demo-linux

# Ship the latest Windows build
echo Shipping for Windows...
make win
./zip.sh game.exe
mv package/build.zip package/demo-win64-x86_64.zip
butler push package/demo-win64-x86_64.zip zaklaus/rurik-framework:demo-windows
# butler


echo Done!