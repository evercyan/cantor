APP := "Cantor"
CREATE_DMG ?= $(shell which create-dmg)

define check
	command -v $(1) 1>/dev/null || $(2)
endef

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: 启动调试模式
dev:
	wails dev

## build: 编译 app 文件
build: clean
	wails build

## clean: 清除编译的程序
clean:
	rm -rf ./build/bin/*
	rm -rf ./Cantor.dmg

## dmg: 生成 dmg 文件
dmg: build
	@@$(call check,create-dmg,brew install create-dmg)
	$(CREATE_DMG) \
		--volname "Cantor" \
		--background "assets/background.png" \
		--window-size 558 367 \
		--icon-size 100 \
		--icon "Cantor.app" 200 190 \
		--hide-extension "Cantor.app" \
		--app-drop-link 370 190 \
		"Cantor.dmg" \
		"build/bin/"

.PHONY: help clean dev build dmg