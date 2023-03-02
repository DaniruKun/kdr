package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const Notice = `
	kdr - Kerbal Debris Removal
	Copyright (C) 2023  DaniruKun
	This program comes with ABSOLUTELY NO WARRANTY.
	This is free software, and you are welcome to redistribute it
	under certain conditions.

`

func main() {
	fmt.Print(Notice)

	saveFilePath := flag.String("i", "", "save file JSON path")
	outputFilePath := flag.String("o", "save-patched.json", "output save file JSON")

	if *saveFilePath == "" {
		log.Fatal("missing required argument: inputFilePath")
	}

	saveFile, err := os.Open(*saveFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer saveFile.Close()

	saveData, _ := io.ReadAll(saveFile)
	var save map[string]any

	json.Unmarshal(saveData, &save)

	vessels := save["Vessels"].([]interface{})
	var nonDebrisVessels []map[string]interface{}

	for _, vessel := range vessels {
		vesselMap := vessel.(map[string]interface{})
		if !vesselMap["IsDebris"].(bool) {
			nonDebrisVessels = append(nonDebrisVessels, vesselMap)
		}
	}

	fmt.Printf("saved %d non-debris vessels", len(nonDebrisVessels))

	// for _, nonDebris := range nonDebrisVessels {
	// 	fmt.Println(nonDebris["AssemblyDefinition"].(map[string]any)["assemblyName"])
	// }

	save["Vessels"] = nonDebrisVessels
	outData, _ := json.Marshal(save)

	err = os.WriteFile(*outputFilePath, outData, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func appDataDir() (dir string) {
	dir, _ = os.UserCacheDir()
	return
}

// return filepath.Join(appDataDir(), "LocalLow", "Intercept Games", "Kerbal Space Program 2", "Saves", "SinglePlayer", "Sandboxxy", "autosave_1.json")
