package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=btc,eth"
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accepts", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", os.Getenv("API_KEY"))
	res, _ := netClient.Do(req)
	byt, _ := ioutil.ReadAll(res.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(byt, &data); err != nil {
		panic(err)
	}
	defer res.Body.Close()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Symbol", "Price"})
	for symbol, v := range data["data"].(map[string]interface{}) {
		for kk, vv := range v.(map[string]interface{}) {
			var price float64
			if vv != nil {
				switch kk {
				case "quote":
					quote := vv.(map[string]interface{})["USD"].(map[string]interface{})
					price = quote["price"].(float64)
					t.AppendRows([]table.Row{
						{symbol, fmt.Sprintf("%f", price)},
					})
					t.AppendSeparator()
				}
			}
		}
	}
	t.Render()
}
