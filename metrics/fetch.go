package metrics

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func cpanelVersion() string {
	out, err := exec.Command("/usr/local/cpanel/cpanel", "-V").CombinedOutput()
	if err != nil {
		log.Println(err)
		return ""
	}
	return strings.Trim(string(out), "\n")
}

func getUsers() []string {
	files := getFilesInDir("/var/cpanel/users")
	return files
}

func getUsersCount(onlySuspended bool) int {
	files := getUsers()
	if onlySuspended {
		return len(matchFilesLine(files, "SUSPENDED=1"))
	}
	return len(files)
}

func getBandwidth(user string) (bw int) {
	var lines []string
	file, err := os.Open("/var/cpanel/bandwidth.cache/" + user)
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txty := scanner.Text()
		lines = append(lines, txty)
	}

	file.Close()
	out := strings.Join(lines, "\n")
	bw, _ = strconv.Atoi(out)
	return
}

func getFTP() (lines []string) {
	file, err := os.Open("/etc/proftpd/passwd.vhosts")
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txty := scanner.Text()
		parts := strings.Split(txty, ":")
		if len(parts) > 0 {
			lines = append(lines, parts[0])
		}
	}
	file.Close()
	return
}

func getPlans() map[string]int {
	var plans = make(map[string]int)
	files := getFilesInDir("/var/cpanel/users")
	matches := matchFilesLine(files, "PLAN=.*")
	for _, m := range matches {
		parts := strings.Split(m, "=")
		if len(parts) > 0 {
			plans[parts[1]]++
		}
	}

	return plans
}

func getSessions() (web int, email int) {
	for _, f := range getFilesInDir("/var/cpanel/sessions/raw") {
		if strings.Contains(f, "@") {
			email++
		} else {
			web++
		}
	}
	return
}

func getRelease() string {
	file, err := os.Open("/etc/cpupdate.conf")
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txty := scanner.Text()
		if strings.Contains(txty, "CPANEL=") {
			parts := strings.Split(txty, "=")
			if len(parts) > 0 {
				return parts[1]
			}
		}
	}

	return ""
}

func getDomains() (domains []string) {
	file, err := os.Open("/etc/userdomains")
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txty := scanner.Text()
		parts := strings.Split(txty, ":")
		if len(parts) > 1 {
			domains = append(domains, parts[0])
		}
	}

	file.Close()
	return
}

func getLicenseInfo() (expireTime string, maxUsers int) {
	file, err := os.Open("/usr/local/cpanel/cpanel.lisc")
	if err != nil {
		log.Printf("failed opening file: %s\n", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txty := scanner.Text()
		if strings.Contains(txty, "license_expire_time:") {
			parts := strings.Split(txty, " ")
			if len(parts) > 0 {
				expireTime = parts[1]
			}
		} else if strings.Contains(txty, "maxusers:") {
			parts := strings.Split(txty, " ")
			if len(parts) > 0 {
				maxUsers, _ = strconv.Atoi(parts[1])
			}
		}
	}

	return
}
