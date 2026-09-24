/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/xml"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// AdminCatalog represents a VMware vCD AdminCatalog XML element.
type AdminCatalog struct {
	XMLName      xml.Name    `xml:"Catalog"`
	Name         string      `xml:"name,attr"`
	ID           string      `xml:"id,attr"`
	Href         string      `xml:"href,attr"`
	Description  string      `xml:"Description,omitempty"`
	IsEnabled    bool        `xml:"IsEnabled,omitempty"`
	Owner        Reference   `xml:"Owner"`
	Org          Reference   `xml:"Org"`
	CatalogItems []Reference `xml:"CatalogItems>CatalogItem"`
	IsPublished  bool        `xml:"IsPublished"`
	IsTrusted    bool        `xml:"IsTrusted"`
}

// AdminCatalogRequest represents the request body for creating or updating a VMware vCD AdminCatalog.
type AdminCatalogRequest struct {
	XMLName          xml.Name    `xml:"Catalog"`
	Name             string      `xml:"name,attr"`
	Description      string      `xml:"Description,omitempty"`
	StorageProfiles  []Reference `xml:"StorageProfiles>StorageProfile,omitempty"`
}

// AdminCatalogs represents the wrapper for a list of catalogs in VMware vCD AdminCatalog API.
type AdminCatalogs struct {
	XMLName  xml.Name       `xml:"Catalogs"`
	Catalogs []AdminCatalog `xml:"Catalog"`
}

// ToModel converts the VMware vCD AdminCatalog XML response to the public model.
func (r *AdminCatalog) ToModel() *types.ModelAdminCatalog {
	return &types.ModelAdminCatalog{
		ID:          r.ID,
		Name:        r.Name,
		Href:        r.Href,
		Description: r.Description,
		IsEnabled:   r.IsEnabled,
		IsPublished: r.IsPublished,
		IsTrusted:   r.IsTrusted,
	}
}
