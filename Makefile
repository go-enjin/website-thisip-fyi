#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2023  The Go-Enjin Authors
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

-include .env

BE_LOCAL_PATH ?= ../be

APP_NAME    := be-thisip-fyi
APP_SUMMARY := thisip.fyi

DENY_DURATION ?= 60

ADD_TAGS_DEFAULTS := true

COMMON_TAGS += driver_fs_embed
COMMON_TAGS += driver_kvs_gocache memory
COMMON_TAGS += log_papertrail
COMMON_TAGS += user_auth_basic
COMMON_TAGS += user_base_htenv
COMMON_TAGS += page_pql
COMMON_TAGS += page_robots
COMMON_TAGS += fs_theme fs_menu fs_content fs_public
COMMON_TAGS += ngrokio

BUILD_TAGS     = production embeds $(COMMON_TAGS)
DEV_BUILD_TAGS = locals $(COMMON_TAGS)

# Custom go.mod locals
GOPKG_KEYS = SET DJHT

# Semantic Enjin Theme
SET_GO_PACKAGE = github.com/go-enjin/semantic-enjin-theme
SET_LOCAL_PATH = ../semantic-enjin-theme

# Go-Enjin gotext package (manual v-enjin.0 releases)
GOXT_GO_PACKAGE = github.com/go-enjin/golang-org-x-text
GOXT_LOCAL_PATH = ../golang-org-x-text

# Go-Enjin times package
DJHT_GO_PACKAGE = github.com/go-enjin/github-com-djherbis-times
DJHT_LOCAL_PATH = ../github-com-djherbis-times

#MAKE_THEME_LOCALES   := true
#MAKE_SOURCE_LOCALES  := true
#MAKE_CONTENT_LOCALES := true

include ./Enjin.mk
