package art

import (
	"bytes"
	"fmt"
	"image"
	"path"
)

// Content helps store and take the content of a variant to or from the
// filesystem.
type Content struct {
	ID        string
	Category  string
	Variation string
	Content   []byte
}

func (s *Content) String() string {
	if s.Variation == VariationOrigin {
		return fmt.Sprintf("%s", s.ID)
	} else {
		return fmt.Sprintf("%s/%s", s.Variation, s.ID)
	}
}

func (s *Content) PathRel() string {
	return path.Join(s.Category, s.String()+".webp")
}

func (s *Content) Check() (*image.Config, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(s.Content))
	if err != nil {
		return nil, fmt.Errorf("cannot decode content file %v: %w", s, err)
	}
	if format != "webp" {
		return nil, fmt.Errorf("the format of content file %v is not webp", s)
	}
	return &config, nil
}
