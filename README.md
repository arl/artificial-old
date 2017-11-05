# ARTificial


## Build dependencies

Dependencies are included as submodules, to build them:

```sh
cd lib\pHash
./configure --disable-video-hash --disable-audio-hash CXXFLAGS="-I$PWD/../CImg" LDFLAGS="-lpthread"
make
```
