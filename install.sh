#!/bin/bash
set -e
release="0.0.0"
localBin="/usr/local/bin/dope"
os=$(uname -s)
arch=$(uname -m)
goOs="$(tr '[:upper:]' '[:lower:]' <<< ${os})"
goArch=""

if [ "${arch}" == "x86_64" ]; then
    goArch="amd64" 
else
    echo "Unsupported arch ${arch}\n"
    exit 1
fi

goBin="dope-${goOs}-${goArch}"
binUrl="https://github.com/offers/dope/releases/download/${release}/${goBin}"
curl -L ${binUrl} > ${localBin} 2>/dev/null
chmod +x ${localBin}

echo Installed dope
