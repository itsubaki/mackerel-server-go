package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TSDBLatest struct {
	TSDBLatest TSDBLatestValue `json:"tsdbLatest"`
}

// [hostId][name]metric_value
type TSDBLatestValue map[string]map[string]float64

type HostMetadata struct {
	HostID    string      `json:"-"`
	Namespace string      `json:"-"`
	Metadata  interface{} `json:"-"`
}

type HostMetadataList struct {
	Metadata []Namespace `json:"metadata"`
}

type Namespace struct {
	Namespace string `json:"namespace"`
}

type MetricNames struct {
	Names []string `json:"names"`
}

type HostRetire struct {
}

type RoleFullNames struct {
	Names []string `json:"roleFullnames"`
}

func (r *RoleFullNames) Roles() map[string][]string {
	roles := make(map[string][]string)
	for i := range r.Names {
		svc := strings.Split(r.Names[i], ":")
		if _, ok := roles[svc[0]]; !ok {
			roles[svc[0]] = []string{}
		}
		roles[svc[0]] = append(roles[svc[0]], svc[1])
	}

	return roles
}

type HostID struct {
	ID string `json:"id"`
}

type HostStatus struct {
	Status string `json:"status"`
}

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

type HostInfo struct {
	Host Host `json:"host"`
}

type Host struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Status           string              `json:"status"`
	Memo             string              `json:"memo"`
	DisplayName      string              `json:"displayName,omitempty"`
	CustomIdentifier string              `json:"customIdentifier,omitempty"`
	CreatedAt        int64               `json:"createdAt"`
	RetiredAt        int64               `json:"-"`
	IsRetired        bool                `json:"isRetired"`
	Roles            map[string][]string `json:"roles"`
	RoleFullNames    []string            `json:"roleFullnames,omitempty"`
	Interfaces       []Interface         `json:"interfaces,omitempty"`
	Checks           []Check             `json:"checks,omitempty"`
	Meta             Meta                `json:"meta"`
}

func (h *Host) Init() {
	h.Roles = make(map[string][]string)
	for i := range h.RoleFullNames {
		svc := strings.Split(h.RoleFullNames[i], ":")
		if _, ok := h.Roles[svc[0]]; !ok {
			h.Roles[svc[0]] = make([]string, 0)
		}

		h.Roles[svc[0]] = append(h.Roles[svc[0]], svc[1])
	}

	if len(h.ID) > 0 {
		return
	}

	sha := sha256.Sum256([]byte(uuid.Must(uuid.NewRandom()).String()))
	hash := hex.EncodeToString(sha[:])

	h.ID = hash[:11]
	h.CreatedAt = time.Now().Unix()
	h.RetiredAt = 0
	h.IsRetired = false
	h.Checks = []Check{}
	if len(h.Status) < 1 {
		h.Status = "working"
	}
}

type Meta struct {
	AgentName     string                 `json:"agent-name"`
	AgentRevision string                 `json:"agent-revision"`
	AgentVersion  string                 `json:"agent-version"`
	BlockDevice   map[string]interface{} `json:"block_device,omitempty"`
	CPU           []interface{}          `json:"cpu,omitempty"`
	FileSystem    map[string]interface{} `json:"filesystem,omitempty"`
	Kernel        map[string]string      `json:"kernel,omitempty"`
	Memory        map[string]string      `json:"memory,omitempty"`
}

type Interface struct {
	Name           string   `json:"name"`
	IpAddress      string   `json:"ipAddress"`
	MacAddress     string   `json:"macAddress"`
	IpV4Addresses  []string `json:"ipv4Addresses"`
	IpV6Addresses  []string `json:"ipv6Addresses"`
	DefaultGateway string   `json:"defaultGateway"`
}

type Check struct {
	Name string `json:"name"`
	Memo string `json:"memo,omitempty"`
}

type MetricValues struct {
	Metrics []MetricValue `json:"metrics"`
}

type MetricValue struct {
	HostID string  `json:"hostId,omitempty"`
	Name   string  `json:"name,omitempty"`
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
}
