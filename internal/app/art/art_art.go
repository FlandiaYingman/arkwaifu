package art

type Art struct {
	ID       string `gorm:"primaryKey;type:text COLLATE numeric;check:id=lower(id)" json:"id"`
	Category string `gorm:"" json:"category"`

	Variants []Variant `gorm:"" json:"variants,omitempty" validate:""`
}

func NewArt(id string, category string) *Art {
	return &Art{
		ID:       id,
		Category: category,
		Variants: nil,
	}
}

const (
	CategoryImage      string = "image"
	CategoryBackground string = "background"
	CategoryItem       string = "item"
	CategoryCharacter  string = "character"
)
