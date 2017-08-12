# 用Elastic Stack来看看祖国的蓝天之数据导入篇

> 原文地址：[这里](https://mp.weixin.qq.com/s?__biz=MzI2NDExNTk5Mg==&mid=2247484189&idx=1&sn=9f6e274e5d65bfa6a6821ecfce4d7d75&chksm=eab0c16addc7487c135a064d074ff10ec31afdcccf003a63ecfe3a6724a86eb71fe79f6a1a65&mpshare=1&scene=1&srcid=05122kHHC9lJpzdRfqsd0X0m&key=c7e325e91b5a6402d0395473982db5fd8387edc15dd017871890bb2ddc6cfa3f25de1ee3615de537b4c0bfa19757aa18e42d7c5fcd803ff7be2600ae95c56835c97b5c3dacedb9017e584e353f4d7c05&ascene=0&uin=MTE2NDAxMTQyMA%3D%3D&devicetype=iMac+MacBookPro12%2C1+OSX+OSX+10.11.6+build(15G1217)&version=12020010&nettype=WIFI&fontScale=100&pass_ticket=blC1ou7VzXwa2KJnvukc3e%2F82JpE4a7vVGXgy6gwWoyaC4qoxPK%2BbBeJpjgBE4az)

## 数据获取

> 原始数据下载地址：[这里](https://github.com/elastic-adventures/air-quality)

下载后查看

```
➜  ~ git clone https://github.com/elastic-adventures/air-quality.git
➜  ~
➜  ~
➜  ~ cd workspace/GIT
➜  GIT ll
total 0
drwxr-xr-x    6 sunfei  staff   204B  4 26 11:19 IdeaProjects
drwxr-xr-x    6 sunfei  staff   204B  8 29  2016 ImageCache
drwxr-xr-x   15 sunfei  staff   510B 11 26 15:12 MarkSomethingDown
drwxr-xr-x    6 sunfei  staff   204B  5 16 14:30 air-quality
drwxr-xr-x    7 sunfei  staff   238B  4  6 12:12 dotsfile
drwxr-xr-x   13 sunfei  staff   442B  9 17  2016 entop
drwxr-xr-x  134 sunfei  staff   4.4K  8 16  2016 glibc
drwxr-xr-x    2 sunfei  staff    68B  1 23 17:59 gpg
drwxr-xr-x   10 sunfei  staff   340B  1 23 14:23 jupyter-notebook-slides
drwxr-xr-x    8 sunfei  staff   272B 10 10  2016 learngit
drwxr-xr-x    9 sunfei  staff   306B  2  9 15:09 pcaphub
drwxr-xr-x   13 sunfei  staff   442B  3 28 11:36 pesticides
drwxr-xr-x   14 sunfei  staff   476B 11 18 17:58 rabbit-stress
drwxr-xr-x   24 sunfei  staff   816B  3 28 18:14 recon_web
drwxr-xr-x   11 sunfei  staff   374B 11 22 16:43 redis_dissector_for_wireshark
drwxr-xr-x   18 sunfei  staff   612B  2  8 15:42 reveal.js
➜  GIT git clone https://github.com/elastic-adventures/air-quality.git
➜  GIT
➜  GIT
➜  GIT ll air-quality
total 24
-rw-r--r--   1 sunfei  staff   522B  5 16 14:20 README.md
drwxr-xr-x  11 sunfei  staff   374B  5 16 14:21 data
-rw-r--r--   1 sunfei  staff   4.1K  5 16 14:20 filebeat.yml
➜  GIT ll air-quality/data
total 9568
-rw-r--r--  1 sunfei  staff   328K  5 16 14:20 Beijing_2008_HourlyPM2.5_created20140325.csv
-rw-r--r--  1 sunfei  staff   563K  5 16 14:20 Beijing_2009_HourlyPM25_created20140709.csv
-rw-r--r--  1 sunfei  staff   558K  5 16 14:20 Beijing_2010_HourlyPM25_created20140709.csv
-rw-r--r--  1 sunfei  staff   558K  5 16 14:20 Beijing_2011_HourlyPM25_created20140709.csv
-rw-r--r--  1 sunfei  staff   567K  5 16 14:20 Beijing_2012_HourlyPM2.5_created20140325.csv
-rw-r--r--  1 sunfei  staff   556K  5 16 14:20 Beijing_2013_HourlyPM2.5_created20140325.csv
-rw-r--r--  1 sunfei  staff   555K  5 16 14:20 Beijing_2014_HourlyPM25_created20150203.csv
-rw-r--r--  1 sunfei  staff   542K  5 16 14:20 Beijing_2015_HourlyPM25_created20160201.csv
-rw-r--r--  1 sunfei  staff   543K  5 16 14:20 Beijing_2016_HourlyPM25_created20170201.csv
➜  GIT
```

> 如何把这些 csv 文件转换为 elasticsearch 中的文档呢

## 数据导入

下面通过 `filebeat` 进行上述 csv 文件内容导入；

配置文件如下

```
filebeat.prospectors:
- input_type: stdin
  # Type to be published in the 'type' field. For Elasticsearch output,
  # the type defines the document type these entries should be stored
  # in. Default: log
  document_type: air_quality
  exclude_lines: ["^A ","^The","^Site","^,"]

output.elasticsearch:
  hosts: ["localhost:9200"]
  
  # Optional protocol and basic auth credentials.
  username: "elastic"   # 使能 x-pack 后需要
  password: "changeme"  # 使能 x-pack 后需要
  
  # Optional index name. The default is "filebeat" plus date
  # and generates [filebeat-]YYYY.MM.DD keys.
  index: "air_quality"
  
  # Optional ingest node pipeline. By default no pipeline will be used.
  #pipeline: ""

output.console:
  pretty: true
```

导入数据

```
filebeat git:(5.3) ✗ cat data_for_test/Beijing_2016_HourlyPM25_created20170201.csv| ./filebeat -e -c filebeat.yml
...
{
  "@timestamp": "2017-05-16T08:32:55.034Z",
  "beat": {
    "hostname": "sunfeideMacBook-Pro.local",
    "name": "sunfeideMacBook-Pro.local",
    "version": "5.3.3"
  },
  "input_type": "stdin",
  "message": "Beijing,PM2.5,12/31/2016 22:00,2016,12,31,22,488,\ufffdg/m\ufffd,1 Hr,Valid",
  "offset": 0,
  "source": "",
  "type": "air_quality"
}
{
  "@timestamp": "2017-05-16T08:32:55.034Z",
  "beat": {
    "hostname": "sunfeideMacBook-Pro.local",
    "name": "sunfeideMacBook-Pro.local",
    "version": "5.3.3"
  },
  "input_type": "stdin",
  "message": "Beijing,PM2.5,12/31/2016 23:00,2016,12,31,23,507,\ufffdg/m\ufffd,1 Hr,Valid",
  "offset": 0,
  "source": "",
  "type": "air_quality"
}

(Ctrl+C)

2017/05/16 08:33:23.208117 metrics.go:39: INFO Non-zero metrics in the last 30s: filebeat.harvester.closed=1 filebeat.harvester.open_files=-1 filebeat.harvester.started=1 libbeat.es.call_count.PublishEvents=176 libbeat.es.publish.read_bytes=86923 libbeat.es.publish.write_bytes=3110824 libbeat.es.published_and_acked_events=8784 libbeat.publisher.published_events=8784 publish.events=8789 registrar.states.cleanup=1 registrar.states.update=5 registrar.writes=5
2017/05/16 08:33:53.207663 metrics.go:34: INFO No non-zero metrics in the last 30s
^C2017/05/16 08:34:17.543907 filebeat.go:226: INFO Stopping filebeat
2017/05/16 08:34:17.543948 crawler.go:90: INFO Stopping Crawler
2017/05/16 08:34:17.543960 crawler.go:100: INFO Stopping 1 prospectors
2017/05/16 08:34:17.543957 prospector.go:137: INFO Prospector channel stopped because beat is stopping.
2017/05/16 08:34:17.543972 prospector.go:180: INFO Prospector ticker stopped
2017/05/16 08:34:17.543993 prospector.go:232: INFO Stopping Prospector: 6063977265577267127
2017/05/16 08:34:17.544005 crawler.go:112: INFO Crawler stopped
2017/05/16 08:34:17.544010 spooler.go:101: INFO Stopping spooler
2017/05/16 08:34:17.544023 registrar.go:291: INFO Stopping Registrar
2017/05/16 08:34:17.544029 registrar.go:248: INFO Ending Registrar
2017/05/16 08:34:17.544657 metrics.go:51: INFO Total non-zero values:  filebeat.harvester.closed=1 filebeat.harvester.open_files=-1 filebeat.harvester.started=1 libbeat.es.call_count.PublishEvents=176 libbeat.es.publish.read_bytes=86923 libbeat.es.publish.write_bytes=3110824 libbeat.es.published_and_acked_events=8784 libbeat.publisher.published_events=8784 publish.events=8789 registrar.states.cleanup=1 registrar.states.update=5 registrar.writes=6
2017/05/16 08:34:17.544677 metrics.go:52: INFO Uptime: 1m24.349095906s
2017/05/16 08:34:17.544684 beat.go:225: INFO filebeat stopped.
```

> 数据已经进入 elasticsearch 了，现在如何把 message 中的数据展开成最终的 json 形式？

## Ingest Node

`Ingest Node` 用于数据写入 elasticsearch 前对其做**预处理**（后续介绍）。

**数据预处理**是通过指定一个 `pipeline` 做到的，一个 `pipeline` 由一组 `processor` 组成。`processor` 是最小的处理单元，比如**删除字段**、**设定字段类型**、**设定 Id**、**大小写转换**等等。我们这里用 `Grok processor` 来解析 message 中的数据。


一个快速模拟数据处理（使用 Grok 匹配）的 api ：

```
curl -u elastic:changeme -XGET http://localhost:9200/_ingest/pipeline/_simulate\?pretty -d '{
  "pipeline": {
    "processors": [
      {
        "grok": {
          "field": "message",
          "patterns": [
            "%{DATA:city},%{DATA:parameter},%{DATA:date},%{NUMBER:year},%{NUMBER:month},%{NUMBER:day},%{NUMBER:hour},%{NUMBER:value},%{DATA:unit},%{DATA:duration},%{WORD:status}"
          ]
        }
      }
    ]
  },
  "docs": [
    {
      "_index": "air_quality",
      "_type": "air_quality",
      "_id": "AasdfliOkdopoCuhya",
      "_score": 1,
      "_source": {
        "@timestamp": "2017-02-05T10:59:51.045Z",
        "beat": {
          "hostname": "",
          "name": "",
          "version": "5.2.0"
        },
        "input_type": "stdin",
        "message": "Beijing,PM2.5,2016-12-20 23:00,2016,12,31,23,507,g/m,1 Hr,Valid",
        "offset": 0,
        "source": "",
        "type": "air_quality"
      }
    }
  ]
}
'
```

返回结果

```
➜  ~
{
  "docs" : [
    {
      "doc" : {
        "_id" : "AasdfliOkdopoCuhya",
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_source" : {
          "date" : "2016-12-20 23:00",
          "offset" : 0,
          "city" : "Beijing",
          "year" : "2016",
          "input_type" : "stdin",
          "source" : "",
          "message" : "Beijing,PM2.5,2016-12-20 23:00,2016,12,31,23,507,g/m,1 Hr,Valid",
          "type" : "air_quality",
          "duration" : "1 Hr",
          "unit" : "g/m",
          "@timestamp" : "2017-02-05T10:59:51.045Z",
          "month" : "12",
          "hour" : "23",
          "parameter" : "PM2.5",
          "beat" : {
            "hostname" : "",
            "name" : "",
            "version" : "5.2.0"
          },
          "day" : "31",
          "value" : "507",
          "status" : "Valid"
        },
        "_ingest" : {
          "timestamp" : "2017-05-16T09:00:04.934Z"
        }
      }
    }
  ]
}
➜  ~
```

![ingest and pipeline - 1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ingest%20and%20pipeline%20-%201.png "ingest and pipeline - 1")

通过下面的 api 创建一个名为 airquality 的 pipeline

```
curl -u elastic:changeme -XPUT http://localhost:9200/_ingest/pipeline/airquality -d '{
  "processors": [
    {
      "grok": {
        "field": "message",
        "patterns": [
          "%{DATA:city},%{DATA:parameter},%{DATA:date},%{NUMBER:year},%{NUMBER:month},%{NUMBER:day},%{NUMBER:hour},%{NUMBER:value},%{DATA:unit},%{DATA:duration},%{WORD:status}"
        ]
      }
    },
    {
      "set": {
        "field": "_id",
        "value": "{{city}}-{{date}}"
      }
    },
    {
      "date": {
        "field": "date",
        "target_field": "@timestamp",
        "formats": [
          "MM/dd/yyyy HH:mm",
          "yyyy-MM-dd HH:mm"
        ],
        "timezone": "Asia/Shanghai"
      }
    },
    {
      "remove": {
        "field": "message"
      }
    },
    {
      "remove": {
        "field": "beat"
      }
    },
    {
      "remove": {
        "field": "input_type"
      }
    },
    {
      "remove": {
        "field": "offset"
      }
    },
    {
      "remove": {
        "field": "source"
      }
    },
    {
      "remove": {
        "field": "date"
      }
    },
    {
      "remove": {
        "field": "year"
      }
    },
    {
      "remove": {
        "field": "month"
      }
    },
    {
      "remove": {
        "field": "day"
      }
    },
    {
      "remove": {
        "field": "hour"
      }
    },
    {
      "remove": {
        "field": "duration"
      }
    },
    {
      "remove": {
        "field": "unit"
      }
    },
    {
      "remove": {
        "field": "type"
      }
    },
    {
      "convert": {
        "field": "value",
        "type": "integer"
      }
    }
  ]
}
'
{"acknowledged":true}%
➜  ~
```

![ingest and pipeline - 2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ingest%20and%20pipeline%20-%202.png "ingest and pipeline - 2")

在 `filebeat.yml` 中，将 pipeline 设定好

```
output.elasticsearch:
   hosts: ["localhost:9200"]
   pipeline: "airquality"   # 增加上面定义的 pipeline
   index: "air_quality"
```

再次导入数据（注意：每次基于下面命令进行导入都会导致数据增加，若想保持结果干净，需要自行删除）

```
➜  filebeat git:(5.3) ✗ cat data_for_test/Beijing_2016_HourlyPM25_created20170201.csv| ./filebeat -e -c filebeat.yml
```

之后重新查询会得到转换后的“干净的”文档

```
➜  ~ curl -u elastic:changeme -i -XGET 'http://localhost:9200/air_quality/air_quality/_search?pretty&q=Beijing'
HTTP/1.1 200 OK
content-type: application/json; charset=UTF-8
content-length: 3875

{
  "took" : 8,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "failed" : 0
  },
  "hits" : {
    "total" : 8784,
    "max_score" : 1.04389255E-4,
    "hits" : [
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/1/2016 8:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-01T08:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 100,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/1/2016 21:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-01T21:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 345,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 7:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T07:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 206,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 12:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T12:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 222,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 15:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T15:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 208,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 16:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T16:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 240,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 18:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T18:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 368,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/2/2016 19:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-02T19:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 493,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/3/2016 2:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-03T02:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 406,
          "status" : "Valid"
        }
      },
      {
        "_index" : "air_quality",
        "_type" : "air_quality",
        "_id" : "Beijing-1/3/2016 3:00",
        "_score" : 1.04389255E-4,
        "_source" : {
          "city" : "Beijing",
          "@timestamp" : "2016-01-03T03:00:00.000+08:00",
          "parameter" : "PM2.5",
          "value" : 261,
          "status" : "Valid"
        }
      }
    ]
  }
}
➜  ~
```

------

## 小节

- 创建 ingest node pipeline 前可以先进行 simulate 模拟；
- ingest node pipeline 是创建在 elasticsearch 上的；
- 当 beat 想要使用 ingest node pipeline 功能时，需要在 `output.elasticsearch` 下配置 `pipeline: "xxx"`（xxx 为所创建的 pipeline 名字）；

