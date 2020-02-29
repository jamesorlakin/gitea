// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"fmt"

	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func addHTTPDeployKeys(x *xorm.Engine) error {
	// HTTPDeployKey see models/http_deploy_key.go
	type HTTPDeployKey struct {
		ID           int64 `xorm:"pk autoincr"`
		RepositoryID int64 `xorm:"UNIQUE(s) INDEX NOT NULL"`
		Name         string
		KeyContent   string             `xorm:"UNIQUE(s) INDEX NOT NULL"`
		CreatedUnix  timeutil.TimeStamp `xorm:"created"`
	}

	if err := x.Sync2(new(HTTPDeployKey)); err != nil {
		return fmt.Errorf("Sync2: %v", err)
	}
	return nil
}
