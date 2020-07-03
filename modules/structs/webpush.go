// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package structs

// WebPushSubscription represents a HTML5 Web Push Subscription used for background notifications.
type WebPushSubscription struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	P256DH   string `json:"p256dh"`
}

// WebPushPayload marks a JSON payload sent in a push event to the JS service worker.
// This is used for background notifications.
type WebPushPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}
