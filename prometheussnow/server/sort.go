package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/alertmanager/template"
)

// type to sort alert by StartsAt time (implement sort interface)
type TimeSortedAlerts []template.Alert

/****needs to chage the function namings 1. NewTimeSortedAlerts, 2.TimeSortedAlerts*/
func NewTimeSortedAlerts(alerts []template.Alert) TimeSortedAlerts {
	result := make(TimeSortedAlerts, 0)
	for _, alert := range alerts {
		result = append(result, alert)
	}
	return result
}

func (a TimeSortedAlerts) Len() int {
	return len(a)
}

func (a TimeSortedAlerts) Less(i, j int) bool {
	return a[i].StartsAt.Before(a[j].StartsAt)
}

func (a TimeSortedAlerts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// check if the passed-in map contains the specified key and returns its value
func keyValue(values template.KV, key string) (error, string) {
	value := ""
	found := false
	// loop through the alert values
	for k, v := range values {
		// if it finds the required key
		if strings.ToLower(key) == strings.ToLower(k) {
			// assign the value of the key
			value = v
			// set the found flag
			found = true
			// exit the loop
			break
		}
	}
	// if a key was found but with no value
	if found && len(value) == 0 {
		return errors.New(fmt.Sprintf("annotation '%s' has no value, check the Prometheus rule has the correct label", key)), value
	}
	// if the key was not found
	if !found {
		return errors.New(fmt.Sprintf("cannot find '%s' annotation in alert", key)), value
	}
	// we have a value
	return nil, value
}

// get the service unique natural key
func uniqueKey(platform string, service string, facet string, status string, location string) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", platform, service, facet, status, strings.Replace(strings.Replace(location, ":", "_", -1), ".", "_", -1))
}
