#!/bin/bash



echo "Installing Go V1.20.3...."
goVersion=go1.20.3
nodeVersion=18.13.0
nvmVersion=v0.39.3



wget https://go.dev/dl/$goVersion.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $goVersion.linux-amd64.tar.gz
cat PATH=$PATH:/usr/local/go/bin >> ~/.profile
cat export GOPATH=$HOME/go >> ~/.profile
cat export PATH=$PATH:$GOROOT/bin:$GOPATH/bin >> ~/.profile



go version



echo "Installing NVM and NodeJs..."



curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash



export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
cat ~/.nvm/nvm.sh >> ~/.profile



echo "Installing NodeJs V18.13.0"



nvm install $nodeVersion



echo "installing homebrew...."



/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"



test -d ~/.linuxbrew && eval "$(~/.linuxbrew/bin/brew shellenv)"
test -d /home/linuxbrew/.linuxbrew && eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
test -r ~/.bash_profile && echo "eval \"\$($(brew --prefix)/bin/brew shellenv)\"" >> ~/.bash_profile
echo "eval \"\$($(brew --prefix)/bin/brew shellenv)\"" >> ~/.profile



echo "installing brew dependencies..."
brew install asdf awscli bash chamber circlelci diffutils direnv entr jq nodenv postgresql pre-commit spellcheck watchman yarn



echo "dependency install complete"
source ~/.profile



echo "Verify install by running the following:"
echo "go version"
echo "nvm version"
echo "brew"
