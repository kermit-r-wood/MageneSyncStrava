BINARY_NAME=MageneSync.exe

.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

.PHONY: sync
sync: build
	./$(BINARY_NAME) sync

.PHONY: auth
auth: build
	./$(BINARY_NAME) auth

.PHONY: check
check: build
	./$(BINARY_NAME) check

.PHONY: status
status: build
	./$(BINARY_NAME) status

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf tmp
