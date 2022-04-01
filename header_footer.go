package generator

import (
	"github.com/creasty/defaults"
	"github.com/signintech/gopdf"
)

// HeaderFooter define header or footer informations on document
type HeaderFooter struct {
	Text       string  `json:"text,omitempty"`
	FontSize   float64 `json:"font_size,omitempty" default:"7"`
	Pagination bool    `json:"pagination,omitempty"`
}

type fnc func()

// ApplyFunc allow user to apply custom func
func (hf *HeaderFooter) ApplyFunc(pdf *gopdf.GoPdf, fn fnc) {
	//TODO pdf.SetHeaderFunc(fn)
}

func (hf *HeaderFooter) applyHeader(doc *Document) error {
	if err := defaults.Set(hf); err != nil {
		return err
	}
	return nil
}

func (hf *HeaderFooter) applyFooter(doc *Document) error {
	if err := defaults.Set(hf); err != nil {
		return err
	}

	return nil
}
