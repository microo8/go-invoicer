package generator

// Address represent an address
type Address struct {
	Address    string `json:"address,omitempty" validate:"required"`
	Address2   string `json:"address_2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
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

	return res
}
