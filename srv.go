package srv

import (
	"fmt"
	"net"
	"strings"
)

func Resolve(address string) ([]string, error) {
	var err error
	var records []*net.SRV

	// _noded
	// _noded._tcp.example.com
	if strings.Count(address, ".") == 0 {
		_, records, err = net.LookupSRV("", "", address)
	} else if strings.Count(address, ".") > 1 {
		parts := strings.SplitN(address, ".", 3)
		_, records, err = net.LookupSRV(
			strings.TrimLeft(parts[0], "_"),
			strings.TrimLeft(parts[1], "_"),
			parts[2],
		)
	} else {
		return nil, fmt.Errorf(
			"SRV record '%s' is malformed, should be "+
				"like as _service or _service._tcp.example.com", address,
		)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"can't resolve SRV record for '%s': %s", address, err,
		)
	}

	var addresses []string
	for _, srvRecord := range records {
		resolved := fmt.Sprintf(
			"%s:%d",
			strings.Trim(srvRecord.Target, "."), srvRecord.Port,
		)

		addresses = append(addresses, resolved)
	}

	return addresses, nil
}
