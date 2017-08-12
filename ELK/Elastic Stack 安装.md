# Elastic Stack 安装

那么问题来了：**如何安装 Elastic stack ？**

## [Getting Started with Beats and the Elastic Stack](https://www.elastic.co/guide/en/beats/libbeat/5.0/getting-started.html)

Looking for an "**ELK tutorial**" that shows how to set up the Elastic stack for Beats? You’ve come to the right place. The topics in this section describe how to install and configure the Elastic stack for Beats.

Beats setup 通常由如下内容构成：

- **Elasticsearch** 用于提供存储和索引，详见 [Installing Elasticsearch](https://www.elastic.co/guide/en/beats/libbeat/5.0/elasticsearch-installation.html) ；
- **Logstash** (optional) 用于往 Elasticsearch 中插入数据，详见 [Installing Logstash](https://www.elastic.co/guide/en/beats/libbeat/5.0/logstash-installation.html) ；
- **Kibana** 用于提供 UI ，详见 [Installing Kibana](https://www.elastic.co/guide/en/beats/libbeat/5.0/kibana-installation.html) ；
- 一种或多种 **Beats** 用于在目标服务器上捕获操作数据内容（operational data），详见 [Installing Beats](https://www.elastic.co/guide/en/beats/libbeat/5.0/installing-beats.html) ；
- **Kibana dashboards** 用于数据可视化；

See the [Elastic Support Matrix](https://www.elastic.co/support/matrix) for information about supported operating systems and product compatibility.

> Note
>> To get started, you can install Elasticsearch and Kibana on a single VM or even on your laptop. The only condition is that the machine must be accessible from the servers you want to monitor. As you add more Beats and your traffic grows, you’ll want to replace the single Elasticsearch instance with a cluster. You’ll probably also want to automate the installation process.


----------

> 以下内容均基于 Mac 版本进行说明；

## Mac 上基于 Homebrew 安装 elasticsearch

### [Installing Elasticsearch](https://www.elastic.co/guide/en/beats/libbeat/5.0/elasticsearch-installation.html)

Elasticsearch is a real-time, distributed storage, search, and analytics engine. It can be used for many purposes, but one context where it excels is indexing streams of semi-structured data, such as logs or decoded network packets.

The binary packages of Elasticsearch have only one dependency: **Java**. The minimum supported version is **Java 8**. To download and install Elasticsearch, use the commands that work with your system

> 以下内容为基于 brew 在 Mac 上安装的过程；

查找资源

```shell
➜  ~ brew search elasticsearch
elasticsearch                                      elasticsearch@1.7
homebrew/versions/elasticsearch24
➜  ~
```

查看 elasticsearch 相关信息（可以看到 Required: java >= 1.8 ✘ ）

```shell
➜  ~ brew info elasticsearch
elasticsearch: stable 5.0.2, HEAD
Distributed search & analytics engine
https://www.elastic.co/products/elasticsearch
Not installed
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/elasticsearch.rb
==> Requirements
Required: java >= 1.8 ✘
==> Caveats
Data:    /usr/local/var/elasticsearch/elasticsearch_sunfei/
Logs:    /usr/local/var/log/elasticsearch/elasticsearch_sunfei.log
Plugins: /usr/local/Cellar/elasticsearch/5.0.2/libexec/plugins/
Config:  /usr/local/etc/elasticsearch/
plugin script: /usr/local/Cellar/elasticsearch/5.0.2/libexec/bin/plugin

To have launchd start elasticsearch now and restart at login:
  brew services start elasticsearch
Or, if you don't want/need a background service you can just run:
  elasticsearch
➜  ~
```

尝试安装 elasticsearch（提示 Java 安装需求不满足）

```shell
➜  ~ brew install elasticsearch
elasticsearch: Java 1.8+ is required to install this formula.JavaRequirement unsatisfied!

You can install with Homebrew-Cask:
  brew cask install java

You can download from:
  http://www.oracle.com/technetwork/java/javase/downloads/index.html
Error: An unsatisfied requirement failed this build.
➜  ~
```

安装 java

```shell
➜  ~ brew cask install java
==> Caveats
This Cask makes minor modifications to the JRE to prevent issues with
packaged applications, as discussed here:

  https://bugs.eclipse.org/bugs/show_bug.cgi?id=411361

If your Java application still asks for JRE installation, you might need
to reboot or logout/login.

Installing this Cask means you have AGREED to the Oracle Binary Code
License Agreement for Java SE at

  https://www.oracle.com/technetwork/java/javase/terms/license/index.html

==> Downloading http://download.oracle.com/otn-pub/java/jdk/8u112-b16/jdk-8u112-macosx-x64.dmg
######################################################################## 100.0%
==> Verifying checksum for Cask java
==> Running installer for java; your password may be necessary.
==> Package installers may write to any location; options such as --appdir are ignored.
Password:
==> installer: Package name is JDK 8 Update 112
==> installer: Installing at base path /
==> installer: The install was successful.
🍺  java was successfully installed!
➜  ~
```

确认 Java 安装成功

```shell
➜  ~ java -version
java version "1.8.0_112"
Java(TM) SE Runtime Environment (build 1.8.0_112-b16)
Java HotSpot(TM) 64-Bit Server VM (build 25.112-b16, mixed mode)
➜  ~
```

重新安装 elasticsearch

```shell
➜  ~ brew install elasticsearch
==> Using the sandbox
==> Downloading https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.0.2.tar.gz
######################################################################## 100.0%
==> Caveats
Data:    /usr/local/var/elasticsearch/elasticsearch_sunfei/
Logs:    /usr/local/var/log/elasticsearch/elasticsearch_sunfei.log
Plugins: /usr/local/Cellar/elasticsearch/5.0.2/libexec/plugins/
Config:  /usr/local/etc/elasticsearch/
plugin script: /usr/local/Cellar/elasticsearch/5.0.2/libexec/bin/plugin

To have launchd start elasticsearch now and restart at login:
  brew services start elasticsearch
Or, if you don't want/need a background service you can just run:
  elasticsearch
==> Summary
🍺  /usr/local/Cellar/elasticsearch/5.0.2: 98 files, 34.8M, built in 1 minute 1 second
➜  ~
```

测试通过

```shell
➜  ~ elasticsearch &
[1] 94255
...
➜  ~ curl http://127.0.0.1:9200
{
  "name" : "EofoeDw",
  "cluster_name" : "elasticsearch_sunfei",
  "cluster_uuid" : "70NsT860TgOwq_yVKuuX5Q",
  "version" : {
    "number" : "5.0.2",
    "build_hash" : "f6b4951",
    "build_date" : "2016-11-24T10:07:18.101Z",
    "build_snapshot" : false,
    "lucene_version" : "6.2.1"
  },
  "tagline" : "You Know, for Search"
}
➜  ~
```

----------

## Mac 上基于 Homebrew 安装 Logstash

### [Installing Logstash (Optional)](https://www.elastic.co/guide/en/beats/libbeat/5.0/logstash-installation.html)

Beats platform 的最简架构为 Elasticsearch + Kibana + Beats （因此理论上可以不安装使用 Logstash）；

> The simplest architecture for the Beats platform setup consists of one or more **Beats**, **Elasticsearch**, and **Kibana**. This architecture is easy to get started with and sufficient for networks with low traffic. It also uses the minimum amount of servers: a single machine running Elasticsearch and Kibana. The Beats insert the transactions directly into the Elasticsearch instance.

Logstash 可以提供额外的数据缓冲和处理功能；当然，你也可以使用 Redis 或 RabbitMQ 干类似的事情；

> If you want to perform additional processing or buffering on the data, however, you’ll want to install Logstash.
> 
> An important advantage to this approach is that you can use Logstash to modify the data captured by Beats in any way you like. You can also use Logstash’s many output plugins to integrate with other systems.

下载安装 Logstash 的步骤如下：

```shell
➜  ~ brew search logstash
logstash
homebrew/versions/logstash24
➜  ~
➜  ~
➜  ~ brew info logstash
logstash: stable 5.0.2, HEAD
Tool for managing events and logs
https://www.elastic.co/products/logstash
Not installed
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/logstash.rb
==> Requirements
Required: java >= 1.8 ✔
==> Caveats
Please read the getting started guide located at:
  https://www.elastic.co/guide/en/logstash/current/getting-started-with-logstash.html
➜  ~
➜  ~
➜  ~ brew cask search logstash
No Cask found for "logstash".
➜  ~
➜  ~
➜  ~ brew install logstash
Updating Homebrew...
==> Auto-updated Homebrew!
Updated 1 tap (homebrew/core).
==> Updated Formulae
pandoc

==> Using the sandbox
==> Downloading https://artifacts.elastic.co/downloads/logstash/logstash-5.0.2.tar.gz
######################################################################## 100.0%
==> Caveats
Please read the getting started guide located at:
  https://www.elastic.co/guide/en/logstash/current/getting-started-with-logstash.html
==> Summary
🍺  /usr/local/Cellar/logstash/5.0.2: 10,424 files, 189.9M, built in 1 minute 19 seconds
➜  ~
```

### Setting Up Logstash

略

### Updating the Beats Input Plugin for Logstash

略

### Starting Logstash

略

----------

## Mac 上基于 Homebrew 安装 Kibana

### [Installing Kibana](https://www.elastic.co/guide/en/beats/libbeat/5.0/kibana-installation.html)

Kibana 基于 Elasticsearch 中的数据进行可视化展示；

> Kibana is a visualization application that gets its data from Elasticsearch. It provides a customizable and user-friendly UI in which you can combine various widget types to create your own dashboards. The dashboards can be easily saved, shared, and linked.

建议 Kibana 和 Elasticsearch 装在一台机器上（方便）；

> For getting started, we recommend installing Kibana on the same server as Elasticsearch, but it is not required. If you install the products on different servers, you’ll need to change the URL (IP:PORT) of the Elasticsearch server in the Kibana configuration file, `config/kibana.yml`, before starting Kibana.

使用如下命令下载并运行 Kibana ：

```shell
➜  ~ brew search kibana
kibana
homebrew/versions/kibana41
homebrew/versions/kibana44
➜  ~
➜  ~
➜  ~ brew info kibana
kibana: stable 5.0.1 (bottled), HEAD
Analytics and search dashboard for Elasticsearch
https://www.elastic.co/products/kibana
Not installed
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/kibana.rb
==> Requirements
Build: git ✔
==> Caveats
Config: /usr/local/etc/kibana/
If you wish to preserve your plugins upon upgrade, make a copy of
/usr/local/opt/kibana/plugins before upgrading, and copy it into the
new keg location after upgrading.

To have launchd start kibana now and restart at login:
  brew services start kibana
Or, if you don't want/need a background service you can just run:
  kibana
➜  ~
➜  ~
➜  ~ brew cask search kibana
No Cask found for "kibana".
➜  ~
➜  ~
➜  ~ brew install kibana
Updating Homebrew...
==> Downloading https://homebrew.bintray.com/bottles/kibana-5.0.1.el_capitan.bottle.1.tar.gz
######################################################################## 100.0%
==> Pouring kibana-5.0.1.el_capitan.bottle.1.tar.gz
==> Using the sandbox
==> Caveats
Config: /usr/local/etc/kibana/
If you wish to preserve your plugins upon upgrade, make a copy of
/usr/local/opt/kibana/plugins before upgrading, and copy it into the
new keg location after upgrading.

To have launchd start kibana now and restart at login:
  brew services start kibana
Or, if you don't want/need a background service you can just run:
  kibana
==> Summary
🍺  /usr/local/Cellar/kibana/5.0.1: 17,514 files, 144M
➜  ~
```

### Launching the Kibana Web Interface

为了启动 Kibana 的 web 接口，只需要在浏览器上访问 5601 端口；例如 http://127.0.0.1:5601 ；

关于 Kibana 的更多细节信息，详见 [Kibana 用户手册](https://www.elastic.co/guide/en/kibana/current/index.html) ；