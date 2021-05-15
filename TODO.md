# Core Spec
Core specification includes anything that are not modules and the spec deem necessary for every client.

[X] API Standards
[X] Server Discovery
[X] Client Auth
	[X] Soft Logout
	[X] User Interactive Authentication API (UIAA)
	[X] Login
	[X] Account Registration and Management
	[X] Adding 3PIDs
	[X] whoami
[X] Capabilities Negotiation
[X] Filtering
[X] Events
	[X] Room Events
	[X] Syncing
	[X] Getting Events for a Room
	[X] Sending Events
	[X] Redaction
[X] Rooms
	[X] Creation
	[X] Room Aliases
	[X] Room Membership
		[X] Joining
		[X] Leaving
		[X] Kick
		[X] Ban
		[X] Unban
	[X] Listing Rooms
[X] User Data
[X] Rate Limiting

# Modules
These are modules that are "optional" but somewhat required to figure out what is going on.

[X] Instant Messaging
[/] Voice over IP
	// Note: Need to provide user-friendly API.
[X] Typing Notifications
[X] Receipts
[ ] Fully Read Marker
[X] Presence (Online/Unavailable/Offline)
[X] Content Repository
[ ] Send-to-Device Messaging
[ ] Device Management
[ ] End-to-end Encryption
	[ ] Key Upload
	[ ] Key Fetch
	[ ] Key Claim
	[ ] Key Verification
	[ ] Key Sharing
	[ ] Curve25519 Message Encryption
	[ ] MegOlm Message Encryption
[ ] History Visibility
[ ] Push Notification
[ ] Third Party Invites
[ ] Server Side Search
[ ] Guest Access
[ ] Room Previews
[ ] Room Tagging
[X] Client Config
[ ] Server Administration
[ ] Event Context
[ ] SSO Client Login
[ ] Direct Messaging
[ ] Ignoring Users
[ ] Sticker Messages
[ ] Reporting Content
[ ] Third Party Networks
[ ] OpenID
[ ] Server ACL
[ ] User, Room and Group mentions
[ ] Room Upgrades
[ ] Server Notices
[ ] Moderation Policy
