package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/gamedata"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/resource"
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
	resVersion, outOfDate, err := c.checkVersion(ctx)
	if err != nil {
		return err
	}
	if outOfDate || c.forceUpdate {
		if c.forceUpdate {
			logrus.Info("Though the local resources are up-to-date, a force update is required.")
		}

		log := logrus.WithFields(logrus.Fields{
			"resVersion": resVersion,
		})

		resLocation := filepath.Join(c.resLocation, resVersion)

		log.Info("Getting gamedata...")
		avgGameData, err := GetAvgGameData(resVersion)
		if err != nil {
			return err
		}

		log.Info("Getting resources...")
		err = GetAvgResources(resVersion, resLocation)
		if err != nil {
			return err
		}

		log.Info("Setting AVG gamedata...")
		err = c.avgService.SetAvgs(resVersion, avgGameData)
		if err != nil {
			return err
		}

		c.forceUpdate = false
		log.Info("Updated resources.")
	}
	return nil
}

func (c *Controller) checkVersion(ctx context.Context) (string, bool, error) {
	log := logrus.WithFields(logrus.Fields{})
	log.Info("Attempt to update resources.")

	rResVersion, err := getLatestResourceResVersion()
	if err != nil {
		return "", false, err
	}
	log.WithFields(logrus.Fields{
		"rResVersion": rResVersion,
	}).Info("Got resource resVersion.")

	gResVersion, gCommitRef, err := getLatestGamedataResVersion()
	if err != nil {
		return "", false, err
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
		return "", false, nil
	}

	remoteResVersion := rResVersion
	localResVersion, err := c.avgService.GetVersion(ctx)
	if err != nil {
		return "", false, err
	}

	log = log.WithFields(logrus.Fields{
		"remoteResVersion": remoteResVersion,
		"localResVersion":  localResVersion,
	})
	// Test whether the resource is up-to-date.
	if remoteResVersion == localResVersion {
		log.Info("The local resources are up-to-date.")
		return remoteResVersion, false, nil
	}
	log.Info("The local resources are out-of-date.")
	return remoteResVersion, true, nil
}

func getLatestResourceResVersion() (string, error) {
	return resource.GetResVersion()
}

func getLatestGamedataResVersion() (string, string, error) {
	return gamedata.FindLatestCommitRef()
}
