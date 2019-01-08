package processing

import (
	"image"

	"github.com/jonasdebeukelaer/goGlitch/processing/effects"
)

// Effect represents a possible effect to apply to an image
type Effect func(image.Image) (image.Image, error)

var (
	EffectLignify      Effect = effects.Lignify      // use lignify effect
	EffectRandomMuddle Effect = effects.RandomMuddle // use random muddle effect
)
