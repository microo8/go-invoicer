package generator

const (
	// Invoice define the "invoice" document type
	Invoice string = "INVOICE"

	// Quotation define the "quotation" document type
	Quotation string = "QUOTATION"

	// DeliveryNote define the "delievry note" document type
	DeliveryNote string = "DELIVERY_NOTE"

	// BaseMargin define base margin used in documents
	BaseMargin float64 = 30

	// BaseMarginTop define base margin top used in documents
	BaseMarginTop float64 = 40

	// HeaderMarginTop define base header margin top used in documents
	HeaderMarginTop float64 = 5

	// MaxPageHeight define the maximum height for a single page
	MaxPageHeight float64 = 900
)

// Cols offsets
const (
	// ItemColUnitPriceOffset ...
	ItemColUnitPriceOffset float64 = PageWidth * 0.4

	// ItemColQuantityOffset ...
	ItemColQuantityOffset float64 = PageWidth * 0.5

	// ItemColTotalHTOffset ...
	ItemColTotalHTOffset float64 = PageWidth * 0.55

	// ItemColDiscountOffset ...
	ItemColDiscountOffset float64 = PageWidth * 0.69

	// ItemColTaxOffset ...
	ItemColTaxOffset float64 = PageWidth * 0.75

	// ItemColTotalTTCOffset ...
	ItemColTotalTTCOffset float64 = PageWidth * 0.85
)

var (
	// BaseTextFontSize define the base font size for text in document
	BaseTextFontSize float64 = 8

	// SmallTextFontSize define the small font size for text in document
	SmallTextFontSize float64 = 7

	// ExtraSmallTextFontSize define the extra small font size for text in document
	ExtraSmallTextFontSize float64 = 6

	// LargeTextFontSize define the large font size for text in document
	LargeTextFontSize float64 = 10
)

const (
	ColumnWidth     = 250
	PageWidth       = 592
	itemFontSize    = 8
	itemTitleMargin = 6
	itemsPaddingTop = 40
	metasFontSize   = 8
	contactMargin   = 3
	totalMargin     = 5
	imageHeight     = 80
)
