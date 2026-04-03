package astronomy

import (
	"sort"
	"testing"
)

func TestGetVisibleConstellations_January(t *testing.T) {
	result := GetVisibleConstellations(1, DefaultLatitude)
	if len(result) == 0 {
		t.Fatal("expected visible constellations in January")
	}

	// January should include winter constellations
	expected := map[string]bool{
		"オリオン座":  true,
		"おうし座":   true,
		"ふたご座":   true,
		"おおいぬ座":  true,
		"こいぬ座":   true,
		"こぐま座":   true,
		"かに座":    true,
		"ペルセウス座": true,
	}

	resultSet := make(map[string]bool)
	for _, name := range result {
		resultSet[name] = true
	}

	for name := range expected {
		if !resultSet[name] {
			t.Errorf("expected %s to be visible in January", name)
		}
	}
}

func TestGetVisibleConstellations_July(t *testing.T) {
	result := GetVisibleConstellations(7, DefaultLatitude)
	if len(result) == 0 {
		t.Fatal("expected visible constellations in July")
	}

	// July should include summer constellations
	expected := map[string]bool{
		"さそり座":   true,
		"こと座":    true,
		"はくちょう座": true,
		"わし座":    true,
	}

	resultSet := make(map[string]bool)
	for _, name := range result {
		resultSet[name] = true
	}

	for name := range expected {
		if !resultSet[name] {
			t.Errorf("expected %s to be visible in July", name)
		}
	}
}

func TestGetVisibleConstellations_Circumpolar(t *testing.T) {
	// こぐま座 should be visible all year at Tokyo latitude
	for month := 1; month <= 12; month++ {
		result := GetVisibleConstellations(month, DefaultLatitude)
		found := false
		for _, name := range result {
			if name == "こぐま座" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected こぐま座 to be visible in month %d", month)
		}
	}
}

func TestGetVisibleConstellations_InvalidMonth(t *testing.T) {
	result := GetVisibleConstellations(0, DefaultLatitude)
	if result != nil {
		t.Error("expected nil for month 0")
	}

	result = GetVisibleConstellations(13, DefaultLatitude)
	if result != nil {
		t.Error("expected nil for month 13")
	}

	result = GetVisibleConstellations(-1, DefaultLatitude)
	if result != nil {
		t.Error("expected nil for month -1")
	}
}

func TestGetVisibleConstellations_HighLatitude_ExcludesSouthern(t *testing.T) {
	// At latitude > 40, southern-only constellations should be excluded
	result := GetVisibleConstellations(5, 50.0)
	for _, name := range result {
		if name == "みなみじゅうじ座" || name == "ケンタウルス座" {
			t.Errorf("expected %s to be excluded at latitude 50", name)
		}
	}
}

func TestGetVisibleConstellations_LowLatitude_IncludesSouthern(t *testing.T) {
	// At latitude <= 40, southern constellations should be included
	result := GetVisibleConstellations(5, 30.0)
	resultSet := make(map[string]bool)
	for _, name := range result {
		resultSet[name] = true
	}

	if !resultSet["ケンタウルス座"] {
		t.Error("expected ケンタウルス座 to be visible at latitude 30")
	}
	if !resultSet["みなみじゅうじ座"] {
		t.Error("expected みなみじゅうじ座 to be visible at latitude 30")
	}
}

func TestAllConstellationNames(t *testing.T) {
	names := AllConstellationNames()
	if len(names) != 30 {
		t.Errorf("expected 30 constellation names, got %d", len(names))
	}

	// Verify no duplicates
	sort.Strings(names)
	for i := 1; i < len(names); i++ {
		if names[i] == names[i-1] {
			t.Errorf("duplicate constellation name: %s", names[i])
		}
	}
}

func TestGetVisibleConstellations_SeasonalTransition(t *testing.T) {
	// Verify that spring and autumn have different dominant constellations
	springResult := GetVisibleConstellations(4, DefaultLatitude)
	autumnResult := GetVisibleConstellations(10, DefaultLatitude)

	springSet := make(map[string]bool)
	for _, name := range springResult {
		springSet[name] = true
	}
	autumnSet := make(map[string]bool)
	for _, name := range autumnResult {
		autumnSet[name] = true
	}

	// Spring should have these but not autumn
	if !springSet["しし座"] {
		t.Error("expected しし座 in spring")
	}
	if autumnSet["しし座"] {
		t.Error("expected しし座 not in autumn")
	}

	// Autumn should have these but not spring
	if !autumnSet["ペガスス座"] {
		t.Error("expected ペガスス座 in autumn")
	}
	if springSet["ペガスス座"] {
		t.Error("expected ペガスス座 not in spring")
	}
}
