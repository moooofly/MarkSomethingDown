
本文针对 Docker for Mac 进行说明，以下内容整理自《[Getting Started with Docker for Mac](https://docs.docker.com/docker-for-mac/)》


----------


# 系统要求

>- Mac must be a 2010 or newer model, with Intel’s hardware support for **memory management unit (MMU) virtualization**; i.e., **Extended Page Tables (EPT)**
>- OS X **10.10.3** Yosemite or newer
>- At least **4GB of RAM**
>- VirtualBox prior to version 4.3.30 must NOT be installed (it is incompatible with Docker for Mac)


# 安装

两种方式：

- 通过 `brew cask install docker` 进行安装（推荐）；
- 登陆到 [Docker for Mac](https://www.docker.com/docker-mac) 官网下载 Docker.dmg 后安装；

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

# 测试

> 如果目标 image 无法在本地找到，那么 Docker 将默认从 [Docker Hub](https://hub.docker.com/) 上进行 pull 拉取；

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

测试步骤：

```shell
docker run -d -p 80:80 --name webserver nginx
docker ps
docker stop webserver       ## stop the container
docker start webserver      ## start the container
```

具体执行过程：

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

> ⚠️ 在一些情况下，很有可能你会希望将一些 image 保留下来（会占用磁盘空间），这样就不再需要再次从 Docker Hub 上进行 pull ；

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

- 可以用于在 Mac 或 Windows 上安装和运行 Docker
- 可以 Provision 和 manage 多个远端 Docker hosts
- 可以 Provision Swarm clusters

Docker Machine 是一种工具，允许你在 virtual hosts 上安装 Docker Engine，允许你通过  `docker-machine` 命令管理这些 hosts ；可以使用 Machine 在你本地 Mac 或 Windows box 中，在你的公司网络中，在你的数据中心中，以及在类似 AWS 或 Digital Ocean 的云提供商中创建 Docker hosts ；

通过 `docker-machine` 命令你可以 start, inspect, stop 以及 restart 一个受管控的 host ，升级指定的 Docker client 和 daemon ，以及配置指定的 Docker client 与你的 host 进行通信；

将 Machine 的 CLI 指向一处于运行状态的，受管控的 host 后，你就可以直接在那个 host 上运行 docker 命令了；例如，运行 `docker-machine env default` 以指向一个名为 default 的 host ，并按照屏幕上提示的信息完成 env 的设置，之后就可以运行 `docker ps` ，`docker run hello-world` 等指令了；

Machine 是 **Docker v1.12** 出现之前，在 Mac 或 Windows 上运行 Docker 的唯一方式；从 beta program 和 Docker v1.12 开始，Docker for Mac 和 Docker for Windows 已经作为 native apps 可用了，并作为更好的使用方式被推荐；


### 为什么要使用 Docker Machine

Docker Machine 令你能够对运行不同 Linux 版本的众多远端 Docker hosts 进行 provision ；     
Docker Machine 允许你将 Docker 运行在 older Mac 或 Windows 系统上；     

Docker Machine 在如下两种场景下被广泛使用：

- I have an older desktop system and want to run Docker on Mac or Windows

> If you work primarily on an older Mac or Windows laptop or desktop that doesn’t meet the requirements for the new Docker for Mac and Docker for Windows apps, then you need Docker Machine in order to “run Docker” (that is, Docker Engine) locally. Installing Docker Machine on a Mac or Windows box with the Docker Toolbox installer provisions a local virtual machine with Docker Engine, gives you the ability to connect it, and run docker commands.

- I want to provision Docker hosts on remote systems

> Docker Engine runs natively on Linux systems. If you have a Linux box as your primary system, and want to run docker commands, all you need to do is download and install Docker Engine. However, if you want an efficient way to provision multiple Docker hosts on a network, in the cloud or even locally, you need Docker Machine.
> 
> Whether your primary system is Mac, Windows, or Linux, you can install Docker Machine on it and use docker-machine commands to provision and manage large numbers of Docker hosts. It automatically creates hosts, installs Docker Engine on them, then configures the docker clients. Each managed host (”machine”) is the combination of a Docker host and a configured client.


## [Docker Engine](https://docs.docker.com/engine/understanding-docker/#/what-is-docker-engine)

Docker Engine 属于 client-server 应用模型，主要由以下组件构成：

- A **server** which is a type of long-running program called a daemon process.
- A **REST API** which specifies interfaces that programs can use to talk to the daemon and instruct it what to do.
- A command line interface (CLI) **client**  that talks to the daemon (through the REST API wrapper).

CLI 使用 Docker 提供的 REST API 来控制 Docker daemon 或与 Docker daemon 进行交互（通过脚本或直接 CLI 命令调用）     
许多其他 Docker 应用同样使用底层 API 和 CLI ；    

daemon 负责创建和管理各种 Docker 对象，如 **images**, **containers**, **networks** 和 **data volumes** ；


## [Docker Engine v.s. Docker Machine](https://docs.docker.com/machine/overview/#/what-s-the-difference-between-docker-engine-and-docker-machine)

当人们谈起 “Docker” 时，他们通常指的是 Docker Engine， 即由 Docker daemon, REST API 和 CLI 构成的 client-server 应用；Docker Engine 接受来自 CLI 的 docker 命令，例如 `docker run <image>`, `docker ps` 和 `docker images` 等；

![Docker Engine](https://docs.docker.com/machine/img/engine.png)

Docker Machine 是一种可以 provisioning 和 managing 你的 Dockerized hosts（运行了 Docker Engine 的 hosts）的工具；典型情况下，你会安装 Docker Machine 在你的本地系统上；Docker Machine 具有其自己的命令行客户端 `docker-machine` 和 Docker Engine 客户端 `docker` ；你可以使用 Machine 来安装 Docker Engine 到一个或多个虚拟系统中；这些虚拟系统可能是本地的（例如当你在 Mac 或 Windows 上的 VirtualBox 中基于 Machine 安装并运行 Docker Engine 时）或者是远端的（例如当你使用 Machine 在云提供商中对 Dockerized hosts 进行 provision 时）；在有些场景中，Dockerized hosts 会被看作受管理的 “machines” ；


![Docker Machine](https://docs.docker.com/machine/img/machine.png)


## [Docker for Mac](https://docs.docker.com/docker-for-mac/)

使用 HyperKit ；


## [Docker for Mac v.s. Docker Toolbox](https://docs.docker.com/docker-for-mac/docker-toolbox/)

> 两者之间的影响和共存问题


## [HyperKit](https://github.com/docker/HyperKit/)

HyperKit 是用于嵌入 hypervisor 能力到你的应用中的一种 toolkit ；其内部包含了基于 [xhyve](https://github.com/mist64/xhyve)/[bhyve](http://bhyve.org/) 提供的完整 hypervisor 功能，而这两者被专门优化用于提供轻量级 VM 和 container 部署；     
HyperKit 被设计成可以同更高层组件进行交互，例如 [VPNKit](https://github.com/docker/vpnkit) 和 [DataKit](https://github.com/docker/datakit) ；     
HyperKit 当前仅支持采用了 Hypervisor.framework 框架的 Mac OS X ；     
HyperKit 是 Docker For Mac 的核心组件；


### 系统要求

- OS X 10.10.3 Yosemite or later
- a 2010 or later Mac (i.e. a CPU that supports EPT)

