package geocode

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/opaquee/EventMapAPI/graph/model"
)

const geo_api_url string = "http://api.positionstack.com"
const forward_geo string = "v1/forward"

type ResponseData struct {
	Data []ResponseDataEntry `json:"data,omitempty"`
}

type ResponseDataEntry struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

func GetLatLng(event *model.Event) error {
	baseURL, err := url.Parse(geo_api_url)
	if err != nil {
		return err
	}

	baseURL.Path += forward_geo

	params := url.Values{}
	params.Add("access_key", os.Getenv("GEO_API_KEY"))
	params.Add("query", event.AddressLine1+", "+event.City+", "+event.State+", "+string(event.Zip))
	params.Add("output", "json")
	params.Add("limit", "1")

	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	res_data := &ResponseData{}

	if err = decoder.Decode(res_data); err != nil {
		return err
	}

	event.Latitude = res_data.Data[0].Latitude
	event.Longitude = res_data.Data[0].Longitude

	return nil
}
