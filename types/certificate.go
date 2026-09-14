/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type ParamsListCertificate struct {
	Name string
}

type ParamsGetCertificate struct {
	ID   string
	Name string
}

type ParamsCreateCertificate struct {
	Name                 string
	Description          string
	Certificate          string
	PrivateKey           string
	PrivateKeyPassphrase string
}

type ParamsUpdateCertificate struct {
	ID          string
	LookupName  string
	Name        string
	Description string
}

type ParamsDeleteCertificate struct {
	ID   string
	Name string
}

type ParamsListCertificateConsumers struct {
	CertificateID   string
	CertificateName string
}

type ParamsAddCertificateConsumer struct {
	CertificateID   string
	CertificateName string
	ConsumerID      string
	ConsumerName    string
}

type ParamsSetCertificateConsumers struct {
	CertificateID   string
	CertificateName string
	Consumers       []ModelEntityReference
}

type ModelListCertificate struct {
	Certificates []ModelGetCertificate
}

type ModelGetCertificate struct {
	ID            string `documentation:"URN of the certificate library item"`
	Name          string `documentation:"Certificate alias"`
	Description   string `documentation:"Certificate description"`
	Certificate   string `documentation:"PEM encoded certificate"`
	ConsumerCount int    `documentation:"Number of consumers of the certificate"`
}

type ModelEntityReference struct {
	ID   string `documentation:"Entity URN"`
	Name string `documentation:"Entity name"`
}

type ModelListEntityReference struct {
	References []ModelEntityReference
}
