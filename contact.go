package generator

import (
	"bytes"
	"image"

	_ "image/jpeg"
	_ "image/png"

	"github.com/signintech/gopdf"
)

// Contact contact a company informations
type Contact struct {
	Name    string   `json:"name,omitempty" validate:"required,min=1,max=256"`
	Logo    *[]byte  `json:"logo,omitempty"` // Logo byte array
	Address *Address `json:"address,omitempty"`
}

func (c *Contact) appendContactTODoc(
	x float64,
	y float64,
	fill bool,
	logoAlign string,
	doc *Document,
) float64 {
	doc.pdf.SetX(x)
	doc.pdf.SetY(y)

	// Logo
	if c.Logo != nil {
		img, _, err := image.Decode(bytes.NewReader(*c.Logo))
		if err != nil {
			panic(err)
		}
		b := img.Bounds()
		imgH, err := gopdf.ImageHolderByBytes(*c.Logo)
		if err != nil {
			panic(err)
		}
		if err := doc.pdf.ImageByHolderWithOptions(
			imgH,
			gopdf.ImageOptions{
				X:    x,
				Y:    y,
				Rect: &gopdf.Rect{W: imageHeight * float64(b.Dx()) / float64(b.Dy()), H: imageHeight},
				Transparency: &gopdf.Transparency{
					Alpha:         0.0,
					BlendModeType: "",
				},
			},
		); err != nil {
			panic(err)
		}
		doc.pdf.SetY(y + imageHeight + 3)
	}

	// Name
	if fill {
		doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
		doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	} else {
		doc.pdf.SetFillColor(255, 255, 255)
		doc.pdf.SetStrokeColor(255, 255, 255)
	}

	// Name rect
	doc.pdf.Rectangle(x, doc.pdf.GetY(), x+ColumnWidth, doc.pdf.GetY()+LargeTextFontSize, "F", 0, 0)

	// Reset x
	doc.pdf.SetX(x + contactMargin)
	// Set name
	doc.pdf.SetFont("Ubuntu", "B", LargeTextFontSize)
	doc.pdf.Cell(nil, c.Name)
	doc.pdf.SetFont("Ubuntu", "", LargeTextFontSize)

	if c.Address != nil {
		// Address rect
		lines := c.Address.lines()
		var addrRectHeight float64 = LargeTextFontSize * float64(len(lines)+1)

		offsetY := doc.pdf.GetY() + LargeTextFontSize + 3
		doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
		doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
		doc.pdf.Rectangle(x, offsetY, x+ColumnWidth, doc.pdf.GetY()+addrRectHeight+contactMargin*2, "F", 0, 0)

		doc.pdf.SetFont("Ubuntu", "", LargeTextFontSize)
		doc.pdf.SetX(x + contactMargin)
		doc.pdf.SetY(offsetY + contactMargin)
		// Set address
		for _, line := range lines {
			doc.pdf.MultiCell(&gopdf.Rect{W: ColumnWidth, H: addrRectHeight}, line)
		}
	}

	return doc.pdf.GetY()
}

func (c *Contact) appendCompanyContactToDoc(doc *Document) float64 {
	return c.appendContactTODoc(BaseMargin, BaseMarginTop, true, "L", doc)
}

func (c *Contact) appendCustomerContactToDoc(doc *Document) float64 {
	return c.appendContactTODoc(PageWidth-BaseMargin-ColumnWidth, BaseMarginTop+45, true, "R", doc)
}
