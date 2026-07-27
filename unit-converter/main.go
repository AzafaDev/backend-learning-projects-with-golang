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
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	fmt.Println("server is running on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		category := r.FormValue("category")
		value, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
		from := r.FormValue("from")
		to := r.FormValue("to")

		switch category {
		case "length":
			result, err := converter.ConvertLength(value, from, to)
			if err != nil {
				fmt.Fprintln(w, "error:", err)
				return
			}
			fmt.Fprintln(w, "result:", result)

		case "weight":
			result, err := converter.ConvertWeight(value, from, to)
			if err != nil {
				fmt.Fprintln(w, "error:", err)
				return
			}
			fmt.Fprintln(w, "result:", result)

		case "temperature":
			result, err := converter.ConvertTemperature(value, from, to)
			if err != nil {
				fmt.Fprintln(w, "error:", err)
				return
			}
			fmt.Fprintln(w, "result:", result)
		}
	}

	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Fprintln(w, "error:", err)
			return
		}
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
