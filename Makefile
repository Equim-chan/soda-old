# Copyright (c) 2017 Equim.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PKG := .

#####################################################################
SHELL := /bin/bash
MAKE := make --no-print-directory

GC := go build
ifdef DEBUG
  GCFLAGS += -N -l
else
  LDFLAGS += -s -w
endif

VERSION := $(shell git describe --dirty --always --tags 2> /dev/null || date +"%y%m%d")
LDFLAGS += -X main.Version=$(VERSION)

ifdef GOOS
  ifeq ($(GOOS),windows)
    GOEXE := .exe
  endif
else ifeq ($(OS),Windows_NT)
  GOEXE := .exe
endif

DIRNAME := $(shell basename $(shell pwd))
TITLE := $(DIRNAME)
ifdef SUFFIX
  TITLE := $(DIRNAME)-$(SUFFIX)
else ifdef GOOS
  ifdef GOARCH
    ifdef GOARM
      TITLE := $(DIRNAME)-$(GOOS)-$(GOARCH)$(GOARM)
    else
      TITLE := $(DIRNAME)-$(GOOS)-$(GOARCH)
    endif
  endif
endif
BIN := $(TITLE)$(GOEXE)

PREFIX := $(GOPATH)/bin

HAVE_UPX := $(shell hash upx 2> /dev/null && echo 1)

define magic
  SUFFIX=$1 $(shell test "$(RELEASE)" && echo "pack")
endef

ifndef VERBOSE
  AT := @
endif

#####################################################################
SRC := $(shell find . -path ./vendor -prune -o -name '*.go')
SRC += Makefile
PLATFORMS := linux-64 linux-32 windows-64 windows-32 macos
PREPARE := vendor

$(BIN): $(PREPARE) $(SRC)
	@test "$(GOOS)" -o "$(GOARCH)" -o "$(GOARM)" && echo -ne "  "; \
	test "$(GOOS)" && echo -ne "\x1b[32mGOOS=$(GOOS) "; \
	test "$(GOARCH)" && echo -ne "\x1b[33mGOARCH=$(GOARCH) "; \
	test "$(GOARM)" && echo -ne "\x1b[35mGOARM=$(GOARM) "; \
	test "$(GOOS)" -o "$(GOARCH)" -o "$(GOARM)" && echo -ne "\n"; \
	echo -e "\x1b[0m  - \x1b[1;36mGC\x1b[0m"
	$(AT)$(GC) -o $@ --gcflags "$(GCFLAGS)" --ldflags "$(LDFLAGS)" $(PKG)

.PHONY: pack
pack: $(TITLE)-$(VERSION).tar.xz

$(TITLE)-$(VERSION).tar.xz: $(BIN)
ifdef HAVE_UPX
	@echo "  - UPX"
  ifdef VERBOSE
	upx -9 $^
  else
	@upx -q9 $^ > /dev/null
  endif
endif
	@echo "  - TAR | XZ"
	$(AT)tar cpf - $^ | xz -c9 - > $@

.PHONY: all
all: $(PLATFORMS)

.PHONY: $(PLATFORMS)
linux-64:
	@$(MAKE) GOOS=linux GOARCH=amd64 $(call magic,$@)
linux-32:
	@$(MAKE) GOOS=linux GOARCH=386 $(call magic,$@)
windows-64:
	@$(MAKE) GOOS=windows GOARCH=amd64 $(call magic,$@)
windows-32:
	@$(MAKE) GOOS=windows GOARCH=386 $(call magic,$@)
macos:
	@$(MAKE) GOOS=darwin GOARCH=amd64 $(call magic,$@)

.PHONY: release
release: RELEASE=1
release: all
	@echo -e "\n\x1b[35m  - SHA256 > $(DIRNAME)-$(VERSION).sha256\x1b[0m"
	$(AT)sha256sum *.tar.xz > $(DIRNAME)-$(VERSION).sha256

vendor:
ifdef VERBOSE
	dep ensure -v
else
	dep ensure
endif

.PHONY: install
install: $(BIN)
	cp -t $(PREFIX) $^

.PHONY: uninstall
uninstall:
	$(RM) $(PREFIX)/$(BIN)

.PHONY: test
test: $(PREPARE)
	go test ./...

.PHONY: bench
bench: $(PREPARE)
	go test --bench . ./...

.PHONY: clean
clean:
	go clean
	$(RM) $(BIN) $(DIRNAME).tar.xz $(foreach t,$(PLATFORMS),$(DIRNAME)-$(t)*) $(DIRNAME)-*.sha256
