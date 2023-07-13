package art

import "fmt"

type Art struct {
	ID       string   `gorm:"primaryKey;type:text COLLATE numeric;check:id=lower(id)" json:"id"`
	Category Category `gorm:"type:art_category" json:"category"`

	Variants []Variant `gorm:"" json:"variants,omitempty" validate:""`
}

func NewArt(id string, category Category) *Art {
	return &Art{
		ID:       id,
		Category: category,
		Variants: nil,
	}
}

type Category string

const (
	CategoryImage      Category = "image"
	CategoryBackground Category = "background"
	CategoryItem       Category = "item"
	CategoryCharacter  Category = "character"
)

func ParseCategory(str string) (Category, error) {
	category := Category(str)
	switch category {
	case
		CategoryImage,
		CategoryBackground,
		CategoryItem,
		CategoryCharacter:
		return category, nil
	default:
		return "", fmt.Errorf("string %q is not a category", str)
	}
}

func MustParseCategory(str string) Category {
	category, err := ParseCategory(str)
	if err != nil {
		panic(err)
	}
	return category
}

func (c *Category) UnmarshalJSON(bytes []byte) error {
	category, err := ParseCategory(string(bytes))
	if err != nil {
		return err
	}
	*c = category
	return nil
}
