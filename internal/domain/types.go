package domain

type RawMaterials struct {
	BoMName       string
	BoMCode       string
	ComponentName string
	ProductID     string
	Description   string
	SKU           string
	Qty           string
}

type BoM struct {
	RawMaterials []RawMaterials
	CostAccount  []CostAccount
}

type CostAccount struct {
	BoMCode string
	Name    string
	Amount  string
}

type BomMap map[string]BoM

type BoMComponents map[string][]RawMaterials
type BoMCostAccount map[string][]CostAccount
