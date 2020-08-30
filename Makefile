OUTPUT_PATH := ./bin
BUILD_FLAGS := -mod=vendor
BUILD_CMD   := go build -o $(OUTPUT_PATH) $(BUILD_FLAGS)

install : clean $(OUTPUT_PATH)/client $(OUTPUT_PATH)/server

PHONY : buildproto

mod :
	@echo update vendor
	@go mod vendor

clean :
	@rm -f $(OUTPUT_PATH)/*

buildproto :
	@protoc -I=./proto --go_out=message --go_opt=paths=source_relative ./proto/*.proto
	@echo done!

$(OUTPUT_PATH)/client : $(shell find ./internal/client/ -name "*.go" -not -name "*_test.go" -print0)
	@ $(BUILD_CMD) github.com/nerored/chat-test-golang/internal/client
	@echo process build client finish

$(OUTPUT_PATH)/server : $(shell find ./internal/server/ -name "*.go" -not -name "*_test.go" -print0)
	@ $(BUILD_CMD) github.com/nerored/chat-test-golang/internal/server
	@echo process build server finish
