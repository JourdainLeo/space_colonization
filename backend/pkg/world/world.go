package world

import (
	ss "ai30-spatial/pkg/solar_system"
)

type World struct {
	SolarSystem ss.SolarSystem
}

func NewWorld() *World {
	return &World{*ss.NewSolarSystem()}
}
