package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	Roblox "github.com/piratepeep/ugc-limited-sniper/roblox"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Roblox struct {
		Accounts []Roblox.Config `toml:"accounts"`
	} `toml:"roblox"`
	ItemIDs struct {
		IDs []int `toml:"ids"`
	} `toml:"item_ids"`
}

var config Config

func main() {
	configFile := "config.toml"

	configTree, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	err = configTree.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	stop := make(chan struct{})
	defer close(stop)

	for i, accountConfig := range config.Roblox.Accounts {
		r := Roblox.New(accountConfig)

		// Get the user info and update the struct property
		userInfo, err := r.GetUserInfo()

		if err != nil {
			fmt.Printf("Cookie #%v is invalid!\n", i+1)
			continue
		}

		r.SetCurrentUser(userInfo.ID, userInfo.Name)

		// Get the XSRF token before we start doing anything
		token, err := r.RefreshXsrfToken()
		if err != nil {
			fmt.Printf("Failed to get the XSRF token for: %s", userInfo.Name)
		}

		r.SetXsrfToken(token)

		go monitorAndBuyItems(r, config.ItemIDs.IDs, stop)
		go updateXsrfToken(r, stop)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}

func monitorAndBuyItems(r *Roblox.Roblox, itemIDs []int, stop <-chan struct{}) {
	// This is the requests per minute that the program can make without proxies. This is a ratelimit imposed by the Roblox API.
	// Do not change this unless you know what you are doing.
	rpm := 20 / len(config.Roblox.Accounts)
	ticker := time.NewTicker(time.Minute / time.Duration(rpm))
	messagePrinted := false
	for {
		select {
		case <-ticker.C:
			assetsInfo, err := r.GetAssetsInfo(itemIDs)
			if err != nil {
				log.Printf("Failed to get asset info: %v", err)
				continue
			}

			if messagePrinted != true {
				userID, username := r.GetCurrentUser()
				itemNamesList := make([]string, 0, len(assetsInfo))
				for _, assetInfo := range assetsInfo {
					itemNamesList = append(itemNamesList, assetInfo.Name)
				}

				log.Printf("Monitoring items for user %s (%d): %s", username, userID, strings.Join(itemNamesList, ", "))

				messagePrinted = true
			}

			for _, assetInfo := range assetsInfo {
				// This property may or may not be reliable.
				if assetInfo.UnitsAvailableForConsumption > 0 {
					// Just be an asshole and buy everything.
					for i := 0; i < 5; i++ {
						go func(assetInfo Roblox.POSTItemDetails) {
							// Attempt to buy the item using the available information
							response, err := r.BuyItem(assetInfo.ID, assetInfo.CollectibleItemId, assetInfo.CreatorTargetId)
							if err != nil {
								log.Printf("Failed to buy item %d: %v", assetInfo.ID, err)
							} else {
								log.Printf("Successfully bought item %s", response.PurchaseResult)
							}
						}(assetInfo)
					}
				}
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func updateXsrfToken(r *Roblox.Roblox, stop <-chan struct{}) {
	// Get the new XSRF token every 30 seconds so we don't need to send an additional request when trying to purchase an item.
	// This saves the program a small amount of time which can be crucial depending on the quantity of items available.
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			token, err := r.RefreshXsrfToken()
			if err != nil {
				log.Printf("Failed to update XSRF token: %v", err)
				continue
			}
			r.SetXsrfToken(token)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
