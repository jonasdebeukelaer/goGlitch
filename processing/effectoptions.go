package processing

import (
	"image"

	"github.com/jonasdebeukelaer/goGlitch/processing/effects"
)

// Effect represents a possible effect to apply to an image
type Effect func(image.Image) (image.Image, error)

type effectOption struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	effect Effect
}

var (
	// EffectList contains all the available effects to the app.
	// it also populates the frontend list
	EffectList = []effectOption{
		effectOption{
			ID:     "lignify",
			Name:   "Lignify",
			effect: effects.Lignify,
		},
		effectOption{
			ID:     "colourise",
			Name:   "Colourise",
			effect: effects.RandomMuddle,
		},
	}

	// EffectMap maps from an effect id to its function, for
	// dynamic function calls from the handler
	EffectMap = func() map[string]Effect {
		var eMap map[string]Effect
		for i := range EffectList {
			if len(eMap) == 0 {
				eMap = map[string]Effect{EffectList[i].ID: EffectList[i].effect}
			} else {
				eMap[EffectList[i].ID] = EffectList[i].effect
			}
		}
		return eMap
	}()
)
