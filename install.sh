if [ $(id -u) != "0" ]; then
echo "You must be the root to run this script" >&2
exit 1
fi

echo "[*]Imp[*] The tools only works on Kali as of now. Dependencies will be installed accordingly"
echo
echo "[*] Ensure that the script is run from the same directory from which install.sh is run"

echo
echo "[*] Updating packages"
apt-get update
echo

##Creating temp directory
echo "[*] Creating temporary directory "subdomainenum_temp" to store temporary data"
mkdir subdomainenum_temp
DEST=${DEST:-subdomainenum_temp}
echo

##Install Golang
echo "[*] Verifying Golang Installation"
hash go 2>/dev/null || { echo "Golang not installed. Installing…"; cd "${DEST}"; wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz; tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz; echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile; echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc; source $HOME/.profile;}
echo

##Install Golang Packages
echo "[*] Installing Golang packages"
go get github.com/fatih/color
go get github.com/common-nighthawk/go-figure
echo

##Install MassDNS
echo "[*] Installing MassDNS"
cd "${DEST}"
git clone https://github.com/blechschmidt/massdns.git
cd massdns
make
cd bin
cp massdns /usr/sbin
echo

##Installing Json processor
echo "[*]Installing JSON processor"
echo
apt-get -y install jq
echo

echo "All dependencies have been installed successfully"
