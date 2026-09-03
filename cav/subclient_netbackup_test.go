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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNetbackupContextData(t *testing.T) {
	n := &netbackup{}
	data := n.ContextData(t.Context())
	require.Equal(t, "", data.OrganizationID)
	require.Equal(t, "", data.SiteID)
}

func TestNetbackupTokenExpiry(t *testing.T) {
	n := &netbackup{}

	// No token set - should be expired
	require.True(t, n.isTokenExpired())

	// Set token with future expiry
	n.setTokens("access-token", "refresh-token", 3600)
	require.False(t, n.isTokenExpired())
	require.Equal(t, "access-token", n.getAccessToken())
	require.Equal(t, "refresh-token", n.getRefreshToken())

	// Set token with past expiry
	n.setTokens("access-token-2", "refresh-token-2", -1)
	require.True(t, n.isTokenExpired())
}

func TestNetbackupSessionPersistence(t *testing.T) {
	n := &netbackup{}
	n.setTokens("access-token", "refresh-token", 3600)
	n.baseURL = "https://backup1.cloudavenue.orange-business.com"

	session := n.getSession()
	require.Equal(t, "access-token", session["accessToken"])
	require.Equal(t, "refresh-token", session["refreshToken"])
	require.Equal(t, "https://backup1.cloudavenue.orange-business.com", session["baseURL"])

	// Restore into new instance
	n2 := &netbackup{}
	err := n2.restoreSession(session)
	require.NoError(t, err)
	require.Equal(t, "access-token", n2.getAccessToken())
	require.Equal(t, "refresh-token", n2.getRefreshToken())
	require.Equal(t, "https://backup1.cloudavenue.orange-business.com", n2.getBaseURL())
}

func TestNetbackupGetTokenMissingEnv(t *testing.T) {
	n := &netbackup{}
	// Ensure env vars are not set
	t.Setenv("NETBACKUP_USERNAME", "")
	t.Setenv("NETBACKUP_PASSWORD", "")

	err := n.getToken(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "NETBACKUP_USERNAME and NETBACKUP_PASSWORD environment variables must be set")
}

func TestNetbackupEnsureToken(t *testing.T) {
	n := &netbackup{}

	// No env vars, should fail
	t.Setenv("NETBACKUP_USERNAME", "")
	t.Setenv("NETBACKUP_PASSWORD", "")

	err := n.ensureToken(t.Context())
	require.Error(t, err)
}

func TestNetbackupSetTokensExpiryBuffer(t *testing.T) {
	n := &netbackup{}

	// Set token expiring in 3 minutes - should be considered expired due to 5-minute buffer
	n.setTokens("access", "refresh", 180)
	require.True(t, n.isTokenExpired(), "token expiring in 3 minutes should be considered expired due to 5-minute buffer")

	// Set token expiring in 10 minutes - should NOT be considered expired
	n.setTokens("access", "refresh", 600)
	require.False(t, n.isTokenExpired(), "token expiring in 10 minutes should not be considered expired")

	// Simulate token that expired 4 minutes ago - should be expired
	n.mu.Lock()
	n.tokenExpiresAt = time.Now().Add(-4 * time.Minute)
	n.mu.Unlock()
	require.True(t, n.isTokenExpired(), "token expired 4 minutes ago should be considered expired")
}
