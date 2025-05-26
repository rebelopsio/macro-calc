package calculator

import (
	"math"
	"testing"
)

func TestCalculateBMR(t *testing.T) {
	tests := []struct {
		name     string
		sex      Sex
		age      int
		heightCM float64
		weightKG float64
		expected float64
	}{
		{
			name:     "Male BMR",
			sex:      Male,
			age:      30,
			heightCM: 180,
			weightKG: 80,
			expected: 1780,
		},
		{
			name:     "Female BMR",
			sex:      Female,
			age:      25,
			heightCM: 165,
			weightKG: 60,
			expected: 1345.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBMR(tt.sex, tt.age, tt.heightCM, tt.weightKG)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("CalculateBMR() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateTDEE(t *testing.T) {
	bmr := 1500.0
	tests := []struct {
		name          string
		activityLevel ActivityLevel
		expected      float64
	}{
		{"Sedentary", Sedentary, 1800},
		{"Lightly Active", LightlyActive, 2062.5},
		{"Moderately Active", ModeratelyActive, 2325},
		{"Very Active", VeryActive, 2587.5},
		{"Extra Active", ExtraActive, 2850},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTDEE(bmr, tt.activityLevel)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("CalculateTDEE() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRoundToNearest5(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"211 rounds to 210", 211, 210},
		{"213 rounds to 215", 213, 215},
		{"215 stays 215", 215, 215},
		{"217.5 rounds to 220", 217.5, 220},
		{"212.4 rounds to 210", 212.4, 210},
		{"212.5 rounds to 215", 212.5, 215},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundToNearest5(tt.input)
			if result != tt.expected {
				t.Errorf("roundToNearest5(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateDailyMacros(t *testing.T) {
	input := DailyInput{
		Sex:           Male,
		Age:           30,
		HeightCM:      180,
		WeightKG:      80,
		ActivityLevel: ModeratelyActive,
		Goal:          Maintain,
		DietType:      Standard,
	}

	result := CalculateDailyMacros(input)

	// Check that values are reasonable
	if result.BMR < 1500 || result.BMR > 2000 {
		t.Errorf("BMR seems unreasonable: %v", result.BMR)
	}

	if result.TDEE < 2000 || result.TDEE > 3500 {
		t.Errorf("TDEE seems unreasonable: %v", result.TDEE)
	}

	if result.ProteinGrams < 100 || result.ProteinGrams > 200 {
		t.Errorf("Protein seems unreasonable: %v", result.ProteinGrams)
	}

	// Check that all macros are rounded to nearest 5
	if int(result.ProteinGrams)%5 != 0 {
		t.Errorf("Protein not rounded to nearest 5: %v", result.ProteinGrams)
	}
	if int(result.CarbsGrams)%5 != 0 {
		t.Errorf("Carbs not rounded to nearest 5: %v", result.CarbsGrams)
	}
	if int(result.FatGrams)%5 != 0 {
		t.Errorf("Fat not rounded to nearest 5: %v", result.FatGrams)
	}

	// Check that macros add up to approximately the total calories
	macroCalories := (result.ProteinGrams * 4) + (result.CarbsGrams * 4) + (result.FatGrams * 9)
	if math.Abs(macroCalories-result.Calories) > 50 { // Increased tolerance due to rounding
		t.Errorf("Macros don't add up to total calories: %v vs %v", macroCalories, result.Calories)
	}
}

func TestDietTypeMacros(t *testing.T) {
	baseInput := DailyInput{
		Sex:           Male,
		Age:           30,
		HeightCM:      180,
		WeightKG:      80,
		ActivityLevel: ModeratelyActive,
		Goal:          Maintain,
	}

	tests := []struct {
		name            string
		dietType        DietType
		expectedCarbPct float64
		expectedFatPct  float64
		maxCarbs        float64
	}{
		{"Keto", Keto, 0.05, 0.75, 50},
		{"Paleo", Paleo, 0.20, 0.50, 999},
		{"Zone", Zone, 0.40, 0.30, 999},
		{"Low Fat", LowFat, 0.55, 0.15, 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := baseInput
			input.DietType = tt.dietType
			result := CalculateDailyMacros(input)

			// Calculate actual percentages
			totalCalories := result.Calories
			carbCalories := result.CarbsGrams * 4
			fatCalories := result.FatGrams * 9
			
			carbPct := carbCalories / totalCalories
			fatPct := fatCalories / totalCalories

			// Check carb percentage (with tolerance for rounding)
			if math.Abs(carbPct-tt.expectedCarbPct) > 0.05 {
				t.Errorf("%s: Carb percentage = %.2f, want ~%.2f", tt.name, carbPct, tt.expectedCarbPct)
			}

			// Check fat percentage (with tolerance for rounding)
			if math.Abs(fatPct-tt.expectedFatPct) > 0.05 {
				t.Errorf("%s: Fat percentage = %.2f, want ~%.2f", tt.name, fatPct, tt.expectedFatPct)
			}

			// Check keto carb limit
			if tt.dietType == Keto && result.CarbsGrams > tt.maxCarbs {
				t.Errorf("Keto carbs = %.0f, want <= %.0f", result.CarbsGrams, tt.maxCarbs)
			}
		})
	}
}

func TestCalculateWeeklyMacros(t *testing.T) {
	input := WeeklyInput{
		DailyInput: DailyInput{
			Sex:           Female,
			Age:           28,
			HeightCM:      165,
			WeightKG:      65,
			ActivityLevel: LightlyActive,
			Goal:          Lose,
			DietType:      Standard,
		},
		DailyActivities: map[string]ActivityLevel{
			"Monday":    ModeratelyActive,
			"Tuesday":   LightlyActive,
			"Wednesday": ModeratelyActive,
			"Thursday":  LightlyActive,
			"Friday":    ModeratelyActive,
			"Saturday":  VeryActive,
			"Sunday":    Sedentary,
		},
	}

	result := CalculateWeeklyMacros(input)

	// Check that we have 7 days
	if len(result.DailyMacros) != 7 {
		t.Errorf("Expected 7 days of macros, got %d", len(result.DailyMacros))
	}

	// Check that average is reasonable
	if result.Average.Calories < 1500 || result.Average.Calories > 2500 {
		t.Errorf("Average calories seem unreasonable: %v", result.Average.Calories)
	}

	// Verify different activity levels produce different results
	mondayMacros := result.DailyMacros["Monday"]
	sundayMacros := result.DailyMacros["Sunday"]
	if mondayMacros.Calories <= sundayMacros.Calories {
		t.Error("More active day should have more calories")
	}
}

func TestPoundsToKg(t *testing.T) {
	tests := []struct {
		name     string
		pounds   float64
		expected float64
	}{
		{"150 pounds", 150, 68.0388},
		{"200 pounds", 200, 90.7184},
		{"0 pounds", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PoundsToKg(tt.pounds)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("PoundsToKg(%v) = %v, want %v", tt.pounds, result, tt.expected)
			}
		})
	}
}

func TestFeetInchesToCm(t *testing.T) {
	tests := []struct {
		name     string
		feet     float64
		inches   float64
		expected float64
	}{
		{"5 feet 10 inches", 5, 10, 177.8},
		{"6 feet 0 inches", 6, 0, 182.88},
		{"5 feet 5 inches", 5, 5, 165.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FeetInchesToCm(tt.feet, tt.inches)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("FeetInchesToCm(%v, %v) = %v, want %v", tt.feet, tt.inches, result, tt.expected)
			}
		})
	}
}