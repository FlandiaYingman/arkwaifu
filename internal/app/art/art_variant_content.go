package art

import (
	"bytes"
	"fmt"
	"image"
	"path"
)

// VariantContent helps store and take the content of a variant to or from the
// filesystem.
type VariantContent struct {
	ArtID     string `params:"id"`
	Variation string `params:"variation"`
	Content   []byte
}

func (s *VariantContent) String() string {
	if s.Variation == VariationOrigin {
		return fmt.Sprintf("%s", s.ArtID)
	} else {
		return fmt.Sprintf("%s.%s", s.ArtID, s.Variation)
	}
}

func (s *VariantContent) PathRel() string {
	return fmt.Sprintf("%v.webp", s)
}

func (s *VariantContent) Path(root string) string {
	return path.Join(root, s.PathRel())
}

func (s *VariantContent) Check() (*image.Config, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(s.Content))
	if err != nil {
		return nil, fmt.Errorf("cannot decode static file %v: %w", s, err)
	}
	if format != "webp" {
		return nil, fmt.Errorf("the format of static file %v is not webp", s)
	}
	return &config, nil
}
