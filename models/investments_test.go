package models

import "testing"

func TestTotalSellValue(t *testing.T) {
	// If I buy at 20x $10, with 5x leverage
	i := Investment{BuyPrice: 1000, Leverage: 5, Amount: 20}

	if i.Pennies() != 1000 {
		t.Error("Pennies is not right for $10")
	}

	// Mongobucks
	if i.TotalBuyValue() != 2 {
		t.Error("Failed to calculate total buy value in mongobucks")
	}

	// last value was 5
	ticker := Ticker{Last: 1500}
	if i.TotalSellValue(ticker) != 8 {
		t.Error("Failed to calculate total sell value in mongobucks")
	}

	if i.ReturnOnInvestment(ticker) != 6 {
		t.Error("failed to calculate return on investment")
	}

	ticker = Ticker{Last: 500}
	if i.TotalSellValue(ticker) != -6 {
		t.Error("Failed to calcualte negative sell value")
	}
}
