.PHONY: deps clean build
MAIN_FUNCTION=main_function
ORCHESTRATOR_FUNCTION=orchestrator_function
WORKER_FUNCTION=worker_function
FUNCTION_PACKAGE=functions
deps:

clean: 
	rm -rf ./bin/*

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(MAIN_FUNCTION)/$(MAIN_FUNCTION) ./$(FUNCTION_PACKAGE)/main/$(MAIN_FUNCTION).go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(ORCHESTRATOR_FUNCTION)/$(ORCHESTRATOR_FUNCTION) ./$(FUNCTION_PACKAGE)/orchestrator/$(ORCHESTRATOR_FUNCTION).go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(WORKER_FUNCTION)/$(WORKER_FUNCTION) ./$(FUNCTION_PACKAGE)/worker/$(WORKER_FUNCTION).go