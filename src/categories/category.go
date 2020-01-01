package categories

type Category struct {
	Id             CategoryId
	Name           string
	NormalizedName string
}

type CategoryId int
