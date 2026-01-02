package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/fumihumi/bt-manage/internal/core"
)

func WriteTSV(w io.Writer, devices []core.Device, withHeader bool) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if withHeader {
		fmt.Fprintln(tw, "Name\tAddress\tType\tRSSI")
	}

	for _, d := range devices {
		rssi := ""
		if d.RSSI != nil {
			rssi = fmt.Sprintf("%d", *d.RSSI)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.Name, d.Address, d.Type, rssi)
	}

	return tw.Flush()
}
