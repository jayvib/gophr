APPNAME := gophr

.PHONY: build
build:
	go build -ldflags "-X main.commitID=`git rev-parse HEAD`" -o $(APPNAME)

.PHONY: run
run: build
	./$(APPNAME)
