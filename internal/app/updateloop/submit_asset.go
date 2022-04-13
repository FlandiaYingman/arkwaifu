package updateloop

import (
	"context"
)

// submitAsset submits the asset data of the given resource version to the AVG service.
//
// Note that submitting asset data is a fully overwrite operation.
func (s *Service) submitAsset(ctx context.Context, resVer ResVersion) error {
	err := s.AssetService.InitNames(ctx)
	if err != nil {
		return err
	}
	err = s.AssetService.PopulateFrom(ctx, s.StaticDir(resVer))
	return err
}
