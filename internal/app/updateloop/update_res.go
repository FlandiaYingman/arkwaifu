package updateloop

import (
	"context"
	_ "image/jpeg" // register jpeg codec
	_ "image/png"  // register png codec
	"regexp"

	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// fullResRE are the regular expressions of resources that need to always be fully retrieved.
	fullResRE = []*regexp.Regexp{
		regexp.MustCompile("^gamedata/excel"),
		regexp.MustCompile("^gamedata/story"),
	}
	// incrementalResRE are the regular expressions of resources that is ok to be retrieved incrementally.
	incrementalResRE = []*regexp.Regexp{
		regexp.MustCompile("^avg/bg"),
		regexp.MustCompile("^avg/imgs"),
		regexp.MustCompile("^avg/characters"),
	}
	// allResRE are the regular expressions which consist of fullResRE and incrementalResRE.
	allResRE []*regexp.Regexp
)

func init() {
	allResRE = append(allResRE, fullResRE...)
	allResRE = append(allResRE, incrementalResRE...)
}

// updateRes gets the resources of newResVer into the corresponding resource directory.
// If oldResVer is empty, it will get full resources; otherwise, it will get only incremental resources.
//
// It is skipped if the corresponding resource directory of newResVer already exists.
//
// It doesn't ensure that the incrementally updated resources are identical to fully updated resources.
// However, it ensures the incrementally updated resources are a superset of fully updated resources.
func (s *Service) updateRes(ctx context.Context, oldResVer ResVersion, newResVer ResVersion) error {
	newResDir := s.ResDir(newResVer)

	exists, err := fileutil.Exists(newResDir)
	if err != nil {
		return errors.Wrapf(err, "cannot check if resDir exists: %s", newResDir)
	}
	if exists && !s.ForceUpdate {
		log.WithFields(log.Fields{
			"newResDir": newResDir,
		}).Info("update res: newResDir already exists; skipping")
		return nil
	}

	err = updateResources(ctx, string(oldResVer), string(newResVer), newResDir)
	if err != nil {
		return errors.Wrapf(err, "cannot update resources: %s", newResDir)
	}
	log.WithFields(log.Fields{
		"oldResVer": oldResVer,
		"newResVer": newResVer,
		"newResDir": newResDir,
	}).Info("updated resources from oldResVer to newResVer, into newResDir")
	return nil
}

// updateResources simply wraps around the package arkres.
func updateResources(ctx context.Context, oldResVer string, newResVer string, newResDir string) error {
	if oldResVer == "" {
		return arkres.Get(ctx, newResVer, newResDir, allResRE...)
	} else {
		err := arkres.Get(ctx, newResVer, newResDir, fullResRE...)
		if err != nil {
			return err
		}
		err = arkres.GetIncrementally(ctx, oldResVer, newResVer, newResDir, incrementalResRE...)
		if err != nil {
			return err
		}
		return nil
	}
}
