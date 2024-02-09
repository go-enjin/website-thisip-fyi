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
CL_LOCAL_PATH ?= ../../go-corelibs

APP_NAME    := be-thisip-fyi
APP_SUMMARY := thisip.fyi
APP_COMMAND := ./cmd/${APP_NAME}

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

BUILD_TAGS     = production embeds $(COMMON_TAGS)
DEV_BUILD_TAGS = locals ngrokio $(COMMON_TAGS)

# Custom go.mod locals
GOPKG_KEYS += _TIMES
GOPKG_KEYS += _SEMANTIC_THEME

AUTO_CORELIBS_KEYS := true

#MAKE_THEME_LOCALES   := true
#MAKE_SOURCE_LOCALES  := true
#MAKE_CONTENT_LOCALES := true

include ./Enjin.mk
