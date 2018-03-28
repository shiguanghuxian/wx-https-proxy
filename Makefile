default:
	@echo 'Usage of make: [ build | linux_build | clean ]'

build: 
	@go build -o ./build/wx-https-proxy ./

linux_build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/wx-https-proxy ./

windows_build: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/wx-https-proxy.exe ./

windows_build_386: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./build/wx-https-proxy.exe ./


clean: 
	@rm -f ./build/wx-https-proxy
	@rm -f ./build/wx-https-proxy.exe

.PHONY: default build linux_build windows_build clean