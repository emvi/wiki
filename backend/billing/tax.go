package billing

import (
	"emviwiki/shared/config"
	"github.com/stripe/stripe-go/v71"
)

var iso31661alpha2EU = []string{
	"BE",
	"BG",
	"DK",
	"DE",
	"EE",
	"FI",
	"FR",
	"GR",
	"IE",
	"IT",
	"HR",
	"LV",
	"LT",
	"LU",
	"MT",
	"NL",
	"AT",
	"PL",
	"PT",
	"RO",
	"SK",
	"SI",
	"ES",
	"SE",
	"CZ",
	"HU",
	"CY",
}

// getTaxId returns the Stripe tax ID for given country ISO3166-1-alpha-2 code and tax number.
// The only tax ID returned is the one for Germany at the moment, if the country is part of the EU.
// Else empty strings are returned (no VAT required).
// DE -> tax for private/business
// EU -> tax for private/no tax for businesses
// worldwide -> no taxes
func getTaxId(country, taxNumber string) string {
	if country == "DE" || (taxNumber == "" && isEUCountry(country)) {
		return config.Get().Stripe.TaxIDDE
	}

	return ""
}

// getTaxIdType returns the Stripe tax ID type for given country ISO3166-1-alpha-2 code.
// If it's not an EU country, type "unknown" will be returned.
func getTaxIdType(country string) string {
	if isEUCountry(country) {
		return string(stripe.TaxIDTypeEUVAT)
	}

	return string(stripe.TaxIDTypeUnknown)
}

// isEUCountry reports wether the country is part of the EU or not.
func isEUCountry(country string) bool {
	for _, code := range iso31661alpha2EU {
		if code == country {
			return true
		}
	}

	return false
}

// getTaxExempt returns the tax exempt for given country ISO3166-1-alpha-2 code.
// This will be "reverse" for EU countries except Germany and "none" for Germany and countries outside the EU.
func getTaxExempt(country string) string {
	if isEUCountry(country) && country != "DE" {
		return string(stripe.CustomerTaxExemptReverse)
	}

	return string(stripe.CustomerTaxExemptNone)
}
