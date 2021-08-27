package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lord3ver/adnslookup/pkg/dnsip"
)

const desc = "A DNS Lookup tool. Concurrent with worker, pool of workers and generator \"load balance\" mechanism."
const version = "1.0.0"
const author = "LordEver (@Lord3ver)"

func main() {

	banner := `
     ___    ____             __                __             
    /   |  / __ \____  _____/ /   ____  ____  / /____  ______ 
   / /| | / / / / __ \/ ___/ /   / __ \/ __ \/ //_/ / / / __ \
  / ___ |/ /_/ / / / (__  ) /___/ /_/ / /_/ / ,< / /_/ / /_/ /
 /_/  |_/_____/_/ /_/____/_____/\____/\____/_/|_|\__,_/ .___/ 
                                                     /_/      
	`

	fmt.Println(banner)
	fmt.Printf("%s\t\n\nVersion:\t%s\nAuthor:\t\t%s\n\n", desc, version, author)

	threadsPtr := flag.Int("t", 25, "Max threads")
	domainPtr := flag.String("d", "", "Target domain")
	domainsFPtr := flag.String("f", "", "Targets file. One per line.")
	outPtr := flag.Bool("out", true, "Print results to stdout")
	outFilePtr := flag.String("outfile", "", "Specify an output file when completed. Create or append if exists.")
	outFileNonePtr := flag.String("outfileNoDNS", "", "Specify an output file for domains with no DNS record. Create or append if exists.")

	flag.Parse()

	if *threadsPtr < 1 {
		fmt.Println("What would you like to do? Threads cannot be 0!")
		flag.PrintDefaults()
		return
	}
	if *outPtr == false && *outFilePtr == "" && *outFileNonePtr == "" {
		fmt.Println("\"out\" is false and both \"outfile\" and \"outfileNoDNS\" are not set, where do you wanna to go?\n\nStdout output enabled.")
		*outPtr = true
	}
	if *domainPtr == "" && *domainsFPtr == "" {
		fmt.Println("At least one option among \"d\" (target domain) and \"f\" (targets file) must be specified.")
		flag.PrintDefaults()
		return
	}
	if *domainPtr != "" && *domainsFPtr != "" {
		fmt.Println("\"d\" (target domain) and \"f\" (targets file) cannot be used together. We will going to considering only the input filename")
		*domainPtr = ""
	}

	var res []string
	var resNone []string

	// Single target
	if *domainPtr != "" {
		res, resNone = dnsip.Lookup(*threadsPtr, []string{*domainPtr}, *outPtr)
	} else {
		targets := getTargets(*domainsFPtr)
		res, resNone = dnsip.Lookup(*threadsPtr, targets, *outPtr)
	}

	if *outFilePtr != "" {
		if len(res) > 0 {
			// Sort results
			sort.Slice(res, func(i, j int) bool {
				switch strings.Compare(strings.Split(res[i], "\t")[1], strings.Split(res[j], "\t")[1]) {
				case -1:
					return true
				case 1:
					return false
				}
				return true
			})

			lookupToFile(*outFilePtr, res)
		}
	}

	if *outFileNonePtr != "" {
		if len(resNone) > 0 {
			// Sort "none" results
			sort.Slice(resNone, func(i, j int) bool {
				switch strings.Compare(resNone[i], resNone[j]) {
				case -1:
					return true
				case 1:
					return false
				}
				return true
			})
			lookupToFile(*outFileNonePtr, resNone)
		}
	}

	fmt.Println("\nDone!")
}

// getTargets read targets from file.
// Returns a slice of targets.
func getTargets(fn string) []string {
	file, err := os.Open(fn)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return lines

}

// lookupToFile save lookup results to a file.
func lookupToFile(fn string, res []string) {
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range res {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}
