package metrics

import (
	"encoding/json"
	"log"
	"math"
	"os/exec"
	"strings"
)

type uapiResponse struct {
	Result struct {
		Data struct {
			InodeLimit    interface{} `json:"inode_limit"`
			MegabytesUsed float64     `json:"megabytes_used"`
			MegabyteLimit interface{} `json:"megabyte_limit"`
			InodesUsed    float64     `json:"inodes_used"`
		} `json:"data"`
	} `json:"result"`
}

func getQuota(user string) (usedMegabyte float64, percentMegabyte float64, usedInode float64, percentInode float64) {
	out := command(strings.TrimSpace(user), "Quota", "get_quota_info")
	var resp uapiResponse
	err := json.Unmarshal(out, &resp)
	if err != nil {
		log.Println("grt quota error:", err, string(out))
		return
	}

	usedMegabyte = resp.Result.Data.MegabytesUsed
	limitMegabyte := parseStringFloat(resp.Result.Data.MegabyteLimit)
	usedInode = resp.Result.Data.InodesUsed
	limitInode := parseStringFloat(resp.Result.Data.InodeLimit)

	if limitMegabyte > 0 {
		percentMegabyte = math.Round((usedMegabyte / limitMegabyte) * 100)
	}
	if limitInode > 0 {
		percentInode = math.Round((usedInode / limitInode) * 100)
	}
	return
}

func command(user string, commands ...string) (out []byte) {
	com := []string{"--user=" + user, "--output=json"}
	com = append(com, commands...)
	out, err := exec.Command("/usr/bin/uapi", com...).CombinedOutput()
	if err != nil {
		log.Println("execute command error:", err)
	}
	return
}
