package usecase

import (
	"context"
	"encoding/json"
	"track-my-package/app/domain"
)

type packageUseCase struct {
	pc domain.PackageConsumer
}

func NewPackageUseCase(pc domain.PackageConsumer) domain.PackageUseCase {
	return &packageUseCase{
		pc: pc,
	}
}
func (p *packageUseCase) TrackByVehicleID(ctx context.Context, vehicleID string) (*domain.Package, error) {
	bytes, err := p.pc.ConsumeByVehicleID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	var res domain.Package
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *packageUseCase) UpdateLocation(dp *domain.Package) error {
	return p.pc.Publish(dp)
}
