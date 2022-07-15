package domain

import "context"

type Package struct {
	From      string `json:"from"`
	To        string `json:"to"`
	VehicleID string `json:"vehicle_id"`
}

type PackageUseCase interface {
	TrackByVehicleID(ctx context.Context, vehicleID string) (*Package, error)
	UpdateLocation(p *Package) error
}

type PackageConsumer interface {
	ConsumeByVehicleID(ctx context.Context, vehicleID string) ([]byte, error)
	Publish(p *Package) error
}
