package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"

	"github.com/blang/semver"
	"github.com/brimstone/whoisx"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func PublicKey() (*ecdsa.PublicKey, error) {
		block, _ := pem.Decode([]byte(`
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEV+xkT70NtOJJRIOfnMw/KAFBIzCw
EPrC6Nz3G8+OEjZUTNMiB/Q598rCaGPtaKYnU6WE4tG6GrH7tVBcubUdfA==
-----END PUBLIC KEY-----`))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse public key: " + err.Error())
	}

	pubKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("PublicKey is not ECDSA")
	}

	return pubKey, nil
}

func Update() *cobra.Command {
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update whoisx",
		Long: `Check github.com/brimstone/whoisx for the latest version of whoisx.
Download it if available and replace the current executable.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := semver.Parse(whoisx.Version)
			if err != nil {
				slog.Warn("Unable to parse version", "version", whoisx.Version)
			}
			// Check for updates
			pubkey, err := PublicKey()
			if err != nil {
				return fmt.Errorf("error occurred while extracting public key: %w", err)
			}

			slog.Info("Checking for updates")

			up, err := selfupdate.NewUpdater(selfupdate.Config{
				Validator: &selfupdate.ECDSAValidator{
					PublicKey: pubkey,
				},
				Filters: []string{
					whoisx.Binary,
				},
			})
			if err != nil {
				return fmt.Errorf("unable to create updater component: %w", err)
			}

			latest, err := up.UpdateSelf(v, "brimstone/whoisx")

			if err != nil {
				return fmt.Errorf("binary update failed: %w", err)
			}

			if latest.Version.Equals(v) {
				// latest version is the same as current version. It means current binary is up to date.
				slog.Info("Current binary is the latest", "version", v)
			} else {
				slog.Info("Successfully updated", "version", latest.Version)
			}

			return nil
		},
	}

	return updateCmd
}
