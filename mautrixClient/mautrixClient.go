package mautrixClient

import (
	"context"
	"errors"
	"fmt"
	"grafana-matrix-alerts/config"
	"grafana-matrix-alerts/logger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"sync"
)

var (
	Client *mautrix.Client
)

func Connect() {

	var err error
	Client, err = mautrix.NewClient(config.Homeserver, "", "")
	if err != nil {
		panic(err)
	}
	Client.Log = logger.Log

	syncer := Client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.StateMember, func(ctx context.Context, evt *event.Event) {
		if evt.GetStateKey() == Client.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
			_, err := Client.JoinRoomByID(ctx, evt.RoomID)
			if err == nil {
				logger.Log.Info().
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Joined room after invite")

			} else {
				logger.Log.Error().Err(err).
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Failed to join room after invite")
			}

			infoMessage := fmt.Sprintf(
				"👋 Hey!<br>The room id is: <code>%s</code>",
				evt.RoomID,
			)

			content := format.RenderMarkdown(infoMessage, true, true)
			_, err = Client.SendMessageEvent(context.TODO(), evt.RoomID, event.EventMessage, &content)
		}
	})

	cryptoHelper, err := cryptohelper.NewCryptoHelper(Client, []byte("meow"), config.Database)
	if err != nil {
		panic(err)
	}

	cryptoHelper.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: config.Username},
		Password:   config.Password,
	}

	err = cryptoHelper.Init(context.TODO())
	if err != nil {
		panic(err)
	}
	// Set the client crypto helper in order to automatically encrypt outgoing messages
	Client.Crypto = cryptoHelper

	logger.Log.Info().Msg("Matrix client connected")
	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)

	// Verify own device using recoveryKey
	_, isVerified, err := GetVerificationStatus(cryptoHelper, Client, context.TODO())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get verification status")
	}

	if !isVerified && config.RecoveryKey != "" {
		err = VerifyWithRecoveryKey(cryptoHelper, context.TODO(), config.RecoveryKey)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to verify device signature")
		}
	}

	go func() {
		err = Client.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	// Ensure cleanup runs on exit
	defer func() {
		cancelSync()
		syncStopWait.Wait()

		err := cryptoHelper.Close()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error closing database")
		}
	}()

	select {}
}
