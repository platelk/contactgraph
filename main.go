package main

import (
	"github.com/platelk/contactgraph/domain/usecases"
	"github.com/platelk/contactgraph/infra/contactstore"
	"github.com/platelk/contactgraph/infra/logger"
	"github.com/platelk/contactgraph/infra/userstore"
	"github.com/platelk/contactgraph/transport/http"
)

func main() {
	// Load configuration from different sources.
	cfg := Load()

	// Initialise log.
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't initialise logger.")
	}
	log.Debug().Interface("config", cfg).Send()

	// Initialise user store
	usrStore := userstore.WithStats(userstore.NewInMemoryImproved())

	// Initialise contact store
	contactStore := contactstore.WithStats(contactstore.NewShardStore(
		contactstore.NewInMemoryMap(),
		contactstore.NewInMemoryMap(),
		contactstore.NewInMemoryMap(),
		contactstore.NewInMemoryMap(),
		contactstore.NewInMemoryMap(),
	))

	// Build http server
	srv := http.NewBuilder(log, cfg.HTTP).
		// /users
		WithV1CreateUser(usecases.SetupCreateUser(log, usrStore)).
		WithV1UpdateUser(usecases.SetupUpdateUser(log, usrStore)).
		WithV1DeleteUser(usecases.SetupDeleteUser(log, usrStore)).
		WithV1SearchUser(usecases.SetupSearchUser(log, usrStore)).
		// /contacts
		WithV1ConnectContact(usecases.SetupConnectContact(log, contactStore)).
		WithV1LookupContact(usecases.SetupLookupContact(log, contactStore, usrStore)).
		WithV1ReverseLookupContact(usecases.SetupReverseLookupContact(log, contactStore, usrStore)).
		WithV1SuggestionContact(usecases.SetupSuggestContact(log, 10, contactStore, usrStore)).
		// /dev
		WithV1DevStats(contactStore, usrStore).
		WithV1DevGenerateData(usrStore, contactStore).
		WithHealthCheck().
		Build()

	// Run HTTP server.
	log.Info().Msg("running http server")
	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("http server didn't end correctly")
	}

	// Finished.
	log.Info().Msg("done")
}
