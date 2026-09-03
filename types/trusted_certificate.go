/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type ParamsListTrustedCertificate struct {
	Name string
}

type ParamsGetTrustedCertificate struct {
	ID   string
	Name string
}

type ParamsCreateTrustedCertificate struct {
	Name        string
	Certificate string
}

type ParamsUpdateTrustedCertificate struct {
	ID          string
	LookupName  string
	Name        string
	Certificate string
}

type ParamsDeleteTrustedCertificate struct {
	ID   string
	Name string
}

type ModelListTrustedCertificate struct {
	Certificates []ModelGetTrustedCertificate
}

type ModelGetTrustedCertificate struct {
	ID          string `documentation:"URN of trusted certificate"`
	Name        string `documentation:"Trusted certificate alias"`
	Certificate string `documentation:"PEM encoded certificate"`
}
