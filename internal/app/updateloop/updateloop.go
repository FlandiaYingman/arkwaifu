package updateloop

import (
	"context"
	"path/filepath"
	"time"

	"github.com/flandiayingman/arkwaifu/internal/app/asset"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("updateloop",
		fx.Provide(
			ProvideConfig,
			NewService,
		),
	)
}

type Service struct {
	MainResDir      string
	MainStaticDir   string
	ForceUpdate     bool
	ForceSubmit     bool
	ForceResVersion string

	AvgService   *avg.Service
	AssetService *asset.Service
}
type ResVersion string

func NewService(
	avgS *avg.Service,
	assetS *asset.Service,
	appConf *config.Config,
	conf *Config,
) *Service {
	return &Service{
		MainResDir:      appConf.ResourceDir,
		MainStaticDir:   appConf.StaticDir,
		ForceUpdate:     conf.ForceUpdate,
		ForceSubmit:     conf.ForceSubmit,
		ForceResVersion: conf.ForceResVersion,
		AvgService:      avgS,
		AssetService:    assetS,
	}
}

func (s *Service) ResVer(ctx context.Context) (ResVersion, ResVersion, error) {
	local, err := s.AvgService.GetVersion(ctx)
	if err != nil {
		return "", "", errors.Wrapf(err, "cannot get local resource version")
	}
	remote, err := arkres.GetResVersion()
	if err != nil {
		return "", "", errors.Wrapf(err, "cannot get remote resource version")
	}
	if s.ForceResVersion != "" {
		log.WithFields(log.Fields{"local": local, "remote": remote, "forceResVersion": s.ForceResVersion}).
			Warn("force res version is set, ignoring res remote version and use the force res version")
		remote = s.ForceResVersion
		defer func() { s.ForceResVersion = "" }()
	}
	return ResVersion(local), ResVersion(remote), nil
}

func (s *Service) VersioningDir(r ResVersion) string {
	return filepath.Join(s.MainResDir, string(r))
}
func (s *Service) ResDir(r ResVersion) string {
	return filepath.Join(s.VersioningDir(r), "res")
}
func (s *Service) StaticDir(r ResVersion) string {
	return filepath.Join(s.VersioningDir(r), "static")
}

// AttemptUpdate attempts to update the resources.
func (s *Service) AttemptUpdate(ctx context.Context) error {
	log.Info("attempt to update the resources")

	current, latest, err := s.ResVer(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve the current and latest version")
	}
	log.WithFields(log.Fields{
		"current": current,
		"latest":  latest,
	}).Debug("retrieved the current and latest versions")

	err = s.update(ctx, current, latest)
	if err != nil {
		return err
	}
	err = s.submit(ctx, current, latest)
	if err != nil {
		return err
	}

	return nil
}

// update updates the resources and static files from the old ResVersion to the new ResVersion.
func (s *Service) update(ctx context.Context, old ResVersion, new ResVersion) error {
	if old == new && !s.ForceUpdate {
		log.WithFields(log.Fields{
			"old":         old,
			"new":         new,
			"ForceUpdate": s.ForceUpdate,
		}).Info("skipping update phase, since the resource versions are the same, and ForceUpdate is not set")
		return nil
	}
	if s.ForceUpdate {
		log.Warn("forcing update since ForceUpdate is set; setting old resVersion to empty to force update")
		old = ""
		defer func() { s.ForceUpdate = false }()
	}
	log.WithFields(log.Fields{
		"old": old,
		"new": new,
	}).Info("updating the resources and static files from old to new")

	begin := time.Now()

	var err error
	err = s.updateRes(ctx, old, new)
	if err != nil {
		return errors.Wrapf(err, "failed to update resources from %v to %v", old, new)
	}
	err = s.updateStatic(ctx, new)
	if err != nil {
		return errors.Wrapf(err, "failed to update static files from %v to %v", old, new)
	}

	elapsed := time.Since(begin)
	log.WithFields(log.Fields{
		"old":     old,
		"new":     new,
		"elapsed": elapsed,
	}).Info("updated the resources and static files from old to new")
	return nil
}

// update submits the resources and static files of the new ResVersion to the services.
func (s *Service) submit(ctx context.Context, old ResVersion, new ResVersion) error {
	if old == new && !s.ForceSubmit {
		log.WithFields(log.Fields{
			"old":         old,
			"new":         new,
			"ForceSubmit": s.ForceSubmit,
		}).Info("skipping submit phase, since the resource versions are the same, and ForceSubmit is not set")
		return nil
	}
	if s.ForceSubmit {
		log.Warn("forcing submit, since ForceSubmit is set")
		defer func() { s.ForceSubmit = false }()
	}
	log.WithFields(log.Fields{
		"new": new,
	}).Info("submitting the resources and static files of new to services")

	begin := time.Now()

	var err error
	err = s.submitAvg(ctx, new)
	if err != nil {
		return errors.Wrapf(err, "failed to submit the resources and static files of %v to AVG service", new)
	}
	err = s.submitAsset(ctx, new)
	if err != nil {
		return errors.Wrapf(err, "failed to submit the resources and static files of %v to Asset service", new)
	}

	elapsed := time.Since(begin)
	log.WithFields(log.Fields{
		"new":     new,
		"elapsed": elapsed,
	}).Info("submitted the resources and static files of new to the services")
	return nil
}
