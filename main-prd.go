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

package thisip_fyi

import (
	"embed"

	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/fs/menu"
	"github.com/go-enjin/be/features/fs/public"
)

//go:embed content/**
var contentFsWWW embed.FS

//go:embed public/**
var publicFs embed.FS

//go:embed menus/**
var menuFsWWW embed.FS

func init() {

	fMenu = menu.New().
		MountEmbedPath("/", "menus", menuFsWWW).
		Make()
	fPublic = public.New().
		MountEmbedPath("/", "public", publicFs).
		Make()
	fContent = content.New().
		MountEmbedPath("/", "content", contentFsWWW).
		AddToIndexProviders(gPagesPqlFeature).
		Make()

	hotReload = false
}
