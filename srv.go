package srv // import "github.com/reconquest/srv-go"

import (
	"fmt"
	"net"
	"strings"
)

var (
	Testing_RecordToResult = map[string][]string{}
)

// Resolve specified address, returns sorted by priority and randomized by
// weight within a priority records list.
func Resolve(record string) ([]string, error) {
	var err error
	var records []*net.SRV

	if testing, ok := Testing_RecordToResult[record]; ok {
		return testing, nil
	}

	// _noded
	// _noded._tcp.example.com
	if strings.Count(record, ".") == 0 {
		_, records, err = net.LookupSRV("", "", record)
	} else if strings.Count(record, ".") > 1 {
		parts := strings.SplitN(record, ".", 3)
		_, records, err = net.LookupSRV(
			strings.TrimLeft(parts[0], "_"),
			strings.TrimLeft(parts[1], "_"),
			parts[2],
		)
	} else {
		return nil, fmt.Errorf(
			"SRV record '%s' is malformed, should be "+
				"like as _service or _service._tcp.example.com", record,
		)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"can't resolve SRV record for '%s': %s", record, err,
		)
	}

	var addresses []string
	for _, srvRecord := range records {
		address := fmt.Sprintf(
			"%s:%d",
			strings.Trim(srvRecord.Target, "."), srvRecord.Port,
		)

		addresses = append(addresses, address)
	}

	return addresses, nil
}
