package eureka_client

import "net"

// GetLocalIP 获取本地ip
func GetLocalIP() (string, bool) {
	addres, err := net.InterfaceAddrs()
	if err != nil {
		return "", false
	}
	for _, address := range addres {
		// check the address type and if it is not a loopback the display it
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), true
			}
		}
	}
	return "", false
}
