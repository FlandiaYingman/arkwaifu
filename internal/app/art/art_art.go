package art

type Art struct {
	ID       string `gorm:"primaryKey;type:text COLLATE numeric" json:"id" validate:""`
	Category string `gorm:"" json:"category" validate:"oneof=image background item character"`

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
