// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package archive

import (
	"go.uber.org/zap"
	"path/filepath"
	"regexp"
	"sort"
)

func (p *IndexPool) listAllReadOnlyIndexes() ([]string, map[string]bool, error) {
	local, err := filepath.Glob(filepath.Join(p.IndexesPath, "??????????.bleve"))
	if err != nil {
		return nil, nil, err
	}

	for _, readOnlyPath := range p.ReadOnlyIndexesPaths {
		more, err := filepath.Glob(filepath.Join(readOnlyPath, "??????????.bleve"))
		if err != nil {
			zlog.Warn("failed listing files in read-only path, continuing", zap.String("path", readOnlyPath), zap.Error(err))
			continue
		}
		local = append(local, more...)
	}

	// dedupe
	seen := map[string]bool{}
	var sorted []string
	for _, el := range local {
		baseDir := toIndexBase(el)
		if seen[baseDir] {
			continue
		}
		sorted = append(sorted, baseDir)
		seen[baseDir] = true
	}

	sort.Strings(sorted)
	return sorted, seen, nil
}

var localPathRE = regexp.MustCompile(`(\d{10})\.bleve`)

func toIndexBase(indexPath string) string {
	match := localPathRE.FindStringSubmatch(indexPath)
	if match == nil {
		return ""
	}
	return match[1]
}
