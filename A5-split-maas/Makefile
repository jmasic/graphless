.PHONY: deps clean build
MAIN_FUNCTION=main_function
ORCHESTRATOR_FUNCTION=orchestrator_function
WORKER_FUNCTION=worker_function
FUNCTION_PACKAGE=functions

OUTPUT_FOLDER="bin"
LOCAL_OUTPUT_FOLDER="bin-local"

deps:

clean:
	rm -rf ./$(OUTPUT_FOLDER)/*

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_FOLDER)/$(MAIN_FUNCTION)/$(MAIN_FUNCTION) ./$(FUNCTION_PACKAGE)/main/$(MAIN_FUNCTION).go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_FOLDER)/$(ORCHESTRATOR_FUNCTION)/$(ORCHESTRATOR_FUNCTION) ./$(FUNCTION_PACKAGE)/orchestrator/$(ORCHESTRATOR_FUNCTION).go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_FOLDER)/$(WORKER_FUNCTION)/$(WORKER_FUNCTION) ./$(FUNCTION_PACKAGE)/worker/$(WORKER_FUNCTION).go


clean-local:
	rm -rf ./$(LOCAL_OUTPUT_FOLDER)/*

build-local:
	go build -ldflags="-s -w" -o $(LOCAL_OUTPUT_FOLDER)/$(MAIN_FUNCTION)/$(MAIN_FUNCTION) ./$(FUNCTION_PACKAGE)/main/$(MAIN_FUNCTION).go
	go build -ldflags="-s -w" -o $(LOCAL_OUTPUT_FOLDER)/$(ORCHESTRATOR_FUNCTION)/$(ORCHESTRATOR_FUNCTION) ./$(FUNCTION_PACKAGE)/orchestrator/$(ORCHESTRATOR_FUNCTION).go
	go build -ldflags="-s -w" -o $(LOCAL_OUTPUT_FOLDER)/$(WORKER_FUNCTION)/$(WORKER_FUNCTION) ./$(FUNCTION_PACKAGE)/worker/$(WORKER_FUNCTION).go
