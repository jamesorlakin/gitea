// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webpush

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/graceful"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/notification/base"
	"code.gitea.io/gitea/modules/queue"
)

type (
	// Web Push service deals with sending push notifications to an external browser-provided notification service.
	webpushService struct {
		base.NullNotifier
		webpushQueue queue.Queue
	}

	issueNotificationOpts struct {
		IssueID              int64
		CommentID            int64
		NotificationAuthorID int64
		ReceiverID           int64 // 0 -- ALL Watcher
		IsPull               bool
	}
)

var (
	_ base.Notifier = &webpushService{}
)

// NewNotifier create a new webpushService notifier
func NewNotifier() base.Notifier {
	ns := &webpushService{}
	ns.webpushQueue = queue.CreateQueue("webpush-service", ns.handle, issueNotificationOpts{})
	return ns
}

func (ns *webpushService) handle(data ...queue.Data) {
	for _, datum := range data {
		opts := datum.(issueNotificationOpts)
		handleNotification(&opts)
	}
}

func (ns *webpushService) Run() {
	graceful.GetManager().RunWithShutdownFns(ns.webpushQueue.Run)
}

func (ns *webpushService) NotifyCreateIssueComment(doer *models.User, repo *models.Repository,
	issue *models.Issue, comment *models.Comment) {
	var opts = issueNotificationOpts{
		IssueID:              issue.ID,
		NotificationAuthorID: doer.ID,
		IsPull:               issue.IsPull,
	}
	if comment != nil {
		opts.CommentID = comment.ID
	}
	_ = ns.webpushQueue.Push(opts)
}

func (ns *webpushService) NotifyNewIssue(issue *models.Issue) {
	_ = ns.webpushQueue.Push(issueNotificationOpts{
		IssueID:              issue.ID,
		NotificationAuthorID: issue.Poster.ID,
		IsPull:               issue.IsPull,
	})
}

func (ns *webpushService) NotifyIssueChangeStatus(doer *models.User, issue *models.Issue, actionComment *models.Comment, isClosed bool) {
	_ = ns.webpushQueue.Push(issueNotificationOpts{
		IssueID:              issue.ID,
		NotificationAuthorID: doer.ID,
		IsPull:               issue.IsPull,
	})
}

func (ns *webpushService) NotifyMergePullRequest(pr *models.PullRequest, doer *models.User) {
	_ = ns.webpushQueue.Push(issueNotificationOpts{
		IssueID:              pr.Issue.ID,
		NotificationAuthorID: doer.ID,
		IsPull:               true,
	})
}

func (ns *webpushService) NotifyNewPullRequest(pr *models.PullRequest) {
	if err := pr.LoadIssue(); err != nil {
		log.Error("Unable to load issue: %d for pr: %d: Error: %v", pr.IssueID, pr.ID, err)
		return
	}
	_ = ns.webpushQueue.Push(issueNotificationOpts{
		IssueID:              pr.Issue.ID,
		NotificationAuthorID: pr.Issue.PosterID,
		IsPull:               true,
	})
}

func (ns *webpushService) NotifyPullRequestReview(pr *models.PullRequest, r *models.Review, c *models.Comment) {
	var opts = issueNotificationOpts{
		IssueID:              pr.Issue.ID,
		NotificationAuthorID: r.Reviewer.ID,
		IsPull:               true,
	}
	if c != nil {
		opts.CommentID = c.ID
	}
	_ = ns.webpushQueue.Push(opts)
}

func (ns *webpushService) NotifyPullRequestPushCommits(doer *models.User, pr *models.PullRequest, comment *models.Comment) {
	var opts = issueNotificationOpts{
		IssueID:              pr.IssueID,
		NotificationAuthorID: doer.ID,
		CommentID:            comment.ID,
		IsPull:               true,
	}
	_ = ns.webpushQueue.Push(opts)
}

func (ns *webpushService) NotifyIssueChangeAssignee(doer *models.User, issue *models.Issue, assignee *models.User, removed bool, comment *models.Comment) {
	if !removed {
		var opts = issueNotificationOpts{
			IssueID:              issue.ID,
			NotificationAuthorID: doer.ID,
			ReceiverID:           assignee.ID,
			IsPull:               issue.IsPull,
		}

		if comment != nil {
			opts.CommentID = comment.ID
		}

		_ = ns.webpushQueue.Push(opts)
	}
}

func (ns *webpushService) NotifyPullReviewRequest(doer *models.User, issue *models.Issue, reviewer *models.User, isRequest bool, comment *models.Comment) {
	if isRequest {
		var opts = issueNotificationOpts{
			IssueID:              issue.ID,
			NotificationAuthorID: doer.ID,
			ReceiverID:           reviewer.ID,
			IsPull:               true,
		}

		if comment != nil {
			opts.CommentID = comment.ID
		}

		_ = ns.webpushQueue.Push(opts)
	}
}
