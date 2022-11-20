//go:build locals || !embeds

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
	"github.com/go-enjin/be/features/fs/locals/content"
	"github.com/go-enjin/be/features/fs/locals/locales"
	"github.com/go-enjin/be/features/fs/locals/menu"
	"github.com/go-enjin/be/features/fs/locals/public"
	"github.com/go-enjin/be/pkg/log"
	"github.com/go-enjin/be/pkg/theme"
)

func init() {
	// locals environment, early startup debug logging
	// log.Config.LogLevel = log.LevelDebug
	// log.Config.Apply()

	fMenu = menu.New().MountPath("menus", "menus").Make()
	fPublic = public.New().MountPath("/", "public").Make()
	fContent = content.New().MountPath("/", "content").Make()
	fLocales = locales.New().Include("locales").Make()

	hotReload = true
}

func thisIpFyiTheme() (t *theme.Theme) {
	var err error
	if t, err = theme.NewLocal("themes/thisip-fyi"); err != nil {
		log.FatalF("error loading local theme: %v", err)
	} else {
		log.DebugF("loaded local theme: %v", t.Name)
	}
	return
}
