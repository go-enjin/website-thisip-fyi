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

	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/log"
	"github.com/go-enjin/be/pkg/page"

	"github.com/go-enjin/website-thisip-fyi/pkg/features/whois"
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

	whois    map[string]*whois.Info
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
	f.whois = make(map[string]*whois.Info)
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
			whoisInfo, nslookup := f.lookupInfo(r.RemoteAddr)
			content += "\n\nnslookup: " + strings.Join(nslookup, ", ")
			content += "\n\nwhois:\n" + whoisInfo.Response
			p.Content = content
			p.Context.SetSpecific("Content", content)
			p.Context.SetSpecific("ContentDisposition", fmt.Sprintf(`attachment; filename="%s.txt"`, r.RemoteAddr))

		default:
			whoisInfo, nslookup := f.lookupInfo(r.RemoteAddr)
			if whoisInfo != nil {
				p.Context.SetSpecific("Whois", whoisInfo)
			}
			p.Context.SetSpecific("LookupAddr", nslookup)
		}
		p.Context.SetSpecific("Title", r.RemoteAddr)
		pg = p
		processed = true
	}
	return
}

func (f *CFeature) lookupInfo(addr string) (info *whois.Info, nslookup []string) {
	var ok bool
	var err error
	if nslookup, ok = f.nslookup[addr]; !ok {
		if nslookup, err = net.LookupAddr(addr); err != nil {
			log.WarnF("error net.LookupAddr: %v - %v", addr, err.Error())
			delete(f.nslookup, addr)
			nslookup = make([]string, 0)
		} else {
			f.nslookup[addr] = nslookup
		}
	}
	if info, ok = f.whois[addr]; !ok {
		if info, err = whois.LookupAndParse(addr); err != nil {
			log.WarnF("error whois.LookupAndParse: %v - %v", addr, err.Error())
			info = nil
		} else {
			f.whois[addr] = info
		}
	}
	return
}