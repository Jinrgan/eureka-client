package eureka_client

import (
	"fmt"
	"strings"
)

const metadataVersionKey = "VERSION"

type Port struct {
	Port    uint16 `json:"$"`
	Enabled string `json:"@enabled"`
}

type SecurePort struct {
	Field1  int    `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

type LeaseInfo struct {
	RenewalIntervalInSecs int   `json:"renewalIntervalInSecs"`
	DurationInSecs        int   `json:"durationInSecs"`
	RegistrationTimestamp int64 `json:"registrationTimestamp"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
	EvictionTimestamp     int   `json:"evictionTimestamp"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
}

type Metadata struct {
	Zone           string `json:"zone"`
	Profile        string `json:"profile"`
	ManagementPort string `json:"management.port,omitempty"`
	Version        string `json:"version,omitempty"`
}

type Instance struct {
	InstanceId                    string                 `json:"instanceId"`
	HostName                      string                 `json:"hostName"`
	App                           string                 `json:"app"`
	IpAddr                        string                 `json:"ipAddr"`
	Status                        string                 `json:"status"`
	OverriddenStatus              string                 `json:"overriddenstatus"`
	Port                          *Port                  `json:"port"`
	SecurePort                    *SecurePort            `json:"securePort,omitempty"`
	CountryId                     int                    `json:"countryId,omitempty"`
	DataCenterInfo                *DataCenterInfo        `json:"dataCenterInfo"`
	LeaseInfo                     *LeaseInfo             `json:"leaseInfo"`
	Metadata                      map[string]interface{} `json:"metadata"`
	HomePageUrl                   string                 `json:"homePageUrl"`
	StatusPageUrl                 string                 `json:"statusPageUrl"`
	HealthCheckUrl                string                 `json:"healthCheckUrl,omitempty"`
	VipAddress                    string                 `json:"vipAddress"`
	SecureVipAddress              string                 `json:"secureVipAddress"`
	IsCoordinatingDiscoveryServer string                 `json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string                 `json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string                 `json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string                 `json:"actionType,omitempty"`
}

type InstanceOption func(ins *Instance)

func NewInstance(app string, port uint16, opts ...InstanceOption) (*Instance, error) {
	app = strings.ToLower(app)
	ip, ok := GetLocalIP()
	if !ok {
		return nil, fmt.Errorf("cannot get local ip")
	}
	url := fmt.Sprintf("http://%s:%d", ip, port)
	ins := &Instance{
		InstanceId:       fmt.Sprintf("%s:%s:%d", ip, app, port),
		HostName:         ip,
		App:              app,
		IpAddr:           ip,
		Status:           "UP",      // TODO: enum
		OverriddenStatus: "UNKNOWN", // TODO: enum
		Port: &Port{
			Port:    port,
			Enabled: "true",
		}, // TODO: bool
		SecurePort: nil,
		CountryId:  0,
		DataCenterInfo: &DataCenterInfo{
			Name:  "MyOwn",
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		LeaseInfo: &LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        15,
		},
		VipAddress:       app,
		SecureVipAddress: app,
		Metadata: map[string]interface{}{
			metadataVersionKey:     "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
		HomePageUrl:   url,
		StatusPageUrl: url + "/info",
	}

	for _, opt := range opts {
		opt(ins)
	}

	return ins, nil
}

func WithIP(ip string) InstanceOption {
	return func(ins *Instance) {
		port := ins.Port.Port
		url := fmt.Sprintf("http://%s:%d", ip, port)

		ins.InstanceId = fmt.Sprintf("%s:%s:%d", ip, ins.App, port)
		ins.HostName = ip
		ins.IpAddr = ip
		ins.HomePageUrl = url
		ins.StatusPageUrl = url + "/info"
	}
}

func WithVersion(v string) InstanceOption {
	return func(ins *Instance) {
		ins.Metadata[metadataVersionKey] = v
	}
}

type Application struct {
	Name     string      `json:"name"`
	Instance []*Instance `json:"instance"`
}

type Applications struct {
	VersionsDelta string         `json:"versions__delta"`
	AppsHashcode  string         `json:"apps__hashcode"`
	Application   []*Application `json:"application"`
}

type GetApplicationsResponse struct {
	Applications *Applications `json:"applications"`
}
