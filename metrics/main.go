package metrics

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
)

type metrics struct {
	bandwidth         bool
	domainsConfigured bool
	ftpAccounts       bool
	sessions          bool
	license           bool
	meta              bool
	cache             map[string]string
}

func New(
	bandwidth bool,
	domainsConfigured bool,
	ftpAccounts bool,
	sessions bool,
	license bool,
	meta bool,
) *metrics {
	return &metrics{
		bandwidth:         bandwidth,
		domainsConfigured: domainsConfigured,
		ftpAccounts:       ftpAccounts,
		sessions:          sessions,
		license:           license,
		meta:              meta,
		cache:             make(map[string]string),
	}
}

type Metrics interface {
	FetchMetrics()
	FetchUpiMetrics()
	GetCache() map[string]string
	GetSortedCacheString() (data string)
}

func (m metrics) FetchMetrics() {
	total := getUsersCount(false)
	suspended := getUsersCount(true)
	m.cache["cpanel_users_total"] = strconv.Itoa(total)
	m.cache["cpanel_users_active"] = strconv.Itoa(total - suspended)
	m.cache["cpanel_users_suspended"] = strconv.Itoa(suspended)
	if m.sessions {
		webSessions, emailSessions := getSessions()
		m.cache["cpanel_sessions_web"] = strconv.Itoa(webSessions)
		m.cache["cpanel_sessions_email"] = strconv.Itoa(emailSessions)
	}
	if m.license {
		expireTime, maxUsers := getLicenseInfo()
		m.cache[fmt.Sprintf("cpanel_license{expire_time=\"%s\",max_usesrs=\"%d\"}", expireTime, maxUsers)] = "0"
	}
	if m.meta {
		m.cache[fmt.Sprintf("cpanel_meta{version=\"%s\",release=\"%s\"}", cpanelVersion(), getRelease())] = "0"
	}
	if m.domainsConfigured {
		m.cache["cpanel_domains_configured"] = strconv.Itoa(len(getDomains()))
	}
	if m.ftpAccounts {
		m.cache["cpanel_ftp_accounts"] = strconv.Itoa(len(getFTP()))
	}
	for p, ct := range getPlans() {
		m.cache[fmt.Sprintf("cpanel_plans{plan=\"%s\"}", p)] = strconv.Itoa(ct)
	}
	return
}

func (m metrics) FetchUpiMetrics() {
	for _, u := range getUsers() {
		us := filepath.Base(u)
		usedMegabyte, percentMegabyte, usedInode, percentInode := getQuota(us)
		m.cache[fmt.Sprintf("cpanel_megabyte_quota_used{user=\"%s\"}", us)] = strconv.FormatFloat(usedMegabyte, 'f', -1, 64)
		m.cache[fmt.Sprintf("cpanel_megabyte_quota_percent{user=\"%s\"}", us)] = strconv.FormatFloat(percentMegabyte, 'f', -1, 64)
		m.cache[fmt.Sprintf("cpanel_inode_quota_used{user=\"%s\"}", us)] = strconv.FormatFloat(usedInode, 'f', -1, 64)
		m.cache[fmt.Sprintf("cpanel_inode_quota_percent{user=\"%s\"}", us)] = strconv.FormatFloat(percentInode, 'f', -1, 64)
		if m.bandwidth {
			m.cache[fmt.Sprintf("cpanel_bandwidth{user=\"%s\"}", us)] = strconv.Itoa(getBandwidth(us))
		}
	}
}

func (m metrics) GetCache() map[string]string {
	return m.cache
}

func (m metrics) GetSortedCacheString() (data string) {
	keys := make([]string, 0, len(m.cache))
	for k, _ := range m.cache {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		data += fmt.Sprintf("%s %s\n", k, m.cache[k])
	}
	return
}
