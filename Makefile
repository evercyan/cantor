APP := "Cantor"

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: 启动调试模式
run:
	wails serve

## cli: 编译可执行文件
cli:
	wails build -d

## app: 编译 app 文件
app: clean
	wails build -p

## clean: 清除编译的程序
clean:
	rm -rf ./build/*

.PHONY: clean run cli app help