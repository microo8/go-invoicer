package generator

import (
	"fmt"

	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
	"github.com/signintech/gopdf"
)

// Item represent a 'product' or a 'service'
type Item struct {
	Name        string    `json:"name,omitempty" validate:"required"`
	Description string    `json:"description,omitempty"`
	UnitCost    string    `json:"unit_cost,omitempty"`
	Quantity    string    `json:"quantity,omitempty"`
	Tax         *Tax      `json:"tax,omitempty"`
	Discount    *Discount `json:"discount,omitempty"`
}

func (i *Item) unitCost() decimal.Decimal {
	unitCost, _ := decimal.NewFromString(i.UnitCost)
	return unitCost
}

func (i *Item) quantity() decimal.Decimal {
	quantity, _ := decimal.NewFromString(i.Quantity)
	return quantity
}

func (i *Item) totalWithoutTax() decimal.Decimal {
	quantity, _ := decimal.NewFromString(i.Quantity)
	price, _ := decimal.NewFromString(i.UnitCost)
	total := price.Mul(quantity)

	return total
}

func (i *Item) totalWithoutTaxAndWithDiscount() decimal.Decimal {
	total := i.totalWithoutTax()

	// Check discount
	if i.Discount != nil {
		dType, dNum := i.Discount.getDiscount()

		if dType == "amount" {
			total = total.Sub(dNum)
		} else {
			// Percent
			toSub := total.Mul(dNum.Div(decimal.NewFromFloat(100)))
			total = total.Sub(toSub)
		}
	}

	return total
}

func (i *Item) totalWithTaxAndDiscount() decimal.Decimal {
	return i.totalWithoutTaxAndWithDiscount().Add(i.taxWithDiscount())
}

func (i *Item) taxWithDiscount() decimal.Decimal {
	result := decimal.NewFromFloat(0)

	if i.Tax == nil {
		return result
	}

	totalHT := i.totalWithoutTaxAndWithDiscount()
	taxType, taxAmount := i.Tax.getTax()

	if taxType == "amount" {
		result = taxAmount
	} else {
		divider := decimal.NewFromFloat(100)
		result = totalHT.Mul(taxAmount.Div(divider))
	}

	return result
}

func (i *Item) appendColTo(options *Options, doc *Document) {
	ac := accounting.Accounting{
		Symbol:    (options.CurrencySymbol),
		Precision: options.CurrencyPrecision,
		Thousand:  options.CurrencyThousand,
		Decimal:   options.CurrencyDecimal,
	}

	// Get base Y (top of line)
	baseY := doc.pdf.GetY()

	// Name
	doc.pdf.SetX(BaseMargin + itemTitleMargin)
	doc.pdf.MultiCell(
		&gopdf.Rect{
			W: ItemColUnitPriceOffset - BaseMargin - itemTitleMargin*2,
			H: itemFontSize * 3,
		},
		i.Name,
	)

	// Description
	if len(i.Description) > 0 {
		doc.pdf.SetX(BaseMargin + itemTitleMargin)

		doc.pdf.SetFont("Ubuntu", "", SmallTextFontSize)
		doc.pdf.SetTextColor(
			doc.Options.GreyTextColor[0],
			doc.Options.GreyTextColor[1],
			doc.Options.GreyTextColor[2],
		)

		doc.pdf.MultiCell(
			&gopdf.Rect{
				W: ItemColUnitPriceOffset - BaseMargin - itemTitleMargin*2,
				H: itemFontSize * 3,
			},
			i.Description,
		)

		// Reset font
		doc.pdf.SetFont("Ubuntu", "", BaseTextFontSize)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
	}

	// Compute line height
	colHeight := doc.pdf.GetY() - baseY

	// Unit price
	doc.pdf.SetY(baseY)
	doc.pdf.SetX(ItemColUnitPriceOffset)
	doc.pdf.Cell(&gopdf.Rect{W: ItemColQuantityOffset - ItemColUnitPriceOffset, H: colHeight}, ac.FormatMoneyDecimal(i.unitCost())) //, "0", 0, "", false, 0, "")

	// Quantity
	doc.pdf.SetX(ItemColQuantityOffset)
	doc.pdf.Cell(&gopdf.Rect{W: ItemColTaxOffset - ItemColQuantityOffset, H: colHeight}, i.quantity().String()) //, "0", 0, "", false, 0, "")

	// Total HT
	doc.pdf.SetX(ItemColTotalHTOffset)
	doc.pdf.Cell(&gopdf.Rect{W: ItemColTaxOffset - ItemColTotalHTOffset, H: colHeight}, ac.FormatMoneyDecimal(i.totalWithoutTax())) //, "0", 0, "", false, 0, "")

	// Discount
	doc.pdf.SetX(ItemColDiscountOffset)
	if i.Discount == nil {
		doc.pdf.Cell(&gopdf.Rect{W: ItemColTotalTTCOffset - ItemColDiscountOffset, H: colHeight}, "--")
	} else {
		// If discount
		discountType, discountAmount := i.Discount.getDiscount()
		var discountTitle string
		var discountDesc string

		if discountType == "percent" {
			discountTitle = fmt.Sprintf("%s %s", discountAmount, "%")
			// get amount from percent
			dCost := i.totalWithoutTax()
			dAmount := dCost.Mul(discountAmount.Div(decimal.NewFromFloat(100)))
			discountDesc = fmt.Sprintf("-%s", ac.FormatMoneyDecimal(dAmount))
		} else {
			discountTitle = fmt.Sprintf("%s %s", discountAmount, "€")
			dCost := i.totalWithoutTax()
			dPerc := discountAmount.Mul(decimal.NewFromFloat(100))
			dPerc = dPerc.Div(dCost)
			// get percent from amount
			discountDesc = fmt.Sprintf("-%s %%", dPerc.StringFixed(2))
		}

		// discount title
		// lastY := doc.pdf.GetY()
		doc.pdf.Cell(&gopdf.Rect{W: ItemColTotalTTCOffset - ItemColDiscountOffset, H: colHeight / 2}, discountTitle)

		// discount desc
		doc.pdf.SetX(ItemColDiscountOffset)
		doc.pdf.SetY(baseY + BaseTextFontSize)
		doc.pdf.SetFont("Ubuntu", "", SmallTextFontSize)
		doc.pdf.SetTextColor(doc.Options.GreyTextColor[0], doc.Options.GreyTextColor[1], doc.Options.GreyTextColor[2])

		doc.pdf.Cell(&gopdf.Rect{W: ItemColTotalTTCOffset - ItemColDiscountOffset, H: colHeight / 2}, discountDesc)

		// reset font and y
		doc.pdf.SetFont("Ubuntu", "", BaseTextFontSize)
		doc.pdf.SetTextColor(
			doc.Options.BaseTextColor[0],
			doc.Options.BaseTextColor[1],
			doc.Options.BaseTextColor[2],
		)
		doc.pdf.SetY(baseY)
	}

	// Tax
	doc.pdf.SetX(ItemColTaxOffset)
	if i.Tax == nil {
		// If no tax
		doc.pdf.Cell(&gopdf.Rect{W: ItemColDiscountOffset - ItemColTaxOffset, H: colHeight}, "--")
	} else {
		// If tax
		taxType, taxAmount := i.Tax.getTax()
		var taxTitle string
		var taxDesc string

		if taxType == "percent" {
			taxTitle = fmt.Sprintf("%s %s", taxAmount, ("%"))
			// get amount from percent
			dCost := i.totalWithoutTaxAndWithDiscount()
			dAmount := dCost.Mul(taxAmount.Div(decimal.NewFromFloat(100)))
			taxDesc = ac.FormatMoneyDecimal(dAmount)
		} else {
			taxTitle = fmt.Sprintf("%s %s", taxAmount, ("€"))
			dCost := i.totalWithoutTaxAndWithDiscount()
			dPerc := taxAmount.Mul(decimal.NewFromFloat(100))
			dPerc = dPerc.Div(dCost)
			// get percent from amount
			taxDesc = fmt.Sprintf("%s %%", dPerc.StringFixed(2))
		}

		// tax title
		// lastY := doc.pdf.GetY()
		doc.pdf.Cell(&gopdf.Rect{W: ItemColDiscountOffset - ItemColTaxOffset, H: colHeight / 2}, taxTitle)

		// tax desc
		doc.pdf.SetX(ItemColTaxOffset)
		doc.pdf.SetY(baseY + BaseTextFontSize)
		doc.pdf.SetFont("Ubuntu", "", SmallTextFontSize)
		doc.pdf.SetTextColor(doc.Options.GreyTextColor[0], doc.Options.GreyTextColor[1], doc.Options.GreyTextColor[2])

		doc.pdf.Cell(&gopdf.Rect{W: ItemColDiscountOffset - ItemColTaxOffset, H: colHeight / 2}, taxDesc)

		// reset font and y
		doc.pdf.SetFont("Ubuntu", "", BaseTextFontSize)
		doc.pdf.SetTextColor(doc.Options.BaseTextColor[0], doc.Options.BaseTextColor[1], doc.Options.BaseTextColor[2])
		doc.pdf.SetY(baseY)
	}

	// TOTAL TTC
	doc.pdf.SetX(ItemColTotalTTCOffset)
	doc.pdf.Cell(&gopdf.Rect{W: 190 - ItemColTotalTTCOffset, H: colHeight}, ac.FormatMoneyDecimal(i.totalWithTaxAndDiscount()))

	// Set Y for next line
	doc.pdf.SetY(baseY + colHeight)
}
