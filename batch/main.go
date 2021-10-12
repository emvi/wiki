package main

import (
	auth "emviwiki/auth/model"
	"emviwiki/batch/balance"
	"emviwiki/batch/invitation"
	"emviwiki/batch/newsletter"
	"emviwiki/batch/notification"
	"emviwiki/batch/registration"
	dashboard "emviwiki/dashboard/model"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	"emviwiki/shared/feed"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	backend "emviwiki/shared/model"
	"emviwiki/shared/server"
	"github.com/emvi/logbuch"
	"log"
	"strings"
)

var (
	// list of all batch execution structs
	batches = map[string]batch{
		"notifications":         {notification.LoadConfig, notification.SendNotificationMails},
		"newsletter":            {newsletter.LoadConfig, newsletter.SendNewsletterMails},
		"cleanup_registrations": {nil, registration.CleanupRegistrations},
		"cleanup_invitations":   {nil, invitation.CleanupInvitations},
		"update_balance":        {balance.LoadConfig, balance.UpdateBalance},
	}
)

// entry point for all batch processes
type batch struct {
	init func()
	exec func()
}

func selectAndStartBatch() {
	batchName := strings.TrimSpace(strings.ToLower(config.Get().Batch.Process))
	var name string
	var exec *batch

	for key, value := range batches {
		if key == batchName {
			name = key
			exec = &value
			break
		}
	}

	if exec == nil {
		log.Printf("No batch selected or name unknown '%s'. Select one of the following batches by setting (EMVI_WIKI_BATCH_PROCESS):\n", batchName)

		for batch := range batches {
			log.Println(batch)
		}
	} else {
		logbuch.Info("Starting batch", logbuch.Fields{"name": name})
		backendDB, authDB, dashboardDB := connectDB()
		backend.SetConnection(backendDB)
		auth.SetConnection(authDB)
		dashboard.SetConnection(backendDB, authDB, dashboardDB)
		defer backendDB.Disconnect()
		defer authDB.Disconnect()
		defer dashboardDB.Disconnect()

		if exec.init != nil {
			exec.init()
		}

		exec.exec()
		logbuch.Info("Finished batch execution", logbuch.Fields{"name": name})
	}
}

func connectDB() (*db.Connection, *db.Connection, *db.Connection) {
	authdb := config.Get().AuthDB
	backenddb := config.Get().BackendDB
	dashboarddb := config.Get().DashboardDB
	return db.NewConnection(db.ConnectionData{
			Host:               backenddb.Host,
			Port:               backenddb.Port,
			User:               backenddb.User,
			Password:           backenddb.Password,
			Schema:             backenddb.Schema,
			SSLMode:            backenddb.SSLMode,
			SSLCert:            backenddb.SSLCert,
			SSLKey:             backenddb.SSLKey,
			SSLRootCert:        backenddb.SSLRootCert,
			MaxOpenConnections: backenddb.MaxOpenConnections,
		}), db.NewConnection(db.ConnectionData{
			Host:               authdb.Host,
			Port:               authdb.Port,
			User:               authdb.User,
			Password:           authdb.Password,
			Schema:             authdb.Schema,
			SSLMode:            authdb.SSLMode,
			SSLCert:            authdb.SSLCert,
			SSLKey:             authdb.SSLKey,
			SSLRootCert:        authdb.SSLRootCert,
			MaxOpenConnections: authdb.MaxOpenConnections,
		}), db.NewConnection(db.ConnectionData{
			Host:               dashboarddb.Host,
			Port:               dashboarddb.Port,
			User:               dashboarddb.User,
			Password:           dashboarddb.Password,
			Schema:             dashboarddb.Schema,
			SSLMode:            dashboarddb.SSLMode,
			SSLCert:            dashboarddb.SSLCert,
			SSLKey:             dashboarddb.SSLKey,
			SSLRootCert:        dashboarddb.SSLRootCert,
			MaxOpenConnections: dashboarddb.MaxOpenConnections,
		})
}

func main() {
	config.Load()
	stdout, stderr := server.ConfigureLogging()
	defer server.CloseLogger(stdout, stderr)
	feed.LoadConfig()
	mail.LoadConfig()
	i18n.LoadConfig()
	selectAndStartBatch()
}
