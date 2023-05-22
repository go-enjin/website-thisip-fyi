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

	"github.com/go-enjin/be/pkg/userbase"
	"github.com/go-enjin/golang-org-x-text/language"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/drivers/kvs/gocache"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/outputs/htmlify"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/features/pages/pql"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	"github.com/go-enjin/be/features/user/auth/basic"
	"github.com/go-enjin/be/features/user/base/htenv"
	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/lang"

	"github.com/go-enjin/website-thisip-fyi/pkg/features/thisip"
)

var (
	fThemes  feature.Feature
	fContent feature.Feature
	fPublic  feature.Feature
	fMenu    feature.Feature

	fCachePagesPql feature.Feature

	hotReload bool
)

func init() {
	fCachePagesPql = gocache.NewTagged(gPagesPqlKvsFeature).AddMemoryCache(gPagesPqlKvsCache).Make()
}

func setup(eb *be.EnjinBuilder) *be.EnjinBuilder {
	eb.SiteName("ThisIp.Fyi").
		SiteTagLine("This IP for your information.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		AddFeature(fCachePagesPql).
		AddFeature(pql.NewTagged("pages-pql").
			SetKeyValueCache(gPagesPqlKvsFeature, gPagesPqlKvsCache).
			Make()).
		AddFeature(formats.New().Defaults().Make()).
		AddFeature(fThemes).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo")
	return eb
}

func features(eb feature.Builder) feature.Builder {
	return eb.
		AddFeature(papertrail.Make()).
		AddFeature(htmlify.New().Make()).
		AddFeature(proxy.New().Enable().Make()).
		AddFeature(htenv.NewTagged("htenv").Make()).
		AddFeature(basic.New().
			AddUserbase("htenv", "htenv", "htenv").
			Ignore(`^/favicon.ico$`).
			Make()).
		AddFeature(robots.New().
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddDisallowed("/").Make(),
			).Make()).
		AddFeature(thisip.New().Make()).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		SetPublicAccess(gPublicActions...).
		HotReload(hotReload)
}

func main() {
	enjin := be.New()
	setup(enjin).SiteTag("TIPFYI").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewPathMode().Make()).
		SiteSupportedLanguages(language.English)
	features(enjin).
		AddFeature(fMenu).
		AddFeature(fPublic).
		AddFeature(fContent)
	if err := enjin.Build().Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

var (
	gPublicActions = userbase.Actions{
		userbase.NewAction("fs-content", "view", "page"),
	}
)

const (
	gFsContentKvsFeature = "fs-content-kvs-feature"
	gFsContentKvsCache   = "fs-content-kvs-cache"
	gPagesPqlKvsFeature  = "pages-pql-kvs-feature"
	gPagesPqlKvsCache    = "pages-pql-kvs-cache"

	main500tmpl = `500 - Internal Server Error`
	main404tmpl = `404 - Not Found`
	main204tmpl = `+++
url = "/"
+++
204 - {{ _ "No Content" }}`
)