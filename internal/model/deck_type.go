package model

// DeckType представляет собой перечисление типов колод.
type DeckType string

const (
	DeckTypeFibonacci         DeckType = "fibonacci"
	DeckTypeModifiedFibonacci DeckType = "modified_fibonacci"
	DeckTypeTShirt            DeckType = "tshirt"
	DeckTypeHydra             DeckType = "hydra"
	DeckTypeClassic           DeckType = "classic"
)

// DeckValues сопоставляет каждый тип колоды с набором значений.
var DeckValues = map[DeckType][]string{
	DeckTypeFibonacci:         {"0", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89"},
	DeckTypeModifiedFibonacci: {"0", "½", "1", "2", "3", "5", "8", "13", "20", "40", "100"},
	DeckTypeTShirt:            {"XS", "S", "M", "L", "XL", "XXL"},
	DeckTypeHydra:             {"?", "1", "2", "3", "5", "8", "13"},
	DeckTypeClassic:           {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
}
