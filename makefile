THIS_OS := unknown
ifeq ($(OS),Windows_NT)
	THIS_OS := windows
	SET_ENV := cmd /c "set
	SET_ENV_TAIL := "
	DEL_FILE := del
endif
ifeq ($(shell uname -s),Linux)
	THIS_OS := linux
	DEL_FILE := rm
endif

local_exe_list := ./local/host.texe ./local/controller.texe

EXE_EXTENSION := .exe
PACKAGE_NAME := undefined

help:
	@echo it is all for $(THIS_OS)

make_local_dir:
	-mkdir ./local

build: build_w build_l

build_l: export EXE_EXTENSION := .elf
build_l: make_local_dir set_target_linux build_controller build_host

set_target_linux: export EXE_EXTENSION := .elf
set_target_linux:
	$(SET_ENV) GOOS=linux$(SET_ENV_TAIL)

build_host: export PACKAGE_NAME := host
build_host: export PACKAGE_BUILDER := go
build_host: ./local/host.elf

build_controller: export PACKAGE_NAME := controller
build_controller: export PACKAGE_BUILDER := go
build_controller: ./local/controller.elf

./local/controller.elf ./local/controller.exe: ./local/controller.texe
./local/host.elf ./local/host.exe: ./local/host.texe

./local/host.texe: ./src/host/*.go
./local/controller.texe: ./src/controller/*.go

$(local_exe_list): (%.texe): *.go
	go build -o ./local/$(PACKAGE_NAME).texe ./src/$(PACKAGE_NAME)/
	cp ./local/$(PACKAGE_NAME).texe ./local/$(PACKAGE_NAME)$(EXE_EXTENSION)
