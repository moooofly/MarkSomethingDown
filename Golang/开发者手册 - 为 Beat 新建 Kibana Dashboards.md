# 开发者手册 - 为 Beat 新建 Kibana Dashboards

> 原文地址：[Developer Guide: Creating New Kibana Dashboards for a Beat](https://www.elastic.co/guide/en/beats/libbeat/5.4/new-dashboards.html)

从 Beats 5.0.0 开始，Kibana dashboards 不再作为 Beat package 的一部分发布了；而是以单独的 `beats-dashboards` package 进行发布；

在 Beats 开发过程中，你可能想要添加或修改 dashboards ，现在可以通过 `import_dashboards` 将现有 Beat 的 dashboards 内容导入到 Kibana 中；

导入完成后，Kibana 会将 **dashboards** 及其对应的所有依赖项，即 **visualizations**, **searches** 和 **index patterns** ，导入到 Elasticsearch 的一个特殊 index 中；默认情况下，这个特殊 index（的名字）为 `.kibana` ，但允许你指定一个不同名字的 index ；

当你完成了对 Kibana dashboards 的变更后，你可以使用 `export_dashboards` 脚本导出对应的 dashboards 以及其所有依赖项，到一个本地目录中；

为了确保 dashboards 兼容于最新版本的 Kibana 和 Elasticsearch ，建议使用 [`beats/testing/environments`](https://github.com/elastic/beats/tree/master/testing/environments) 虚拟环境测试针对 Kibana dashboards 的 import, create 和 export ；

如下 topics 提供了关于 import 和 Beats dashboards 使用的详细信息：

- Importing Existing Beat Dashboards
- Building Your Own Beat Dashboards
- Generating the Beat Index Pattern
- Exporting New and Modified Beat Dashboards
- Archiving Your Beat Dashboards
- Sharing Your Beat Dashboards


> ⚠️ 如下内容的具体含义：
> 
> - `.kibana`
> - dashboards 和 visualizations, searches 和 index patterns 的关系


## Importing Existing Beat Dashboards

可以使用 `import_dashboards` 脚本为 Beat 导入所有 dashboards 和 **index pattern** ，包括相关的 dependencies ，例如 **visualizations** 和 **searches** ；

`import_dashboards` 脚本在 [`beats/libbeat/dashboards`](https://github.com/elastic/beats/tree/5.4/libbeat/dashboards) 中提供；若使用 Beat package 则在 `scripts` 目录下；

> ⚠️ master 分支内容有所不同！

存在很多常见的、需要导入 dashboards 的场景：

- Users who are **getting started with Beats** may want to import dashboards and/or the index pattern for a single Beat.
- Community Beats developers may want to **import dashboards for development to use** as a starting point for new dashboards.

### Import Dashboards and/or the Index Pattern for a Single Beat

Using the `import_dashboards` script from the Beat package, you can import the **dashboards** and the **index pattern** to Elasticsearch running on localhost for a single Beat (eg. Metricbeat):

- 基于**本地目录**内容导入：

```
./scripts/import_dashboards -dir kibana/metricbeat
```

- 基于本地 **zip** 归档文件导入：

```
./scripts/import_dashboards -file metricbeat-dashboards-1.1.zip
```

- 直接下载官方提供的 **zip** 归档文件导入，例如 https://artifacts.elastic.co/downloads/beats/beats-dashboards/beats-dashboards-5.4.0.zip ，具体版本取决于代码中定义的版本信息：

```
./scripts/import_dashboards
```

- 直接指定 github 上的 **zip** 归档文件地址进行下载导入：

```
./scripts/import_dashboards -url https://github.com/monicasarbu/metricbeat-dashboards/archive/v1.1.zip
```

如果你不指定 archive 的位置，则默认使用官方 zip 归档文件，其中包含了 official Beats 提供的 **index pattern** 和 **dashboards** ；

- 针对单个 Beat (例如 Metricbeat) 仅导入 **index-pattern** 时使用：

```
./scripts/import_dashboards -only-index
```

- 针对单个 Beat (例如 Metricbeat) 仅导入 **dashboards** 以及 **visualizations** 和 **searches** 时使用：

```
./scripts/import_dashboards -only-dashboards
```

### Import Dashboards for Development

对于 Beats 开发者来说，直接在 [`beats/libbeat/dashboards`](https://github.com/elastic/beats/tree/5.4/libbeat/dashboards) 目录下运行 `import_dashboards` 脚本更为简单；但首先需要自行编译出该脚本：

```
cd beats/libbeat/dashboards
make
```

> ⚠️ master 分支内容有所不同！

之后，你就可以指定 `-beat` 导入上面提及的各种文件了；

例如，通过如下命令可以导入 Metricbeat 的 dashboards 和 visualizations, searches ，以及 Metricbeat index pattern ：

```
beats/libbeat/dashboards/import_dashboards -beat metricbeat
```

> 执行上述命令若提示 "zip: not a valid zip file" 之类的错误，则可以改为执行如下命令
> 
> ```
> ./import_dashboards -url https://artifacts.elastic.co/downloads/beats/beats-dashboards/beats-dashboards-5.3.2.zip
> ```

对于上面的例子，你必须指定 `-beat metricbeat` ，否则该脚本会导入全部 Beats 的 dashboards ；

你还可以利用 Beat 的 github 仓库中提供的 Makefile 文件导入 dashboards ；如果 Elasticsearch 运行在本地，则你可以在 Beat 仓库内（合适位置）运行如下命令：

```
make import-dashboards
```

如果 Elasticsearch 运行在其他主机上，你可以使用 `ES_URL` 变量：

```
ES_URL="http://192.168.3.206:9200" make import-dashboards
```


### Command Line Options

`import_dashboards` 脚本接受如下命令行选项；

```
./import_dashboards -h
```

- **`-beat <beatname>`**
The Beat name. The Beat name is required when importing from a zip archive. When using `import_dashboards` from the Beat package, this option is set automatically with the name of the Beat. **When running the script from source, the default value is ""**, so you need to set this option in order to install the `index pattern` and the `dashboards` for a single Beat. Otherwise the script imports the `index pattern` and the `dashboards` for all Beats.

- **`-cacert <certificate_authority>`**
The `Certificate Authority` to use for server verification.

- **`-cert <client_certificate>`**
The `certificate` to use for SSL client authentication. The certificate **must be** in `PEM` format.

- **`-dir <local_dir>`**
Local directory that contains the subdirectories: `dashboard`, `visualization`, `search`, and `index-pattern`. The default value is the current directory.

- **`-es <elasticsearch_url>`**
The Elasticsearch URL. The default value is `http://localhost:9200`.

- **`-file <local_archive>`**
**Local zip archive** with the dashboards. The archive can contain Kibana dashboards for a single Beat or for multiple Beats.

- **`-i <elasticsearch_index>`**
You should **only use** this option if you want **to change the `index pattern` name** that’s used by default. For example, if the default is `metricbeat-*`, you can change it to `custombeat-*`.

- **`-insecure`**
If specified, "insecure" SSL connections are allowed.

- **`-k <kibana_index>`**
The Elasticsearch index pattern where Kibana saves its configuration. The default value is `.kibana`.

- **`-key <client_key>`**
The client certificate key. The key must be in `PEM` format.

- **`-only-dashboards`**
If specified, then only the `dashboards`, along with their `visualizations` and `searches`, are imported. The index pattern is not imported. By default, this is **false**.

- **`-only-index`**
If specified, then only the `index pattern` is imported. The dashboards, along with their visualizations and searches, are not imported. By default, this is **false**.

- **`-pass <password>`**
The password for authenticating the connection to Elasticsearch by using Basic Authentication. By default no username and password are used.

- **`-snapshot`**
Using `-snapshot` will import the snapshot dashboards build for the current version. This is mainly useful when running a snapshot Beat build for testing purpose.

> Note
> When using `-snapshot`, `-url` will be ignored.

- **`-url <zip_url>`**
Zip archive with the dashboards, available online. The archive can contain Kibana dashboards for a single Beat or for multiple Beats.

- **`-user <username>`**
The username for authenticating the connection to Elasticsearch by using Basic Authentication. By default no username and password are used.

### Structure of the Dashboards Archive

zip 归档文件中包含了至少一个 Beat 的 dashboards 内容；其中 index pattern, dashboards, visualizations 和 searches 分别位于针对每一种 Beat 设置的不同目录中；例如，官方 zip 归档文件（beats-dashboards-5.4.0）具有如下目录结构：

```
  metricbeat/
    dashboard/
    search/
    visualization/
    index-pattern/
  packetbeat/
    dashboard/
    search/
    visualization/
    index-pattern/
  filebeat/
    index-pattern/
  winlogbeat/
    dashboard/
    search/
    visualization/
    index-pattern/
```

## Building Your Own Beat Dashboards
## Generating the Beat Index Pattern
## Exporting New and Modified Beat Dashboards
## Archiving Your Beat Dashboards
## Sharing Your Beat Dashboards