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

func TestCalculateDailyMacros(t *testing.T) {
	input := DailyInput{
		Sex:           Male,
		Age:           30,
		HeightCM:      180,
		WeightKG:      80,
		ActivityLevel: ModeratelyActive,
		Goal:          Maintain,
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

	// Check that macros add up to approximately the total calories
	macroCalories := (result.ProteinGrams * 4) + (result.CarbsGrams * 4) + (result.FatGrams * 9)
	if math.Abs(macroCalories-result.Calories) > 10 {
		t.Errorf("Macros don't add up to total calories: %v vs %v", macroCalories, result.Calories)
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