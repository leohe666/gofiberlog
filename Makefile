.PHONY: start build


ENV				   = local
SECOND			   = 00
RELEASE_VERSION    = v1.0
PKG_BIN 		   = fiberlog
PKG_NAME 		   = log
CONF_NAME		   = config.yaml
NOW_TIME		   = $(shell date '+%Y%m%d%H%M')
PKG_USER		   = $(shell git config user.name)
GIT_HASH           = $(shell git rev-parse --short HEAD)
GIT_COUNT 		   = $(shell git rev-list --all --count)
GIT_CURRENT_BRANCH = $(shell git symbolic-ref --short HEAD)
BASE_SECURITY_CODE = KbgrC6bA8GjozOuEFJqBXjdaVDkGuWOD2lmuW3uM3W4Ypp6LuplEwZVEqAzvhFob
RELEASE_TAG        = $(GIT_CURRENT_BRANCH)_$(RELEASE_VERSION).$(GIT_COUNT)_$(GIT_HASH)_$(PKG_USER)
PKG_FULL_NAME      = $(PKG_NAME)_$(RELEASE_TAG)_$(NOW_TIME)$(SECOND)
SECURITY_CODE      = $(BASE_SECURITY_CODE)$(PKG_FULL_NAME)


clean:
	@rm -rf ./$(PKG_BIN)


doc:
	@swag init --exclude=tmp


run:
	@rm -rf ./tmp/run
	@air -c .air.toml


gobuild:
	@go clean
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=$(RELEASE_TAG)" -o ./$(PKG_BIN)


utilsFileName = $(wildcard utils)
dbUpgradeFileName = $(wildcard dbupgrade)
init:
	@cd ./tmp/pkg && rm -rf ./*.tar.bz2 && rm -rf ./*.zip
	@go generate


build: clean gitpull init gobuild pkgbin
gitpull:
	git pull
pkgbin:
	@if [ "${ENV}" != "prod" ]; then tar jcf ./tmp/pkg/$(PKG_FULL_NAME).tar.bz2 --exclude=*.log --exclude=*.log.gz --exclude=*.xls* --exclude=pkg --exclude=main* --exclude=dev.yaml --exclude=local.yaml --exclude=prod.yaml --exclude=test.yaml --exclude=config.go  tmp $(PKG_BIN) $(CONF_NAME) && rm -rf ./$(PKG_BIN) ; else tar jcf ./tmp/pkg/$(PKG_FULL_NAME).tar.bz2 --exclude=*.log --exclude=*.log.gz --exclude=*.xls* --exclude=pkg --exclude=main* --exclude=dev.yaml --exclude=local.yaml --exclude=prod.yaml --exclude=test.yaml --exclude=config.go tmp $(PKG_BIN) && rm -rf ./$(PKG_BIN); fi
	@cd ./tmp/pkg && zip $(PKG_FULL_NAME).zip $(PKG_FULL_NAME).tar.bz2 && rm -rf ./*.tar.bz2
