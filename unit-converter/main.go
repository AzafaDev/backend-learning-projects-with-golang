package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"unit-converter/converter"
)

type Data struct {
	ActiveCategory string
	Units          []string
	Value          float64
	From           string
	To             string
	Result         float64
	Error          string
}

func main() {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler(tmpl))

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var data Data
			category := r.FormValue("category")
			value, err := strconv.ParseFloat(r.FormValue("value"), 64)
			if err != nil {
				data.Error = err.Error()
				tmpl.Execute(w, data)
				return
			}
			from := r.FormValue("from")
			to := r.FormValue("to")

			data.ActiveCategory = category

			result, units, convertErr := convert(category, value, from, to)
			data.Units = units
			data.Value = value
			data.From = from
			data.To = to
			if convertErr != nil {
				data.Error = convertErr.Error()
			} else {
				data.Result = result
			}
			tmpl.Execute(w, data)
		}

		if r.Method == "GET" {

			category := r.URL.Query().Get("category")
			if category == "" {
				category = "length"
			}

			var data Data

			data.ActiveCategory = category
			switch category {
			case "length":
				data.Units = converter.LengthUnits()
			case "weight":
				data.Units = converter.WeightUnits()
			case "temperature":
				data.Units = converter.TemperatureUnits()
			}

			if err := tmpl.Execute(w, data); err != nil {
				fmt.Fprintln(w, "error:", err)
				return
			}
		}
	}

}

func convert(category string, value float64, from, to string) (float64, []string, error) {
	switch category {
	case "length":
		result, err := converter.ConvertLength(value, from, to)
		if err != nil {
			return 0, converter.LengthUnits(), err
		}
		return result, converter.LengthUnits(), nil

	case "weight":
		result, err := converter.ConvertWeight(value, from, to)
		if err != nil {
			return 0, converter.WeightUnits(), err
		}
		return result, converter.WeightUnits(), nil

	case "temperature":
		result, err := converter.ConvertTemperature(value, from, to)
		if err != nil {
			return 0, converter.TemperatureUnits(), err
		}
		return result, converter.TemperatureUnits(), nil

	default:
		return 0, converter.LengthUnits(), fmt.Errorf("kategori tidak dikenali: %s", category)
	}
}
