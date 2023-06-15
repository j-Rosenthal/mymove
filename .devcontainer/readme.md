# Devcontainer Setup

In order to make this application more modular for those working on any type of machine we can utilize [Devcontainers](https://containers.dev/). This will allow anyone to setup their environment with the most minimal amount of dependencies.

# VsCode
First you will have to download [VSCode](https://code.visualstudio.com/) in order to take advantage of the toolset.

Also it is recommended to install the extensions that are included in the `extensions.json` inside the `.vscode` folder

# Docker
Next you need to install [Docker Desktop](https://www.docker.com/products/docker-desktop/) on your machine.


## Windows Only Steps
If you are on a Windows machine it is recommended you do these additional steps.

### WSL2
In order to keep file endings in sync between the main project and the project inside of the docker container it is recommended to use a linux environment. (**NOTE** You can modify your git config on windows to do the same thing, however all projects going forward would have the different line endings. You also must do this before you clone the project.)

1. Open a powershell terminal in administrator mode and enter the following:
```PowerShell
wsl --install
```
This command will enable the features necessary to run WSL and install the Ubuntu distribution of Linux. ([This default distribution can be changed](https://learn.microsoft.com/en-us/windows/wsl/basic-commands#install)).

Afer you have WSL2 installed and Ubuntu you may open the Ubuntu Window and clone the repository.

Navigate to the project in WSL and type `code .` this will open VSCode within the current directory.

# Open in Dev Container

Now that you have the project cloned and opened in VSCode and the DevContainer Extension is installed all you have to do is type: 

`Ctrl + Shift + p` you should then see the command window appear at the top of your window. (It should be prefixed with `>`)
In there type:
`Dev Containers: Reopen in Container`.

The VSCode window will reload and then open again inside the dev container running through the build steps defined in the devcontainer.json, docker-compose.yml, and Dockerfile
