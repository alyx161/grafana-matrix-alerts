package mautrixClient

import (
	"context"
	"fmt"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
)

func VerifyWithRecoveryKey(cryptoHelper *cryptohelper.CryptoHelper, ctx context.Context, recoveryKey string) error {
	keyID, keyData, err := cryptoHelper.Machine().SSSS.GetDefaultKeyData(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default SSSS key data: %w", err)
	}
	key, err := keyData.VerifyRecoveryKey(keyID, recoveryKey)
	if err != nil {
		return err
	}
	err = cryptoHelper.Machine().FetchCrossSigningKeysFromSSSS(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to fetch cross-signing keys from SSSS: %w", err)
	}
	err = cryptoHelper.Machine().SignOwnDevice(ctx, cryptoHelper.Machine().OwnIdentity())
	if err != nil {
		return fmt.Errorf("failed to sign own device: %w", err)
	}
	err = cryptoHelper.Machine().SignOwnMasterKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to sign own master key: %w", err)
	}
	return nil
}

func GetVerificationStatus(cryptoHelper *cryptohelper.CryptoHelper, client *mautrix.Client, ctx context.Context) (hasKeys, isVerified bool, err error) {
	pubkeys := cryptoHelper.Machine().GetOwnCrossSigningPublicKeys(ctx)
	if pubkeys != nil {
		hasKeys = true
		isVerified, err = cryptoHelper.Machine().CryptoStore.IsKeySignedBy(
			ctx, client.UserID, cryptoHelper.Machine().GetAccount().SigningKey(), client.UserID, pubkeys.SelfSigningKey,
		)
		if err != nil {
			err = fmt.Errorf("failed to check if current device is signed by own self-signing key: %w", err)
		}
	}
	return
}
