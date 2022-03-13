package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/asset"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/data"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type Controller struct {
	resLocation string
	forceUpdate bool
	avgService  *avg.Service
}

func NewController(avgService *avg.Service, conf *config.Config) *Controller {
	return &Controller{conf.ResourceLocation, conf.ForceUpdate, avgService}
}

func (c *Controller) UpdateResources() error {
	ctx := context.Background()
	localResVer, remoteResVer, err := c.checkVersion(ctx)
	if err != nil {
		return err
	}
	if c.forceUpdate {
		localResVer = ""
		logrus.Info("Because a force updated is specified, set the local resVersion to empty.")
	}

	if localResVer != remoteResVer {
		log := logrus.WithFields(logrus.Fields{
			"localResVer":  localResVer,
			"remoteResVer": remoteResVer,
		})

		resLocation := filepath.Join(c.resLocation, remoteResVer)

		log.Info("Getting gamedata...")
		avgGameData, err := GetAvgGameData(remoteResVer)
		if err != nil {
			return err
		}

		log.Info("Getting resources...")
		err = GetAvgResources(ctx, localResVer, remoteResVer, resLocation)
		if err != nil {
			return err
		}

		log.Info("Setting AVG gamedata...")
		err = c.avgService.SetAvgs(remoteResVer, avgGameData)
		if err != nil {
			return err
		}

		c.forceUpdate = false
		log.Info("Updated resources.")
	}
	return nil
}

func (c *Controller) checkVersion(ctx context.Context) (string, string, error) {
	log := logrus.WithFields(logrus.Fields{})
	log.Info("Attempt to update resources.")

	rResVersion, err := getLatestResourceResVersion()
	if err != nil {
		return "", "", err
	}
	log.WithFields(logrus.Fields{
		"rResVersion": rResVersion,
	}).Info("Got resource resVersion.")

	gResVersion, gCommitRef, err := getLatestGamedataResVersion()
	if err != nil {
		return "", "", err
	}
	log.WithFields(logrus.Fields{
		"gResVersion": gResVersion,
		"gCommitRef":  gCommitRef,
	}).Info("Got gamedata resVersion and commitRef.")

	if rResVersion != gResVersion {
		log.WithFields(logrus.Fields{
			"rResVersion": rResVersion,
			"gResVersion": gResVersion,
		}).Info("The remote resources are updating.")
		return "", "", nil
	}

	remoteResVersion := rResVersion
	localResVersion, err := c.avgService.GetVersion(ctx)
	if err != nil {
		return "", "", err
	}

	log = log.WithFields(logrus.Fields{
		"remoteResVersion": remoteResVersion,
		"localResVersion":  localResVersion,
	})
	// Test whether the resource is up-to-date.
	if remoteResVersion == localResVersion {
		log.Info("The local resources are up-to-date.")
	} else {
		log.Info("The local resources are out-of-date.")
	}
	return localResVersion, remoteResVersion, nil
}

func getLatestResourceResVersion() (string, error) {
	return asset.GetLatestResVersion()
}

func getLatestGamedataResVersion() (string, string, error) {
	return data.FindLatestCommitRef()
}
