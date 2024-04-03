package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	udmPath = "./templates/udm/"
	data    = []string{
		"HTMX is a great alternative to React etc.",
		"It offers a different way of loading data into your frontend web UI.",
		"It might be especially interesting for server-side developers who are not so familiar with frontend development.",
		"But - as you will see - it's actually also a very promising alternative to React, Angular etc.",
		"You just have to be open for a diffent mental model.",
		"When using HTMX you typically write way less frontend JavaScript code.",
		"You also don't need to manage any frontend state.",
		"Though you can always add extra JS code if needed.",
		"And you can also combine HTMX with other libraries like AlpineJS or integrate it into React apps etc.",
	}
	goals              = []string{}
	availableLocations = []Location{
		{
			ID:    "1",
			Title: "Forest Waterfall",
			Image: Image{
				Src: "forest-waterfall.jpg",
				Alt: "A tranquil forest with a cascading waterfall amidst greenery.",
			},
			Lat: 44.5588,
			Lon: -80.344,
		},
		{
			ID:    "2",
			Title: "Sahara Desert Dunes",
			Image: Image{
				Src: "desert-dunes.jpg",
				Alt: "Golden dunes stretching to the horizon in the Sahara Desert.",
			},
			Lat: 25.0,
			Lon: 0.0,
		},
		{
			ID:    "3",
			Title: "Himalayan Peaks",
			Image: Image{
				Src: "majestic-mountains.jpg",
				Alt: "The sun setting behind snow-capped peaks of majestic mountains.",
			},
			Lat: 27.9881,
			Lon: 86.925,
		},
		{
			ID:    "4",
			Title: "Caribbean Beach",
			Image: Image{
				Src: "caribbean-beach.jpg",
				Alt: "Pristine white sand and turquoise waters of a Caribbean beach.",
			},
			Lat: 18.2208,
			Lon: -66.5901,
		},
		{
			ID:    "5",
			Title: "Ancient Grecian Ruins",
			Image: Image{
				Src: "ruins.jpg",
				Alt: "Historic ruins standing tall against the backdrop of the Grecian sky.",
			},
			Lat: 37.9715,
			Lon: 23.7257,
		},
		{
			ID:    "6",
			Title: "Amazon Rainforest Canopy",
			Image: Image{
				Src: "rainforest.jpg",
				Alt: "Lush canopy of a rainforest, teeming with life.",
			},
			Lat: -3.4653,
			Lon: -62.2159,
		},
		{
			ID:    "7",
			Title: "Northern Lights",
			Image: Image{
				Src: "northern-lights.jpg",
				Alt: "Dazzling display of the Northern Lights in a starry sky.",
			},
			Lat: 64.9631,
			Lon: -19.0208,
		},
		{
			ID:    "8",
			Title: "Japanese Temple",
			Image: Image{
				Src: "japanese-temple.jpg",
				Alt: "Ancient Japanese temple surrounded by autumn foliage.",
			},
			Lat: 34.9949,
			Lon: 135.785,
		},
		{
			ID:    "9",
			Title: "Great Barrier Reef",
			Image: Image{
				Src: "great-barrier-reef.jpg",
				Alt: "Vibrant coral formations of the Great Barrier Reef underwater.",
			},
			Lat: -18.2871,
			Lon: 147.6992,
		},
		{
			ID:    "10",
			Title: "Parisian Streets",
			Image: Image{
				Src: "parisian-streets.jpg",
				Alt: "Charming streets of Paris with historic buildings and cafes.",
			},
			Lat: 48.8566,
			Lon: 2.3522,
		},
		{
			ID:    "11",
			Title: "Grand Canyon",
			Image: Image{
				Src: "grand-canyon.jpg",
				Alt: "Expansive view of the deep gorges and ridges of the Grand Canyon.",
			},
			Lat: 36.1069,
			Lon: -112.1129,
		},
		{
			ID:    "12",
			Title: "Venetian Canals",
			Image: Image{
				Src: "venetian-canals.jpg",
				Alt: "Glistening waters of the Venetian canals as gondolas glide by at sunset.",
			},
			Lat: 45.4408,
			Lon: 12.3155,
		},
		{
			ID:    "13",
			Title: "Taj Mahal",
			Image: Image{
				Src: "taj-mahal.jpg",
				Alt: "The iconic Taj Mahal reflecting in its surrounding waters during sunrise.",
			},
			Lat: 27.1751,
			Lon: 78.0421,
		},
		{
			ID:    "14",
			Title: "Kerala Backwaters",
			Image: Image{
				Src: "kerala-backwaters.jpg",
				Alt: "Tranquil waters and lush greenery of the Kerala backwaters.",
			},
			Lat: 9.4981,
			Lon: 76.3388,
		},
		{
			ID:    "15",
			Title: "African Savanna",
			Image: Image{
				Src: "african-savanna.jpg",
				Alt: "Wild animals roaming freely in the vast landscapes of the African savanna.",
			},
			Lat: -2.153,
			Lon: 34.6857,
		},
		{
			ID:    "16",
			Title: "Victoria Falls",
			Image: Image{
				Src: "victoria-falls.jpg",
				Alt: "The powerful cascade of Victoria Falls, a natural wonder between Zambia and Zimbabwe.",
			},
			Lat: -17.9243,
			Lon: 25.8572,
		},
		{
			ID:    "17",
			Title: "Machu Picchu",
			Image: Image{
				Src: "machu-picchu.jpg",
				Alt: "The historic Incan citadel of Machu Picchu illuminated by the morning sun.",
			},
			Lat: -13.1631,
			Lon: -72.545,
		},
		{
			ID:    "18",
			Title: "Amazon River",
			Image: Image{
				Src: "amazon-river.jpg",
				Alt: "Navigating the waters of the Amazon River, surrounded by dense rainforest.",
			},
			Lat: -3.4653,
			Lon: -58.38,
		},
	}
	interestingLocations = []Location{}
)

type Image struct {
	Src string
	Alt string
}

type Location struct {
	ID    string
	Title string
	Image Image
	Lat   float64
	Lon   float64
}

func (app *application) showAppPageHandler(w http.ResponseWriter, r *http.Request) {
	file1 := udmPath + "main.tmpl"
	file2 := udmPath + "info.tmpl"
	tmpl := template.Must(template.ParseFiles(file1, file2))

	if err := tmpl.ExecuteTemplate(w, "main", envelope{"Data": data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showInfoPageHandler(w http.ResponseWriter, r *http.Request) {

	file := udmPath + "info.tmpl"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, envelope{"Data": data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) postNoteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fv := r.PostFormValue("note")

	data = append([]string{fv}, data...)

	http.Redirect(w, r, "/udm/v1", http.StatusSeeOther) // Method 2

	// Method 1
	// file := udmPath + "info.tmpl"
	// tmpl := template.Must(template.ParseFiles(file))

	// if err := tmpl.ExecuteTemplate(w, "info", envelope{"Data": data}); err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }
}

func (app *application) showGoalPageHandler(w http.ResponseWriter, r *http.Request) {
	file1 := udmPath + "main.tmpl"
	file2 := udmPath + "goal.tmpl"
	tmpl := template.Must(template.ParseFiles(file1, file2))

	if err := tmpl.ExecuteTemplate(w, "main", envelope{"Goals": goals}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) postGoalHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fv := r.PostFormValue("goal")

	goals = append(goals, fv)

	http.Redirect(w, r, "/udm/v2/goals", http.StatusSeeOther)
}

func (app *application) deleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	gid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if gid == 0 {
		goals = goals[1:]
	} else {
		goals = remove(goals, gid)
	}

	http.Redirect(w, r, "/udm/v2/goals", http.StatusSeeOther)
}

func remove[T any](slice []T, index int64) []T {
	if index < 0 || index >= int64(len(slice)) {
		return slice
	}

	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

var IsInterested = true

func (app *application) showPlacePageHandler(w http.ResponseWriter, r *http.Request) {
	file1 := udmPath + "main.tmpl"
	file2 := udmPath + "location.tmpl"
	file3 := udmPath + "deloc.tmpl"
	tmpl := template.Must(template.ParseFiles(file1, file2, file3))

	if err := tmpl.ExecuteTemplate(
		w,
		"main",
		envelope{
			"AvailableLocations":   availableLocations,
			"InterestingLocations": interestingLocations,
			"IsInterested":         IsInterested}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) postPlaceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	as := make([]string, 0)

	for _, v := range body {
		as = append(as, string(v))
	}

	lid := strings.Join(as[11:], "")

	var location Location
	for i, loc := range availableLocations {
		if loc.ID == lid {
			location = loc
			availableLocations = remove(availableLocations, int64(i))
			break
		}
	}

	interestingLocations = append(interestingLocations, location)

	file := udmPath + "location.tmpl"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.ExecuteTemplate(
		w,
		"location",
		envelope{"Location": location, "IsInterested": false}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// http.Redirect(w, r, "/udm/v3/places", http.StatusSeeOther)
	// if err = app.writeJSON(
	// 	w,
	// 	http.StatusAccepted,
	// 	envelope{"InterestingLocations": interestingLocations}, nil); err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }
}

func (app *application) deletePlaceHandler(w http.ResponseWriter, r *http.Request) {
	gid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var (
		index    int64
		location Location
	)

	for i, loc := range interestingLocations {
		s := strconv.FormatInt(gid, 10)
		if loc.ID == s {
			index = int64(i)
		}
	}

	location = interestingLocations[index]

	interestingLocations = remove(interestingLocations, index)

	availableLocations = append(availableLocations, location)

	file := udmPath + "deloc.tmpl"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.ExecuteTemplate(w, "deloc", interestingLocations); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/udm/v3/places", http.StatusSeeOther)
}
