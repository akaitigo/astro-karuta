package astronomy

// DefaultLatitude is the default observation latitude (Tokyo).
const DefaultLatitude = 35.68

// visibilityTable maps constellation names to their visible months (1-12).
// Each entry is a slice of month numbers when the constellation is best visible.
var visibilityTable = map[string][]int{
	"オリオン座":      {11, 12, 1, 2, 3},
	"おおぐま座":      {2, 3, 4, 5, 6},
	"さそり座":       {5, 6, 7, 8, 9},
	"こと座":        {6, 7, 8, 9, 10},
	"はくちょう座":     {6, 7, 8, 9, 10},
	"わし座":        {6, 7, 8, 9, 10},
	"おうし座":       {10, 11, 12, 1, 2},
	"ふたご座":       {11, 12, 1, 2, 3},
	"しし座":        {2, 3, 4, 5, 6},
	"おとめ座":       {3, 4, 5, 6, 7},
	"みずがめ座":      {8, 9, 10, 11, 12},
	"うお座":        {9, 10, 11, 12, 1},
	"おひつじ座":      {9, 10, 11, 12, 1},
	"かに座":        {1, 2, 3, 4, 5},
	"てんびん座":      {4, 5, 6, 7, 8},
	"いて座":        {6, 7, 8, 9, 10},
	"やぎ座":        {8, 9, 10, 11, 12},
	"カシオペヤ座":     {9, 10, 11, 12, 1},
	"ペルセウス座":     {10, 11, 12, 1, 2},
	"アンドロメダ座":    {9, 10, 11, 12, 1},
	"こぐま座":       {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
	"おおいぬ座":      {11, 12, 1, 2, 3},
	"こいぬ座":       {12, 1, 2, 3, 4},
	"ぎょしゃ座":      {10, 11, 12, 1, 2},
	"りゅう座":       {5, 6, 7, 8, 9},
	"ケンタウルス座":    {3, 4, 5, 6, 7},
	"みなみじゅうじ座":   {3, 4, 5, 6, 7},
	"ヘルクレス座":     {5, 6, 7, 8, 9},
	"ペガスス座":      {8, 9, 10, 11, 12},
	"うしかい座":      {3, 4, 5, 6, 7},
}

// GetVisibleConstellations returns the names of constellations visible
// in the given month at the specified latitude.
// Latitude adjusts visibility: constellations at extreme southern declinations
// are excluded for high northern latitudes and vice versa.
func GetVisibleConstellations(month int, latitude float64) []string {
	if month < 1 || month > 12 {
		return nil
	}

	// Southern-hemisphere-only constellations that are hard to see
	// from latitudes above +40 degrees.
	southernOnly := map[string]bool{
		"みなみじゅうじ座": true,
		"ケンタウルス座":  true,
	}

	var result []string
	for name, months := range visibilityTable {
		if !containsMonth(months, month) {
			continue
		}
		// At latitudes above 40 N, southern-only constellations are not visible.
		if latitude > 40.0 && southernOnly[name] {
			continue
		}
		result = append(result, name)
	}
	return result
}

// AllConstellationNames returns all constellation names in the visibility table.
func AllConstellationNames() []string {
	names := make([]string, 0, len(visibilityTable))
	for name := range visibilityTable {
		names = append(names, name)
	}
	return names
}

func containsMonth(months []int, target int) bool {
	for _, m := range months {
		if m == target {
			return true
		}
	}
	return false
}
