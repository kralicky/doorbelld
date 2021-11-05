package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gen2brain/beeep"
	"github.com/kralicky/doorbelld/pkg/hue"
	"github.com/kralicky/doorbelld/pkg/unifi"
)

func main() {
	viper.SetConfigName("doorbelld")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
	url, err := url.Parse(viper.GetString("unifi.endpoint"))
	if err != nil {
		log.Fatal(err)
	}
	cfg := unifi.Config{
		Endpoint: url,
	}
	hueCfg := hue.Config{
		Endpoint: viper.GetString("hue.endpoint"),
		Username: viper.GetString("hue.username"),
	}
	token, err := unifi.Login(cfg,
		viper.GetString("unifi.username"), viper.GetString("unifi.password"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		listener, err := unifi.NewUpdateListener(context.Background(), cfg, token)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		for {
			update, open := <-listener.C
			if !open {
				break
			}
			log.WithField("kind", update.ActionFrame.Action).
				WithField("key", update.ActionFrame.ModelKey).
				Debug(string(update.DataFrame))
			if update.ActionFrame.ModelKey == unifi.ModelKeyEvent &&
				update.ActionFrame.Action == unifi.ActionKindAdd {
				event := unifi.EventDataFrame{}
				err := json.Unmarshal(update.DataFrame, &event)
				if err != nil {
					log.WithError(err).Warn("Could not unmarshal event data frame")
					continue
				}
				if event.Type == "ring" {
					log.Info("** DOORBELL IS RINGING **")
					go func() {
						if err := beeep.Notify("Doorbell is ringing", "", ""); err != nil {
							log.Error(err)
						}
					}()
					go func() {
						if err := beeep.Beep(783.99, 450); err != nil {
							log.Error(err)
						}
						if err := beeep.Beep(659.25, 550); err != nil {
							log.Error(err)
						}
					}()
					for i := 0; i < 3; i++ {
						if err := hue.AlertAllLights(hueCfg); err != nil {
							log.Error(err)
						}
						time.Sleep(1 * time.Second)
					}
				}
			}
		}
	}
}
