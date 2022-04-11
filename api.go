package eureka_client

type Port struct {
	Field1  int    `json:"$"`
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
	InstanceId                    string          `json:"instanceId"`
	HostName                      string          `json:"hostName"`
	App                           string          `json:"app"`
	IpAddr                        string          `json:"ipAddr"`
	Status                        string          `json:"status"`
	OverriddenStatus              string          `json:"overriddenstatus"`
	Port                          *Port           `json:"port"`
	SecurePort                    *SecurePort     `json:"securePort"`
	CountryId                     int             `json:"countryId"`
	DataCenterInfo                *DataCenterInfo `json:"dataCenterInfo"`
	LeaseInfo                     *LeaseInfo      `json:"leaseInfo"`
	Metadata                      *Metadata       `json:"metadata"`
	HomePageUrl                   string          `json:"homePageUrl"`
	StatusPageUrl                 string          `json:"statusPageUrl"`
	HealthCheckUrl                string          `json:"healthCheckUrl"`
	VipAddress                    string          `json:"vipAddress"`
	SecureVipAddress              string          `json:"secureVipAddress"`
	IsCoordinatingDiscoveryServer string          `json:"isCoordinatingDiscoveryServer"`
	LastUpdatedTimestamp          string          `json:"lastUpdatedTimestamp"`
	LastDirtyTimestamp            string          `json:"lastDirtyTimestamp"`
	ActionType                    string          `json:"actionType"`
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