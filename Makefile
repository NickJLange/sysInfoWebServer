
GO = /usr/local/bin/go

TOPLEVEL = /Users/njl/dev/
BINARY = github.com/NickJLange/sysInfoWebServer/

TARGETARCH = linux:arm darwin:amd64 windows:386

INSTALLPATH=$(TOPLEVEL)/pkg/


help:
	echo "There is no help at this time."

darwin:
	cd $(TOPLEVEL)/src/$(BINARY)
	mkdir -p $(INSTALLPATH)/darwin_amd64/$(BINARY)
	env GOOS=darwin GOARCH=amd64 $(GO) build -o $(INSTALLPATH)/darwin_amd64/$(BINARY)/sysInfoWebServer

linux:
	cd $(TOPLEVEL)/src/$(BINARY)
	mkdir -p $(INSTALLPATH)/linux_arm/$(BINARY)
	env GOOS=linux GOARCH=arm $(GO) build -o $(INSTALLPATH)/linux_arm/$(BINARY)/sysInfoWebServer

windows:
		cd $(TOPLEVEL)/src/$(BINARY)
		mkdir -p $(INSTALLPATH)/windows_386/$(BINARY)
		env GOOS=windows GOARCH=386 $(GO) build  -o $(INSTALLPATH)/windows_386/$(BINARY)/sysInfoWebServer


clean:
	rm -rf ./_obj

all:
	echo "that's all folks!	"
