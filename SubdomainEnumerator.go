package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

const (
	// LogPrefix - for verbose console logging
	LogPrefix = "[*] "
)

var (
	//initialize logvar for package level scope
	logvar *log.Logger
	green  func(...interface{}) string
)

func main() {
	//ASCII Text
	figure.NewFigure("Subdomain Enumerator", "small", true).Print()
	fmt.Printf("\n")

	//Set up logging
	logvar = log.New(os.Stdout, LogPrefix, log.Ltime)
	green = color.New(color.FgGreen).SprintFunc()

	//Define flag for input
	domain := flag.String("domain", "", "Domain to be enumerated for subdomains e.g. yahoo.com")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Please input the domain:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	wildcard(*domain)
	amass(*domain)
	crtsh(*domain)
	dnsbrute(*domain)
	fmt.Printf("\n")
	logvar.Println(green("Script execution has been completed. Please check the results"))
}

//Perform wildcard check
func wildcard(d string) {

	logvar.Println("Checking wildcard configuration for:", d)

	//creating a string of command to be run
	arg := fmt.Sprintf("%s%s%s", "dig @1.1.1.1 A,CNAME {test321123,testingforwildcard,plsdontgimmearesult}.", d, " +short | wc -l")
	out, err := exec.Command("bash", "-c", arg).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//split the string output to exclude newline and get an int
	output := strings.Split(string(out), "\n")
	if output[0] == "0" {
		logvar.Printf(green("No wildcard misconfiguration detected\n"))
	} else {
		logvar.Printf("Possible wildcard misconfiguration ! Please evaluate further.\n")
	}

	return
}

func amass(d string) {
	logvar.Printf("Running Amass for subdomain enumeration for TLD: %s", d)

	file := "subdomainenum_temp/subdomains_temp.txt"

	//creating a string of commands to be run
	cmd := exec.Command("amass", "enum", "-o", file, "-d", d)

	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Amass execution failed with %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n")
	logvar.Printf(green("Amass ran successfully\n"))
	return
}

func crtsh(domain string) {
	logvar.Println("Querying SSL certificates for subdomains for TLD:", domain)

	arg1 := fmt.Sprintf("%s%s%s", "curl \"https://crt.sh/?q=%.", domain, "&output=json\"")
	arg2 := ` | jq '.[].name_value' | sed 's/\"//g' | sed 's/\*\.//g' | sort -u`
	crtarg := fmt.Sprintf("%s%s", arg1, arg2)

	cmd, err := exec.Command("bash", "-c", crtarg).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//appending data to subdomains_temp.txt file
	f, err := os.OpenFile("subdomainenum_temp/subdomains_temp.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	_, err = f.Write(cmd)
	if err != nil {
		fmt.Println(err)
	}

	logvar.Printf(green("Subdomains have been extracted from SSL certificates successfully\n"))
	return
}

func dnsbrute(domain string) {
	logvar.Printf("Initiating DNS bruteforcing\n")
	logvar.Printf("Downloading top 1 million DNS records wordlist from Seclists")

	//downloading wordlist and saving to a file
	response, err := http.Get("https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-110000.txt")
	if err != nil {
		println(err)
	}
	defer response.Body.Close()

	file, err := os.Create("subdomainenum_temp/subdomains-top1million-110000.txt")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	defer file.Close()

	io.Copy(file, response.Body)

	//creating a string of commands to be run
	cmd := fmt.Sprintf("%s%s%s", `sed 's/$/.`, domain, "/' subdomainenum_temp/subdomains-top1million-110000.txt")

	filename := "subdomainenum_temp/subdomains-top1million-110000-wordlist.txt"
	domainList, err := os.Create(filename)
	if err != nil {
		println(err)
		os.Exit(1)
	}

	//creating wordlist using the downloaded list of domains
	out := exec.Command("bash", "-c", cmd)

	//redirecting standard output to file
	out.Stdout = domainList
	err = out.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//not using defer as file needs to be reopened for reading as output is needed as []bytes
	domainList.Close()
	logvar.Printf("Subdomain wordlist %s has been created\n", filename)

	//opening already created subdomain file for appending the wordlist
	subdomainFile, err := os.OpenFile("subdomainenum_temp/subdomains_temp.txt", os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer subdomainFile.Close()

	//opening and reading subdomains-top1million-110000-wordlist.txt file to return []bytes
	wordListfile, err := os.Open("subdomainenum_temp/subdomains-top1million-110000-wordlist.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer wordListfile.Close()

	data, err := ioutil.ReadAll(wordListfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//writing []bytes from reading subdomains-top1million-110000-wordlist.txt to subdomains_temp.txt
	_, err = subdomainFile.Write(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logvar.Printf(green("Data from all sources merged in subdomainenum_temp/subdomains_temp.txt\n"))

	logvar.Printf("Running MassDNS to find online domains")

	arg := "massdns -r subdomainenum_temp/massdns/lists/resolvers.txt -t A -o S -w subdomainenum_temp/subdomains_with_metadata.txt subdomainenum_temp/subdomains_temp.txt"
	massdnsCmd := fmt.Sprintf("%s", arg)

	massdnsOut := exec.Command("bash", "-c", massdnsCmd)

	err = massdnsOut.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logvar.Printf(green("MassDNS has been run successfully"))

	//Filtering to keep all the data(to be developed in future releases): sed -e 's/[^ ] / /' -e 's/m./m/2' subdomains_with_metadata.txt
	// awk '{print $1}' sc.txt | sed 's/.$//' | sort -u
	uniqueDomains := `awk '{print $1}' subdomainenum_temp/subdomains_with_metadata.txt | sed 's/.$//' | sort -u > subdomains.txt`
	count := `awk '{print $1}' subdomainenum_temp/subdomains_with_metadata.txt | sed 's/.$//' | sort -u | wc -l`

	//counting number of unique domains found
	countCmd, err := exec.Command("bash", "-c", count).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	domainCount := string(countCmd)

	logvar.Println(green("Total number of unique domains have been discovered: "), domainCount)

	//finding unique domains
	_, err = exec.Command("bash", "-c", uniqueDomains).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logvar.Printf(green("List of domains with metadata has been saved at subdomainenum_temp/subdomains_with_metadata.txt and unique subdomains have been saved at subdomains.txt"))
	return
}
