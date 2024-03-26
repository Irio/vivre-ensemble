package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func saveCourses(filename string, courses []Course) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}
	defer file.Close()

	jsonBytes, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err)
	}
	_, err = file.Write(jsonBytes)
	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
	}

	return nil
}

func fileMD5(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
		return "", err
	}

	hash := md5.Sum(data)
	md5String := hex.EncodeToString(hash[:])

	return md5String, nil
}

func doesNewFileHaveChanges() (bool, error) {
	files, err := os.ReadDir(DATA_DIR)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	// only .gitkeep and new file
	if len(files) <= 2 {
		return true, nil
	}

	file1 := DATA_DIR + files[len(files)-1].Name()
	file2 := DATA_DIR + files[len(files)-2].Name()

	fmt.Println("Files:\n\t" + file1)
	fmt.Println("\t" + file2)

	file1Md5, err := fileMD5(file1)
	if err != nil {
		return false, err
	}
	file2Md5, err := fileMD5(file2)
	if err != nil {
		return false, err
	}

	return file1Md5 != file2Md5, nil
}
