package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseUrl string = "https://swp.webspeiseplan.de/index.php"

func getMenuData() (string, error) {
	params := url.Values{}
	params.Add("token", "55ed21609e26bbf68ba2b19390bf7961")
	params.Add("model", "menu")
	params.Add("location", "9601")
	params.Add("languagetype", "2")
	params.Add("_", fmt.Sprintf("%d", time.Now().UnixMilli()))

	req, err := http.NewRequest(http.MethodGet, baseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Referer", "https://sqp.webspeiseplan.de/Menu")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(resBody), nil
}

func main() {
	fmt.Println(getMenuData())
}
