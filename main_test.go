package main

import (
	"testing"
)

func TestCalculatePoints(t *testing.T) {
	type args struct {
		receipt Receipt
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "Target",
			args: args{
				receipt: Receipt{
					Retailer:     "Target",
					PurchaseDate: "2022-01-01",
					PurchaseTime: "13:01",
					Items: []Item{
						{"Mountain Dew 12PK", "6.49"},
						{"Emils Cheese Pizza", "12.25"},
						{"Knorr Creamy Chicken", "1.26"},
						{"Doritos Nacho Cheese", "3.35"},
						{"   Klarbrunn 12-PK 12 FL OZ  ", "12.00"},
					},
					Total: "35.35",
				},
			},
			want: 28,
		},
		{name: "M&M Corner Market",
			args: args{
				receipt: Receipt{
					Retailer:     "M&M Corner Market",
					PurchaseDate: "2022-03-20",
					PurchaseTime: "14:33",
					Items: []Item{
						{"Gatorade", "2.25"},
						{"Gatorade", "2.25"},
						{"Gatorade", "2.25"},
						{"Gatorade", "2.25"},
					},
					Total: "9.00",
				},
			},
			want: 109,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePoints(tt.args.receipt); got != tt.want {
				t.Errorf("CalculatePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
