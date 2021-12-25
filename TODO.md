# Core Spec
Core specification includes anything that are not modules and the spec deem necessary for every client.

- [X] API Standards
- [X] Server Discovery
- [X] Client Auth
	- [X] Soft Logout
	- [X] User Interactive Authentication API (UIAA)
	- [X] Login
	- [X] Account Registration and Management
	- [X] Adding 3PIDs
	- [X] whoami
- [X] Capabilities Negotiation
- [X] Filtering
- [X] Events
	- [X] Room Events
	- [X] Syncing
	- [X] Getting Events for a Room
	- [X] Sending Events
	- [X] Redaction
- [X] Rooms
	- [X] Creation
	- [X] Room Aliases
	- [X] Room Membership
		- [X] Joining
		- [X] Leaving
		- [X] Kick
		- [X] Ban
		- [X] Unban
	- [X] Listing Rooms
- [X] User Data
- [X] Rate Limiting

# Modules
These are modules that are "optional" but somewhat required to figure out what is going on.

- [X] Instant Messaging
- [/] Voice over IP
	// Note: Need to provide user-friendly API.
- [X] Typing Notifications
- [X] Receipts
- [ ] Fully Read Marker
- [X] Presence (Online/Unavailable/Offline)
- [X] Content Repository
- [X] Send-to-Device Messaging
- [ ] Device Management
- [ ] End-to-end Encryption
- [ ] Secrets
- [X] History Visibility
- [ ] Push Notification
- [ ] Third Party Invites
- [ ] Server Side Search
- [X] Guest Access
- [ ] Room Previews
- [X] Room Tagging
- [X] Client Config
- [ ] Server Administration
- [ ] Event Context
- [X] SSO Client Login
- [X] Direct Messaging
- [X] Ignoring Users
- [ ] Sticker Messages
- [ ] Reporting Content
- [ ] Third Party Networks
- [ ] OpenID
- [ ] Server ACL
- [X] User, Room and Group mentions
- [X] Room Upgrades
- [ ] Server Notices
- [ ] Moderation Policy

# Migration
These are the changelogs of v1.1. They should be marked off when they're checked to be implemented.

## Breaking Changes
- [ ] MSC 2687: `curve25519-hkdf-sha256` for SAS verification
- [ ] MSC 3139: `m.key.verification.ready` and `m.key.verification.done`

## Deprecations
- [X] MSC 3199: Remove starting verification without `m.key.verification.ready`

## New Endpoints
- [ ] MSC 2387 + MSC 2639: `/room_keys/*`
- [ ] MSC 2536: `POST /keys/device_signing/upload` and `POST /keys/signatures/upload`
- [ ] MSC 3154 + MSC 3254: `/knock`
- [ ] MSC 3163: `/login/sso/redirect/{idpId}`

## Removed Endpoints
- [X] MSC 2609: `m.login.oauth2` and `m.login.token` User-Interactive Auth API

## Backwards Compatible Changes
- [ ] MSC 2399: Advise recipients about withholding keys
- [ ] MSC 2536: Cross-signing property to `POST /keys/query`
- [ ] MSC 2597 + MSC 3151: Secure Secret Storage and Sharing
- [ ] MSC 2709: `device_id` parameter to login fallback
- [ ] MSC 2728: SAS Emojis
- [ ] MSC 2795: `reason` on membership events
- [ ] MSC 2796: `M_NOT_FOUND` on push rule endpoints
- [ ] MSC 2807: Content reporting API: `reason`/`score` now optional
- [ ] MSC 2808: Guest may get list of members of a room
- [ ] MSC 3098: Support for spoilers
- [ ] MSC 3100: `<details>` and `<summary>` now in HTML subset
- [ ] MSC 3139 + MSC 3150: Key verification using in-room messages
- [ ] MSC 3147: SSSS for cross-signing and key backup
- [ ] MSC 3149: Key verification using QR code
- [ ] MSC 3163: Multiple SSO providers
- [ ] MSC 3166: `device_id` on `/account/whoami`
- [ ] MSC 3169: Identity server discovery failure results in `FAIL_PROMPT`
- [X] MSC 3421: Re-version to be `v3` instead of `r0`
