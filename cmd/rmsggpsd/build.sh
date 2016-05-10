#!/bin/bash

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_linux_amd64 ./rmsggpsd.go
env GOOS=linux GOARCH=386 GO386=sse2 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_linux_x86 ./rmsggpsd.go
env GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_linux_arm6 ./rmsggpsd.go
env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_linux_arm7 ./rmsggpsd.go
env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_linux_arm64 ./rmsggpsd.go

env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_win64.exe ./rmsggpsd.go
env GOOS=windows GOARCH=386 GO386=sse2 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_win32.exe ./rmsggpsd.go

env GOOS=darwin GOARCH=386 GO386=sse2 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_darwin_x86 ./rmsggpsd.go
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_darwin_arm64 ./rmsggpsd.go
env GOOS=darwin GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_darwin_arm ./rmsggpsd.go
env GOOS=darwin GOARCH=arm64 GOARM=7 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_darwin_arm64 ./rmsggpsd.go

env GOOS=freebsd GOARCH=arm CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_freebsd_arm ./rmsggpsd.go
env GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_darwin_amd64 ./rmsggpsd.go
env GOOS=linux GOARCH=386 GO386=sse2 CGO_ENABLED=0 go build -ldflags="-s -w -d" -o rmsggpsd_freebsd_x86 ./rmsggpsd.go
