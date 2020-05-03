package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/leoluk/perflib_exporter/perflib"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	counter := kingpin.Arg("object", "Which Performance Counter to inspect").Required().String()
	kingpin.Parse()

	nametable := perflib.QueryNameTable("Counter 009") // English
	index := strconv.Itoa(int(nametable.LookupIndex(*counter)))
	if index == "0" {
		fmt.Printf("Counter %q not found\n", *counter)
		os.Exit(1)
	}
	fmt.Printf("Querying for %s (counter index %s)\n", *counter, index)

	t := time.Now()
	objects, err := perflib.QueryPerformanceData(index)
	d := time.Since(t)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		os.Exit(1)
	}
	if len(objects) == 0 {
		fmt.Println("No result")
		os.Exit(1)
	}
	fmt.Printf("Query took %v\n", d)

	var obj *perflib.PerfObject
	for _, o := range objects {
		if o.Name == *counter {
			obj = o
			break
		}
	}
	if obj == nil {
		fmt.Printf("Counter %q was not returned by query\n", *counter)
		os.Exit(1)
	}
	if len(obj.Instances) == 0 {
		fmt.Printf("No instances of %q found\n", *counter)
		os.Exit(1)
	}
	i := obj.Instances[0]

	fmt.Printf("First instance (of %d): %s\n", len(obj.Instances), i.Name)
	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Name", "Value", "Type"})
	for _, c := range i.Counters {
		w.Append([]string{c.Def.Name, strconv.Itoa(int(c.Value)), typeNameMapping[c.Def.CounterType]})
	}
	w.Render()
}

var typeNameMapping = map[uint32]string{
	0x00000000: "PERF_COUNTER_RAWCOUNT_HEX",
	0x00000100: "PERF_COUNTER_LARGE_RAWCOUNT_HEX",
	0x00000b00: "PERF_COUNTER_TEXT",
	0x00010000: "PERF_COUNTER_RAWCOUNT",
	0x00010100: "PERF_COUNTER_LARGE_RAWCOUNT",
	0x00012000: "PERF_DOUBLE_RAW",
	0x00400400: "PERF_COUNTER_DELTA",
	0x00400500: "PERF_COUNTER_LARGE_DELTA",
	0x00410400: "PERF_SAMPLE_COUNTER",
	0x00450400: "PERF_COUNTER_QUEUELEN_TYPE",
	0x00450500: "PERF_COUNTER_LARGE_QUEUELEN_TYPE",
	0x00550500: "PERF_COUNTER_100NS_QUEUELEN_TYPE",
	0x00650500: "PERF_COUNTER_OBJ_TIME_QUEUELEN_TYPE",
	0x10410400: "PERF_COUNTER_COUNTER",
	0x10410500: "PERF_COUNTER_BULK_COUNT",
	0x20020400: "PERF_RAW_FRACTION",
	0x20020500: "PERF_LARGE_RAW_FRACTION",
	0x20410500: "PERF_COUNTER_TIMER",
	0x20470500: "PERF_PRECISION_SYSTEM_TIMER",
	0x20510500: "PERF_100NSEC_TIMER",
	0x20570500: "PERF_PRECISION_100NS_TIMER",
	0x20610500: "PERF_OBJ_TIME_TIMER",
	0x20670500: "PERF_PRECISION_OBJECT_TIMER",
	0x20c20400: "PERF_SAMPLE_FRACTION",
	0x21410500: "PERF_COUNTER_TIMER_INV",
	0x21510500: "PERF_100NSEC_TIMER_INV",
	0x22410500: "PERF_COUNTER_MULTI_TIMER",
	0x22510500: "PERF_100NSEC_MULTI_TIMER",
	0x23410500: "PERF_COUNTER_MULTI_TIMER_INV",
	0x23510500: "PERF_100NSEC_MULTI_TIMER_INV",
	0x30020400: "PERF_AVERAGE_TIMER",
	0x30240500: "PERF_ELAPSED_TIME",
	0x40000200: "PERF_COUNTER_NODATA",
	0x40020500: "PERF_AVERAGE_BULK",
	0x40030401: "PERF_SAMPLE_BASE",
	0x40030402: "PERF_AVERAGE_BASE",
	0x40030403: "PERF_RAW_BASE",
	0x40030500: "PERF_PRECISION_TIMESTAMP",
	0x40030503: "PERF_LARGE_RAW_BASE",
	0x42030500: "PERF_COUNTER_MULTI_BASE",
	0x80000000: "PERF_COUNTER_HISTOGRAM_TYPE",
}
