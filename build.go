package generator

import (
	"bytes"
	"fmt"
	"time"

	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	"github.com/signintech/gopdf"
)

// Build pdf document from data provided
func (doc *Document) Build() (*gopdf.GoPdf, error) {
	// Validate document data
	err := doc.Validate()
	if err != nil {
		return nil, err
	}

	// Build base doc
	doc.pdf.SetMargins(BaseMargin, BaseMarginTop, BaseMargin, 0)
	doc.pdf.SetX(10)
	doc.pdf.SetY(10)
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)

	// Set header
	if doc.Header != nil {
		err = doc.Header.applyHeader(doc)

		if err != nil {
			return nil, err
		}
	}

	// Set footer
	if doc.Footer != nil {
		err = doc.Footer.applyFooter(doc)

		if err != nil {
			return nil, err
		}
	}

	// Add first page
	doc.pdf.AddPage()

	// Load font
	doc.pdf.SetFont("Ubuntu", "", 12)

	// Appenf document title
	doc.appendTitle()

	// Appenf document metas (ref & version)
	doc.appendMetas()

	// Append company contact to doc
	companyBottom := doc.Company.appendCompanyContactToDoc(doc)

	// Append customer contact to doc
	customerBottom := doc.Customer.appendCustomerContactToDoc(doc)

	if customerBottom > companyBottom {
		doc.pdf.SetX(10)
		doc.pdf.SetY(customerBottom)
	} else {
		doc.pdf.SetX(10)
		doc.pdf.SetY(companyBottom)
	}

	// Append description
	doc.appendDescription()

	// Append items
	doc.appendItems()

	// Check page height (total bloc height = 30, 45 when doc discount)
	offset := doc.pdf.GetY() + 30
	if doc.Discount != nil {
		offset += 15
	}
	if offset > MaxPageHeight {
		doc.pdf.AddPage()
	}

	// Append notes
	doc.appendNotes()

	// Append total
	doc.appendTotal()

	// Append payment term
	doc.appendPaymentTerm()

	return doc.pdf, nil
}

const (
	titleFontSize = 14
	titleMargin   = 6
)

func (doc *Document) appendTitle() {
	title := doc.typeAsString()

	// Draw rect
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth, BaseMarginTop, PageWidth-BaseMargin, BaseMarginTop+titleFontSize+titleMargin, "F", 0, 0)

	// Set x y
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.SetY(BaseMarginTop + titleMargin/2)

	// Draw text
	doc.pdf.SetFont("Ubuntu", "", titleFontSize)
	doc.pdf.CellWithOption(&gopdf.Rect{W: ColumnWidth, H: titleFontSize}, title, gopdf.CellOption{Align: gopdf.Center})
}

func (doc *Document) appendMetas() {
	// Append ref
	refString := fmt.Sprintf("%s: %s", doc.Options.TextRefTitle, doc.Ref)
	const (
		top = BaseMarginTop + titleFontSize + titleMargin + 1
	)

	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.SetY(top)
	doc.pdf.SetFont("Ubuntu", "", metasFontSize)
	doc.pdf.CellWithOption(&gopdf.Rect{W: ColumnWidth, H: metasFontSize}, refString, gopdf.CellOption{Align: gopdf.Right})

	// Append version
	if len(doc.Version) > 0 {
		versionString := fmt.Sprintf("%s: %s", doc.Options.TextVersionTitle, doc.Version)
		doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
		doc.pdf.SetY(top + metasFontSize)
		doc.pdf.SetFont("Ubuntu", "", metasFontSize)
		doc.pdf.CellWithOption(&gopdf.Rect{W: ColumnWidth, H: metasFontSize}, versionString, gopdf.CellOption{Align: gopdf.Right})
	}

	// Append date
	date := time.Now().Format("02/01/2006")
	if len(doc.Date) > 0 {
		date = doc.Date
	}
	dateString := fmt.Sprintf("%s: %s", doc.Options.TextDateTitle, date)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.SetY(top + metasFontSize*2)
	doc.pdf.SetFont("Ubuntu", "", metasFontSize)
	doc.pdf.CellWithOption(&gopdf.Rect{W: ColumnWidth, H: metasFontSize}, dateString, gopdf.CellOption{Align: gopdf.Right})
}

func (doc *Document) appendDescription() {
	if len(doc.Description) > 0 {
		doc.pdf.SetY(doc.pdf.GetY() + 10)
		doc.pdf.SetFont("Ubuntu", "", 10)
		doc.pdf.MultiCell(&gopdf.Rect{W: 190, H: 5}, doc.Description)
	}
}

func (doc *Document) drawsTableTitles() {
	// Draw rec
	doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.Rectangle(
		BaseMargin,
		doc.pdf.GetY(),
		PageWidth-BaseMargin,
		doc.pdf.GetY()+itemFontSize+itemTitleMargin,
		"F",
		0,
		0,
	)

	// Draw table titles
	doc.pdf.SetX(BaseMargin)
	doc.pdf.SetY(doc.pdf.GetY() + itemTitleMargin/2)
	doc.pdf.SetFont("Ubuntu", "B", itemFontSize)

	// Name
	doc.pdf.SetX(BaseMargin + itemTitleMargin)
	doc.pdf.Cell(nil, doc.Options.TextItemsNameTitle)

	// Unit price
	doc.pdf.SetX(ItemColUnitPriceOffset)
	doc.pdf.Cell(
		&gopdf.Rect{
			W: ItemColQuantityOffset - ItemColUnitPriceOffset,
			H: itemFontSize + itemTitleMargin,
		},
		doc.Options.TextItemsUnitCostTitle,
	)

	// Quantity
	doc.pdf.SetX(ItemColQuantityOffset)
	doc.pdf.Cell(nil, doc.Options.TextItemsQuantityTitle)

	// Total HT
	doc.pdf.SetX(ItemColTotalHTOffset)
	doc.pdf.Cell(nil, doc.Options.TextItemsTotalHTTitle)

	// Tax
	doc.pdf.SetX(ItemColTaxOffset)
	doc.pdf.Cell(nil, doc.Options.TextItemsTaxTitle)

	// Discount
	doc.pdf.SetX(ItemColDiscountOffset)
	doc.pdf.Cell(nil, doc.Options.TextItemsDiscountTitle)

	// TOTAL TTC
	doc.pdf.SetX(ItemColTotalTTCOffset)
	doc.pdf.Cell(nil, doc.Options.TextItemsTotalTTCTitle)
}

func (doc *Document) appendItems() {
	doc.pdf.SetY(doc.pdf.GetY() + itemsPaddingTop)
	doc.drawsTableTitles()

	doc.pdf.SetX(BaseMargin)
	doc.pdf.SetY(doc.pdf.GetY() + itemFontSize + itemTitleMargin)
	doc.pdf.SetFont("Ubuntu", "", itemFontSize)

	for i := 0; i < len(doc.Items); i++ {
		item := doc.Items[i]

		// Check item tax
		if item.Tax == nil {
			item.Tax = doc.DefaultTax
		}

		// Append to pdf
		item.appendColTo(doc.Options, doc)

		if doc.pdf.GetY() > MaxPageHeight {
			// Add page
			doc.pdf.AddPage()
			doc.drawsTableTitles()
			doc.pdf.SetFont("Ubuntu", "", itemFontSize)
		}

		//doc.pdf.SetX(10)
		doc.pdf.SetY(doc.pdf.GetY() + 6)
	}
}

func (doc *Document) appendNotes() {
	if len(doc.Notes) == 0 {
		return
	}

	currentY := doc.pdf.GetY()

	doc.pdf.SetFont("Ubuntu", "", 9)
	doc.pdf.SetX(BaseMargin)
	doc.pdf.SetMarginRight(100)
	doc.pdf.SetY(currentY + 10)

	doc.pdf.MultiCell(
		&gopdf.Rect{W: PageWidth - BaseMargin*2 - ColumnWidth, H: MaxPageHeight * 0.3},
		doc.Notes,
	)

	doc.pdf.SetMarginRight(BaseMargin)
	doc.pdf.SetY(currentY)
}

func (doc *Document) appendTotal() {
	ac := accounting.Accounting{
		Symbol:    (doc.Options.CurrencySymbol),
		Precision: doc.Options.CurrencyPrecision,
		Thousand:  doc.Options.CurrencyThousand,
		Decimal:   doc.Options.CurrencyDecimal,
	}

	// Get total (without tax)
	total, _ := decimal.NewFromString("0")

	for _, item := range doc.Items {
		total = total.Add(item.totalWithoutTaxAndWithDiscount())
	}

	// Apply document discount
	totalWithDiscount := decimal.NewFromFloat(0)
	if doc.Discount != nil {
		discountType, discountNumber := doc.Discount.getDiscount()

		if discountType == "amount" {
			totalWithDiscount = total.Sub(discountNumber)
		} else {
			// Percent
			toSub := total.Mul(discountNumber.Div(decimal.NewFromFloat(100)))
			totalWithDiscount = total.Sub(toSub)
		}
	}

	// Tax
	totalTax := decimal.NewFromFloat(0)
	if doc.Discount == nil {
		for _, item := range doc.Items {
			totalTax = totalTax.Add(item.taxWithDiscount())
		}
	} else {
		discountType, discountAmount := doc.Discount.getDiscount()
		discountPercent := discountAmount
		if discountType == "amount" {
			// Get percent from total discounted
			discountPercent = discountAmount.Mul(decimal.NewFromFloat(100)).Div(totalWithDiscount)
		}

		for _, item := range doc.Items {
			if item.Tax != nil {
				taxType, taxAmount := item.Tax.getTax()
				if taxType == "amount" {
					// If tax type is amount, juste add amount to tax
					totalTax = totalTax.Add(taxAmount)
				} else {
					// Else, remove doc discount % from item total without tax and item discount
					itemTotal := item.totalWithoutTaxAndWithDiscount()
					toSub := discountPercent.Mul(itemTotal).Div(decimal.NewFromFloat(100))
					itemTotalDiscounted := itemTotal.Sub(toSub)

					// Then recompute tax on itemTotalDiscounted
					itemTaxDiscounted := taxAmount.Mul(itemTotalDiscounted).Div(decimal.NewFromFloat(100))

					totalTax = totalTax.Add(itemTaxDiscounted)
				}
			}
		}
	}

	// finalTotal
	totalWithTax := total.Add(totalTax)
	if doc.Discount != nil {
		totalWithTax = totalWithDiscount.Add(totalTax)
	}

	doc.pdf.SetY(doc.pdf.GetY() + 10)
	doc.pdf.SetFont("Ubuntu", "", LargeTextFontSize)
	doc.pdf.SetTextColor(
		doc.Options.BaseTextColor[0],
		doc.Options.BaseTextColor[1],
		doc.Options.BaseTextColor[2],
	)

	// Draw TOTAL HT title
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth, doc.pdf.GetY(), PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		doc.Options.TextTotalNoTax,
		gopdf.CellOption{Align: gopdf.Middle | gopdf.Right},
	)

	// Draw TOTAL HT amount
	doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY(), PageWidth-BaseMargin, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth/2 + totalMargin)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		ac.FormatMoneyDecimal(total),
		gopdf.CellOption{Align: gopdf.Middle},
	)

	if doc.Discount != nil {
		baseY := doc.pdf.GetY() + LargeTextFontSize + totalMargin*2

		// Draw DISCOUNTED title
		doc.pdf.SetY(baseY)
		doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
		doc.pdf.SetStrokeColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
		doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth, doc.pdf.GetY(), PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY()+LargeTextFontSize+5+totalMargin*2, "F", 0, 0)

		// title
		doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
		doc.pdf.SetY(baseY + totalMargin)
		doc.pdf.CellWithOption(
			&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize},
			doc.Options.TextTotalDiscounted,
			gopdf.CellOption{Align: gopdf.Right},
		)

		// description
		doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
		doc.pdf.SetY(baseY + 7.5 + totalMargin)
		doc.pdf.SetFont("Ubuntu", "", BaseTextFontSize)
		doc.pdf.SetTextColor(doc.Options.GreyTextColor[0], doc.Options.GreyTextColor[1], doc.Options.GreyTextColor[2])

		var descString bytes.Buffer
		discountType, discountAmount := doc.Discount.getDiscount()
		if discountType == "percent" {
			descString.WriteString("-")
			descString.WriteString(discountAmount.String())
			descString.WriteString(" % / -")
			descString.WriteString(ac.FormatMoneyDecimal(total.Sub(totalWithDiscount)))
		} else {
			descString.WriteString("-")
			descString.WriteString(ac.FormatMoneyDecimal(discountAmount))
			descString.WriteString(" / -")
			descString.WriteString(
				discountAmount.Mul(decimal.NewFromFloat(100)).Div(total).StringFixed(2),
			)
			descString.WriteString(" %")
		}

		doc.pdf.SetY(baseY + 9.5 + totalMargin)
		doc.pdf.CellWithOption(
			&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: BaseTextFontSize + 2},
			descString.String(),
			gopdf.CellOption{Align: gopdf.Right},
		)

		doc.pdf.SetFont("Ubuntu", "", LargeTextFontSize)
		doc.pdf.SetTextColor(doc.Options.BaseTextColor[0], doc.Options.BaseTextColor[1], doc.Options.BaseTextColor[2])

		// Draw DISCOUNT amount
		doc.pdf.SetY(baseY)
		doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth/2 + totalMargin)
		doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
		doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
		doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY()-totalMargin, PageWidth-BaseMargin, doc.pdf.GetY()+LargeTextFontSize+5+totalMargin*2, "F", 0, 0)
		doc.pdf.CellWithOption(
			&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
			ac.FormatMoneyDecimal(totalWithDiscount),
			gopdf.CellOption{Align: gopdf.Middle},
		)
		doc.pdf.SetY(doc.pdf.GetY() + LargeTextFontSize + 5)
	} else {
		doc.pdf.SetY(doc.pdf.GetY() + LargeTextFontSize)
	}

	// Draw TAX title
	doc.pdf.SetY(doc.pdf.GetY() + totalMargin*2)
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth, doc.pdf.GetY(), PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		doc.Options.TextTotalTax,
		gopdf.CellOption{Align: gopdf.Middle | gopdf.Right},
	)

	// Draw TAX amount
	doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY(), PageWidth-BaseMargin, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth/2 + totalMargin)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		ac.FormatMoneyDecimal(totalTax),
		gopdf.CellOption{Align: gopdf.Middle},
	)

	// Draw TOTAL TTC title
	doc.pdf.SetY(doc.pdf.GetY() + totalMargin*2)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
	doc.pdf.SetY(doc.pdf.GetY() + LargeTextFontSize)
	doc.pdf.SetFillColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.DarkBgColor[0], doc.Options.DarkBgColor[1], doc.Options.DarkBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth, doc.pdf.GetY(), PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		doc.Options.TextTotalWithTax,
		gopdf.CellOption{Align: gopdf.Middle | gopdf.Right},
	)

	// Draw TOTAL TTC amount
	doc.pdf.SetFillColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.SetStrokeColor(doc.Options.GreyBgColor[0], doc.Options.GreyBgColor[1], doc.Options.GreyBgColor[2])
	doc.pdf.Rectangle(PageWidth-BaseMargin-ColumnWidth/2, doc.pdf.GetY(), PageWidth-BaseMargin, doc.pdf.GetY()+LargeTextFontSize+totalMargin*2, "F", 0, 0)
	doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth/2 + totalMargin)
	doc.pdf.CellWithOption(
		&gopdf.Rect{W: ColumnWidth/2 - totalMargin, H: LargeTextFontSize + totalMargin*2},
		ac.FormatMoneyDecimal(totalWithTax),
		gopdf.CellOption{Align: gopdf.Middle},
	)
}

func (doc *Document) appendPaymentTerm() {
	if len(doc.PaymentTerm) > 0 {
		paymentTermString := fmt.Sprintf(
			"%s: %s",
			doc.Options.TextPaymentTermTitle,
			doc.PaymentTerm,
		)
		doc.pdf.SetY(doc.pdf.GetY() + LargeTextFontSize + 5 + totalMargin*2)

		doc.pdf.SetX(PageWidth - BaseMargin - ColumnWidth)
		doc.pdf.SetFont("Ubuntu", "B", LargeTextFontSize)
		doc.pdf.CellWithOption(
			&gopdf.Rect{W: ColumnWidth, H: LargeTextFontSize},
			paymentTermString,
			gopdf.CellOption{Align: gopdf.Right},
		)
	}
}
