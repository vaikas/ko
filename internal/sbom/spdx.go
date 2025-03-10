// Copyright 2021 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sbom

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

const dateFormat = "2006-01-02T15:04:05Z"

func GenerateSPDX(koVersion string, date time.Time, mod []byte) ([]byte, error) {
	bi := &BuildInfo{}
	if err := bi.UnmarshalText(mod); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplInfo{
		BuildInfo: *bi,
		Date:      date.Format(dateFormat),
		KoVersion: koVersion,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type tmplInfo struct {
	BuildInfo
	Date, UUID, KoVersion string
}

// TODO: use k8s.io/release/pkg/bom
var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"dots": func(s string) string { return strings.ReplaceAll(s, "/", ".") },
}).Parse(`SPDXVersion: SPDX-2.2
DataLicense: CC0-1.0
SPDXID: SPDXRef-DOCUMENT
DocumentName: {{ .BuildInfo.Main.Path }}
DocumentNamespace: http://spdx.org/spdxpackages/{{ .BuildInfo.Main.Path }}
Creator: Tool: ko {{ .KoVersion }}
Created: {{ .Date }}

##### Package representing {{ .BuildInfo.Main.Path }}

PackageName: {{ .BuildInfo.Main.Path }}
SPDXID: SPDXRef-Package-{{ .BuildInfo.Main.Path | dots }}
PackageSupplier: Organization: {{ .BuildInfo.Main.Path }}
PackageDownloadLocation: https://{{ .BuildInfo.Main.Path }}
FilesAnalyzed: false
PackageHomePage: https://{{ .BuildInfo.Main.Path }}
PackageLicenseConcluded: NOASSERTION
PackageLicenseDeclared: NOASSERTION
PackageCopyrightText: NOASSERTION
PackageLicenseComments: NOASSERTION
PackageComment: NOASSERTION

Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-Package-{{ .BuildInfo.Main.Path | dots }}
{{ range .Deps }}
Relationship: SPDXRef-Package-{{ $.Main.Path | dots }} DEPENDS_ON SPDXRef-Package-{{ .Path | dots }}-{{ .Version }}{{ end }}

{{ range .Deps }}
##### Package representing {{ .Path }}

PackageName: {{ .Path }}
SPDXID: SPDXRef-Package-{{ .Path | dots }}-{{ .Version }}
PackageVersion: {{ .Version }}
PackageSupplier: Organization: {{ .Path }}
PackageDownloadLocation: https://proxy.golang.org/{{ .Path }}/@v/{{ .Version }}.zip
FilesAnalyzed: false
PackageChecksum: SHA256: {{ .Sum }}
PackageLicenseConcluded: NOASSERTION
PackageLicenseDeclared: NOASSERTION
PackageCopyrightText: NOASSERTION
PackageLicenseComments: NOASSERTION
PackageComment: NOASSERTION

{{ end }}
`))
