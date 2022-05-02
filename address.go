package generator

// Address represent an address
type Address struct {
	Address    string `json:"address,omitempty" validate:"required"`
	Address2   string `json:"address_2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	BusinessID string
	TaxID      string
	VAT        string
	IBAN       string
	BankName   string
}

// ToString output address as string
// Line break are added for new lines
func (a *Address) lines() []string {
	res := []string{
		a.Address,
	}
	if len(a.Address2) > 0 {
		res = append(res, a.Address2)
	}
	if len(a.PostalCode) > 0 {
		res = append(res, a.PostalCode+" "+a.City)
	}
	if len(a.Country) > 0 {
		res = append(res, a.Country)
	}
	if len(a.BusinessID) > 0 {
		res = append(res, a.BusinessID)
	}
	if len(a.TaxID) > 0 {
		res = append(res, a.TaxID)
	}
	if len(a.VAT) > 0 {
		res = append(res, a.VAT)
	}
	if len(a.IBAN) > 0 {
		res = append(res, a.IBAN)
	}
	if len(a.BankName) > 0 {
		res = append(res, a.BankName)
	}

	return res
}
