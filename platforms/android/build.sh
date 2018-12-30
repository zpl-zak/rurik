# This script compiles the game into a shared library, which can then be used within an Android native project.

# Generate a standalone Android toolchain by using the tool provided by Android NDK:
# <android-ndk>/build/tools/make_standalone_toolchain.py --install-dir=/opt/android-toolchain --arch=arm --api 24 --stl=libc++

: ${ANDROID_NDK_HOME:=/opt/android-toolchain}
: ${ANDROID_SYSROOT:=${ANDROID_NDK_HOME}/sysroot}
: ${ANDROID_HOME:=/opt/android-studio}
export ANDROID_NDK_HOME ANDROID_HOME ANDROID_SYSROOT
export PATH=${ANDROID_NDK_HOME}/bin:${GRADLE_HOME}/bin:${PATH}

CC=arm-linux-androideabi-gcc \
CGO_CFLAGS="-I${ANDROID_SYSROOT}/usr/include --sysroot=${ANDROID_SYSROOT}" \
CGO_LDFLAGS="-L${ANDROID_SYSROOT}/usr/lib --sysroot=${ANDROID_SYSROOT}" \
CGO_ENABLED=1 GOOS=android GOARCH=arm \
go build -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libgame.so" \
-o=android/libs/armeabi-v7a/libgame.so ../../src/demo/*.go

rm -rf android/assets
rm -rf android/build/outputs
mkdir -p android/assets/assets
cp -r ../../assets android/assets

unset ANDROID_NDK_HOME

./gradlew assembleDebug