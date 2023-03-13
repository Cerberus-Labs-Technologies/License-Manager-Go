package main

import (
	"fmt"
	"github.com/Cerberus-Labs-Technologies/License-Manager-Go/license"
	"github.com/Cerberus-Labs-Technologies/License-Manager-Go/rest"
	"log"
	"strconv"
	"strings"
	"time"
)

type LicenseManager struct {
	*license.Config
	License license.License
}

func (lm *LicenseManager) generateLogList(validity bool) []string {
	logs := []string{"",
		"LicenseManager v1.0.0",
		"Made by Cerberus-Labs.tech",
		"Authored by Kelvin Bill",
		"License Key: " + lm.Config.LicenseKey,
		"Product ID: " + strconv.Itoa(lm.Config.ProductId),
	}
	licensedToMsg := "Licensed to: "
	if validity {
		licensedToMsg += lm.License.Username
	} else {
		licensedToMsg += strconv.Itoa(lm.Config.UserId)
	}
	logs = append(logs, licensedToMsg)
	validityMsg := "License is: "
	if validity {
		validityMsg += "Valid"
	} else {
		validityMsg += "Invalid"
	}

	logs = append(logs, validityMsg)
	return logs
}

func (lm *LicenseManager) printLogs(logs []string) {
	printRectangle(logs)
}

func (lm *LicenseManager) validate(validCallback func(), nonValidCallback func()) {
	validity := lm.isValid()

	if validity {
		validCallback()
	} else {
		nonValidCallback()
	}

	logs := lm.generateLogList(validity)
	printRectangle(logs)
	lm.startLicenseChecker(nonValidCallback)
}

func (lm *LicenseManager) startLicenseChecker(nonValidCallback func()) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if !lm.isValid() {
					nonValidCallback()
				}
			}
		}
	}()
}

func (lm *LicenseManager) isValid() bool {
	licenseServerUrl := "https://backend.cerberus-labs.tech/api/v1/license/" + lm.Config.LicenseKey
	l, err := rest.MakeHttpGetRequest(licenseServerUrl)
	if err != nil {
		return false
	}
	lm.License = l
	log.Println("License: " + lm.License.License)
	if lm.License.ProductId != lm.Config.ProductId {
		return false
	}
	if lm.License.UserId != lm.Config.UserId {
		return false
	}
	if !lm.License.Active {
		return false
	}
	if !lm.License.Permanent && lm.License.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

func printRectangle(texts []string) {
	maxLen := 0
	for _, text := range texts {
		if len(text) > maxLen {
			maxLen = len(text)
		}
	}
	maxLen += 14
	fmt.Println(strings.Repeat("#", maxLen+2))
	for i, text := range texts {
		paddingSize := maxLen - len(text) - 2
		paddingStart := paddingSize / 2
		paddingEnd := paddingSize - paddingStart
		switch i {
		case 0:
			fmt.Printf("# %s%s%s #\n", strings.Repeat(" ", paddingStart), text, strings.Repeat(" ", paddingEnd))
		case len(texts) - 1:
			fmt.Printf("#%s%s%s #\n", strings.Repeat(" ", paddingStart), text, strings.Repeat(" ", paddingEnd+1))
		default:
			fmt.Printf("# %s%s%s #\n", strings.Repeat(" ", paddingStart), text, strings.Repeat(" ", paddingEnd))
		}
	}
	fmt.Println(strings.Repeat("#", maxLen+2))
}
