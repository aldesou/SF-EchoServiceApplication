rmdir /s /q out

mkdir out
mkdir out\EchoService
mkdir out\EchoService\Code

copy sfpkg\ApplicationManifest.xml out\ApplicationManifest.xml
copy sfpkg\EchoService\ServiceManifest.xml out\EchoService\ServiceManifest.xml

cd code

go build -o ..\out\EchoService\Code\EchoService.exe

cd ..