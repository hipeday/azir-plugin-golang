package properties

type defaultProperties struct {
	Language string `json:"language"`
}

type DefaultProperty struct {
	Property
	Configure defaultProperties `json:"configure"`
}
