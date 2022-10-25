package go_scstjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ScstHandlerDriver string
type ScstTargetDriverType string

const (
	SCST_VDISK_BLOCKIO       ScstHandlerDriver    = "vdisk_blockio"
	SCST_TARGET_DRIVER_ISCSI ScstTargetDriverType = "iscsi"
)

type ScstDevice struct {
	Name       string `json:"name"`
	Filename   string `json:"filename"`
	NvCache    bool   `json:"nv_cache"`
	Rotational bool   `json:"rotational"`
	Blocksize  int    `json:"blocksize,omitempty"`
	T10VendId  string `json:"t10_vend_id"`
}

type ScstHandler struct {
	Driver  ScstHandlerDriver `json:"driver"`
	Devices []ScstDevice      `json:"devices"`
}

type ScstIscsiTargetDriver struct {
	Driver  ScstTargetDriverType `json:"driver"`
	Enabled bool                 `json:"enabled"`
	Targets []ScstIscsiTarget    `json:"targets"`
}

type ScstTargetAcl struct {
	Name       string             `json:"name"`
	Luns       []ScstTargetAclLun `json:"luns"`
	Initiators []string           `json:"initiators"`
}

type ScstTargetAclLun struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ScstIscsiTarget struct {
	Name           string          `json:"name"`
	Enabled        bool            `json:"enabled"`
	QueuedCommands int             `json:"queued_commands"`
	AllowedPortal  string          `json:"allowed_portal"`
	Groups         []ScstTargetAcl `json:"groups"`
}
type ScstConfig struct {
	Handlers      []ScstHandler           `json:"handlers"`
	TargetDrivers []ScstIscsiTargetDriver `json:"target_drivers"`
}

func (d *ScstDevice) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	temp = append(temp, fmt.Sprintf("%sDEVICE %s {", ind, d.Name))
	temp = append(temp, fmt.Sprintf("%s filename %s", ind+ind, d.Filename))
	if d.NvCache {
		temp = append(temp, fmt.Sprintf("%s nv_cache %d", ind+ind, 1))
	} else {
		temp = append(temp, fmt.Sprintf("%s nv_cache %d", ind+ind, 0))
	}
	if d.Rotational {
		temp = append(temp, fmt.Sprintf("%s rotational %d", ind+ind, 1))
	} else {
		temp = append(temp, fmt.Sprintf("%s rotational %d", ind+ind, 0))
	}
	if d.Blocksize > 0 {
		temp = append(temp, fmt.Sprintf("%s blocksize %d", ind+ind, d.Blocksize))
	}
	temp = append(temp, fmt.Sprintf("%s t10_vend_id %s", ind+ind, d.T10VendId))
	temp = append(temp, fmt.Sprintf("%s}\n", ind))
	res = strings.Join(temp, "\n")
	return
}

func (h *ScstHandler) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	temp = append(temp, fmt.Sprintf("HANDLER %s {", h.Driver))
	for _, d := range h.Devices {
		temp = append(temp, d.ToString(ind))
	}
	temp = append(temp, "}\n")
	res = strings.Join(temp, "\n")
	return
}

func (t *ScstIscsiTarget) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	temp = append(temp, fmt.Sprintf("%sTARGET %s {", ind, t.Name))
	if t.Enabled {
		temp = append(temp, fmt.Sprintf("%senabled %d", ind+ind, 1))
	} else {
		temp = append(temp, fmt.Sprintf("%senabled %d", ind+ind, 0))
	}
	temp = append(temp, fmt.Sprintf("%sQueuedCommands %d", ind+ind, t.QueuedCommands))
	temp = append(temp, fmt.Sprintf("%sallowed_portal %s", ind+ind, t.AllowedPortal))
	for _, g := range t.Groups {
		temp = append(temp, fmt.Sprintf("%s%s", ind, g.ToString(ind)))
	}
	temp = append(temp, fmt.Sprintf("%s}\n", ind))
	res = strings.Join(temp, "\n")
	return
}

func (l *ScstTargetAclLun) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	res = fmt.Sprintf("%sLUN %d %s\n", ind, l.Id, l.Name)
	return
}

func (a *ScstTargetAcl) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	temp = append(temp, fmt.Sprintf("%sGROUP %s {", ind, a.Name))
	for _, l := range a.Luns {
		temp = append(temp, fmt.Sprintf("%s%s", ind, l.ToString(ind+ind)))
	}
	for _, i := range a.Initiators {
		temp = append(temp, fmt.Sprintf("%sINITIATOR %s", ind+ind+ind, i))
	}
	temp = append(temp, fmt.Sprintf("%s}\n", ind))
	res = strings.Join(temp, "\n")
	return
}

func (d *ScstIscsiTargetDriver) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	temp = append(temp, fmt.Sprintf("TARGET_DRIVER %s {", d.Driver))
	if d.Enabled {
		temp = append(temp, fmt.Sprintf("%s enabled 1", ind))
	} else {
		temp = append(temp, fmt.Sprintf("%s enabled 0", ind))
	}
	for _, t := range d.Targets {
		temp = append(temp, t.ToString(ind))
	}
	temp = append(temp, "}")
	res = strings.Join(temp, "\n")
	return
}

func (c *ScstConfig) ToString(indent ...string) (res string) {
	ind := strings.Join(indent, "")
	temp := []string{}
	for _, h := range c.Handlers {
		temp = append(temp, h.ToString(ind))
	}
	for _, d := range c.TargetDrivers {
		temp = append(temp, d.ToString(ind))
	}
	res = strings.Join(temp, "\n")
	return
}

func NewScstConfFromFile(filePath string) (conf *ScstConfig, err error) {
	var (
		jsonData []byte
	)
	if jsonData, err = os.ReadFile(filePath); err != nil {
		err = fmt.Errorf("unable to read scst json config: %w", err)
	} else {
		if err = json.Unmarshal(jsonData, conf); err != nil {
			err = fmt.Errorf("unable to unmarshal scst json data: %w", err)
		}
	}

	return
}
