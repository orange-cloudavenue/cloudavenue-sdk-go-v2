/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"encoding/json"
	"os"
)

type cachedSessions struct {
	Sessions map[string]map[string]string `json:"sessions"`
}

func (c *client) storeSessionsToCache(passphrase, path string) error {
	cs := cachedSessions{
		Sessions: map[string]map[string]string{},
	}

	for _, value := range c.clientsInitialized {
		cs.Sessions[value.getID()] = value.getCredential().getSession()
	}

	csJSON, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	csEncrypted, err := encryptSessions([]byte(passphrase), string(csJSON))
	if err != nil {
		return err
	}

	return writeGobFile(path, csEncrypted)
}

func (c *client) restoreSessionsFromCache(passphrase, path string) error {
	c.cachePassphrase = passphrase
	c.cachePath = path

	// check if cache file exist. If not exist ignore
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.logger.Warn("Cache file does not exist, skipping restoration", "path", path)
		return nil
	}

	var csEncrypted string

	err := readGobFile(path, &csEncrypted)
	if err != nil {
		return err
	}

	csJSON, err := decryptSessions([]byte(passphrase), csEncrypted)
	if err != nil {
		return err
	}

	var cs cachedSessions
	if err := json.Unmarshal([]byte(csJSON), &cs); err != nil {
		return err
	}

	for id, session := range cs.Sessions {
		for _, value := range c.clientsInitialized {
			if value.getID() == id {
				if err := value.getCredential().restoreSession(session); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
