rm -rf package
mkdir package
cp -R assets package/assets
cp build/rurik.exe package/game.exe

cd package
7za a build.zip assets/ game.exe -tzip -mx7 -xr!*.aseprite
cd ..