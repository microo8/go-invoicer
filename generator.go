// Package generator allows you to easily generate invoices, delivery notes and quotations in GoLang.
package generator

import (
	_ "embed"

	"github.com/creasty/defaults"
	"github.com/signintech/gopdf"
)

//go:embed Ubuntu-L.ttf
var ubuntuTTF []byte

// New return a new documents with provided types and defaults
func New(docType string, options *Options) (*Document, error) {
	_ = defaults.Set(options)

	doc := &Document{
		Options: options,
		Type:    docType,
	}

	doc.pdf = &gopdf.GoPdf{}
	doc.pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	doc.pdf.AddTTFFontData("Ubuntu", ubuntuTTF)

	return doc, nil
}
