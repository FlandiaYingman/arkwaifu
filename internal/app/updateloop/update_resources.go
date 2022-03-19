package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	_ "image/jpeg" // register jpeg codec
	_ "image/png"  // register png codec
	"regexp"
)

var (
	incrementalResRE = []*regexp.Regexp{
		regexp.MustCompile("^avg/bg"),
		regexp.MustCompile("^avg/imgs"),
		regexp.MustCompile("^gamedata/story"),
	}
	fullResRE = []*regexp.Regexp{
		regexp.MustCompile("^gamedata/excel"),
	}
)

// updateResources retrieves the incremental resources from oldResVer to newResVer.
//
// It doesn't ensure that the updated resources is identical to freshly retrieved resources.
// i.e., the updated resources will always contain the freshly retrieved resources, but not vice-versa.
func updateResources(ctx context.Context, oldResVer string, newResVer string, newResDir string) error {
	if oldResVer == "" {
		var re []*regexp.Regexp
		re = append(re, fullResRE...)
		re = append(re, incrementalResRE...)
		return arkres.Get(ctx, newResVer, newResDir, re...)
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
