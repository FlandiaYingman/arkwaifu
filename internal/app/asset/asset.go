package asset

import "fmt"

type Variant string

var (
	Img      = Variant("img")
	Timg     = Variant("timg")
	Variants = map[string]Variant{
		"img":  Img,
		"timg": Timg,
	}
	VariantExts = map[Variant]string{
		Img:  ".webp",
		Timg: ".webp",
	}
)

func ParseVariant(str string) (Variant, error) {
	variant, ok := Variants[str]
	if !ok {
		return "", fmt.Errorf("str %v cannot be parsed into variant", str)
	}
	return variant, nil
}

type Kind string

var (
	Image       = Kind("images")
	Backgrounds = Kind("backgrounds")
	Kinds       = map[string]Kind{
		"images":      Image,
		"backgrounds": Backgrounds,
	}
)

func ParseKind(str string) (Kind, error) {
	kind, ok := Kinds[str]
	if !ok {
		return "", fmt.Errorf("str %v cannot be parsed into kind", str)
	}
	return kind, nil
}
