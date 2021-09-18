package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"io"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func main() {

	title,stickers := getStickerTitleAndUrl(GetLineStickerURL())

	cmd := exec.Command("mkdir" , GetCurPath() + "/" + title)
	err := cmd.Run()

	if(err != nil){
		log.Fatal(err)
	}

	for i,v := range stickers{
		DownloadFile(fmt.Sprintf("%s/%d.png",title,i),v)
	}
}

func GetCurPath() string {
	cur,_ := os.Getwd()
	return cur
}

func GetLineStickerURL() string {
	fmt.Print("Enter Line Sticker URL: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return ""
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	return input

}

func getStickerTitleAndUrl(url string) (title string, stickers []string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	title = doc.Find(".mdCMN38Item01Txt").Text()

	// Find the review items
	doc.Find(".mdCMN09LiInner").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		url,_ := s.Find(".mdCMN09Image").Attr("style")
		url = strings.ReplaceAll(url,"background-image:url(","")
		url = strings.ReplaceAll(url,";compress=true);","")
		stickers = append(stickers, url)
	})

	return title,stickers
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
