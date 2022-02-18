package service

import (
	"arkwaifu/internal/app/entity"
	"arkwaifu/internal/app/repo"
	"context"
)

type AvgService struct {
	resVersionRepo *repo.ResVersionRepo
	avgGroupRepo   *repo.AvgGroupRepo
	avgRepo        *repo.AvgRepo
}

func NewAvgService(rvRepo *repo.ResVersionRepo, agRepo *repo.AvgGroupRepo, aRepo *repo.AvgRepo) *AvgService {
	return &AvgService{
		resVersionRepo: rvRepo,
		avgGroupRepo:   agRepo,
		avgRepo:        aRepo,
	}
}

func (s *AvgService) GetResVersion(ctx context.Context) (string, error) {
	return s.resVersionRepo.GetResVersion(ctx)
}

//UpsertAvgs updates all
func (s *AvgService) UpsertAvgs(ctx context.Context, resVersion string, groups []entity.AvgGroup) error {
	return s.avgGroupRepo.Atomic(ctx, func(agr *repo.AvgGroupRepo) error {
		return s.avgRepo.Atomic(ctx, func(ar *repo.AvgRepo) error {
			return s.resVersionRepo.Atomic(ctx, func(rvr *repo.ResVersionRepo) error {
				for _, group := range groups {
					err := agr.UpsertAvgGroup(ctx, group)
					if err != nil {
						return err
					}
					for _, avg := range group.Avgs {
						err = ar.UpsertAvg(ctx, *avg)
						if err != nil {
							return err
						}
					}
				}
				return rvr.UpsertResVersion(ctx, resVersion)
			})
		})
	})
}
