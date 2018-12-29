package main

import "testing"

func Testmian(t *testing.T) {
	tables := []struct {
		x string
		n []string
	}{
		{"casesTest/01.eml", [2018 2018]},
		{"casesTest/02.eml", [2018 2019]},
		{"casesTest/03.eml", [2018 2018]},
		{"casesTest/04.eml", [2018 2019]},
		{"casesTest/05.eml", [2018 2018]},
		{"casesTest/06.eml", [2018 2018]},
		{"casesTest/07.eml", [2018 2018]},
		{"casesTest/08.eml", [2018 2019]},
		{"casesTest/09.eml", [2019 2019]},
		{"casesTest/10.eml", [2018 2019]},
		{"casesTest/11.eml", [2018 2018]},
		{"casesTest/12.eml", [2019 2019]},
		{"casesTest/13.eml", [2018 2019]},
	}

	for _, table := range tables {
		total := Sum(table.x, table.y)
		if total != table.n {
			t.Errorf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", table.x, table.y, total, table.n)
		}
	}
}
