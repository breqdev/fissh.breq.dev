package timezone

import (
	"log"
	"net"
	"os"

	"github.com/ipinfo/go/v2/ipinfo"
)

func LookupTimezone(ipAddress string) string {
	token := os.Getenv("IPINFO_TOKEN")

	client := ipinfo.NewClient(nil, nil, token)

	info, err := client.GetIPInfo(net.ParseIP(ipAddress))
	if err != nil {
		log.Fatal(err)
	}

	return info.Timezone
}
