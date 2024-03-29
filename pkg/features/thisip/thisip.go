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

	"github.com/urfave/cli/v2"

	"github.com/go-corelibs/x-text/message"
	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/log"
	beNet "github.com/go-enjin/be/pkg/net"

	"github.com/go-enjin/website-thisip-fyi/pkg/whois"
)

var (
	_ Feature     = (*CFeature)(nil)
	_ MakeFeature = (*CFeature)(nil)
)

const Tag feature.Tag = "pages-thisip"

type Feature interface {
	feature.Feature
	feature.PageTypeProcessor
}

type MakeFeature interface {
	Make() Feature
}

type CFeature struct {
	feature.CFeature

	cli   *cli.Context
	enjin feature.Internals

	whois    map[string]*whois.Info
	nslookup map[string][]string
}

func New() MakeFeature {
	return NewTagged(Tag)
}

func NewTagged(tag feature.Tag) MakeFeature {
	f := new(CFeature)
	f.Init(f)
	f.PackageTag = Tag
	f.FeatureTag = tag
	return f
}

func (f *CFeature) Init(this interface{}) {
	f.CFeature.Init(this)
	f.whois = make(map[string]*whois.Info)
	f.nslookup = make(map[string][]string)
}

func (f *CFeature) Make() Feature {
	return f
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

func (f *CFeature) PageTypeNames() (names []string) {
	names = []string{"thisip"}
	return
}

func (f *CFeature) ProcessRequestPageType(r *http.Request, p feature.Page) (pg feature.Page, redirect string, processed bool, err error) {
	// reqArgv := site.GetRequestArgv(r)
	if p.Type() == "thisip" {
		userAgent := r.UserAgent()
		switch {

		case strings.HasPrefix(userAgent, "curl/"):
			p = p.Copy()
			p.SetLayout("none")
			p.Context().SetSpecific("Layout", "none")
			p.Context().SetSpecific("ContentType", "text/plain; charset=utf-8")
			p.SetContent(r.RemoteAddr)
			p.Context().SetSpecific("Content", r.RemoteAddr)

		case strings.HasPrefix(userAgent, "Wget/") || strings.HasPrefix(userAgent, "wget/"):
			p = p.Copy()
			p.SetLayout("none")
			p.Context().SetSpecific("Layout", "none")
			p.Context().SetSpecific("ContentType", "text/plain; charset=utf-8")
			content := "address: " + r.RemoteAddr
			whoisInfo, nslookup := f.lookupInfo(r.RemoteAddr, r)
			content += "\n\nnslookup: " + strings.Join(nslookup, ", ")
			if whoisInfo != nil {
				content += "\n\nwhois:\n" + whoisInfo.Response
			}
			p.SetContent(content)
			p.Context().SetSpecific("Content", content)
			p.Context().SetSpecific("ContentDisposition", fmt.Sprintf(`attachment; filename="%s.txt"`, r.RemoteAddr))

		default:
			whoisInfo, nslookup := f.lookupInfo(r.RemoteAddr, r)
			if whoisInfo != nil {
				p.Context().SetSpecific("Whois", whoisInfo)
			}
			p.Context().SetSpecific("LookupAddr", nslookup)
		}
		p.Context().SetSpecific("Title", r.RemoteAddr)
		pg = p
		processed = true
	}
	return
}

func (f *CFeature) lookupInfo(addr string, r *http.Request) (info *whois.Info, nslookup []string) {
	var ok bool
	var err error

	printer := message.GetPrinter(r)

	if beNet.IsNetIpPrivate(net.ParseIP(addr)) {
		info = nil
		nslookup = append(nslookup, printer.Sprintf("(local/private address)"))
		return
	}

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
