#!/bin/bash

go test -cover -race emviwiki/auth/client
go test -cover -race emviwiki/auth/jwt
go test -cover -race emviwiki/auth/pages
go test -cover -race emviwiki/auth/user

go test -cover -race emviwiki/backend/article
go test -cover -race emviwiki/backend/article/history
go test -cover -race emviwiki/backend/article/schema
go test -cover -race emviwiki/backend/article/util
go test -cover -race emviwiki/backend/articlelist
go test -cover -race emviwiki/backend/billing
go test -cover -race emviwiki/backend/bookmark
go test -cover -race emviwiki/backend/client
go test -cover -race emviwiki/backend/content
go test -cover -race emviwiki/backend/context
go test -cover -race emviwiki/backend/feed
go test -cover -race emviwiki/backend/lang
go test -cover -race emviwiki/backend/member
go test -cover -race emviwiki/backend/newsletter
go test -cover -race emviwiki/backend/observe
go test -cover -race emviwiki/backend/organization
go test -cover -race emviwiki/backend/perm
go test -cover -race emviwiki/backend/pinned
go test -cover -race emviwiki/backend/prosemirror
go test -cover -race emviwiki/backend/search
go test -cover -race emviwiki/backend/support
go test -cover -race emviwiki/backend/tag
go test -cover -race emviwiki/backend/user
go test -cover -race emviwiki/backend/usergroup

go test -cover -race emviwiki/batch/balance
go test -cover -race emviwiki/batch/invitation
go test -cover -race emviwiki/batch/newsletter
go test -cover -race emviwiki/batch/notification
go test -cover -race emviwiki/batch/registration

go test -cover -race emviwiki/shared/auth
go test -cover -race emviwiki/shared/config
go test -cover -race emviwiki/shared/content
go test -cover -race emviwiki/shared/db
go test -cover -race emviwiki/shared/feed
go test -cover -race emviwiki/shared/i18n
go test -cover -race emviwiki/shared/mail
go test -cover -race emviwiki/shared/model
go test -cover -race emviwiki/shared/rest
go test -cover -race emviwiki/shared/tpl
go test -cover -race emviwiki/shared/util

go test -cover -race emviwiki/website/sitemap
