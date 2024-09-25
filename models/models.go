package models

type Record struct {
	MongoID          ID               `json:"_id"`
	PaymentMethods   []PaymentMethods `json:"payment_methods"`
	Origin           string           `json:"origin"`
	Name             string           `json:"name"`
	ID               string           `json:"id"`
	NactCode         string           `json:"nact_code"`
	EcheckCode       string           `json:"echeck_code"`
	Type             string           `json:"type"`
	Combo            []any            `json:"combo"`
	CanBuyForOther   bool             `json:"can_buy_for_other"`
	Status           Status           `json:"status"`
	Category         Category         `json:"category"`
	SubscriptionType SubscriptionType `json:"subscription_type"`
	Rule             Rule             `json:"rule"`
	Validity         Validity         `json:"validity"`
	Value            Value            `json:"value"`
	Size             Size             `json:"size"`
	Cost             Cost             `json:"cost"`
	Freebies         []any            `json:"freebies"`
	OfferCategories  []any            `json:"offer_categories"`
}
type ID struct {
	Oid string `json:"$oid"`
}
type PaymentMethods struct {
	ChargingSystem string `json:"charging_system"`
	Icon           string `json:"icon"`
	ID             int    `json:"id"`
	Links          any    `json:"links"`
	Name           string `json:"name"`
	PaymentSource  string `json:"payment_source"`
	Type           string `json:"type"`
	Value          string `json:"value"`
}
type Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
type SubscriptionType struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Name  string `json:"name"`
}
type Product struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ID          int    `json:"id"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
}
type Category struct {
	SubscriptionType SubscriptionType `json:"subscription_type"`
	Description      string           `json:"description"`
	Product          Product          `json:"product"`
	Name             string           `json:"name"`
	ID               int              `json:"id"`
}
type Rule struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	RuleID      int    `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	ID          string `json:"id"`
}
type Validity struct {
	DisplayName  string `json:"display_name"`
	DisplayValue string `json:"display_value"`
	Label        string `json:"label"`
	Unit         string `json:"unit"`
	Value        int    `json:"value"`
	BundleType   string `json:"bundle_type"`
}
type Value struct {
	BundleType   string `json:"bundle_type"`
	DisplayName  string `json:"display_name"`
	DisplayValue int    `json:"display_value"`
	Label        string `json:"label"`
	Restriction  string `json:"restriction"`
	Unit         string `json:"unit"`
	Value        int    `json:"value"`
}
type Size struct {
	BundleType   string `json:"bundle_type"`
	DisplayName  string `json:"display_name"`
	DisplayValue int    `json:"display_value"`
	Label        string `json:"label"`
	Restriction  string `json:"restriction"`
	Unit         string `json:"unit"`
	Value        int    `json:"value"`
}
type Cost struct {
	Currency     string `json:"currency"`
	DisplayName  string `json:"display_name"`
	DisplayValue int    `json:"display_value"`
	Label        string `json:"label"`
	Unit         string `json:"unit"`
	Value        int    `json:"value"`
}

type Amount struct {
	Cost     float64 `json:"amount"`
	Currency string  `json:"currency"`
}
type SheetConfig struct {
	Columns map[string]int `yaml:"columns"`
}

type Config struct {
	Sheets map[string]SheetConfig `yaml:"HSDP-OFFERS.xlsx"`
}
