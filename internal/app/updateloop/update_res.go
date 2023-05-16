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
	// avgResourceRegexp are the regular expressions of resources that need to always be fully retrieved.
	//
	// Note that the regexps start with 'gamedata', however, the naming of this variable is correct (starts with 'avg').
	avgResourceRegexp = []*regexp.Regexp{
		regexp.MustCompile("^gamedata/excel"),
		regexp.MustCompile("^gamedata/story"),
	}
	// assetResourceRegexp are the regular expressions of resources that is ok to be retrieved incrementally.
	assetResourceRegexp = []*regexp.Regexp{
		regexp.MustCompile("^avg/bg"),
		regexp.MustCompile("^avg/imgs"),
		regexp.MustCompile("^avg/characters"),
	}
	// allResRE are the regular expressions which consist of avgResourceRegexp and assetResourceRegexp.
	allResRE []*regexp.Regexp
)

func init() {
	allResRE = append(allResRE, avgResourceRegexp...)
	allResRE = append(allResRE, assetResourceRegexp...)
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
func updateResources(ctx context.Context, oldResVer string, newResVer string, newResDir string) (err error) {
	if oldResVer == "" {
		err = arkres.GetFromHGAPI(ctx, newResVer, newResDir, assetResourceRegexp...)
		if err != nil {
			return err
		}
		err = arkres.GetFromAGDAPI(ctx, newResVer, newResDir, avgResourceRegexp...)
		if err != nil {
			return err
		}
		return nil
	} else {
		err = arkres.GetFromHGAPIIncrementally(ctx, oldResVer, newResVer, newResDir, assetResourceRegexp...)
		if err != nil {
			return err
		}
		err = arkres.GetFromAGDAPI(ctx, newResVer, newResDir, avgResourceRegexp...)
		if err != nil {
			return err
		}
		return nil
	}
}
