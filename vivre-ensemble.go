package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/andybrewer/mack"
)

const PHONE_NUMBER = "+352 123-45678"
const DATA_DIR = "/Users/irio/Workspace/vivre-ensemble/data/"

type Course struct {
	Type             string  `json:"lu.etat.cgie.rest.dto:type"`
	SchoolYear       string  `json:"anneeScolaire"`
	CategoryCode     string  `json:"codeCategorie"`
	LanguageCode     string  `json:"codeLangue"`
	CourseDate       string  `json:"dateCours"`
	DurationHours    float32 `json:"duree"`
	Duration         string  `json:"dureeLib"`
	Instructor       string  `json:"formateur"`
	CourseSchedule   string  `json:"horaireCours"`
	ID               string  `json:"id"`
	ModuleID         string
	CourseID         string
	CourseInstanceID int
	Title            string `json:"intitule"`
	IsOnline         bool   `json:"isOnline"`
	TrainingLocation string `json:"lieuFormation"`
	SubjectTitle     string `json:"matiereIntitule"`
	RemainingPlaces  int    `json:"nbPlacesRestantes"`
	TrainingCity     string `json:"villeFormation"`
}

type ByID []Course

func (a ByID) Len() int      { return len(a) }
func (a ByID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool {
	if a[i].ModuleID != a[j].ModuleID {
		return a[i].ModuleID < a[j].ModuleID
	}
	if a[i].CourseID != a[j].CourseID {
		return a[i].CourseID < a[j].CourseID
	}
	return a[i].CourseInstanceID < a[j].CourseInstanceID
}

func (c Course) String() string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("Available seats:\t%d\n", c.RemainingPlaces))
	builder.WriteString(fmt.Sprintf("Course:\t\t\t%s %s\n", c.ID, c.Title))
	builder.WriteString(fmt.Sprintf("Date:\t\t\t%s at %s\n", c.CourseDate, c.CourseSchedule))
	builder.WriteString(fmt.Sprintf("Instructor:\t\t%s\n", c.Instructor))
	builder.WriteString(fmt.Sprintf("Training site:\t\t%s\n", c.TrainingLocation))

	return builder.String()

	// s, _ := json.MarshalIndent(c, "", "\t")
	// return string(s)
}

func fetchCourses() ([]Course, error) {
	resp, err := http.Get("https://ssl.education.lu/gicea-wsrest/nat/getCoursList")
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var courses []Course
	if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
		return nil, err
	}

	return courses, nil
}

func preprocessCourses(courses []Course) []Course {
	for i, c := range courses {
		ids := strings.Split(c.ID, "-")
		courses[i].ModuleID = ids[0]
		courses[i].CourseID = ids[1]
		courseInstanceID, err := strconv.Atoi(ids[2])
		if err != nil {
			log.Fatalf("Error parsing ID: %s", err)
		}
		courses[i].CourseInstanceID = courseInstanceID
	}
	sort.Sort(ByID(courses))

	selectedCourses := []Course{}
	for _, c := range courses {
		if c.LanguageCode == "EN" {
			selectedCourses = append(selectedCourses, c)
		}
	}
	return selectedCourses
}

func main() {
	fmt.Println(strings.Repeat("-", 80))

	timestamp := time.Now().Format(time.DateTime)
	fmt.Println("Timestamp:", timestamp)

	shouldShortenOutput := slices.Contains(os.Args, "--short")
	shouldSave := !slices.Contains(os.Args, "--no-save")

	courses, err := fetchCourses()
	if err != nil {
		log.Fatalf("Error fetching courses: %s", err)
		return
	}
	courses = preprocessCourses(courses)

	if shouldSave {
		saveCourses(DATA_DIR+timestamp+".json", courses)

		hasChanges, err := doesNewFileHaveChanges()
		if err != nil {
			log.Fatalf("Error checking for changes: %s", err)
			return
		}
		if !hasChanges {
			fmt.Println("Summary: No changes")
			// return
		}
	}

	remainingPlaces := make(map[string]int)
	for _, c := range courses {
		predicate := true
		// predicate = predicate && c.IsOnline
		// predicate = predicate && c.LanguageCode == "EN"
		predicate = predicate && c.RemainingPlaces > 0
		// predicate = predicate && c.TrainingLocation == "DISTANCE DISTANCE"
		// predicate = predicate && c.TrainingCity == "DISTANCE"
		if predicate {
			remainingPlaces[c.TrainingCity] += 1
			if !shouldShortenOutput {
				fmt.Print(c.String())
				fmt.Println(strings.Repeat("-", 80))
			}
		}
	}

	fmt.Println("Remaining places:")
	for city, n := range remainingPlaces {
		fmt.Printf("\t%s: %d\n", city, n)
	}

	if remainingPlaces["DISTANCE"] > 0 {
		mack.Tell("Messages", `send "There are distance learning courses. https://ssl.education.lu/ve-portal" to buddy "`+PHONE_NUMBER+`"`)
		mack.Notify("There are distance learning courses on Vivre Ensemble", "Check the website for details.")
		// } else {
		// 	mack.Notify("There are updates to Vivre Ensemble", "Check the website for details.")
	}
}
