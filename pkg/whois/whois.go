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

package whois

import (
	"regexp"
	"strings"

	"github.com/likexian/whois"
)

var (
	rxEmptyLine    = regexp.MustCompile(`^\s*$`)
	rxCommentLine  = regexp.MustCompile(`^\s*#`)
	rxKeyValueLine = regexp.MustCompile(`^\s*([^:]+):\s*(.+?)\s*$`)
	rxNetworkLine  = regexp.MustCompile(`^\s*network:(.+?):(.+?)\s*$`)
	rxWhenTimeLine = regexp.MustCompile(`^\s*;;\s*(WHEN|Query\s*time):\s*(.+?)\s*$`)
)

type Record struct {
	Key   string
	Value string
}

type Info struct {
	Response string
	Comments string
	records  []interface{}
	Lookup   map[string]string
	Groups   []map[string]string
}

func (info *Info) Value(key string) (value string, ok bool) {
	for _, record := range info.records {
		if kv, found := record.(Record); found {
			if ok = kv.Key == key; ok {
				value = kv.Value
				return
			}
		}
	}
	return
}

func LookupAndParse(addr string) (info *Info, err error) {
	var response string
	if response, err = LookupIP(addr); err != nil {
		return
	}
	info = ParseResponse(response)
	return
}

func LookupIP(addr string) (response string, err error) {
	response, err = whois.Whois(addr)
	return
}

func ParseResponse(response string) (info *Info) {
	info = new(Info)
	info.Response = response
	info.Lookup = make(map[string]string)
	var lastGroup map[string]string
	var lastComment string
	var consumeComment bool
	var consumeNetwork bool

	setKeyValue := func(key, value string) {
		if rxEmptyLine.MatchString(value) {
			value = ""
		}
		if lastGroup == nil {
			lastGroup = make(map[string]string)
		}
		lastGroup[key] = value
		info.Lookup[key] = value
		info.records = append(info.records, Record{
			Key:   key,
			Value: value,
		})
	}

	for _, line := range strings.Split(response, "\n") {
		switch {

		case rxEmptyLine.MatchString(line):
			if consumeComment {
				if info.Comments == "" {
					info.Comments = lastComment
				}
				info.records = append(info.records, lastComment)
				lastComment = ""
				consumeComment = false
			}
			if lastGroup != nil {
				info.Groups = append(info.Groups, lastGroup)
				lastGroup = nil
			}

		case rxWhenTimeLine.MatchString(line):
			m := rxWhenTimeLine.FindAllStringSubmatch(line, 1)
			setKeyValue(m[0][1], m[0][2])

		case strings.HasPrefix(line, "%") && !consumeNetwork:
			consumeNetwork = true
			info.records = append(info.records, line)

		case strings.HasPrefix(line, "%ok"):
			consumeNetwork = false
			info.records = append(info.records, line)

		case rxCommentLine.MatchString(line):
			consumeComment = true
			lastComment += line + "\n"

		case consumeNetwork && rxNetworkLine.MatchString(line):
			m := rxNetworkLine.FindAllStringSubmatch(line, 1)
			setKeyValue(m[0][1], m[0][2])

		case !consumeNetwork && rxKeyValueLine.MatchString(line):
			m := rxKeyValueLine.FindAllStringSubmatch(line, 1)
			setKeyValue(m[0][1], m[0][2])

		}

	}
	if lastComment != "" {
		info.records = append(info.records, lastComment)
	}
	return
}
