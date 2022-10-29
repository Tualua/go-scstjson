package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	scst "github.com/Tualua/go-scstjson"
)

func main() {
	ScstConfig := scst.ScstConfig{}
	if scstConfFile, err := os.ReadFile("/etc/scst.conf"); err != nil {
		log.Fatalln("Unable to open SCST config")
	} else {
		scstConf := string(scstConfFile)
		//Check config
		if strings.Contains(scstConf, "HANDLER") {
			if strings.Contains(scstConf, "TARGET_DRIVER") {
				hndBlockAddr := strings.Index(scstConf, "HANDLER")
				tgtDrvBlockAddr := strings.Index(scstConf, "TARGET_DRIVER")
				hndBlock := ""
				tgtDrvBlock := ""
				if hndBlockAddr < tgtDrvBlockAddr {
					hndBlock = scstConf[hndBlockAddr:tgtDrvBlockAddr]
					tgtDrvBlock = scstConf[tgtDrvBlockAddr:]
				} else {
					hndBlock = scstConf[tgtDrvBlockAddr:hndBlockAddr]
					tgtDrvBlock = scstConf[hndBlockAddr:]
				}

				hndBlock = strings.ReplaceAll(hndBlock, "\t", "")
				tgtDrvBlock = strings.ReplaceAll(tgtDrvBlock, "\t", "")

				block := scst.ScstHandler{}
				hndBlockSplit := strings.Split(hndBlock, "DEVICE")
				block.Driver = scst.SCST_VDISK_BLOCKIO
				for _, d := range hndBlockSplit[1:] {
					dev := scst.ScstDevice{}
					data := strings.Split(d, "\n")
					dev.Name = strings.TrimSpace(strings.Split(data[0], "{")[0])
					for _, s := range data[1:] {
						s = strings.TrimSpace(s)
						if strings.Contains(s, "filename") {
							dev.Filename = strings.Split(s, " ")[1]
						}
						if strings.Contains(s, "nv_cache") {
							if strings.Split(s, " ")[1] == "0" {
								dev.NvCache = false
							} else {
								dev.NvCache = true
							}

						}
						if strings.Contains(s, "rotational") {
							if strings.Split(s, " ")[1] == "0" {
								dev.Rotational = false
							} else {
								dev.Rotational = true
							}

						}
						if strings.Contains(s, "t10_vend_id") {
							dev.T10VendId = strings.Split(s, " ")[1]
						}
						if strings.Contains(s, "blocksize") {
							bs, _ := strconv.Atoi(strings.Split(s, " ")[1])
							dev.Blocksize = bs
						}
					}
					//					log.Println(data)
					block.Devices = append(block.Devices, dev)
				}
				ScstConfig.Handlers = append(ScstConfig.Handlers, block)
				tgtDrvBlockSplit := strings.Split(tgtDrvBlock, "TARGET ")
				driver := scst.ScstIscsiTargetDriver{}
				driver.Driver = scst.SCST_TARGET_DRIVER_ISCSI
				driver.Enabled = true
				for _, t := range tgtDrvBlockSplit[1:] {
					data := strings.Split(t, "\n")
					tgt := scst.ScstIscsiTarget{}
					tgt.Name = strings.Split(data[0], " ")[0]
					for _, s := range data {
						s = strings.TrimSpace(s)
						if strings.Contains(s, "QueuedCommands") {
							qc, _ := strconv.Atoi(strings.Split(s, " ")[1])
							tgt.QueuedCommands = qc
						}
						if strings.Contains(s, "allowed_portal") {
							tgt.AllowedPortal = strings.Split(s, " ")[1]
						}
						if strings.Contains(s, "enabled") {
							if strings.Split(s, " ")[1] == "0" {
								tgt.Enabled = false
							} else {
								tgt.Enabled = true
							}
						}
						if strings.Contains(s, "GROUP") {
							tgt.Groups = append(tgt.Groups, scst.ScstTargetAcl{Name: "allowed_ini"})
						}

						if strings.Contains(s, "LUN") {
							tgt.Groups[0].Luns = append(tgt.Groups[0].Luns, scst.ScstTargetAclLun{Id: 0, Name: strings.Split(s, " ")[2]})
						}

						if strings.Contains(s, "INITIATOR") {
							tgt.Groups[0].Initiators = append(tgt.Groups[0].Initiators, strings.Split(s, " ")[1])
						}

					}
					driver.Targets = append(driver.Targets, tgt)

				}
				ScstConfig.TargetDrivers = append(ScstConfig.TargetDrivers, driver)
			} else {
				log.Fatalln("Missing TARGET_DRIVER block")
			}
		} else {
			log.Fatalln("Missing HANDLER block")
		}

	}

	j, _ := json.MarshalIndent(ScstConfig, "", "    ")
	os.WriteFile("scst.json", j, 0644)

}
