package res

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Variant struct {
	Key      string
	FileName func(base string) string
	Location func(base string, version string) string
}

var (
	Raw Variant = Variant{
		Key: "raw",
		FileName: func(base string) string {
			return fmt.Sprintf("%s.png", base)
		},
		Location: func(base string, version string) string {
			return filepath.Join(base, version, "raw")
		},
	}
	Thumbnail Variant = Variant{
		Key: "raw",
		FileName: func(base string) string {
			return fmt.Sprintf("%s.webp", base)
		},
		Location: func(base string, version string) string {
			return filepath.Join(base, version, "thumbnail")
		},
	}
)

func VariantFromString(str string) (Variant, error) {
	str = strings.ToLower(str)
	switch str {
	case "raw":
		return Raw, nil
	case "thumbnail":
		return Thumbnail, nil
	default:
		return Variant{}, fmt.Errorf("Variant %v not supported", str)
	}
}
