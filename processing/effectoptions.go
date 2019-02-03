package processing

import (
	"image"

	"github.com/jonasdebeukelaer/goGlitch/processing/effects"
)

// Effect represents a possible effect to apply to an image
type Effect func(image.Image, map[string]string) (image.Image, error)

type effectOption struct {
	Key    string `json:"effect_key"`
	Name   string `json:"name"`
	effect Effect
}

var (
	// EffectList contains all the available effects to the app.
	// it also populates the frontend list
	EffectList = []effectOption{
		effectOption{
			Key:    "li",
			Name:   "Lignify",
			effect: effects.Lignify,
		},
		effectOption{
			Key:    "co",
			Name:   "Colourise",
			effect: effects.RandomMuddle,
		},
		effectOption{
			Key:    "sp",
			Name:   "Split colours",
			effect: effects.SplitColours,
		},
		effectOption{
			Key:    "td",
			Name:   "Threedee",
			effect: effects.Threedee,
		},
	}

	// EffectMap maps from an effect id to its function, for
	// dynamic function calls from the handler
	EffectMap = func() map[string]Effect {
		var eMap map[string]Effect
		for i := range EffectList {
			if len(eMap) == 0 {
				eMap = map[string]Effect{EffectList[i].Key: EffectList[i].effect}
			} else {
				eMap[EffectList[i].Key] = EffectList[i].effect
			}
		}
		return eMap
	}()
)
