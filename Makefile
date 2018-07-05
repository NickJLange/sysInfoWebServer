
GO = go

TOPLEVEL ?= /Users/njl/dev/
BINARY = github.com/NickJLange/sysInfoWebServer/

TARGETARCH = linux:arm darwin:amd64 windows:386

INSTALLPATH=$(TOPLEVEL)/pkg/

vendorparty:
	for i in NickJLange/tickdatabase NickJLange/monitoring ; do \
	 mkdir -p vendor/github.com/$$i; rsync -arp --exclude=.git $(TOPLEVEL)/src/github.com/$$i/ vendor/github.com/$$i/ ; \
	 done
	for i in bmizerany/perks/ ; do \
	  mkdir -p vendor/github.com/$$i;  rsync -arp --exclude=.git  $(TOPLEVEL)/vendor/src/github.com/$$i/ ./vendor/github.com/$$i/; done

docker: vendorparty
	docker build .

help:
	@echo "There is no help at this time."

darwin:
	cd $(TOPLEVEL)/src/$(BINARY)
	mkdir -p $(INSTALLPATH)/darwin_amd64/$(BINARY)
	env GOOS=darwin GOARCH=amd64 $(GO) build -o $(INSTALLPATH)/darwin_amd64/$(BINARY)/sysInfoWebServer

linux:
	cd $(TOPLEVEL)/src/$(BINARY)
	mkdir -p $(INSTALLPATH)/linux_arm/$(BINARY)
	ls -ld $(INSTALLPATH)/linux_arm/$(BINARY)
	env GOOS=linux GOARCH=arm $(GO) build -o $(INSTALLPATH)/linux_arm/$(BINARY)/sysInfoWebServer
	ls -ld $(INSTALLPATH)/linux_arm/$(BINARY)/sysInfoWebServer

windows:
		cd $(TOPLEVEL)/src/$(BINARY)
		mkdir -p $(INSTALLPATH)/windows_386/$(BINARY)
		env GOOS=windows GOARCH=386 $(GO) build  -o $(INSTALLPATH)/windows_386/$(BINARY)/sysInfoWebServer


clean:
	rm -rf ./_obj

all:
	@echo "that's all folks!	"
