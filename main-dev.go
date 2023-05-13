//go:build !production

// Copyright (c) 2022  The Go-Enjin Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/fs/menu"
	"github.com/go-enjin/be/features/fs/public"
	"github.com/go-enjin/be/features/fs/themes"
	"github.com/go-enjin/be/pkg/log"
)

func init() {
	// locals environment, early startup debug logging
	log.Config.LogLevel = log.LevelDebug
	log.Config.Apply()

	fMenu = menu.New().
		MountLocalPath("/", "menus").
		Make()
	fPublic = public.New().
		MountLocalPath("/", "public").
		Make()
	fContent = content.New().
		MountLocalPath("/", "content").
		AddToIndexProviders("pages-pql").
		SetKeyValueCache(gFsContentKvsFeature, gFsContentKvsCache).
		Make()

	fThemes = themes.New().
		AddTheme(semantic.Theme()).
		LocalTheme("themes/thisip-fyi").
		SetTheme("thisip-fyi").
		Make()

	hotReload = true
}
