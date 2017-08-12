# Mac 上基于 Homebrew 安装 Beats

## [Installing Beats](https://www.elastic.co/guide/en/beats/libbeat/5.4/installing-beats.html)

在成功安装配置 **Elastic stack** 后，就可以开始 Beat 部分的处理了；

每一种 Beat 都是独立的、可安装的产品；针对不同类型的 Beat 请参考相应的文档：

- [Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-getting-started.html)
- [Metricbeat](https://www.elastic.co/guide/en/beats/metricbeat/5.4/metricbeat-getting-started.html)
- [Filebeat](https://www.elastic.co/guide/en/beats/filebeat/5.4/filebeat-getting-started.html)
- [Winlogbeat](https://www.elastic.co/guide/en/beats/winlogbeat/5.4/winlogbeat-getting-started.html)


----------

## Packetbeat 安装

### [Getting Started With Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-getting-started.html) 

动嘴不如动手，使用诸如 Packetbeat 这样到网络抓包分析工具的最好方式就是实践；

使用 Packetbeat 前，需要搞定的内容（即 Elastic Stack）：

- **Elasticsearch** 进行数据的存储和索引查询；
- **Kibana** 提供 UI 供查询和展示；
- **Logstash** 用于插入数据到 Elasticsearch 中（可选）；

> Elastic Stack 安装详见 [Getting Started with Beats and the Elastic Stack](https://www.elastic.co/guide/en/beats/libbeat/5.4/getting-started.html) ；

在完成 Elastic Stack 安装后，可以参考下面的内容进行 Packetbeat 的安装、配置和运行：


- [Step 1: Installing Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-installation.html)
- [Step 2: Configuring Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/configuring-packetbeat.html)
- [Step 3: Loading the Index Template in Elasticsearch](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-template.html)
- [Step 4: Starting Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-starting.html)
- [Step 5: Loading Sample Kibana Dashboards](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-sample-dashboards.html)
- [Command Line Options](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-command.html)
- [Directory Layout](https://www.elastic.co/guide/en/beats/packetbeat/5.4/directory-layout.html)


----------


### [Step 1: Installing Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-installation.html)

请使用合适的命令下载和安装 Packetbeat 到你的应用服务器上；

> 如果你使用 Apt 或 Yum ，你可以 [install Packetbeat from our repositories](https://www.elastic.co/guide/en/beats/libbeat/5.4/setup-repositories.html) ，这样更容易进行最新版本的升级；


### [Step 2: Configuring Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/configuring-packetbeat.html)

针对 Packetbeat 的配置需要对配置文件进行编辑；对于 rpm 和 deb 安装方式来说，配置文件位于 `/etc/packetbeat/packetbeat.yml` ；对于 mac 和 win 来说，请在你解压归档文件的相应目录下查找；

配置 Packetbeat 方法如下：

1. 选择进行通信捕获的网络接口
  - On **Linux**: Packetbeat supports capturing all messages sent or received by the server on which Packetbeat is installed. For this, use `any` as the device:
    ```
    packetbeat.interfaces.device: any
    ```
    - On **OS X**, capturing from the `any` device doesn’t work. You would typically use either `lo0` or `en0` depending on which traffic you want to capture.

2. 在 `protocols` 段，配置各种 ports 信息以允许 Packetbeat 捕获相应的协议包；你过你使用了非标准 ports ，也要在这里进行添加；否则，仅使用默认值就足够了；
3. 设置 IP 地址和 port 以便 Packetbeat 和 Elasticsearch 进行连接；

为了对配置文件进行测试，先切换到 Packetbeat 二进制文件所在目录，然后使用如下选项在前台运行 Packetbeat 可执行程序： `sudo ./packetbeat -configtest -e` ；请确保配置文件能够被 Packetbeat 访问到（可以参考 [Directory Layout](https://www.elastic.co/guide/en/beats/packetbeat/5.4/directory-layout.html)）；若基于 DEB 或 RPM 包进行的安装，则运行 `sudo ./packetbeat.sh -configtest -e` ；

```shell
➜  ~ packetbeat -configtest -e
2016/12/01 08:09:59.828010 beat.go:264: INFO Home path: [/usr/local/Cellar/packetbeat/5.0.1] Config path: [/usr/local/etc/packetbeat] Data path: [/usr/local/var/packetbeat] Logs path: [/usr/local/var/log/packetbeat]
2016/12/01 08:09:59.828336 logp.go:219: INFO Metrics logging every 30s
2016/12/01 08:09:59.828345 beat.go:174: INFO Setup Beat: packetbeat; Version: 5.0.1
2016/12/01 08:09:59.828701 output.go:167: INFO Loading template enabled. Reading template file: /usr/local/etc/packetbeat/packetbeat.template.json
2016/12/01 08:09:59.831385 output.go:178: INFO Loading template enabled for Elasticsearch 2.x. Reading template file: /usr/local/etc/packetbeat/packetbeat.template-es2x.json
2016/12/01 08:09:59.833392 client.go:107: INFO Elasticsearch url: http://localhost:9200
2016/12/01 08:09:59.833441 outputs.go:106: INFO Activated elasticsearch as output plugin.
2016/12/01 08:09:59.833528 publish.go:291: INFO Publisher name: sunfeideMacBook-Pro.local
2016/12/01 08:09:59.833639 async.go:63: INFO Flush Interval set to: 1s
2016/12/01 08:09:59.833650 async.go:64: INFO Max Bulk Size set to: 50
2016/12/01 08:09:59.833799 procs.go:91: INFO Process matching disabled
2016/12/01 08:09:59.833871 protos.go:89: INFO registered protocol plugin: thrift
2016/12/01 08:09:59.833878 protos.go:89: INFO registered protocol plugin: cassandra
2016/12/01 08:09:59.833881 protos.go:89: INFO registered protocol plugin: mysql
2016/12/01 08:09:59.833883 protos.go:89: INFO registered protocol plugin: http
2016/12/01 08:09:59.833886 protos.go:89: INFO registered protocol plugin: memcache
2016/12/01 08:09:59.833889 protos.go:89: INFO registered protocol plugin: mongodb
2016/12/01 08:09:59.833892 protos.go:89: INFO registered protocol plugin: nfs
2016/12/01 08:09:59.833894 protos.go:89: INFO registered protocol plugin: pgsql
2016/12/01 08:09:59.833897 protos.go:89: INFO registered protocol plugin: redis
2016/12/01 08:09:59.833900 protos.go:89: INFO registered protocol plugin: amqp
2016/12/01 08:09:59.833902 protos.go:89: INFO registered protocol plugin: dns
Config OK
➜  ~
```


### [Step 3: Loading the Index Template in Elasticsearch](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-template.html)

在 Elasticsearch 中，**index template** 定义了如何分析 field 的 settings 和 mappings 内容；

> In Elasticsearch, [index templates](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/Elasticsearch%20%E4%B8%AD%E7%9A%84%20Index%20Pattern%20%E5%92%8C%20Index%20Template.md) are used to define **settings** and **mappings** that determine how fields should be analyzed.

对于 Packetbeat 来说，推荐使用的 **index template** 会伴随 Packetbeat 一起被安装；如果你能够接受 `packetbeat.yml` 中针对 template 加载的默认配置，那么 Packetbeat 已经能做到在成功连接 Elasticsearch 之后自动加载 template 的功能了（参考上面的打印输出信息）；如果 template 已经存在（于 Elasticsearch 中），则不会被覆盖，除非你通过配置要求 Packetbeat 进行覆盖；

> The **recommended** index template file for Packetbeat is installed by the Packetbeat packages. If you accept the default configuration for template loading in the `packetbeat.yml` config file, Packetbeat loads the template automatically after successfully connecting to Elasticsearch. If the template already exists, it’s not overwritten unless you configure Packetbeat to do so.

如果你想要**去使能 template 自动加载功能**，或者打算**加载自己的 template**，则需要变更 Packetbeat 配置文件中的 template 加载设置；当去使能 template 自动加载功能后，你需要手动加载 template ；

> If you want to disable automatic template loading, or you want to load your own template, you can change the settings for template loading in the Packetbeat configuration file. If you choose to disable automatic template loading, you need to load the template manually. For more information, see:
> 
> - **Configuring Template Loading** - supported for Elasticsearch output only
> - **Loading the Template Manually** - required for Logstash output

#### Configuring Template Loading

- 推荐的默认 template 文件为 `packetbeat.template.json` ，要求 Elasticsearch 这个 output 一定要被启用；
- 加载不同 template 时需要变更 `template.name` 和 `template.path` 的配置内容；

> By default, Packetbeat automatically loads the recommended template file, `packetbeat.template.json`, if Elasticsearch output is enabled. You can configure packetbeat to load a different template by adjusting the `template.name` and `template.path` options in `packetbeat.yml` file:
> 
> ```
> output.elasticsearch:
>   hosts: ["localhost:9200"]
>   template.name: "packetbeat"
>   template.path: "packetbeat.template.json"
>   template.overwrite: false
> ```

- 默认情况下，已存在 template 不会被覆盖；
- 若不想要自动加载功能，则直接将 template 相关配置全部注释掉；
- 使用 Logstash output 时不支持自动加载功能；

> By default, if a template already exists in the index, it is not overwritten. To overwrite an existing template, set `template.overwrite: true` in the configuration file.
> 
> To disable automatic template loading, comment out the template part under the Elasticsearch output.
> 
> The options for auto loading the template are not supported if you are using the Logstash output.


#### Loading the Template Manually

手动加载 template 的办法（注意：`_template` 为关键字，`packetbeat` 为目标 template 名字）；

第一步：添加名为 `packetbeat` 的 `index template` ；

> If you disable automatic template loading, you need to run the following command to load the template:
> 
> ```
> curl -H 'Content-Type: application/json' -XPUT 'http://localhost:9200/_template/packetbeat' -d@packetbeat.template.json
> ```
> 
> where `localhost:9200` is the IP and port where Elasticsearch is listening.

第二步，删除名为 `packetbeat-*` 的 index pattern ，以强制 Kibana 基于 `index template` 内容重新索引所有 documents ；

> If you’ve already used Packetbeat to index data into Elasticsearch, the index may contain old documents. After you load the index template, you can delete the old documents from packetbeat-* to force Kibana to look at the newest documents. Use this command:
> 
> ```
> curl -XDELETE 'http://localhost:9200/packetbeat-*'
> ```


### [Step 4: Starting Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-starting.html)

可以通过如下命令运行 Packetbeat ：

```shell
sudo ./packetbeat -e -c packetbeat.yml -d "publish"
```

#### Testing the Packetbeat Installation

Packetbeat 已经就绪了抓包了；可以通过如下简单的 HTTP 请求进行测试：

```shell
curl http://www.elastic.co/ > /dev/null
```

之后，可以通过如下命令在 Elasticsearch 中确认相应的数据是否存在：

```shell
curl -XGET 'http://localhost:9200/packetbeat-*/_search?pretty'
```

请确保使用正确的 Elasticsearch 实例地址信息替代 `localhost:9200` ；上述命令将会获取到当前 HTTP transaction 到相关信息；


### [Step 5: Loading Sample Kibana Dashboards](https://www.elastic.co/guide/en/beats/packetbeat/5.4/packetbeat-sample-dashboards.html)

为方便针对抓包内容进行应用性能分析，官方提供了 Packetbeat dashboards 样例；但官方提供 dashboards 是作为样例供参考的，建议用户自己[定制](https://www.elastic.co/guide/en/kibana/current/dashboard.html)相应的 dashboards 以满足需求；

![packetbeat-statistics](https://www.elastic.co/guide/en/beats/packetbeat/5.4/images/packetbeat-statistics.png "packetbeat-statistics")

#### Importing the Dashboards

与 Packetbeat 一同打包的 `scripts/import_dashboards` 脚本（基于 go 程序编译出来的），会将预定义的 dashboards 导入到 Elasticsearch 中，以便在 kibana 中提供可视化展示（visualizations），还会为 Packetbeat 配置搜索（search）能力；该脚本会创建一个用于 Packetbeat 的、名为 `packetbeat-*` 的 **index pattern** ；

下面将描述如何导入（用于） Packetbeat 的 (Kibana) dashboards ；

你可能想要导入不止一个 dashboards 以方便不同 Beat 使用，或者可能会想要对导入选项进行定制；在 Beats Platform Reference 的 [Importing Existing Beat Dashboards](https://www.elastic.co/guide/en/beats/libbeat/5.4/index.html) 章节能够看到完整的命令行选项信息；

为了导入用于 Packetbeat 的 Kibana dashboards ，需要执行如下内容；

在 Packetbeat 的安装目录中，运行 `import_dashboards` 脚本；
```shell
./scripts/import_dashboards
```

实际执行情况如下（这里使用官方提供的 mac 版本打包文件进行的演示）

```
➜  WGET cd packetbeat-5.0.2-darwin-x86_64
➜  packetbeat-5.0.2-darwin-x86_64 ll
total 22720
-rwxr-xr-x  1 sunfei  staff    11M 11 24 18:00 packetbeat
-rw-r--r--  1 sunfei  staff    38K 11 24 18:00 packetbeat.full.yml
-rw-r--r--  1 sunfei  staff    44K 11 24 18:00 packetbeat.template-es2x.json
-rw-r--r--  1 sunfei  staff    36K 11 24 18:00 packetbeat.template.json
-rw-r--r--  1 sunfei  staff   5.2K 11 24 18:00 packetbeat.yml
drwxr-xr-x  4 sunfei  staff   136B 12  1 17:04 scripts
➜  packetbeat-5.0.2-darwin-x86_64 cd scripts
➜  scripts ll
total 23656
-rwxr-xr-x  1 sunfei  staff    12M 11 24 18:00 import_dashboards
-rwxr-xr-x  1 sunfei  staff    14K 11 24 17:59 migrate_beat_config_1_x_to_5_0.py
➜  scripts ./import_dashboards
Create temporary directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597
Downloading https://artifacts.elastic.co/downloads/beats/beats-dashboards/beats-dashboards-5.0.2.zip
Unzip archive  /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597
/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/filebeat
/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/metricbeat
/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat
Import directory  /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat/index-pattern
Import index to /.kibana/index-pattern/packetbeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat/index-pattern/packetbeat.json
...
Import search  /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat/search/Thrift-transactions.json
Import vizualization  /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat/visualization/Top-Thrift-RPC-calls-with-errors.json
Import search  /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/packetbeat/search/Thrift-errors.json
/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp586211597/beats-dashboards-5.0.2/winlogbeat
➜  scripts
```

默认情况下，该脚本认为 Elasticsearch 运行在 `127.0.0.1:9200` ；可以使用 `-es` 选项指定一个不同的地址，例如：
```shell
./scripts/import_dashboards -es http://192.168.33.60:9200
```

可以通过 `-user` 选项指定用于 Elasticsearch 鉴权的用户名和密码；例如：

```shell
./scripts/import_dashboards -es https://xyz.found.io -user user -pass password 
```
> Specify the username and password as options.
```shell
./scripts/import_dashboards -es https://xyz.found.io -user admin -pass $(cat ~/pass-file) 
```
> Use a file to avoid polluting the bash history with the password.

#### Opening the Dashboards in Kibana

在导入了 dashboards 之后，可以在浏览器上访问 5601 访问 Kibana 的 web 接口；例如 `http://127.0.0.1:5601` ；

在 **Discover** 页面上，请确保名为 `packetbeat-*` 的预定义 **index pattern** 被选中，以便查看 Packetbeat 数据信息；

![](https://www.elastic.co/guide/en/beats/packetbeat/5.4/images/kibana-created-indexes.png)

如果 Kibana 显示 "*No default index pattern*" 警告，那么你就必须 select 或 create 一个 **index pattern** 才能继续使用；为了解决这个问题，可以选中预定义的 **index pattern** ，即 `packetbeat-*`，并将其设置为默认值；

为了打开已经加载到 Elasticsearch 中的 dashboards ，只需到 Dashboard 页面上选择你想要打开的 dashboard 即可；

![](https://www.elastic.co/guide/en/beats/packetbeat/5.4/images/kibana-navigation-vis.png)


----------


## Command Line Options

The following command line options are available for Packetbeat. To use these options, you need to start Packetbeat in the foreground.

> Run `./packetbeat -h` to see the full list of options from the command line.


### PacketBeat 专有选项

Packetbeat 支持到命令行选项：

- **`-I <file>`**

从事先保存好的 pcap 文件中读取抓包数据到 Packetbeat 中；

- **`-O`**

通过 *Enter* 控制每次读取一个 packet ； 

- **`-devices`**

查看当前主机上可以进行 sniffing 的设备列表；

- **`-dump <file>`**

将捕获到的网络数据包保存到目标文件中；

- **`-l <n>`**

反复读取 pcap 文件 `n` 次；和 `-I` 选项组合使用；若想实现无限循环效果，使用 *0* 值；

- **`-t`**

Read the packets from the pcap file as fast as possible without sleeping. 和 `-I` 选项组合使用；

- **`-waitstop <n>`**

退出前，等待额外 `n` 秒时间（需要确认）；

### 其它选项

以下 libbeat 支持的选项同样可用于 Packetbeat ：

- **`-E <setting>=<value>`**

该选项可以用于覆盖 config 文件中相应的设置内容；例如：

```shell
sudo ./packetbeat -c packetbeat.yml -E name=mybeat
```

- **`-N`**

Disable the publishing of events to the defined output. 

- **`-c <file>`**

为 Beat 指定配置文件位置；

- **`-configtest`**

测试配置文件内容的正确性，测完后退出；

- **`-cpuprofile <output file>`**

Write CPU profile data to the specified file. This option is useful for troubleshooting the Beat.

- **`-d <selectors>`**

使能针对选项 selectors 的调试功能；针对 selectors 设置，可以指定逗号分隔的 selector 列表；或者可以使用 `-d "*"` 使能针对全部 selectors 的调试功能；例如 `-d "publish"` 将会显示全部和 "publish" 相关的消息；

- **`-e`**

输出内容写到 stderr ，同时去使能 syslog/file 输出；

- **`-httpprof [<host>]:<port>`**

Start http server for profiling. This option is useful for troubleshooting and profiling the Beat.

- **`-memprofile <output file>`**

Write memory profile data to the specified output file. This option is useful for troubleshooting the Beat.

- **`-path.config`**

设置配置文件的默认位置（例如 Elasticsearch template 文件）；

- **`-path.data`**

设置数据文件位置的默认位置；

- **`-path.home`**

设置其他文件的默认位置；

- **`-path.logs`**

设置日志文件的默认位置；

- **`-v`**

使能 verbose 输出，以便显示 INFO 级别的消息；

- **`-version`**

Display the Beat version and exit.