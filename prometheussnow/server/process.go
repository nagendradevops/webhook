package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/prometheus/alertmanager/template"
	"github.com/rs/zerolog/log"
)

type Snow struct {
	conf *Config
}

//Create Incident function
func createIncident(data map[string]interface{}, snowurl string, snowusername string, snowpassword string) error {
	tableEndpoint := "/api/now/v1/table/incident"

	// Convert data to JSON format
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", snowurl+tableEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(snowusername, snowpassword)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// process the received alerts
func processAlerts(data template.Alerts) error {
	snowurl := "https://dev183753.service-now.com/"
	snowusername := "aes.creator"
	snowpassword := "@q4RXR2iiK!v"
	// sort the incoming alerts by StartsAt time
	alerts := NewTimeSortedAlerts(data)
	sort.Sort(alerts)
	log.Debug().Msgf("* alerts have been sorted")

	// loops through tha alerts
	for _, alert := range alerts {
		// extract values from the alert
		v, err := values(alert)
		// if the alert did not have all required information
		if err != nil {
			// stops any processing
			return err
		}
		// write debug info about successful data extraction
		log.Debug().Msgf("* extracted values for alert with fingerprint '%s'", alert.Fingerprint)
		log.Debug().Msgf("* platform value = '%s'", v["platform"])
		log.Debug().Msgf("* service value = '%s'", v["service"])
		log.Debug().Msgf("* facet value = '%s'", v["facet"])

		location := v["location"]
		// if the alert does not have a specific location
		if len(location) == 0 {
			// fill the blank
			location = "_"
		}
		// build the natural key for the service item
		serviceKey := uniqueKey(v["platform"], v["service"], v["facet"], v["status"], location)
		log.Debug().Msgf("* service key '%s' created", serviceKey)

		data := map[string]interface{}{
			"Urgency":           "3-Low",
			"Caller":            "Creator User",
			"short_description": fmt.Sprintf("recording event for => %s:%s:%s:%s:%s", serviceKey, v["service"], v["facet"], v["status"], location),
		}

		err = createIncident(data, snowurl, snowusername, snowpassword)
		if err != nil {
			fmt.Println("Error creating incident in ServiceNow:", err)
			return err
		}

	}
	return nil
}

// extract values from the alert
func values(alert template.Alert) (values map[string]string, err error) {
	var result = make(map[string]string)

	err, result["platform"] = keyValue(alert.Labels, "platform")
	if err != nil {
		return result, err
	}
	err, result["service"] = keyValue(alert.Labels, "service")
	if err != nil {
		return result, err
	}
	err, result["status"] = keyValue(alert.Labels, "status")
	if err != nil {
		return result, err
	}
	err, result["description"] = keyValue(alert.Labels, "description")
	if err != nil {
		return result, err
	}
	err, result["facet"] = keyValue(alert.Labels, "facet")
	if err != nil {
		return result, err
	}
	// add any annotations
	for key, value := range alert.Annotations {
		// if the annotation key is not a label (e.g. not in the result map already)
		if result[key] == "" {
			// add the annotation value
			result[key] = value
		} else {
			// skip and issue a warning
			log.Warn().Msgf("skipping annotation '%s' as a label with such name was found", key)
		}
	}
	return result, nil
}
