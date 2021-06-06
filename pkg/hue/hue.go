package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/lucasb-eyer/go-colorful"
)

func ColorToHueHsb(color colorful.Color) (int, int, int) {
	h, s, l := color.Hsl()
	return int(math.Floor(65535 * h / 360)),
		int(math.Floor(s * 255)),
		int(math.Floor(l * 255))
}

func GetLights(cfg Config) (map[string]Light, error) {
	response, err := http.Get(fmt.Sprintf("%s/api/%s/lights", cfg.Endpoint, cfg.Username))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}
	result := map[string]Light{}
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, response.Body); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ConfigureLight(cfg Config, id string, state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/api/%s/lights/%s/state", cfg.Endpoint, cfg.Username, id),
		bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		log.Error(resp.Status)
	}
	x, _ := io.ReadAll(resp.Body)
	fmt.Println(string(x))
	return nil
}

func AlertAllLights(cfg Config) error {
	lights, err := GetLights(cfg)
	if err != nil {
		return err
	}
	for k, v := range lights {
		v.State.Alert = "select"
		lights[k] = v
		if err := ConfigureLight(cfg, k, State{
			On:    true,
			Alert: "select",
		}); err != nil {
			return err
		}
	}
	return nil
}
