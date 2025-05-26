package handlers

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rebelopsio/macro-calc/internal/calculator"
	"github.com/rebelopsio/macro-calc/internal/templates"
)

func Index(c echo.Context) error {
	component := templates.Layout("Macro Calculator")
	handler := templ.Handler(component)
	return handler.Component.Render(c.Request().Context(), c.Response())
}

func ShowCalculator(c echo.Context) error {
	component := templates.CalculatorForm()
	handler := templ.Handler(component)
	return handler.Component.Render(c.Request().Context(), c.Response())
}

func Calculate(c echo.Context) error {
	// Parse form values
	sex := calculator.Sex(c.FormValue("sex"))
	age, _ := strconv.Atoi(c.FormValue("age"))
	unitSystem := c.FormValue("unit_system")
	
	var height, weight float64
	
	if unitSystem == "imperial" {
		// Convert imperial to metric
		feet, _ := strconv.ParseFloat(c.FormValue("feet"), 64)
		inches, _ := strconv.ParseFloat(c.FormValue("inches"), 64)
		height = calculator.FeetInchesToCm(feet, inches)
		
		weightLbs, _ := strconv.ParseFloat(c.FormValue("weight_lbs"), 64)
		weight = calculator.PoundsToKg(weightLbs)
	} else {
		// Already in metric
		height, _ = strconv.ParseFloat(c.FormValue("height"), 64)
		weight, _ = strconv.ParseFloat(c.FormValue("weight"), 64)
	}
	
	activity := calculator.ActivityLevel(c.FormValue("activity"))
	goal := calculator.Goal(c.FormValue("goal"))
	advanced := c.FormValue("advanced") == "on"
	
	// Get diet type, default to standard if not provided
	dietType := calculator.DietType(c.FormValue("diet_type"))
	if dietType == "" {
		dietType = calculator.Standard
	}

	if advanced {
		// Handle weekly calculation
		weeklyInput := calculator.WeeklyInput{
			DailyInput: calculator.DailyInput{
				Sex:           sex,
				Age:           age,
				HeightCM:      height,
				WeightKG:      weight,
				ActivityLevel: activity,
				Goal:          goal,
				DietType:      dietType,
			},
			DailyActivities: make(map[string]calculator.ActivityLevel),
		}

		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, day := range days {
			if dayActivity := c.FormValue(day + "-activity"); dayActivity != "" {
				weeklyInput.DailyActivities[day] = calculator.ActivityLevel(dayActivity)
			}
		}

		result := calculator.CalculateWeeklyMacros(weeklyInput)
		component := templates.WeeklyMacroResults(result)
		handler := templ.Handler(component)
		return handler.Component.Render(c.Request().Context(), c.Response())
	} else {
		// Handle daily calculation
		input := calculator.DailyInput{
			Sex:           sex,
			Age:           age,
			HeightCM:      height,
			WeightKG:      weight,
			ActivityLevel: activity,
			Goal:          goal,
			DietType:      dietType,
		}

		result := calculator.CalculateDailyMacros(input)
		component := templates.MacroResults(result)
		handler := templ.Handler(component)
		return handler.Component.Render(c.Request().Context(), c.Response())
	}
}