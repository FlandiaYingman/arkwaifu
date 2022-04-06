package updateloop

const (
	KindImages      = "images"
	KindBackgrounds = "backgrounds"

	VariantImg        = "img"
	VariantTimg       = "timg"
	VariantRealEsrgan = "real-esrgan"
	VariantRealCugan  = "real-cugan"
)

var (
	AcceptableAssetKinds = []string{
		KindImages,
		KindBackgrounds,
	}
	AcceptableAssetVariants = []string{
		VariantImg,
		VariantTimg,
		VariantRealEsrgan,
		VariantRealCugan,
	}
)
