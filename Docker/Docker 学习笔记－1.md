

本文针对 Docker for Mac 进行说明，以下内容整理自《[Getting Started with Docker for Mac](https://docs.docker.com/docker-for-mac/)》


----------


# 系统要求

>- Mac must be a 2010 or newer model, with Intel’s hardware support for **memory management unit (MMU) virtualization**; i.e., **Extended Page Tables (EPT)**
>- OS X **10.10.3** Yosemite or newer
>- At least **4GB of RAM**
>- VirtualBox prior to version 4.3.30 must NOT be installed (it is incompatible with Docker for Mac)


# 安装

## 构成 Docker for Mac 的组件

Docker for Mac 安装后，将得到如下组件
- Docker Engine
- Docker CLI client
- Docker Compose
- Docker Machine

安装后确认

```shell
docker --version
docker-compose --version
docker-machine --version
docker ps
```

#测试

> 如果 image 无法在本地找到，那么 Docker 将默认从 [Docker Hub](https://hub.docker.com/) 上进行 pull 拉取；

## 测试一（hello world）

```shell
➜  ~ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
c04b14da8d14: Pull complete
Digest: sha256:0256e8a36e2070f7bf2d0b0763dbabdd67798512411de4cdcf9431a1feb60fd9
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker Hub account:
 https://hub.docker.com

For more examples and ideas, visit:
 https://docs.docker.com/engine/userguide/

➜  ~
```

该例子描述了 `docker run hello-world` 背后的行为：
- 首先由 Docker client 与 Docker daemon 建立联系；
- 其次 Docker daemon 尝试从 Docker Hub 上拉取名为 "hello-world" 的 image ；
- 再次 Docker daemon 会基于该 image 创建一个新的 container 用于运行输出 "Hello from Docker!" 的可执行程序；
- 最后 Docker daemon 会将该输出内容 stream 给 Docker client ，即发送到你使用的 terminal 上；

## 测试二（nginx）

启动一个 Docker 化的 web 服务器（基于 nginx 的 webserver 将使用 80 端口运行于 container 中）

测试步骤
```shell
docker run -d -p 80:80 --name webserver nginx
docker ps
docker stop webserver       ## stop the container
docker start webserver      ## start the container
```

具体执行过程
```erlang
➜  Docker docker run -d -p 80:80 --name webserver nginx
Unable to find image 'nginx:latest' locally
latest: Pulling from library/nginx
8ad8b3f87b37: Pull complete
c6b290308f88: Pull complete
f8f1e94eb9a9: Pull complete
Digest: sha256:aa5ac743d65e434c06fff5ceaab6f35cc8519d80a5b6767ed3bdb330f47e4c31
Status: Downloaded newer image for nginx:latest
0b5f46136d5f65573847071aba6b3b27d6a9195cfb9a071854ba2e6a3527ac2a
➜  Docker
➜  Docker
➜  Docker docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED              STATUS              PORTS                         NAMES
0b5f46136d5f        nginx               "nginx -g 'daemon off"   About a minute ago   Up About a minute   0.0.0.0:80->80/tcp, 443/tcp   webserver
➜  Docker
➜  Docker
➜  Docker docker stop webserver
webserver
➜  Docker
➜  Docker docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
➜  Docker
➜  Docker
➜  Docker docker start webserver
webserver
➜  Docker
➜  Docker
➜  Docker docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                         NAMES
0b5f46136d5f        nginx               "nginx -g 'daemon off"   2 minutes ago       Up 2 seconds        0.0.0.0:80->80/tcp, 443/tcp   webserver
➜  Docker
```

移除 container ，但不移除 nginx 的 image ；

```shell
docker rm -f webserver      ## stop and remove the running container
```

列出本地所有 image ；

```shell
docker images
```

> ⚠️ 在一些情况下，很有可能你会希望将一些 image 保留下来，这样就不再需要再次从 Docker Hub 上进行 pull ；

实际结果

```shell
➜  Docker docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
nginx               latest              4a88d06e26f4        5 days ago          183.5 MB
hello-world         latest              c54a2cc56cbb        11 weeks ago        1.848 kB
➜  Docker
```

移除不再需要的 image ；

```shell
docker rmi <imageID>|<imageName>
```

通过 console 命令卸载 Docker 

```shell
/Applications/Docker.app/Contents/MacOS/Docker --uninstall
```

> 你可能想要通过命令行接口执行 uninstall 动作，因为存在 app 本身运行不正常，进而无法通过菜单命令直接 uninstall 的情况；


# Bash Completion on OS X With Brew

> ⚠️ 本段内容针对 bash 使用；

```shell
brew install bash-completion
```

根据提示，将如下内容添加到 ~/.bash_profile 文件中

```shell
if [ -f $(brew --prefix)/etc/bash_completion ]; then
    . $(brew --prefix)/etc/bash_completion
fi
```

```shell
brew tap homebrew/completions
```

为了激活 bash completion ，如下文件需要被拷贝或符号链接到你的 bash_completion.d 目录中；

```shell
cd /usr/local/etc/bash_completion.d
ln -s /Applications/Docker.app/Contents/Resources/etc/docker.bash-completion
ln -s /Applications/Docker.app/Contents/Resources/etc/docker-machine.bash-completion
ln -s /Applications/Docker.app/Contents/Resources/etc/docker-compose.bash-completion
```


# Q&A


## [Docker Toolbox](https://docs.docker.com/toolbox/overview/)

Docker Toolbox 是一种 installer ，用于在不满足 Docker for Mac 和 Docker for Windows apps 运行要求的、older Mac 和 Windows 系统上，快速安装和启动一套 Docker 环境的工具；

Toolbox 中包含如下 Docker 工具：
- **Docker Machine** for running `docker-machine` commands
- **Docker Engine** for running the `docker` commands
- **Docker Compose** for running the `docker-compose` commands
- **Kitematic**, the Docker GUI
- a **shell** preconfigured for a Docker command-line environment
- Oracle **VirtualBox**

> Docker Toolbox 使用的是 Oracle **Virtual Box**，而不是 **HyperKit**；因此允许不满足上面的系统要求；


## [Docker Machine](https://docs.docker.com/machine/overview/)

Docker Machine 的用处：
- 在 Mac 或 Windows 上安装和运行
- Provision 和 manage 多个远端 Docker hosts
- Provision Swarm clusters

Docker Machine is a tool that lets you install Docker Engine on virtual hosts, and manage the hosts with docker-machine commands. You can use Machine to create Docker hosts on your local Mac or Windows box, on your company network, in your data center, or on cloud providers like AWS or Digital Ocean.

Using docker-machine commands, you can start, inspect, stop, and restart a managed host, upgrade the Docker client and daemon, and configure a Docker client to talk to your host.

Machine 是 Docker v1.12 出现之前，在 Mac 或 Windows 上运行 Docker 的唯一方式；从 beta program 和 Docker v1.12 开始，Docker for Mac 和 Docker for Windows 已经作为 native apps 可用了，并作为更好的使用方式被推荐；

Docker Machine 具有如下两种广泛使用的用例：
- I have an older desktop system and want to run Docker on Mac or Windows
- I want to provision Docker hosts on remote systems


## [Docker Engine](https://docs.docker.com/engine/understanding-docker/#/what-is-docker-engine)

Docker Engine 属于 client-server 应用模型，主要由以下组件构成：
- A **server** which is a type of long-running program called a daemon process.
- A **REST API** which specifies interfaces that programs can use to talk to the daemon and instruct it what to do.
- A command line interface (CLI) **client**.

CLI 使用 Docker 提供的 REST API 来控制 Docker daemon 或与 Docker daemon 进行交互（通过脚本或直接 CLI 命令调用）     
许多其他 Docker 应用同样使用底层 API 和 CLI ；    

daemon 负责创建和管理各种 Docker 对象，如 **images**, **containers**, **networks** 和 **data volumes** ；


## [Docker Engine v.s. Docker Machine](https://docs.docker.com/machine/overview/#/what-s-the-difference-between-docker-engine-and-docker-machine)


## [Docker for Mac]()

使用 HyperKit ；


## [Docker for Mac v.s. Docker Toolbox](https://docs.docker.com/docker-for-mac/docker-toolbox/)
> 两者之间的影响和共存问题




