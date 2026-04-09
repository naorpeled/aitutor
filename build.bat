@echo off
REM AITutor build script
REM Converted from Makefile

if "%VERSION%"=="" set VERSION=dev

if "%1"=="build" (
    echo Building aitutor with version %VERSION%...
    go build -ldflags "-X main.version=%VERSION%" -o aitutor.exe .
    goto :eof
)

if "%1"=="run" (
    echo Running aitutor...
    go run .
    goto :eof
)

if "%1"=="install" (
    echo Installing aitutor with version %VERSION%...
    go install -ldflags "-X main.version=%VERSION%" .
    goto :eof
)

if "%1"=="clean" (
    echo Cleaning up...
    if exist aitutor.exe del aitutor.exe
    goto :eof
)

if "%1"=="vet" (
    echo Running go vet...
    go vet ./...
    goto :eof
)

echo Usage: build.bat {build^|run^|install^|clean^|vet}
echo Set VERSION environment variable to override default 'dev'
exit /b 1