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

	"github.com/go-enjin/be/features/pages/formats/html"
	"github.com/go-enjin/be/features/pages/formats/tmpl"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	auth "github.com/go-enjin/be/features/restrict/basic-auth"

	"github.com/go-enjin/be/features/outputs/htmlify"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/pkg/lang"
	"github.com/go-enjin/golang-org-x-text/language"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/pkg/feature"

	"github.com/go-enjin/website-thisip-fyi/pkg/features/thisip"
)

var (
	fLocales feature.Feature
	fContent feature.Feature
	fPublic  feature.Feature
	fMenu    feature.Feature

	hotReload bool
)

const (
	main500tmpl = `500 - Internal Server Error`
	main404tmpl = `404 - Not Found`
	main204tmpl = `+++
url = "/"
+++
204 - {{ _ "No Content" }}`
)

func setup(eb *be.EnjinBuilder) *be.EnjinBuilder {
	eb.SiteName("ThisIp.Fyi").
		SiteTagLine("This IP for your information.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		AddFeature(proxy.New().Enable().Make()).
		AddFeature(
			formats.New().
				AddFormat(html.New().Make()).
				AddFormat(tmpl.New().Make()).
				Make(),
		).
		AddTheme(semantic.SemanticEnjinTheme()).
		AddTheme(thisIpFyiTheme()).
		SetTheme("thisip-fyi").
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo")
	return eb
}

func features(eb feature.Builder) feature.Builder {
	return eb.
		AddFeature(auth.New().EnableEnv(true).Make()).
		AddFeature(papertrail.Make()).
		AddFeature(robots.New().
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddDisallowed("/").Make(),
			).Make()).
		AddFeature(thisip.New().Make()).
		AddFeature(htmlify.New().Make()).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
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
		AddFeature(fContent).
		AddFeature(fLocales)
	if err := enjin.Build().Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}