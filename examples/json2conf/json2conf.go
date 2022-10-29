package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	scst "github.com/Tualua/go-scstjson"
)

func main() {

	if ScstConf, err := scst.NewScstConfFromFile("scst.json"); err != nil {
		fmt.Println(err)
	} else {
		uniqMapDevs := make(map[string]scst.ScstDevice)
		uniqDevNames := []string{}
		fmt.Println(len(ScstConf.Handlers[0].Devices))
		for _, d := range ScstConf.Handlers[0].Devices {
			if val, ok := uniqMapDevs[d.Name]; ok {
				log.Println(val.Name)
			} else {
				uniqMapDevs[d.Name] = d
				uniqDevNames = append(uniqDevNames, d.Name)
			}
		}
		var uniqDevs []scst.ScstDevice
		//sort.Strings(uniqDevNames)
		sort.Slice(uniqDevNames, func(i, j int) bool {
			if strings.Contains(uniqDevNames[i], "-vm") {
				if strings.Contains(uniqDevNames[j], "-vm") {
					vmi, _ := strconv.Atoi(strings.Split(uniqDevNames[i], "-vm")[1])
					vmj, _ := strconv.Atoi(strings.Split(uniqDevNames[j], "-vm")[1])
					if vmi == vmj {
						return uniqDevNames[i] < uniqDevNames[j]
					} else {
						return vmi < vmj
					}

				} else {
					return uniqDevNames[i] < uniqDevNames[j]
				}
			} else {
				return uniqDevNames[i] < uniqDevNames[j]
			}
		})
		for _, v := range uniqDevNames {
			uniqDevs = append(uniqDevs, uniqMapDevs[v])
		}
		ScstConf.Handlers[0].Devices = uniqDevs
		fmt.Println(len(ScstConf.Handlers[0].Devices))

		uniqMapTargets := make(map[string]scst.ScstIscsiTarget)
		uniqTargetNames := []string{}
		fmt.Println(len(ScstConf.TargetDrivers[0].Targets))
		for _, d := range ScstConf.TargetDrivers[0].Targets {
			if val, ok := uniqMapTargets[d.Name]; ok {
				log.Println(val.Name)
			} else {
				uniqMapTargets[d.Name] = d
				uniqTargetNames = append(uniqTargetNames, d.Name)
			}
		}
		var uniqTargets []scst.ScstIscsiTarget
		sort.Slice(uniqTargetNames, func(i, j int) bool {
			if strings.Contains(uniqTargetNames[i], "-vm") {
				if strings.Contains(uniqTargetNames[j], "-vm") {
					vmi, _ := strconv.Atoi(strings.Split(uniqTargetNames[i], "-vm")[1])
					vmj, _ := strconv.Atoi(strings.Split(uniqTargetNames[j], "-vm")[1])
					if vmi == vmj {
						return uniqTargetNames[i] < uniqTargetNames[j]
					} else {
						return vmi < vmj
					}
				} else {
					return uniqTargetNames[i] < uniqTargetNames[j]
				}
			} else {
				return uniqTargetNames[i] < uniqTargetNames[j]
			}
		})
		for _, t := range uniqTargetNames {
			uniqTargets = append(uniqTargets, uniqMapTargets[t])
		}
		ScstConf.TargetDrivers[0].Targets = uniqTargets
		fmt.Println(len(ScstConf.TargetDrivers[0].Targets))

		//ScstConfText := ScstConf.ToString("    ")
		//fmt.Println(ScstConfText)
		os.WriteFile("scst.conf", []byte(ScstConf.ToString("    ")), 0644)
		jsonDedup, _ := json.MarshalIndent(ScstConf, "", "    ")
		os.WriteFile("scst_dedup.json", jsonDedup, 0644)
	}

}
