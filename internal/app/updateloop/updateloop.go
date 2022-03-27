package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/asset"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type Controller struct {
	ResLocation string
	ForceUpdate bool

	avgService   *avg.Service
	assetService *asset.Service
}

func NewController(avgService *avg.Service, assetService *asset.Service, conf *config.Config) *Controller {
	return &Controller{
		ResLocation:  conf.ResourceLocation,
		ForceUpdate:  conf.ForceUpdate,
		avgService:   avgService,
		assetService: assetService,
	}
}

func (c *Controller) AttemptUpdate(ctx context.Context) error {
	log.Info("attempt to update resources")

	if c.ForceUpdate {
		err := c.avgService.Reset(ctx)
		if err != nil {
			return err
		}
		c.ForceUpdate = false
		log.Info("force update flag was set, reset service to force update; force update flag reset")
	}

	localResVer, err := c.avgService.GetVersion(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	remoteResVer, err := arkres.GetResVersion()
	if err != nil {
		return errors.WithStack(err)
	}

	log := log.WithFields(log.Fields{
		"localResVer":  localResVer,
		"remoteResVer": remoteResVer,
	})
	log.Info("retrieved the local & remote resource versions")

	// only the first true case would be executed
	if localResVer != remoteResVer {
		log.Info("updating, since local & remote resource versions are not the same")
		begin := time.Now()
		err = c.doUpdate(ctx, localResVer, remoteResVer)
		elapsed := time.Since(begin)
		log.WithField("elapsed", elapsed).Info("updated from localResVer to remoteResVer")
		if err != nil {
			return err
		}
	}
	return nil
}

// doUpdate updates the resources.
func (c *Controller) doUpdate(ctx context.Context, oldResVer string, newResVer string) error {
	var err error
	err = c.retrieveResources(ctx, oldResVer, newResVer)
	if err != nil {
		return err
	}
	err = c.processStatics(ctx, newResVer)
	if err != nil {
		return err
	}
	err = c.submitUpdate(ctx, newResVer)
	if err != nil {
		return err
	}
	return nil
}

// retrieveResources retrieves the resources into the resource directory.
// If oldResVer is empty, it retrieves full resources; otherwise, it retrieves only incremental resources.
//
// It is skipped if the corresponding resource directory already exists.
func (c *Controller) retrieveResources(ctx context.Context, oldResVer string, newResVer string) error {
	resDir := filepath.Join(c.ResLocation, newResVer, "res")
	if fileutil.Exists(resDir) {
		log.WithFields(log.Fields{"resDir": resDir}).
			Info("retrieving resources: resDir already exists; skipping")
		return nil
	}
	err := updateResources(ctx, oldResVer, newResVer, resDir)
	if err != nil {
		return errors.WithStack(err)
	}
	log.WithFields(log.Fields{
		"oldResVer": oldResVer,
		"newResVer": newResVer,
		"resDir":    resDir,
	}).Info("updated resources from oldResVer to newResVer, into resDir")
	return nil
}

// processStatics processes the resources in the resource directory to static files.
//
// It is skipped if the corresponding resource directory already exists.
func (c *Controller) processStatics(ctx context.Context, resVer string) error {
	resDir := filepath.Join(c.ResLocation, resVer, "res")
	staticDir := filepath.Join(c.ResLocation, resVer, "static")
	if fileutil.Exists(staticDir) {
		log.WithFields(log.Fields{"resDir": resDir}).
			Info("processing statics: staticDir already exists; skipping")
		return nil
	}
	err := os.MkdirAll(staticDir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	err = updateStatics(ctx, resDir, staticDir)
	if err != nil {
		return errors.WithStack(err)
	}
	log.WithFields(log.Fields{
		"resDir":    resDir,
		"staticDir": staticDir,
	}).Info("processed statics from resDir to staticDir")
	return nil
}

// submitUpdate submits the gamedata from resources and static files.
//
// Note that submitting the gamedata from resources will fully override existing;
// but submitting the static files will only override incrementally.
func (c *Controller) submitUpdate(ctx context.Context, resVer string) error {
	resDir := filepath.Join(c.ResLocation, resVer, "res")
	staticDir := filepath.Join(c.ResLocation, resVer, "static")
	mainStaticDir := filepath.Join(c.ResLocation, "static")

	err := fileutil.LowercaseAll(staticDir)
	if err != nil {
		return err
	}

	err = fileutil.CopyAllFileContent(staticDir, mainStaticDir)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.updateDatabase(ctx, resVer, resDir)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.updateAssetDatabase(ctx, mainStaticDir)
	if err != nil {
		return errors.WithStack(err)
	}

	log.WithFields(log.Fields{
		"resVer": resVer,
		"resDir": resDir,
	}).Info("submitted update from resDir, version resVer")
	return nil
}
