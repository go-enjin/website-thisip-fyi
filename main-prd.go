//go:build production

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
	"embed"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/fs/menu"
	"github.com/go-enjin/be/features/fs/public"
	"github.com/go-enjin/be/features/fs/themes"
)

//go:embed content/**
var contentFsWWW embed.FS

//go:embed public/**
var publicFs embed.FS

//go:embed menus/**
var menuFsWWW embed.FS

//go:embed themes/**
var themeFs embed.FS

func init() {
	fMenu = menu.New().
		MountEmbedPath("/", "menus", menuFsWWW).
		Make()
	fPublic = public.New().
		MountEmbedPath("/", "public", publicFs).
		Make()
	fContent = content.New().
		MountEmbedPath("/", "content", contentFsWWW).
		AddToIndexProviders("pages-pql").
		SetKeyValueCache(gFsContentKvsFeature, gFsContentKvsCache).
		Make()

	fThemes = themes.New().
		AddTheme(semantic.Theme()).
		EmbedTheme("thisip-fyi", "themes/thisip-fyi", themeFs).
		SetTheme("thisip-fyi").
		Make()

	hotReload = false
}