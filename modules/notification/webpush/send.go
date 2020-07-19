// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webpush

import (
	"strconv"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"

	"github.com/unknwon/i18n"
)

func handleNotification(opts *issueNotificationOpts) {
	issue, err := models.GetIssueByID(opts.IssueID)
	if err != nil {
		log.Error("unable to get issue %d: %v", opts.IssueID, err)
		return
	}

	userIDs, err := models.GetEligibleNotificationParticipants(issue, opts.NotificationAuthorID, opts.ReceiverID)
	if err != nil {
		log.Error("unable to get user IDs for web push notifications: %v", err)
		return
	}

	// If there's a comment ID, send a link pointing to that comment
	var anchorLink string
	if opts.CommentID != 0 {
		anchorLink = "#issuecomment-" + strconv.FormatInt(opts.CommentID, 10)
	}

	for userID := range userIDs {
		user, err := models.GetUserByID(userID)
		if err != nil {
			log.Error("unable to get user %d: %v", userID, err)
			continue
		}

		var notificationText string
		if issue.IsPull {
			notificationText = i18n.Tr(user.Language, "pushnotification.activity_pr", issue.Index, issue.Title)
		} else {
			notificationText = i18n.Tr(user.Language, "pushnotification.activity_issue", issue.Index, issue.Title)
		}
		notificationPayload := &structs.WebPushNotificationPayload{
			Title: setting.AppName + " - " + issue.Repo.MustOwner().Name + "/" + issue.Repo.Name,
			Text:  notificationText,
			URL:   issue.HTMLURL() + anchorLink,
		}

		err = models.SendWebPushNotificationToUser(userID, notificationPayload)
		if err != nil {
			log.Error("error sending web push notification to user %d: %v", userID, err)
		}
	}
}
