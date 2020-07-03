// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func addWebPushSubscriptionTable(x *xorm.Engine) error {
	type WebPushSubscription struct {
		ID     int64 `xorm:"pk autoincr"`
		UserID int64 `xorm:"INDEX UNIQUE(s)"`

		Endpoint string `xorm:"UNIQUE(s)"`
		Auth     string
		P256DH   string

		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}
	return x.Sync2(new(WebPushSubscription))
}
