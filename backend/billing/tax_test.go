package billing

import (
	"emviwiki/shared/config"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestGetTaxId(t *testing.T) {
	config.Get().Stripe.TaxIDDE = "tax-id"
	input := []struct {
		country   string
		taxNumber string
	}{
		{"DE", ""},
		{"DE", "DE123456789"},
		{"FR", ""},
		{"FR", "DE123456789"},
		{"AF", ""},
		{"AF", "AF123456789"},
	}
	expected := []string{
		"tax-id",
		"tax-id",
		"tax-id",
		"",
		"",
		"",
	}

	for i, in := range input {
		if taxId := getTaxId(in.country, in.taxNumber); taxId != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], taxId)
		}
	}
}

func TestGetTaxIdType(t *testing.T) {
	if getTaxIdType("DE") != string(stripe.TaxIDTypeEUVAT) {
		t.Fatal("EU type must have been returned")
	}

	if getTaxIdType("JP") != string(stripe.TaxIDTypeUnknown) {
		t.Fatal("Unknown type must have been returned")
	}
}

func TestGetTaxExempt(t *testing.T) {
	if getTaxExempt("DE") != string(stripe.CustomerTaxExemptNone) {
		t.Fatal("Tax exempt must be none for Germany")
	}

	if getTaxExempt("JP") != string(stripe.CustomerTaxExemptNone) {
		t.Fatal("Tax exempt must be none for countries outside the EU")
	}

	if getTaxExempt("ES") != string(stripe.CustomerTaxExemptReverse) {
		t.Fatal("Tax exempt must be reverse for countries part of the EU")
	}
}
