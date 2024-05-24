package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"

	"github.com/gotloaf/earthly"
	_ "github.com/gotloaf/earthly/build/docs"
)

// @title Earthly API
// @version 1.0
// @description API for generating Earth images.

// @license.name MIT
// @license.url https://opensource.org/license/mit

// @BasePath /
func main() {
	fmt.Println("Starting earthly server")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", RootGenerate)

	// Swagger endpoint
	r.Mount("/swagger", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", r)
}

// RootGenerate godoc
//
//		@Summary                 Generate Earth Image
//		@Description             Generates an Earth image based on the provided configuration.
//		@Tags                    earth
//		@Accept                  json
//		@Produce                 image/png
//	    @Param strict            query boolean  false "Whether to produce errors instead of falling back for parameters that couldn't be parsed"
//		@Param size              query int      false "Image size" default(512) minimum(16) maximum(2048)
//		@Param latitude          query number   false "Latitude" default(0.0) minimum(-90.0) maximum(90.0)
//		@Param longitude         query number   false "Longitude" default(0.0) minimum(-180.0) maximum(180.0)
//		@Param roll              query number   false "Roll" default(0.0) minimum(-180.0) maximum(180.0)
//		@Param zoom              query number   false "Zoom" default(1.0) minimum(0.01) maximum(4.0)
//		@Success 200 {file}      image/png
//	    @Failure 400 {object}    main.RootGenerate.strictModeError "Only occurs in strict mode. Returns description of fields that failed the strict check and why."
//		@Failure 500 {object}    main.RootGenerate.internalServerError
//		@Router / [get]
func RootGenerate(w http.ResponseWriter, r *http.Request) {
	config, errors := requestToConfiguration(r)

	// Strict mode checks
	type strictModeError struct {
		Error string `json:"error"`
	}

	if slices.Contains([]string{"true", "1"}, strings.ToLower(r.URL.Query().Get("strict"))) && len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(strictModeError{
			Error: strings.Join(errors, "\n"),
		})
		return
	}

	// Load images
	type internalServerError struct {
		Error string `json:"error"`
	}

	var image_format = "1x.png"
	if float64(config.Size)*config.Radius > 1024 {
		image_format = "2x.png"
	}

	image_file, err := os.ReadFile(fmt.Sprintf("build/equirectangular/earth_%s", image_format))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerError{Error: "An internal error occurred."})
		return
	}

	output := config.Generate(earthly.EarthlyBuffers{
		Earth: image_file,
	}, false)

	if output == nil || output.Len() == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerError{Error: "An internal error occurred."})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(output.Bytes())
}

func requestToConfiguration(r *http.Request) (earthly.EarthlyConfig, []string) {
	var errors []string

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 512
		errors = append(errors, "param `size` could not be parsed")
	}
	if size < 16 {
		size = 16
		errors = append(errors, "param `size` is not within bounds [16-2048]")
	} else if size > 2048 {
		size = 2048
		errors = append(errors, "param `size` is not within bounds [16-2048]")
	}

	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil {
		latitude = 0.0
		errors = append(errors, "param `latitude` could not be parsed")
	}

	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil {
		longitude = 0.0
		errors = append(errors, "param `longitude` could not be parsed")
	}

	var rawRoll = r.URL.Query().Get("roll")
	roll, err := strconv.ParseFloat(rawRoll, 64)
	if err != nil {
		roll = 0.0
		if len(rawRoll) > 0 {
			errors = append(errors, "param `roll` could not be parsed")
		}
	}

	var rawZoom = r.URL.Query().Get("zoom")
	zoom, err := strconv.ParseFloat(rawZoom, 64)
	if err != nil {
		zoom = 1.0
		if len(rawZoom) > 0 {
			errors = append(errors, "param `zoom` could not be parsed")
		}
	}

	return earthly.EarthlyConfig{
		Size:       size,
		Background: []uint8{0, 0, 0, 0},
		Latitude:   latitude,
		Longitude:  longitude,
		Roll:       roll,
		Halo:       true,
		Radius:     zoom,
	}, errors
}
