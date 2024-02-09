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

//go:build (fs_theme && (drivers_fs_local || drivers_fs || drivers || locals) && !drivers_fs_embed && !embeds) || all

package themes

import (
	"path/filepath"
	"runtime"

	"github.com/go-enjin/be/features/fs/themes"
)

func Theme() themes.Feature {
	_, fn, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(fn), Name)
	return themes.
		NewTagged(Name).
		LocalTheme(path).
		Make()
}
