# Subdomain Enumerator
Subdomain Enumerator is a simple automation of the subdomain recon process which utilized output from multiple tools to gather a list of subdomain for a TLD. It is aimed to reduce the effort of:

1) Running multiple tools
2) Merging outputs
3) Sanitization

The tool might be a bit rusty as this is my first such attempt but it does the work.

## Tool Workflow
The tool follows the following process for subdomain enumeration:

1) Perform wildcard configuration check for the domain
2) Run Amass
3) Query SSL repositories
4) Create a wordlist for DNS bruteforcing
5) Perform DNS bruteforcing using MassDNS
6) Merge the results 
7) Perform sanitization

## Tool Details
The script is written in Go; so it is possible to cross compile it easily for other environments

## Getting Started | Tool Dependencies
Clone the repository using `git clone`. Then, you need to run `./install.sh`, which will install the dependecies for the tool.

**Note:** You need to be root to install the dependencies using `install.sh`

## Usage
Ensure that the script is run from the same directory from which `install.sh` is run

```
Usage: go run SubdomainENumerator.go -domain [domain_to_be_enumerated]
   -domain string
       Domain to be enumerated for subdomains e.g. yahoo.com
   -h  --help   
       Shows help
```
### Credits
Thanks to [Noobhax](https://medium.com/@noobhax/my-recon-process-dns-enumeration-d0e288f81a8a) for the recon process which sets base for this script.
