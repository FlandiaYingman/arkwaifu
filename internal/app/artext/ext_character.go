package artext

import (
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

func parseCharacterID(characterID string) (base string, body int, face int, err error) {
	// Regex Pattern like: {BASE}#{FACE}${BODY}
	pattern := regexp.MustCompile(`(\w+)#(\d+)\$(\d+)`)
	matches := pattern.FindStringSubmatch(characterID)
	if matches == nil {
		return "", 0, 0, errors.Errorf("invalid character ID %s", characterID)
	}

	base = matches[1]
	face, err = strconv.Atoi(matches[2])
	body, err = strconv.Atoi(matches[3])

	if err != nil {
		return "", 0, 0, errors.Errorf("invalid character ID %s: %w", characterID, err)
	}
	return base, body, face, nil
}

// GetSiblingsOfCharacterArt gets the sibling characters of the specified character.
//
// Sibling Characters: the characters that have the same base name as the specified character.
func (s *Service) GetSiblingsOfCharacterArt(characterID string) (siblings []*art.Art, err error) {
	base, _, _, err := parseCharacterID(characterID)
	if err != nil {
		return nil, err
	}
	return s.art.SelectArtsByIDLike(fmt.Sprintf("%s%%", base))
}

func (c *Controller) GetSiblingsOfCharacterArt(ctx *fiber.Ctx) error {
	characterID := ctx.Params("id")
	siblings, err := c.service.GetSiblingsOfCharacterArt(characterID)
	if err != nil {
		return err
	}
	return ctx.JSON(siblings)
}
