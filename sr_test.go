package chronoscalc

import "testing"

func TestAggregateToWeekly(t *testing.T) {
	rows := []Candle{
		{Date: "2024-01-01", Open: 100, High: 110, Low: 95, Close: 105, Volume: 10},
		{Date: "2024-01-02", Open: 105, High: 120, Low: 100, Close: 115, Volume: 20},
		{Date: "2024-01-08", Open: 115, High: 125, Low: 110, Close: 120, Volume: 15},
	}
	weekly := AggregateToWeekly(rows)
	if len(weekly) != 2 {
		t.Fatalf("harusnya 2 minggu, dapat %d", len(weekly))
	}
	if weekly[0].High != 120 {
		t.Errorf("High minggu 1 = %v, mau 120", weekly[0].High)
	}
	if weekly[0].Low != 95 {
		t.Errorf("Low minggu 1 = %v, mau 95", weekly[0].Low)
	}
	if weekly[0].Volume != 30 {
		t.Errorf("Volume minggu 1 = %v, mau 30", weekly[0].Volume)
	}
}

func TestComputeSRNotEnoughData(t *testing.T) {
	rows := []Candle{{Date: "2024-01-01", Close: 100}}
	result := ComputeSupportResistance(rows, 0.005)
	if result.Error == "" {
		t.Error("harusnya error kalau data kurang")
	}
}
