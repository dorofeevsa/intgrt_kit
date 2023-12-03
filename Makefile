export PROJECT_ROOT             := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
export PROJECT_PLUGINS          := $(PROJECT_ROOT)/cmd
export PROJECT_BIN              := $(PROJECT_ROOT)/bin
export GO               		?= go

.PHONY: build
build:
	@echo "Building "$*
	@mkdir -p $(PROJECT_ROOT)/bin
	@${GO} build -o $(PROJECT_ROOT)/bin $(PROJECT_ROOT)/example/afick_example.go

PHONY: afick_cntr
afick_cntr:
	sudo rm -f example_folder/test_*
	sudo rm -r db_folder
	mkdir -p db_folder
	@dlv --listen=:2349 --headless=true --api-version=2 --accept-multiclient exec $(PROJECT_BIN)/afick_example

PHONY: test
test:
	go test $(PROJECT_ROOT)/pkg/intgrt_afick

PHONY: iteration
iteration: build afick_cntr
