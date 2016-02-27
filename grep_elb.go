package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseLine(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, err
}

func parseElblogfile(t1 *bool, t2 *bool, t3 *bool, threshold *int, debug *bool, verbose *bool, filename string) {

	var time1, time2, time3 = "-", "-", "-"
	lines, err := parseLine(filename)

	if err != nil {
		fmt.Println(err)
	} else {
		for _, line := range lines {
			items := strings.Split(line, " ")
			date, elb, clientip, referrer, requestProcessingTime, backendProcessingTime, responseProcessingTime, elbstatus, backendstatus, rbytes, sbytes, verb, url := items[0], items[1], items[2], items[3], items[4], items[5], items[6],
				items[7], items[8], items[9], items[10], items[11], items[12]

			req, err := strconv.ParseFloat(requestProcessingTime, 64)
			if *t1 && err == nil && (req > float64(*threshold)) {
				time1 = requestProcessingTime
				fmt.Println("-------------1111")
			}

			bkend, err := strconv.ParseFloat(backendProcessingTime, 64)
			if *t2 && err == nil && (bkend > float64(*threshold)) {
				time2 = backendProcessingTime
				fmt.Println("-------------2222")
			}

			resp, err := strconv.ParseFloat(responseProcessingTime, 64)
			if *t3 && err == nil && (resp > float64(*threshold)) {
				time3 = responseProcessingTime
				fmt.Println("-------------3333")
			}

			if (time1 != "-") || (time2 != "-") || (time3 != "-") {
				if *verbose {
					fmt.Println(date, elb, clientip, referrer, requestProcessingTime, backendProcessingTime, responseProcessingTime, elbstatus, backendstatus, rbytes, sbytes, verb, url)
				} else {
					fmt.Println("date:", date, ", t1:", time1, ", t2:", time2, ", t3:", time3, ", url:", url, ", backend:", backendstatus)
				}
			}

		}
	}

}

func main() {
	t1 := flag.Bool("i", false, "Compare time of elb internal processing time (sec). Defaults false, if none specified.")
	t2 := flag.Bool("b", false, "Compare time of backend processing time (sec). Defaults false, if none specified.")
	t3 := flag.Bool("r", false, "Compare time of response processing time (sec). Defaults false, if none specified.")
	t := flag.Int("t", 1, "Time threshold to compare against.")
	d := flag.Bool("d", false, "Display debug information.")
	v := flag.Bool("v", false, "Display whole log statement in results.")
	flag.Parse()

	if *d == true {
		fmt.Println("internal:", *t1, ", backend:", *t2, ", response:", *t3, ", threshold:", *t, ", debug:", *d, ", verbose: ", *v)
		for _, fn := range flag.Args() {
			fmt.Println("Filepath: ", fn)
		}
	}

	for _, fn := range flag.Args() {
		parseElblogfile(t1, t2, t3, t, d, v, fn)
	}

}