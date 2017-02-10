# dope
## About
Dope manages software distributed as Docker images. It can:
* Install Docker images and make them available as a command in your shell. No more adding an alias for every image's command.
* Notify you when a new image version is available, and update it for you.
* Automatically the right Docker command from the image (if it includes a .dope.json file).

## Install
```
\curl -sSL https://raw.githubusercontent.com/offers/dope/master/install.sh | sudo bash
```

Add `~/.dope/bin` to your PATH

## Usage
```bash
$ dope install my_registry/my_org/my_repo # install from your private registry

$ my_repo # run the docker image

$ dope update my_repo # pull the highest tag

$ dope uninstall my_repo # delete docker image and shell command

$ dope list # show installed packages

$ dope self-update # install the latest version of dope

$ dope check my_other_repo # check if an update is available
```
