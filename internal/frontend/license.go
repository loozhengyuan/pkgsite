// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package frontend

import (
	"context"
	"net/url"

	"golang.org/x/discovery/internal"
	"golang.org/x/discovery/internal/license"
)

// License contains information used for a single license section.
type License struct {
	*license.License
	Anchor string
}

// LicensesDetails contains license information for a package or module.
type LicensesDetails struct {
	Licenses []License
}

// LicenseMetadata contains license metadata that is used in the package
// header.
type LicenseMetadata struct {
	Type   string
	Anchor string
}

// fetchPackageLicensesDetails fetches license data for the package version specified by
// path and version from the database and returns a LicensesDetails.
func fetchPackageLicensesDetails(ctx context.Context, ds DataSource, pkg *internal.VersionedPackage) (*LicensesDetails, error) {
	dsLicenses, err := ds.GetPackageLicenses(ctx, pkg.Path, pkg.ModulePath, pkg.VersionInfo.Version)
	if err != nil {
		return nil, err
	}
	return &LicensesDetails{Licenses: transformLicenses(dsLicenses)}, nil
}

// transformLicenses transforms license.License into a License
// by adding an anchor field.
func transformLicenses(dbLicenses []*license.License) []License {
	licenses := make([]License, len(dbLicenses))
	for i, l := range dbLicenses {
		licenses[i] = License{
			Anchor:  licenseAnchor(l.FilePath),
			License: l,
		}
	}
	return licenses
}

// transformLicenseMetadata transforms license.Metadata into a LicenseMetadata
// by adding an anchor field.
func transformLicenseMetadata(dbLicenses []*license.Metadata) []LicenseMetadata {
	var mds []LicenseMetadata
	for _, l := range dbLicenses {
		anchor := licenseAnchor(l.FilePath)
		for _, typ := range l.Types {
			mds = append(mds, LicenseMetadata{
				Type:   typ,
				Anchor: anchor,
			})
		}
	}
	return mds
}

// licenseAnchor returns the anchor that should be used to jump to the specific
// license on the licenses page.
func licenseAnchor(filePath string) string {
	return url.QueryEscape(filePath)
}
