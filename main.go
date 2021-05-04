package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DataPreview struct {
	Id           string `json:"id"`
	StickerType	 string `json:"type"`
	PopupUrl     string `json:"popupUrl"`
	StaticUrl string `json:"staticUrl"`
	AnimationUrl string `json:"animationUrl"`
	SoundUrl     string `json:"soundUrl"`
}

// Function to check if path exist, else create it
func checkPath(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, 0755); if err != nil {
			return err
		}
	}

	return nil
}

// Function to extract required data into a struct from a HTML string
func stickerCodeExtractor(inputHtml string) ([]DataPreview, error) {
	var rgx = regexp.MustCompile(`(data-preview='.*?')`)
	tmpExtracted := rgx.FindAllStringSubmatch(inputHtml, -1)

	output := make([]DataPreview, len(tmpExtracted))

	for i := 0; i < len(tmpExtracted); i++ {
		res := tmpExtracted[i][1]

		// Remove data-preview=' ... '
		res = strings.Split(res, "'")[1]
		res = strings.ReplaceAll(res, "&quot;", "\"")
		res = strings.ReplaceAll(res, " ", "")

		var dp DataPreview
		err := json.Unmarshal([]byte(res), &dp); if err != nil {
			return nil, err
		}

		output[i] = dp
	}

	return output, nil
}

// Function to download image based on either PopupUrl or AnimationUrl
func downloadImage(dp DataPreview, wg *sync.WaitGroup) error {

	var getUrl string
	switch dp.StickerType {
	case "animation", "animation_sound":
		getUrl = dp.AnimationUrl
	case "static":
		getUrl = dp.StaticUrl
	case "popup", "popup_sound":
		getUrl = dp.PopupUrl
	default:
		getUrl = ""
	}

	if getUrl != "" {
		fmt.Println("Downloading: ", getUrl)

		dt := time.Now()
		//Format YYYYMMDDhhmm
		filename := dt.Format("200601021504") + "-" + dp.Id

		downloadResp, err := http.Get(getUrl); if err != nil {
			return err
		}

		defer downloadResp.Body.Close()

		// Create the file
		out, err := os.Create("./output/" + filename + ".png"); if err != nil {
			return err
		}

		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, downloadResp.Body); if err != nil {
			return err
		}

		// Execute conversion to GIF if it is not static
		if dp.StickerType != "static" {
			err = exec.Command("./bin/apng2gif.exe",
				"./output/" + filename + ".png", "./output-gif/" + filename + ".gif").Run(); if err != nil {
				return err
			}
		}
	}

	defer wg.Done()
	return nil
}

// Scrap and download stickers
func scrap(scrapUrl string) error {
	// Ensure directory exists
	err := checkPath("output"); if err != nil {
		return err
	}
	err = checkPath("output-gif"); if err != nil {
		return err
	}

	resp, err := http.Get(scrapUrl); if err != nil {
		return err
	}

	defer resp.Body.Close()

	respHtml, err := ioutil.ReadAll(resp.Body); if err != nil {
		return err
	}

	result, err := stickerCodeExtractor(string(respHtml)); if err != nil {
		return err
	}

	fmt.Println("Sticker type: ", result[0].StickerType)

	// Wait for download to complete
	var wg sync.WaitGroup
	for i := 0; i < len(result); i++ {
		wg.Add(1)
		go downloadImage(result[i], &wg)
	}

	wg.Wait()
	return nil
}

// Entrypoint
func main(){
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter Line Stickershop URL")
		inputUrl, err := consoleReader.ReadString('\n'); if err != nil {
			log.Fatal(err)
		}

		// Check if input has at least a line store format
		if strings.Contains(inputUrl, "https://store.line.me") {
			inputUrl = strings.Replace(inputUrl, "\r\n", "", -1)
			err := scrap(inputUrl); if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Invalid format")
		}
	}
}
