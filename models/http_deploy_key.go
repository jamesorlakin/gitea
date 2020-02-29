// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
)

// HTTPDeployKey is a form of read-only UUID deploy key used for HTTP authentication to a single repo.
// This is separate to SSH-based deploy keys, defined in ssh_key.go
type HTTPDeployKey struct {
	ID           int64 `xorm:"pk autoincr"`
	RepositoryID int64 `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Name         string
	KeyContent   string             `xorm:"UNIQUE(s) INDEX NOT NULL"`
	CreatedUnix  timeutil.TimeStamp `xorm:"created"`
}
