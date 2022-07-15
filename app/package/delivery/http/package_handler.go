package http

import (
	"context"
	"net/http"
	"track-my-package/app/domain"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type PackageHandler struct {
	upgrader websocket.Upgrader
	PUseCase domain.PackageUseCase
}

func NewPackageHandler(e *echo.Echo, pu domain.PackageUseCase) {
	handler := PackageHandler{
		upgrader: websocket.Upgrader{},
		PUseCase: pu,
	}
	e.GET("/packages/track/:vehicleID", handler.TrackByVehicleID)
	e.POST("/packages/location/:vehicleID", handler.UpdateLocation)
}

func (p *PackageHandler) TrackByVehicleID(c echo.Context) error {
	ws, err := p.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, _, err = ws.ReadMessage()
		if err != nil {
			cancel()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			ws.Close()
			return nil
		default:
			p, err := p.PUseCase.TrackByVehicleID(ctx, c.Param("vehicleID"))
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			err = ws.WriteJSON(p)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

func (p *PackageHandler) UpdateLocation(c echo.Context) error {
	pac := &domain.Package{}
	if err := c.Bind(pac); err != nil {
		return err
	}
	p.PUseCase.UpdateLocation(pac)
	return c.JSON(http.StatusCreated, pac)
}
