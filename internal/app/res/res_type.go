package res

import (
	"fmt"
	"path/filepath"
)

type ResourceType string

const (
	Raw       ResourceType = "raw"
	Webp      ResourceType = "webp"
	Thumbnail ResourceType = "thumbnail"
)

func (r ResourceType) FileName(base string) string {
	switch r {
	case Raw:
		return fmt.Sprintf("%s.png", base)
	case Webp, Thumbnail:
		return fmt.Sprintf("%s.webp", base)
	default:
		panic(fmt.Errorf("ResourceType %v not supported", r))
	}
}

func (r ResourceType) Location(baseLocation string, resVersion string) string {
	switch r {
	case Raw, Webp, Thumbnail:
		return filepath.Join(baseLocation, resVersion, string(r))
	default:
		panic(fmt.Errorf("ResourceType %v not supported", r))
	}
}

func ResTypeFromString(str string) (ResourceType, error) {
	if str == "" {
		str = string(Webp)
	}
	resType := ResourceType(str)
	switch resType {
	case Raw, Webp, Thumbnail:
		return resType, nil
	default:
		return "", fmt.Errorf("ResourceType %v not supported", str)
	}
}
