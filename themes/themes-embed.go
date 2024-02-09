// Copyright (c) 2024  The Go-Enjin Authors
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

//go:build fs_theme && !all && (drivers_fs_embed || drivers_fs || drivers || embeds)

package themes

import (
	"embed"

	"github.com/go-enjin/be/features/fs/themes"
)

//go:embed thisip-fyi/**
var themeFS embed.FS

func Theme() themes.Feature {
	return themes.
		NewTagged(Name).
		EmbedTheme(Name, themeFS).
		Make()
}