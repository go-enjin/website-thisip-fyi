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
	"fmt"
	"os"

	"github.com/go-enjin/golang-org-x-text/language"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/drivers/kvs/gocache"
	"github.com/go-enjin/be/features/pages/pql"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/lang"
	"github.com/go-enjin/be/presets/defaults"
	"github.com/go-enjin/website-thisip-fyi/pkg/features/thisip"
)

var (
	gPublicActions = feature.Actions{
		feature.NewAction("enjin", "view", "page"),
		feature.NewAction("fs-content", "view", "page"),
	}
)

const (
	gPagesPqlFeature    = "pages-pql"
	gPagesPqlKvsFeature = "pages-pql-kvs-feature"
	gPagesPqlKvsCache   = "pages-pql-kvs-cache"
)

var (
	fThemes  feature.Feature
	fContent feature.Feature
	fPublic  feature.Feature
	fMenu    feature.Feature

	fListener feature.ServiceListener

	hotReload bool
)

func main() {
	enjin := be.New().
		SiteTag("TIPFYI").
		SiteName("ThisIp.Fyi").
		SiteTagLine("This IP for your information.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewPathMode().Make()).
		SiteSupportedLanguages(language.English).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo").
		AddPreset(defaults.New().SetListener(fListener).Make()).
		AddFeature(gocache.NewTagged(gPagesPqlKvsFeature).AddMemoryCache(gPagesPqlKvsCache).Make()).
		AddFeature(pql.NewTagged(gPagesPqlFeature).
			SetKeyValueCache(gPagesPqlKvsFeature, gPagesPqlKvsCache).
			Make()).
		AddFeature(fThemes).
		AddFeature(robots.New().
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddDisallowed("/").Make(),
			).Make()).
		AddFeature(thisip.New().Make()).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		SetPublicAccess(gPublicActions...).
		HotReload(hotReload).
		AddFeature(fMenu).
		AddFeature(fPublic).
		AddFeature(fContent)
	if err := enjin.Build().Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}