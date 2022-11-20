#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2022  The Go-Enjin Authors
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

#: uncomment to echo instead of execute
#CMD=echo

ENJIN_MK_VERSION = v0.1.7

SHELL = /bin/bash

DEBUG ?= false

APP_NAME    ?= $(shell basename `pwd`)
APP_SUMMARY ?= Go-Enjin
ENV_PREFIX  ?= BE
APP_PREFIX  ?= ${USER}

LISTEN        ?=
PORT          ?= 3334
DEBUG         ?= false
STRICT        ?= false
DOMAIN        ?=
DENY_DURATION ?= 86400

BUILD_TAGS ?=
DEV_BUILD_TAGS ?= ${BUILD_TAGS}
GOPKG_KEYS ?=

CLEAN      ?= ${APP_NAME}
DIST_CLEAN ?=

GOLANG ?= 1.18.5
GO_MOD ?= 1018
NODEJS ?=

RELEASE_BUILD ?= false

GO_ENJIN_PKG ?= github.com/go-enjin/be

BE_PATH ?= $(call _be_local_path)

ENJENV_BIN  ?= $(shell which enjenv)
ENJENV_EXE  ?= $(call _enjenv_bin,${ENJENV_BIN})
ENJENV_URL  ?= https://github.com/go-enjin/enjenv-heroku-buildpack/raw/trunk/bin/enjenv
ENJENV_PKG  ?= github.com/go-enjin/enjenv/cmd/enjenv@latest
ENJENV_DIR_NAME ?= .enjenv
ENJENV_DIR ?= ${ENJENV_DIR_NAME}
ENJENV_PATH ?= $(call _enjenv_path)

UNTAGGED_VERSION ?= v0.0.0

VERSION ?= $(call _be_version)
RELEASE ?= $(call _be_release)

define _check_make_target =
$(shell \
	if (make -n "$(1)" 2>&1) | head -1 | grep -q "No rule to make target"; then \
		echo "false"; \
	else \
		echo "true"; \
	fi)
endef

define _be_local_path =
$(shell \
	if [ "${BE_LOCAL_PATH}" != "" -a -d "${BE_LOCAL_PATH}" ]; then \
		echo "${BE_LOCAL_PATH}"; \
	elif [ "${GOPATH}" != "" ]; then \
		if [ -d "${GOPATH}/src/${GO_ENJIN_PKG}" ]; then \
			echo "${GOPATH}/src/${GO_ENJIN_PKG}"; \
		fi; \
	fi)
endef

define _be_version =
$(shell ${ENJENV_EXE} git-tag --untagged ${UNTAGGED_VERSION})
endef

define _be_release =
$(shell ${ENJENV_EXE} rel-ver)
endef

define _enjenv_bin =
$(shell \
	if [ "${ENJENV_BIN}" != "" -a -x "${ENJENV_BIN}" ]; then \
		echo "${ENJENV_BIN}"; \
	else \
		if [ "$(1)" != "" ]; then \
			echo "$(1)"; \
		elif [ -d .bin -a -x .bin/enjenv ]; then \
			echo "${PWD}/.bin/enjenv"; \
		fi; \
	fi)
endef

define _enjenv_path =
$(shell \
	if [ -x "${ENJENV_EXE}" ]; then \
		${ENJENV_EXE}; \
	elif [ -d "./${ENJENV_DIR}" ]; then \
		echo "${PWD}/${ENJENV_DIR}"; \
	fi)
endef

define _clean =
for thing in $(1); do \
	if [ -d "$${thing}" ]; then \
		rm -rfv "$${thing}"; \
	elif [ -f "$${thing}" ]; then \
		rm -fv "$${thing}"; \
	fi; \
done
endef

define _build_tags =
$(shell if [ "${RELEASE_BUILD}" == "true" ]; then \
		if [ "${BUILD_TAGS}" != "" ]; then \
			echo "-tags ${BUILD_TAGS}"; \
		fi; \
	elif [ "${DEV_BUILD_TAGS}" != "" ]; then \
		echo "-tags ${DEV_BUILD_TAGS}"; \
	fi)
endef

define _build_label =
$(shell \
	if [ "${RELEASE_BUILD}" == "true" ]; then \
		echo "# Building release"; \
	else \
		echo "# Building debug"; \
	fi)
endef

define _build_args =
$(shell \
	if [ "${RELEASE_BUILD}" == "true" ]; then \
		echo "--optimize"; \
	fi)
endef

define _validate_extra_pkgs =
$(if ${GOPKG_KEYS},$(foreach key,${GOPKG_KEYS},$(shell \
	if [ "$($(key)_GO_PACKAGE)" == "" \
			-o "$($(key)_LOCAL_PATH)" == "" \
			-o ! -d "$($(key)_LOCAL_PATH)" ]; \
	then \
		echo "echo \"# $(key)_GO_PACKAGE and/or $(key)_LOCAL_PATH not found\"; false;"; \
	fi \
)))
endef

define _make_go_local =
echo "# go.mod local: $(1)"; \
${CMD} ${ENJENV_EXE} go-local "$(1)" "$(2)"
endef

define _make_go_unlocal =
echo "# go.mod unlocal $(1)"; \
${CMD} ${ENJENV_EXE} go-unlocal "$(1)"
endef

define _make_extra_pkgs =
$(if ${GOPKG_KEYS},$(foreach key,${GOPKG_KEYS},$($(key)_GO_PACKAGE)@latest),)
endef

define _make_extra_locals =
$(call _validate_extra_pkgs) \
$(if ${GOPKG_KEYS},$(foreach key,${GOPKG_KEYS},$(call _make_go_local,$($(key)_GO_PACKAGE),$($(key)_LOCAL_PATH));))
endef

define _make_extra_unlocals =
$(call _validate_extra_pkgs) \
$(if ${GOPKG_KEYS},$(foreach key,${GOPKG_KEYS},$(call _make_go_unlocal,$($(key)_GO_PACKAGE));))
endef

define _env_build_vars =
	BE_APP_NAME="${APP_NAME}" \
	BE_SUMMARY="${APP_SUMMARY}" \
	BE_ENV_PREFIX="${ENV_PREFIX}" \
	BE_VERSION="${VERSION}" \
	BE_RELEASE="${RELEASE}"
endef

define _env_run_vars =
	${ENV_PREFIX}_DEBUG="${DEBUG}" \
	${ENV_PREFIX}_LISTEN="${LISTEN}" \
	${ENV_PREFIX}_PORT="${PORT}" \
	${ENV_PREFIX}_PREFIX="${APP_PREFIX}" \
	${ENV_PREFIX}_STRICT="${STRICT}" \
	${ENV_PREFIX}_DOMAIN="${DOMAIN}" \
	${ENV_PREFIX}_DENY_DURATION="${DENY_DURATION}"
endef

define _is_nodejs_tag =
if [ "$(1)" != "" -a -d "$(1)" ]; then \
	if [ ! -f "$(1)/package.json" ]; then \
		echo "# $(1)/package.json not found"; \
		false; \
	fi; \
else \
	echo "# directory \"$(1)\" not found"; \
	false; \
fi
endef

define _yarn_tag_install =
if ${ENJENV_EXE} features has yarn-$(1)--install 2> /dev/null; then \
	${CMD} ${ENJENV_EXE} yarn-$(1)--install; \
fi
endef

define _yarn_run =
if ${ENJENV_EXE} features has yarn 2> /dev/null; then \
	cd $(1) > /dev/null; \
	${CMD} ${ENJENV_EXE} yarn -- $(2) 2> /dev/null || true; \
	cd - > /dev/null; \
else \
	echo "# yarn feature not found"; \
	false; \
fi
endef

define _yarn_run_script =
if ${ENJENV_EXE} features has yarn-$(1)-$(2) 2> /dev/null; then \
	${CMD} ${ENJENV_EXE} yarn-$(1)-$(2); \
else \
	echo "# yarn-$(1)-$(2) script not found"; \
	false; \
fi
endef

define _golang_nancy_installed =
if ${ENJENV_EXE} features has golang--setup-nancy 2> /dev/null; then \
	${CMD} ${ENJENV_EXE} golang setup-nancy; \
fi
endef

define _list_package_json =
$(shell ls */package.json 2> /dev/null || true)
endef

define _list_node_paths =
$(shell ls */package.json 2> /dev/null | while read P; do dirname $${P}; done)
endef

define _enjenv_present =
$(shell if [ "${ENJENV_EXE}" != "" -a -x "${ENJENV_EXE}" ]; then echo "present"; fi)
endef

define _go_present =
$(shell if ${ENJENV_EXE} features has go 2> /dev/null; then echo "present"; fi)
endef

define _yarn_present =
$(shell if ${ENJENV_EXE} features has yarn 2> /dev/null; then echo "present"; fi)
endef

define _deps_present =
$(shell \
	if [ "$(call _enjenv_present)" = "present" \
			-a "$(call _go_present)" = "present" \
			-a "$(call _yarn_present)" = "present" \
	]; then \
		echo "deps-present"; \
	fi)
endef

define _help_nodejs_audit =
	@echo "  audit-$(1)	runs enjenv $(1)-audit-report"
endef

define _help_nodejs_yarn =
	@if [ "${ENJENV_EXE}" != "" -a -x "${ENJENV_EXE}" ]; then \
		for script in $(shell ${ENJENV_EXE} -h | grep yarn-$(1)- | awk '{print $$1}'); do \
			echo "  $${script}"; \
		done; \
	fi
endef

.PHONY: all help tidy local unlocal be-update clean dist-clean build dev release run

help: NODE_PATHS = $(call _list_node_paths)
help:
	@echo "${APP_NAME} - ${APP_SUMMARY}"
	@echo
	@echo "version: ${VERSION}"
	@echo "release: ${RELEASE}"
	@echo
	@echo "usage: make <target> [targets...]"
	@echo
	@echo "Cleanup Targets:"
	@echo "  clean         remove built binary"
	@echo "  dist-clean    clean and purge local .bin ${ENJENV_DIR}"
	@echo
	@echo "Compile Targets:"
	@echo "  build         build a debug-able version of ${APP_NAME}"
	@echo "  release       build a release version of ${APP_NAME}"
	@echo
	@echo "Runtime Targets:"
	@echo "  dev           set for DEBUG mode and run ./${APP_NAME}"
	@echo "  run           execute ./${APP_NAME}"
ifneq ($(call _deps_present),)
	@echo
	@echo "Auditing Targets:"
endif
ifneq ($(call _go_present),)
	@echo "  audit		runs enjenv go-audit-report"
endif
	@echo
	@echo "helper targets:"
	@echo "  help          this screen of output"
	@echo "  tidy          run go mod tidy"
	@echo "  local         go mod -replace for ${GO_ENJIN_PKG}"
	@echo "  unlocal       go mod -dropreplace for ${GO_ENJIN_PKG}"
	@echo "  be-update     go clean and get -u ${GO_ENJIN_PKG}"
ifneq ($(call _list_package_json),)
	@$(foreach TAG,${NODE_PATHS},$(call _help_nodejs_audit,${TAG}))
	@echo
	@echo "There are nodejs package.json files present in this project"
	@echo "and enjenv has convenient wrappers for a number of extra"
	@echo "targets that Enjin.mk provides. These follow a simple naming"
	@echo "convention: yarn-<dir>-<script-name> where when targetted"
	@echo "with make, will run 'yarn <script-name>' in <dir> for you."
	@echo
	@echo "Yarn Targets:"
	@$(foreach TAG,${NODE_PATHS},$(call _help_nodejs_yarn,${TAG}))
endif
	@echo
	@echo "#############################################################"
	@echo
	@echo "This is a Go-Enjin project which depends upon enjenv to"
	@echo "manage the Golang and Nodejs SDKs and build this Enjin."
	@echo
	@echo "Go-Enjin projects use make as their build system and in"
	@echo "particular, they include Enjin.mk (which is where this"
	@echo "help screen lives)."
	@echo
	@echo "Enjin.mk will ensure that any dependencies that are not"
	@echo "already present, will be included just before they're"
	@echo "actually needed. For your information, these are summarized"
	@echo "into two categories, the managed ones Enjin.mk will handle"
	@echo "and the assumptions are up to the developer(s) managing this"
	@echo "project."
	@echo
	@echo "Managed Dependencies:"
	@echo "  enjenv        the Go-Enjin environment utility"
	@echo "  nancy         for auditing go code using Sonatype nancy"
	@echo "  nodejs        enjenv managed nodejs environment"
	@echo "  yarn          preferred nodejs package manager"
	@echo
	@echo "Assumed Dependencies:"
	@echo "  os            some flavour of linux"
	@echo "  arch          only amd64 supported (for now)"
	@echo "  make          GNU make is installed"
	@echo "  coreutils     GNU coreutils is installed"
	@echo
	@echo "Enjin.mk Version: ${ENJIN_MK_VERSION}"

_enjenv: gobin=$(shell which go)
_enjenv:
	@if [ "${ENJENV_EXE}" = "" ]; then \
		if [ "${gobin}" = "" ]; then \
			echo "# downloading enjenv..."; \
			wget -q -c ${ENJENV_URL}; \
			chmod +x ./enjenv; \
			echo "# installing enjenv..."; \
			if [ "${GOPATH}" = "" ]; then \
				mkdir -v .bin; \
				mv -v ./enjenv ./.bin/; \
				export ENJENV_EXE=$(strip ${PWD}/.bin/enjenv); \
			else \
				mv -v ./enjenv ${GOPATH}/bin/enjenv; \
				export ENJENV_EXE=$(strip ${GOPATH}/bin/enjenv); \
			fi; \
		else \
			echo "# go install enjenv..."; \
			go install ${ENJENV_PKG}; \
			export ENJENV_EXE=$(strip ${GOPATH}/bin/enjenv); \
		fi; \
	fi; \
	if [ -x "${ENJENV_EXE}" ]; then \
		echo "# using ${ENJENV_EXE}"; \
	else \
		echo "# enjenv not found ${ENJENV_EXE}"; \
		false; \
	fi

_golang: _enjenv
	@if ${ENJENV_EXE} features has --not golang--build 2> /dev/null; then \
		if [ "${GOLANG}" != "" ]; then \
			${CMD} ${ENJENV_EXE} golang init --golang "${GOLANG}"; \
		else \
			${CMD} ${ENJENV_EXE} golang init; \
		fi; \
		${CMD} ${ENJENV_EXE} golang setup-nancy; \
		${CMD} ${ENJENV_EXE} write-scripts; \
	fi

_nodejs: _enjenv
	@if ${ENJENV_EXE} features has --not yarn 2> /dev/null; then \
		if [ "${NODEJS}" != "" ]; then \
			${CMD} ${ENJENV_EXE} nodejs init --nodejs "${NODEJS}"; \
		else \
			${CMD} ${ENJENV_EXE} nodejs init; \
		fi; \
		${CMD} ${ENJENV_EXE} write-scripts; \
	fi

tidy: _golang
	@if [ ${GO_MOD} -le 1017 ]; then \
		echo "# go mod tidy -go=1.16 && go mod tidy -go=1.17"; \
		source "${ENJENV_PATH}/activate" \
			&& ${CMD} go mod tidy -go=1.16 \
			&& ${CMD} go mod tidy -go=1.17; \
	else \
		echo "# go mod tidy"; \
		source "${ENJENV_PATH}/activate" && go mod tidy; \
	fi

local: _golang
	@$(call _make_extra_locals)
	@$(call _make_go_local,${GO_ENJIN_PKG},${BE_PATH})

unlocal: _golang
	@$(call _make_extra_unlocals)
	@$(call _make_go_unlocal,${GO_ENJIN_PKG})

be-update: PKG_LIST = ${GO_ENJIN_PKG}@latest $(call _make_extra_pkgs)
be-update: _golang
	@$(call _validate_extra_pkgs)
	@echo "# go get ${PKG_LIST}"
	@source "${ENJENV_PATH}/activate" \
		&& ${CMD} GOPROXY=direct go get \
			$(call _build_tags) \
			${PKG_LIST}

clean:
ifdef override_clean
	@$(call override_clean)
else
	@$(call _clean,${CLEAN})
endif

dist-clean: clean
ifdef override_dist_clean
	@$(call override_dist_clean)
else
	@if [ -d "./${ENJENV_DIR}" ]; then \
		ENJENV_PATH="${PWD}/${ENJENV_DIR}" \
			${CMD} ${ENJENV_EXE} clean --force; \
	fi
	@$(call _clean,${DIST_CLEAN})
endif

build: _golang
ifdef override_build
	@$(call override_build)
else
	@echo "$(call _build_label): ${VERSION}, ${RELEASE}"
	@${CMD} \
		$(call _env_build_vars) \
		${ENJENV_EXE} golang build \
			$(call _build_args) \
			-- -v $(call _build_tags)
	@if [ -x "./${APP_NAME}" ]; then \
		echo "# produced: ${APP_NAME}"; \
		sha256sum ./${APP_NAME}; \
	fi
endif

release: RELEASE_BUILD="true"
release: build

RUN_ARGV ?=

define _run =
	@if [ ! -x "${APP_NAME}" ]; then \
		echo "${APP_NAME} not found or not executable"; \
		false; \
	fi
	@echo "# running ${APP_NAME} -- ${RUN_ARGV}"
	@${CMD} \
		$(call _env_run_vars) \
		./${APP_NAME} ${RUN_ARGV}
endef

run:
ifdef override_run
	@$(call override_run)
else
	@$(call _run)
endif

con-%: export CON_TAG=$(patsubst con-%,%,$@)
con-%: export RUN_ARGV=console ${CON_TAG}
con-%:
ifdef override_run
	@$(call override_run)
else
	@$(call _run)
endif

dev: DEBUG=true
dev: run

_audit: _golang
	@$(call _golang_nancy_installed)
	@${CMD} ${ENJENV_EXE} go-audit-report --tags=all

ifneq ($(call _go_present),)
audit: _audit
endif

ifneq ($(call _list_package_json),)
audit-%: TAG = $(patsubst audit-%,%,$@)
audit-%: _enjenv
	@$(call _is_nodejs_tag,${TAG})
	@$(call _yarn_tag_install,${TAG})
	@${CMD} ${ENJENV_EXE} ${TAG}-audit-report

yarn-%: FIRST = $(patsubst yarn-%,%,$@)
yarn-%: SECOND = $(subst -, ,${FIRST})
yarn-%: TAG = $(word 1,${SECOND})
yarn-%: OP = $(patsubst ${TAG}-%,%,${FIRST})
yarn-%:
	@if [ "${TAG}" = "" -o "${OP}" = "" -o "${TAG}" = "${OP}" ]; then \
		echo "# invalid yarn target: $@"; \
		false; \
	fi
	@$(call _is_nodejs_tag,${TAG})
	@$(call _yarn_tag_install,${TAG})
	@$(call _yarn_run_script,${TAG},${OP})
endif

HEROKU_GIT_REMOTE ?= heroku
HEROKU_SRC_BRANCH ?= trunk
HEROKU_DST_BRANCH ?= main

heroku-push:
	@${CMD} git push ${HEROKU_GIT_REMOTE} ${HEROKU_SRC_BRANCH}:${HEROKU_DST_BRANCH}

heroku-logs:
	@${CMD} heroku logs --tail

# this requires Term::ANSIColor, will error if not present,
# use `make build dev` instead
build-dev-run: build
	@( make dev 2>&1 ) | perl -p -e 'use Term::ANSIColor qw(colored);while (my $$line = <>) {print STDOUT process_line($$line)."\n";}exit(0);sub process_line {my ($$line) = @_;chomp($$line);if ($$line =~ m!^\[(\d+\-\d+\.\d+)\]\s+([A-Z]+)\s+(.+?)\s*$$!) {my ($$datestamp, $$level, $$message) = ($$1, $$2, $$3);my $$colour = "white";if ($$level eq "ERROR") {$$colour = "bold white on_red";} elsif ($$level eq "INFO") {$$colour = "green";} elsif ($$level eq "DEBUG") {$$colour = "yellow";}my $$out = "[".colored($$datestamp, "blue")."]";$$out .= " ".colored($$level, $$colour);if ($$level eq "DEBUG") {$$out .= "\t";if ($$message =~ m!^(.+?)\:(\d+)\s+\[(.+?)\]\s+(.+?)\s*$$!) {my ($$file, $$ln, $$tag, $$info) = ($$1, $$2, $$3, $$4);$$out .= colored($$file, "bright_blue");$$out .= ":".colored($$ln, "blue");$$out .= " [".colored($$tag, "bright_blue")."]";$$out .= " ".colored($$info, "bold cyan");} else {$$out .= $$message;}} elsif ($$level eq "ERROR") {$$out .= "\t".colored($$message, $$colour);} elsif ($$level eq "INFO") {$$out .= "\t".colored($$message, $$colour);} else {$$out .= "\t".$$message;}return $$out;}return $$line;}'

release-dev-run: release
	@( make dev 2>&1 ) | perl -p -e 'use Term::ANSIColor qw(colored);while (my $$line = <>) {print STDOUT process_line($$line)."\n";}exit(0);sub process_line {my ($$line) = @_;chomp($$line);if ($$line =~ m!^\[(\d+\-\d+\.\d+)\]\s+([A-Z]+)\s+(.+?)\s*$$!) {my ($$datestamp, $$level, $$message) = ($$1, $$2, $$3);my $$colour = "white";if ($$level eq "ERROR") {$$colour = "bold white on_red";} elsif ($$level eq "INFO") {$$colour = "green";} elsif ($$level eq "DEBUG") {$$colour = "yellow";}my $$out = "[".colored($$datestamp, "blue")."]";$$out .= " ".colored($$level, $$colour);if ($$level eq "DEBUG") {$$out .= "\t";if ($$message =~ m!^(.+?)\:(\d+)\s+\[(.+?)\]\s+(.+?)\s*$$!) {my ($$file, $$ln, $$tag, $$info) = ($$1, $$2, $$3, $$4);$$out .= colored($$file, "bright_blue");$$out .= ":".colored($$ln, "blue");$$out .= " [".colored($$tag, "bright_blue")."]";$$out .= " ".colored($$info, "bold cyan");} else {$$out .= $$message;}} elsif ($$level eq "ERROR") {$$out .= "\t".colored($$message, $$colour);} elsif ($$level eq "INFO") {$$out .= "\t".colored($$message, $$colour);} else {$$out .= "\t".$$message;}return $$out;}return $$line;}'
