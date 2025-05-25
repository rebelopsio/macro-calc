package calculator

import (
	"math"
)

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

type ActivityLevel string

const (
	Sedentary       ActivityLevel = "sedentary"       // Little or no exercise
	LightlyActive   ActivityLevel = "lightly_active"  // Light exercise/sports 1-3 days/week
	ModeratelyActive ActivityLevel = "moderately_active" // Moderate exercise/sports 3-5 days/week
	VeryActive      ActivityLevel = "very_active"     // Hard exercise/sports 6-7 days a week
	ExtraActive     ActivityLevel = "extra_active"    // Very hard exercise/sports & physical job
)

type Goal string

const (
	Maintain Goal = "maintain"
	Lose     Goal = "lose"
	Gain     Goal = "gain"
)

type DailyInput struct {
	Sex           Sex
	Age           int
	HeightCM      float64
	WeightKG      float64
	ActivityLevel ActivityLevel
	Goal          Goal
}

type WeeklyInput struct {
	DailyInput
	DailyActivities map[string]ActivityLevel // Map of day name to activity level
}

type MacroResult struct {
	Calories      float64
	ProteinGrams  float64
	CarbsGrams    float64
	FatGrams      float64
	BMR           float64
	TDEE          float64
}

type WeeklyMacroResult struct {
	DailyMacros map[string]MacroResult
	Average     MacroResult
}

// Conversion functions
func PoundsToKg(pounds float64) float64 {
	return pounds * 0.453592
}

func FeetInchesToCm(feet, inches float64) float64 {
	totalInches := (feet * 12) + inches
	return totalInches * 2.54
}

func activityMultiplier(level ActivityLevel) float64 {
	switch level {
	case Sedentary:
		return 1.2
	case LightlyActive:
		return 1.375
	case ModeratelyActive:
		return 1.55
	case VeryActive:
		return 1.725
	case ExtraActive:
		return 1.9
	default:
		return 1.2
	}
}

func goalCalorieAdjustment(goal Goal) float64 {
	switch goal {
	case Lose:
		return -500 // 500 calorie deficit for ~1 lb/week loss
	case Gain:
		return 300  // 300 calorie surplus for lean gaining
	default:
		return 0
	}
}

func CalculateBMR(sex Sex, age int, heightCM, weightKG float64) float64 {
	var bmr float64
	
	if sex == Male {
		bmr = (10 * weightKG) + (6.25 * heightCM) - (5 * float64(age)) + 5
	} else {
		bmr = (10 * weightKG) + (6.25 * heightCM) - (5 * float64(age)) - 161
	}
	
	return bmr
}

func CalculateTDEE(bmr float64, activityLevel ActivityLevel) float64 {
	return bmr * activityMultiplier(activityLevel)
}

func CalculateMacros(tdee float64, goal Goal, weightKG float64) MacroResult {
	adjustedCalories := tdee + goalCalorieAdjustment(goal)
	
	// Protein: 0.8-1g per pound of body weight (1.8-2.2g per kg)
	proteinMultiplier := 2.0
	if goal == Gain {
		proteinMultiplier = 2.2
	}
	proteinGrams := weightKG * proteinMultiplier
	proteinCalories := proteinGrams * 4
	
	// Fat: 25-30% of total calories
	fatPercentage := 0.25
	if goal == Maintain {
		fatPercentage = 0.30
	}
	fatCalories := adjustedCalories * fatPercentage
	fatGrams := fatCalories / 9
	
	// Carbs: Remaining calories
	remainingCalories := adjustedCalories - proteinCalories - fatCalories
	carbsGrams := remainingCalories / 4
	
	return MacroResult{
		Calories:     math.Round(adjustedCalories),
		ProteinGrams: math.Round(proteinGrams),
		CarbsGrams:   math.Round(carbsGrams),
		FatGrams:     math.Round(fatGrams),
		TDEE:         math.Round(tdee),
	}
}

func CalculateDailyMacros(input DailyInput) MacroResult {
	bmr := CalculateBMR(input.Sex, input.Age, input.HeightCM, input.WeightKG)
	tdee := CalculateTDEE(bmr, input.ActivityLevel)
	result := CalculateMacros(tdee, input.Goal, input.WeightKG)
	result.BMR = math.Round(bmr)
	return result
}

func CalculateWeeklyMacros(input WeeklyInput) WeeklyMacroResult {
	result := WeeklyMacroResult{
		DailyMacros: make(map[string]MacroResult),
	}
	
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	
	var totalCalories, totalProtein, totalCarbs, totalFat float64
	
	for _, day := range days {
		activityLevel := input.DailyActivities[day]
		if activityLevel == "" {
			activityLevel = input.ActivityLevel // Use default if not specified
		}
		
		bmr := CalculateBMR(input.Sex, input.Age, input.HeightCM, input.WeightKG)
		tdee := CalculateTDEE(bmr, activityLevel)
		macros := CalculateMacros(tdee, input.Goal, input.WeightKG)
		macros.BMR = math.Round(bmr)
		
		result.DailyMacros[day] = macros
		
		totalCalories += macros.Calories
		totalProtein += macros.ProteinGrams
		totalCarbs += macros.CarbsGrams
		totalFat += macros.FatGrams
	}
	
	// Calculate averages
	result.Average = MacroResult{
		Calories:     math.Round(totalCalories / 7),
		ProteinGrams: math.Round(totalProtein / 7),
		CarbsGrams:   math.Round(totalCarbs / 7),
		FatGrams:     math.Round(totalFat / 7),
	}
	
	return result
}