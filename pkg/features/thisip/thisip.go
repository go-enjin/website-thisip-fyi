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

package thisip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/likexian/whois"

	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/page"
)

var (
	_ MakeFeature               = (*CFeature)(nil)
	_ feature.PageTypeProcessor = (*CFeature)(nil)
)

const Tag feature.Tag = "PagesThisIp"

type Feature interface {
	feature.Feature
}

type CFeature struct {
	feature.CFeature

	cli   *cli.Context
	enjin feature.Internals

	whois    map[string]string
	nslookup map[string][]string

	sync.RWMutex
}

type MakeFeature interface {
	Make() Feature
}

func New() MakeFeature {
	f := new(CFeature)
	f.Init(f)
	return f
}

func (f *CFeature) Make() Feature {
	return f
}

func (f *CFeature) Init(this interface{}) {
	f.CFeature.Init(this)
	f.whois = make(map[string]string)
	f.nslookup = make(map[string][]string)
}

func (f *CFeature) Tag() (tag feature.Tag) {
	tag = Tag
	return
}

func (f *CFeature) Setup(enjin feature.Internals) {
	f.enjin = enjin
}

func (f *CFeature) Startup(ctx *cli.Context) (err error) {
	f.cli = ctx
	return
}

func (f *CFeature) ProcessRequestPageType(r *http.Request, p *page.Page) (pg *page.Page, redirect string, processed bool, err error) {
	// reqArgv := site.GetRequestArgv(r)
	if p.Type == "thisip" {
		userAgent := r.UserAgent()
		switch {

		case strings.HasPrefix(userAgent, "curl/"):
			p = p.Copy()
			p.Layout = "none"
			p.Context.SetSpecific("Layout", "none")
			p.Context.SetSpecific("ContentType", "text/plain; charset=utf-8")
			p.Content = r.RemoteAddr
			p.Context.SetSpecific("Content", r.RemoteAddr)

		case strings.HasPrefix(userAgent, "Wget/") || strings.HasPrefix(userAgent, "wget/"):
			p = p.Copy()
			p.Layout = "none"
			p.Context.SetSpecific("Layout", "none")
			p.Context.SetSpecific("ContentType", "text/plain; charset=utf-8")
			content := "address: " + r.RemoteAddr
			if values, ok := f.nslookup[r.RemoteAddr]; ok {
				content += "\n" + fmt.Sprintf("domains: %v", strings.Join(values, ", "))
			} else if names, ee := net.LookupAddr(r.RemoteAddr); ee == nil {
				content += "\n" + fmt.Sprintf("domains: %v", strings.Join(names, ", "))
				f.nslookup[r.RemoteAddr] = names
			}
			if value, ok := f.whois[r.RemoteAddr]; ok {
				content += "\nwhois:" + value
			} else if results, ee := whois.Whois(r.RemoteAddr); ee == nil {
				content += "\nwhois:" + results
				f.whois[r.RemoteAddr] = results
			}
			p.Content = content
			p.Context.SetSpecific("Content", content)
			p.Context.SetSpecific("ContentDisposition", fmt.Sprintf(`attachment; filename="%s.txt"`, r.RemoteAddr))

		default:
			if values, ok := f.nslookup[r.RemoteAddr]; ok {
				p.Context.SetSpecific("LookupAddr", values)
			} else if names, ee := net.LookupAddr(r.RemoteAddr); ee != nil {
				p.Context.SetSpecific("LookupAddrError", ee.Error())
			} else {
				p.Context.SetSpecific("LookupAddr", names)
				f.nslookup[r.RemoteAddr] = names
			}
			if value, ok := f.whois[r.RemoteAddr]; ok {
				p.Context.SetSpecific("Whois", value)
			} else if results, ee := whois.Whois(r.RemoteAddr); ee != nil {
				p.Context.SetSpecific("WhoisError", ee.Error())
			} else {
				p.Context.SetSpecific("Whois", results)
				f.whois[r.RemoteAddr] = results
			}
		}
		p.Context.SetSpecific("Title", r.RemoteAddr)
		pg = p
		processed = true
	}
	return
}