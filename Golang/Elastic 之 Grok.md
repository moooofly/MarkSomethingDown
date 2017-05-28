# Elastic 之 Grok

## Grok 本意

以下内容取自[维基百科](https://en.wikipedia.org/wiki/Grok#In_computer_programmer_culture)：

> Uses of the word in the decades after the 1960s are more concentrated in computer culture, such as a 1984 appearance in InfoWorld: "There isn't any software! Only different internal states of hardware. It's all hardware! It's a shame programmers don't grok that better."
> 
> When you claim to "grok" some knowledge or technique, you are asserting that you have not merely learned it in a detached instrumental way but that it has become part of you, part of your identity. For example, to say that you "know" Lisp is simply to assert that you can code in it if necessary — but to say you "grok" LISP is to claim that you have deeply entered the world-view and spirit of the language, with the implication that it has transformed your view of programming. Contrast zen, which is a similar supernatural understanding experienced as a single brief flash.

----------

## Ingest node & Grok

### [Ingest Node](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/ingest.html)

你可以使用 ingest node 在 documents 被“索引”前对其进行预处理 ；预处理行为包括：对 bulk 和 index 请求进行解析、实施（数据）转换、返回 documents 给相应的 index 或 bulk APIs 调用；

> You can use ingest node to pre-process documents before the actual indexing takes place. This pre-processing happens by an ingest node that intercepts bulk and index requests, applies the transformations, and then passes the documents back to the index or bulk APIs.

可以在任意 node 上使能 ingest 功能，也可以创建专用的 ingest nodes ；Ingest 功能默认在所有 node 上都被使能；如果想去使能某个 node 上的 ingest 功能，可以在 `elasticsearch.yml` 如下配置：

```
node.ingest: false
```

为了在“索引”前进行预处理，你需要[定义一个 pipeline](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/pipeline.html) 用于指定一系列 [processors](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/ingest-processors.html) ；每一个 processor 都以某种方式对 document 进行转换；

> To pre-process documents before indexing, you [define a pipeline](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/pipeline.html) that specifies a series of [processors](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/ingest-processors.html). Each processor transforms the document in some way. For example, you may have a pipeline that consists of one processor that removes a field from the document followed by another processor that renames a field.

为了使用 pipeline ，你可以简单的在 index 或 bulk 请求中指定 `pipeline` 参数，以告知目标 ingest node 需要使用哪种 pipeline ；例如

```
PUT my-index/my-type/my-id?pipeline=my_pipeline_id
{
  "foo": "bar"
}
```

关于 pipelines 的创建、添加和删除，详见 [Ingest APIs](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/ingest-apis.html) ；


### Pipeline Definition

pipeline 本质上就是一组 [processors](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/ingest-processors.html) 定义；processors 的执行顺序与其声明顺序相同；每一个 pipeline 都是由两种主要的 field 构成：一个 `description` 和一组 `processors` ：

```
{
  "description" : "...",
  "processors" : [ ... ]
}
```

`description` 用于保存关于当前 pipeline 的描述信息；`processors` 用于定义需要被顺序执行的一组 processors ；

### Ingest APIs

如下 ingest APIs 用于管理 pipelines ：

- [Put Pipeline API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/put-pipeline-api.html) 用于添加或更新 pipeline ；
- [Get Pipeline API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/get-pipeline-api.html) 用于返回指定的 pipeline ；
- [Delete Pipeline API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/delete-pipeline-api.html) 用于删除一个 pipeline ；
- [Simulate Pipeline API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/simulate-pipeline-api.html) 用于模拟（simulate）针对 pipeline 的调用；


### Put Pipeline API

**put pipeline API** 用于在 cluster 中添加 pipelines 和更新已存在的 pipelines ；

```
PUT _ingest/pipeline/my-pipeline-id
{
  "description" : "describe pipeline",
  "processors" : [
    {
      "set" : {
        "field": "foo",
        "value": "bar"
      }
    }
  ]
}
```

> 注意：**put pipeline API** 同样会令（instructs）全部 ingest nodes 重新加载其 in-memory representation of pipelines ，因此 pipeline 变更是立即生效的；

### Get Pipeline API

**get pipeline API** 用于基于 ID 返回 pipelines ；该 API 总是返回 a local reference of the pipeline ；

```
GET _ingest/pipeline/my-pipeline-id
```

应答示例：

```
{
  "my-pipeline-id" : {
    "description" : "describe pipeline",
    "processors" : [
      {
        "set" : {
          "field" : "foo",
          "value" : "bar"
        }
      }
    ]
  }
}
```

对于每个返回的 pipeline ，都会有 `source` 和 `version` 信息被返回；其中 `version` 可用于获知指定 node 上持有的 pipeline 版本是什么；你也可以一次性指定多个 IDs 以便返回多个 pipeline 信息；Wildcards 同样可以使用；

#### Pipeline Versioning

Pipelines 可以选择添加一个 `version` 号，可以为任意 integer 值，以方便外部系统进行管理；`version` field 完全是可选的，仅用于外部管理目的；若想要重置 `version` ，可以简单的使用不带 `version` 信息的 pipeline 进行替换；

```
PUT _ingest/pipeline/my-pipeline-id
{
  "description" : "describe pipeline",
  "version" : 123,
  "processors" : [
    {
      "set" : {
        "field": "foo",
        "value": "bar"
      }
    }
  ]
}
```

若想检查 `version` 信息，你可以使用 `filter_path` 限制将应答内容仅包含 `version` 信息，即 [filter responses](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/common-options.html#common-options-response-filtering) ；

```
GET /_ingest/pipeline/my-pipeline-id?filter_path=*.version
```

上述命令将产生一个小巧且便于解析的应答；

```
{
  "my-pipeline-id" : {
    "version" : 123
  }
}
```

### Delete Pipeline API

**delete pipeline API** 可以基于 ID 删除 pipelines ，或者使用 wildcard 进行匹配删除 (`my-*`, `*`) ；

```
DELETE _ingest/pipeline/my-pipeline-id
```

### Simulate Pipeline API

**simulate pipeline API** 用于针对请求体中给出的 documents 集合执行指定的 pipeline 处理；

你既可以指定一个已存在的 pipeline 对目标 documents 进行处理，也可以直接在请求体中给出 pipeline 定义再对目标 documents 进行处理；

在请求体中定义 pipeline 进行 documents 处理的 simulate 请求示例如下：

```
POST _ingest/pipeline/_simulate
{
  "pipeline" : {
    // pipeline definition here
  },
  "docs" : [
    { /** first document **/ },
    { /** second document **/ },
    // ...
  ]
}
```

基于已存在 pipeline 进行 documents 处理的 simulate 请求示例如下：

```
POST _ingest/pipeline/my-pipeline-id/_simulate
{
  "docs" : [
    { /** first document **/ },
    { /** second document **/ },
    // ...
  ]
}
```

完整示例：

```
POST _ingest/pipeline/_simulate
{
  "pipeline" :
  {
    "description": "_description",
    "processors": [
      {
        "set" : {
          "field" : "field2",
          "value" : "_value"
        }
      }
    ]
  },
  "docs": [
    {
      "_index": "index",
      "_type": "type",
      "_id": "id",
      "_source": {
        "foo": "bar"
      }
    },
    {
      "_index": "index",
      "_type": "type",
      "_id": "id",
      "_source": {
        "foo": "rab"
      }
    }
  ]
}
```
应答：

```
{
   "docs": [
      {
         "doc": {
            "_id": "id",
            "_ttl": null,
            "_parent": null,
            "_index": "index",
            "_routing": null,
            "_type": "type",
            "_timestamp": null,
            "_source": {
               "field2": "_value",
               "foo": "bar"
            },
            "_ingest": {
               "timestamp": "2016-01-04T23:53:27.186+0000"
            }
         }
      },
      {
         "doc": {
            "_id": "id",
            "_ttl": null,
            "_parent": null,
            "_index": "index",
            "_routing": null,
            "_type": "type",
            "_timestamp": null,
            "_source": {
               "field2": "_value",
               "foo": "rab"
            },
            "_ingest": {
               "timestamp": "2016-01-04T23:53:27.186+0000"
            }
         }
      }
   ]
}
```

#### Viewing Verbose Results

你可以使用 **simulate pipeline API** 来确认每一个 processor 是如何对通过 pipeline 的 ingest document 产生影响的；为了能够看到 simulate 过程中每个 processor 产生的中间结果，你可以使用 `verbose` 参数；

示例如下：

```
POST _ingest/pipeline/_simulate?verbose
{
  "pipeline" :
  {
    "description": "_description",
    "processors": [
      {
        "set" : {
          "field" : "field2",
          "value" : "_value2"
        }
      },
      {
        "set" : {
          "field" : "field3",
          "value" : "_value3"
        }
      }
    ]
  },
  "docs": [
    {
      "_index": "index",
      "_type": "type",
      "_id": "id",
      "_source": {
        "foo": "bar"
      }
    },
    {
      "_index": "index",
      "_type": "type",
      "_id": "id",
      "_source": {
        "foo": "rab"
      }
    }
  ]
}
```

应答：

```
{
   "docs": [
      {
         "processor_results": [
            {
               "tag": "processor[set]-0",
               "doc": {
                  "_id": "id",
                  "_ttl": null,
                  "_parent": null,
                  "_index": "index",
                  "_routing": null,
                  "_type": "type",
                  "_timestamp": null,
                  "_source": {
                     "field2": "_value2",
                     "foo": "bar"
                  },
                  "_ingest": {
                     "timestamp": "2016-01-05T00:02:51.383+0000"
                  }
               }
            },
            {
               "tag": "processor[set]-1",
               "doc": {
                  "_id": "id",
                  "_ttl": null,
                  "_parent": null,
                  "_index": "index",
                  "_routing": null,
                  "_type": "type",
                  "_timestamp": null,
                  "_source": {
                     "field3": "_value3",
                     "field2": "_value2",
                     "foo": "bar"
                  },
                  "_ingest": {
                     "timestamp": "2016-01-05T00:02:51.383+0000"
                  }
               }
            }
         ]
      },
      {
         "processor_results": [
            {
               "tag": "processor[set]-0",
               "doc": {
                  "_id": "id",
                  "_ttl": null,
                  "_parent": null,
                  "_index": "index",
                  "_routing": null,
                  "_type": "type",
                  "_timestamp": null,
                  "_source": {
                     "field2": "_value2",
                     "foo": "rab"
                  },
                  "_ingest": {
                     "timestamp": "2016-01-05T00:02:51.384+0000"
                  }
               }
            },
            {
               "tag": "processor[set]-1",
               "doc": {
                  "_id": "id",
                  "_ttl": null,
                  "_parent": null,
                  "_index": "index",
                  "_routing": null,
                  "_type": "type",
                  "_timestamp": null,
                  "_source": {
                     "field3": "_value3",
                     "field2": "_value2",
                     "foo": "rab"
                  },
                  "_ingest": {
                     "timestamp": "2016-01-05T00:02:51.384+0000"
                  }
               }
            }
         ]
      }
   ]
}
```

### Accessing Data in Pipelines

pipeline 中的 processors 对所有“穿过”自身的 documents 具有 read 和 write 权限；processors 能够访问 document 中 source 里的所有 fields ，以及 document 自身的 metadata fields ；

#### Accessing Fields in the Source

针对 source 中的 field 进行访问是很直接的，示例如下：

```
{
  "set": {
    "field": "my_field"
    "value": 582.1
  }
}
```

除此之外，对于 source 中的 fields 总是可以通过 `_source` 前缀进行访问：

```
{
  "set": {
    "field": "_source.my_field"
    "value": 582.1
  }
}
```

#### Accessing Metadata Fields

对 metadata fields 的访问和上面是一样的；原因在于 Elasticsearch 不允许 source 中的 fields 和 metadata fields 同名；

在下面的示例中，设置 document 的 `_id` metadata field 为 1 ：

```
{
  "set": {
    "field": "_id"
    "value": "1"
  }
}
```

允许被 processor 访问的 metadata fields 包括：`_index`, `_type`, `_id`, `_routing`, `_parent` ；

#### Accessing Ingest Metadata Fields

除了 metadata fields 和 source fields 之外，ingest 还会添加 ingest metadata 到目标 documents 中；这些 metadata properties 可以在 `_ingest` key 下被访问；当前，ingest 会添加 ingest timestamp 到 `_ingest.timestamp` key 下；ingest timestamp 对应的是 Elasticsearch 接收到 index 或 bulk request 开始预处理 document 的时间点；

任何 processor 都能够在 document 处理过程中添加 ingest-related metadata ；Ingest metadata 是 transient 的，因此当 document 被 pipeline 处理完毕后就会丢掉；因此，ingest metadata 不能被索引；

在如下示例中添加了一个名为 `received` 的 field ，其值为 ingest timestamp ：

```
{
  "set": {
    "field": "received"
    "value": "{{_ingest.timestamp}}"
  }
}
```

和 Elasticsearch 的 metadata fields 不同，ingest 的 metadata field 名 `_ingest` 可以用作 document 的 source 中的 field 名字，需要使用 `_source._ingest` 进行引用；否则，`_ingest` 将被理解为 ingest metadata field ；

#### Accessing Fields and Metafields in Templates

许多 processor settings 还支持模版（templating）功能；支持模版的 Settings 可以设置零个或多个 template snippets ；每个 template snippet 都以 `{{` 开始，以 `}}` 结束；针对 templates 中的 fields 和 metafields 的访问，与访问标准的 processor field 是一样的；

如下示例中添加了一个名为 `field_c` 的 field ，其值为 `field_a` 和 `field_b` 的组合；

```
{
  "set": {
    "field": "field_c"
    "value": "{{field_a}} {{field_b}}"
  }
}
```

如下示例使用了 source 中 `geoip.country_iso_code` field 的值，将其设置为 document 的 index ：

```
{
  "set": {
    "field": "_index"
    "value": "{{geoip.country_iso_code}}"
  }
}
```

### Handling Failures in Pipelines

略

### Processors

所有的 processors 都以如下方式进行定义：

```
{
  "PROCESSOR_NAME" : {
    ... processor configuration options ...
  }
}
```

每一个 processor 都会定义自身的配置参数，但全部 processors 都能够声明 `tag` 和 `on_failure` fields ；这些 fields 是可选的；

A `tag` is simply a string identifier of the specific instantiation of a certain processor in a pipeline. The `tag` field does not affect the processor’s behavior, but is very useful for bookkeeping and tracing errors to specific processors.

See [Handling Failures in Pipelines](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/handling-failure-in-pipelines.html) to learn more about the `on_failure` field and error handling in pipelines.

The [node info API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/cluster-nodes-info.html#ingest-info) can be used to figure out what processors are available in a cluster. The node info API will provide a per node list of what processors are available.

Custom processors must be installed on all nodes. The put pipeline API will fail if a processor specified in a pipeline doesn’t exist on all nodes. If you rely on custom processor plugins make sure to mark these plugins as mandatory by adding `plugin.mandatory` setting to the `config/elasticsearch.yml` file, for example:

```
plugin.mandatory: ingest-attachment,ingest-geoip
```

A node will not start if either of these plugins are not available.

The [node stats API](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/cluster-nodes-stats.html#ingest-stats) can be used to fetch ingest usage statistics, globally and on a per pipeline basis. Useful to find out which pipelines are used the most or spent the most time on preprocessing.


----------


### [Grok Processor](https://www.elastic.co/guide/en/elasticsearch/reference/5.4/grok-processor.html)

用于从单独的 text field 中提取 structured fields ；

> Extracts structured fields out of a single text field within a document. You choose which field to extract matched fields from, as well as the grok pattern you expect will match. A **grok pattern** is like a regular expression that supports aliased expressions that can be reused.

该工具适用于各种 log 的解析；该 processor 默认提供了超过 120 种可用的 pattern ；

> This tool is perfect for syslog logs, apache and other webserver logs, mysql logs, and in general, any log format that is generally written for humans and not computer consumption. This processor comes packaged with over [120 reusable patterns](https://github.com/elastic/elasticsearch/tree/master/modules/ingest-common/src/main/resources/patterns).

在线调试 grok 的工具；

> If you need help building patterns to match your logs, you will find the http://grokdebug.herokuapp.com and http://grokconstructor.appspot.com/ applications quite useful!

#### Grok Basics

Grok 以正则表达式为基础，因此任何正则表达式都是合法的 grok 表达式；

> Grok sits on top of regular expressions, so any regular expressions are valid in grok as well. The regular expression library is **Oniguruma**, and you can see the full supported regexp syntax [on the Onigiruma site](https://github.com/kkos/oniguruma/blob/master/doc/RE).

Grok 通过支持针对已有 patterns 进行命名的方式简化正则表达式的编写难度；

> Grok works by leveraging this regular expression language to allow naming existing patterns and combining them into more complex patterns that match your fields.

grok pattern 复用语法有三种形式：

- `%{SYNTAX:SEMANTIC}`
- `%{SYNTAX}`
- `%{SYNTAX:SEMANTIC:TYPE}`

> The syntax for reusing a grok pattern comes in three forms: `%{SYNTAX:SEMANTIC}`, `%{SYNTAX}`, `%{SYNTAX:SEMANTIC:TYPE}`.

`SYNTAX` 对应用于匹配 text 的 pattern 名字；

> The `SYNTAX` is the name of the pattern that will match your text. For example, `3.44` will be matched by the `NUMBER` pattern and `55.3.244.1` will be matched by the `IP` pattern. The syntax is how you match. `NUMBER` and `IP` are both patterns that are provided within the default patterns set.

`SEMANTIC` 作为从 text 中匹配的内容的标识符；

> The `SEMANTIC` is the identifier you give to the piece of text being matched. For example, `3.44` could be the duration of an event, so you could call it simply `duration`. Further, a string `55.3.244.1` might identify the `client` making a request.

`TYPE` 作为针对指定 field 进行 cast 的目标类型；当前仅支持 `int` 和 `float` 类型；

> The `TYPE` is the type you wish to cast your named field. `int` and `float` are currently the only types supported for coercion.

例如，针对如下内容进行匹配：

```
3.44 55.3.244.1
```

可以使用如下 Grok 表达式进行处理；

```
%{NUMBER:duration} %{IP:client}
```

#### Using the Grok Processor in a Pipeline

##### Table 20. Grok Options

| Name | Required | Default | Description |
| --- | --- | --- | --- | 
|field | yes | - | The field to use for grok expression parsing |
| patterns | yes | - | An ordered list of grok expression to match and extract named captures with. Returns on the first expression in the list that matches. |
| pattern_definitions | no | - | A map of pattern-name and pattern tuples defining custom patterns to be used by the current processor. Patterns matching existing names will override the pre-existing definition. |
| trace_match | no | false | when true, _ingest._grok_match_index will be inserted into your matched document’s metadata with the index into the pattern found in patterns that matched. |
| ignore_missing | no | false | If true and field does not exist or is null, the processor quietly exits without modifying the document |


示例如下：

```
{
  "message": "55.3.244.1 GET /index.html 15824 0.043"
}
```

grok pattern 为：

```
%{IP:client} %{WORD:method} %{URIPATHPARAM:request} %{NUMBER:bytes} %{NUMBER:duration}
```

如下为 pipeline 示例，基于 Grok 处理上述 document ：

```
{
  "description" : "...",
  "processors": [
    {
      "grok": {
        "field": "message",
        "patterns": ["%{IP:client} %{WORD:method} %{URIPATHPARAM:request} %{NUMBER:bytes} %{NUMBER:duration}"]
      }
    }
  ]
}
```

该 pipeline 会将捕获到的内容按照定义的命名插入到 document 中：

```
{
  "message": "55.3.244.1 GET /index.html 15824 0.043",
  "client": "55.3.244.1",
  "method": "GET",
  "request": "/index.html",
  "bytes": 15824,
  "duration": "0.043"
}
```

#### Custom Patterns and Pattern Files

Grok processor 带有预先打包好的基础 pattern 集合；这些 patterns 可能无法满足你所有的需求；Pattern 的格式很简单，每一个 entry 描述都由名字和 pattern 本身构成；

你可以在 `pattern_definitions` 选项下添加自定义 patterns 到 processor 定义中；

如下示例展示了如何指定自定义 pattern 的 pipeline ：

```
{
  "description" : "...",
  "processors": [
    {
      "grok": {
        "field": "message",
        "patterns": ["my %{FAVORITE_DOG:dog} is colored %{RGB:color}"]
        "pattern_definitions" : {
          "FAVORITE_DOG" : "beagle",
          "RGB" : "RED|GREEN|BLUE"
        }
      }
    }
  ]
}
```

#### Providing Multiple Match Patterns

有些时候一种 pattern 无法满足某个 field 的所有潜在可能的结构；Let’s assume we want to match all messages that contain your favorite pet breeds of either cats or dogs. 一种方式是提供两种不同的 patterns 用于匹配；另外一种是实现超级复杂的表达式以实现相同的 `or` 行为；

`_simulate` 示例如下：

```
POST _ingest/pipeline/_simulate
{
  "pipeline": {
  "description" : "parse multiple patterns",
  "processors": [
    {
      "grok": {
        "field": "message",
        "patterns": ["%{FAVORITE_DOG:pet}", "%{FAVORITE_CAT:pet}"],
        "pattern_definitions" : {
          "FAVORITE_DOG" : "beagle",
          "FAVORITE_CAT" : "burmese"
        }
      }
    }
  ]
},
"docs":[
  {
    "_source": {
      "message": "I love burmese cats!"
    }
  }
  ]
}
```

得到：

```
{
  "docs": [
    {
      "doc": {
        "_type": "_type",
        "_index": "_index",
        "_id": "_id",
        "_source": {
          "message": "I love burmese cats!",
          "pet": "burmese"
        },
        "_ingest": {
          "timestamp": "2016-11-08T19:43:03.850+0000"
        }
      }
    }
  ]
}
```

两种 patterns 都能子在匹配后成功设置 field `pet` 的值，但如果我们想要 trace 到底是哪种 patterns 完成的匹配并进行了 fields 提取呢？此时可以使用 `trace_match` 参数实现；针对相同 pipeline 但设置了 `"trace_match": true` 的示例输出如下：

```
{
  "docs": [
    {
      "doc": {
        "_type": "_type",
        "_index": "_index",
        "_id": "_id",
        "_source": {
          "message": "I love burmese cats!",
          "pet": "burmese"
        },
        "_ingest": {
          "_grok_match_index": "1",
          "timestamp": "2016-11-08T19:43:03.850+0000"
        }
      }
    }
  ]
}
```

在上面的应答中，你能看到 "_grok_match_index" 为 "1" ，表明完成匹配的是第二个 pattern（从 0 开始计数）；

该 trace metadata 使能了针对那种 patterns 完成匹配的调试功能；相应信息被保存在 ingest metadata 中，因此不会被索引到；


----------


## [Grok 正则捕获](http://www.ctolib.com/docs/sfile/ELKstack-guide-cn/logstash/plugins/filter/grok.html)

Grok 是 Logstash 最重要的插件；你可以在 grok 里预定义好命名正则表达式，在稍后（grok 参数或者其他正则表达式里）引用它；

### 正则表达式语法

你可以在 grok 里写标准的正则：

```
\s+(?<request_time>\d+(?:\.\d+)?)\s+
```

### Grok 表达式语法

Grok 支持把预定义的 grok 表达式 写入到文件中，官方提供的预定义 grok 表达式见[这里](https://github.com/logstash-plugins/logstash-patterns-core/tree/master/patterns)。

特点在于：

- 用普通的正则表达式来定义一个 grok 表达式；
- 通过打印赋值格式（sprintf format），用前面定义好的 grok 表达式来定义另一个 grok 表达式；

举例

```
USERNAME [a-zA-Z0-9._-]+
USER %{USERNAME}
```

grok 表达式的打印赋值格式的完整语法如下：

```
%{PATTERN_NAME:capture_name:data_type}
```

> 建议每个人都要使用 [Grok Debugger](http://grokdebug.herokuapp.com/) 来调试自己的 grok 表达式；


----------

## [grok-patterns](https://github.com/logstash-plugins/logstash-patterns-core/blob/master/patterns/grok-patterns)

```
USERNAME [a-zA-Z0-9._-]+
USER %{USERNAME}
EMAILLOCALPART [a-zA-Z][a-zA-Z0-9_.+-=:]+
EMAILADDRESS %{EMAILLOCALPART}@%{HOSTNAME}
INT (?:[+-]?(?:[0-9]+))
BASE10NUM (?<![0-9.+-])(?>[+-]?(?:(?:[0-9]+(?:\.[0-9]+)?)|(?:\.[0-9]+)))
NUMBER (?:%{BASE10NUM})
BASE16NUM (?<![0-9A-Fa-f])(?:[+-]?(?:0x)?(?:[0-9A-Fa-f]+))
BASE16FLOAT \b(?<![0-9A-Fa-f.])(?:[+-]?(?:0x)?(?:(?:[0-9A-Fa-f]+(?:\.[0-9A-Fa-f]*)?)|(?:\.[0-9A-Fa-f]+)))\b

POSINT \b(?:[1-9][0-9]*)\b
NONNEGINT \b(?:[0-9]+)\b
WORD \b\w+\b
NOTSPACE \S+
SPACE \s*
DATA .*?
GREEDYDATA .*
QUOTEDSTRING (?>(?<!\\)(?>"(?>\\.|[^\\"]+)+"|""|(?>'(?>\\.|[^\\']+)+')|''|(?>`(?>\\.|[^\\`]+)+`)|``))
UUID [A-Fa-f0-9]{8}-(?:[A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}
# URN, allowing use of RFC 2141 section 2.3 reserved characters
URN urn:[0-9A-Za-z][0-9A-Za-z-]{0,31}:(?:%[0-9a-fA-F]{2}|[0-9A-Za-z()+,.:=@;$_!*'/?#-])+

# Networking
MAC (?:%{CISCOMAC}|%{WINDOWSMAC}|%{COMMONMAC})
CISCOMAC (?:(?:[A-Fa-f0-9]{4}\.){2}[A-Fa-f0-9]{4})
WINDOWSMAC (?:(?:[A-Fa-f0-9]{2}-){5}[A-Fa-f0-9]{2})
COMMONMAC (?:(?:[A-Fa-f0-9]{2}:){5}[A-Fa-f0-9]{2})
IPV6 ((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?
IPV4 (?<![0-9])(?:(?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5]))(?![0-9])
IP (?:%{IPV6}|%{IPV4})
HOSTNAME \b(?:[0-9A-Za-z][0-9A-Za-z-]{0,62})(?:\.(?:[0-9A-Za-z][0-9A-Za-z-]{0,62}))*(\.?|\b)
IPORHOST (?:%{IP}|%{HOSTNAME})
HOSTPORT %{IPORHOST}:%{POSINT}

# paths
PATH (?:%{UNIXPATH}|%{WINPATH})
UNIXPATH (/([\w_%!$@:.,+~-]+|\\.)*)+
TTY (?:/dev/(pts|tty([pq])?)(\w+)?/?(?:[0-9]+))
WINPATH (?>[A-Za-z]+:|\\)(?:\\[^\\?*]*)+
URIPROTO [A-Za-z]([A-Za-z0-9+\-.]+)+
URIHOST %{IPORHOST}(?::%{POSINT:port})?
# uripath comes loosely from RFC1738, but mostly from what Firefox
# doesn't turn into %XX
URIPATH (?:/[A-Za-z0-9$.+!*'(){},~:;=@#%&_\-]*)+
#URIPARAM \?(?:[A-Za-z0-9]+(?:=(?:[^&]*))?(?:&(?:[A-Za-z0-9]+(?:=(?:[^&]*))?)?)*)?
URIPARAM \?[A-Za-z0-9$.+!*'|(){},~@#%&/=:;_?\-\[\]<>]*
URIPATHPARAM %{URIPATH}(?:%{URIPARAM})?
URI %{URIPROTO}://(?:%{USER}(?::[^@]*)?@)?(?:%{URIHOST})?(?:%{URIPATHPARAM})?

# Months: January, Feb, 3, 03, 12, December
MONTH \b(?:[Jj]an(?:uary|uar)?|[Ff]eb(?:ruary|ruar)?|[Mm](?:a|ä)?r(?:ch|z)?|[Aa]pr(?:il)?|[Mm]a(?:y|i)?|[Jj]un(?:e|i)?|[Jj]ul(?:y)?|[Aa]ug(?:ust)?|[Ss]ep(?:tember)?|[Oo](?:c|k)?t(?:ober)?|[Nn]ov(?:ember)?|[Dd]e(?:c|z)(?:ember)?)\b
MONTHNUM (?:0?[1-9]|1[0-2])
MONTHNUM2 (?:0[1-9]|1[0-2])
MONTHDAY (?:(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9])

# Days: Monday, Tue, Thu, etc...
DAY (?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)

# Years?
YEAR (?>\d\d){1,2}
HOUR (?:2[0123]|[01]?[0-9])
MINUTE (?:[0-5][0-9])
# '60' is a leap second in most time standards and thus is valid.
SECOND (?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)
TIME (?!<[0-9])%{HOUR}:%{MINUTE}(?::%{SECOND})(?![0-9])
# datestamp is YYYY/MM/DD-HH:MM:SS.UUUU (or something like it)
DATE_US %{MONTHNUM}[/-]%{MONTHDAY}[/-]%{YEAR}
DATE_EU %{MONTHDAY}[./-]%{MONTHNUM}[./-]%{YEAR}
ISO8601_TIMEZONE (?:Z|[+-]%{HOUR}(?::?%{MINUTE}))
ISO8601_SECOND (?:%{SECOND}|60)
TIMESTAMP_ISO8601 %{YEAR}-%{MONTHNUM}-%{MONTHDAY}[T ]%{HOUR}:?%{MINUTE}(?::?%{SECOND})?%{ISO8601_TIMEZONE}?
DATE %{DATE_US}|%{DATE_EU}
DATESTAMP %{DATE}[- ]%{TIME}
TZ (?:[APMCE][SD]T|UTC)
DATESTAMP_RFC822 %{DAY} %{MONTH} %{MONTHDAY} %{YEAR} %{TIME} %{TZ}
DATESTAMP_RFC2822 %{DAY}, %{MONTHDAY} %{MONTH} %{YEAR} %{TIME} %{ISO8601_TIMEZONE}
DATESTAMP_OTHER %{DAY} %{MONTH} %{MONTHDAY} %{TIME} %{TZ} %{YEAR}
DATESTAMP_EVENTLOG %{YEAR}%{MONTHNUM2}%{MONTHDAY}%{HOUR}%{MINUTE}%{SECOND}

# Syslog Dates: Month Day HH:MM:SS
SYSLOGTIMESTAMP %{MONTH} +%{MONTHDAY} %{TIME}
PROG [\x21-\x5a\x5c\x5e-\x7e]+
SYSLOGPROG %{PROG:program}(?:\[%{POSINT:pid}\])?
SYSLOGHOST %{IPORHOST}
SYSLOGFACILITY <%{NONNEGINT:facility}.%{NONNEGINT:priority}>
HTTPDATE %{MONTHDAY}/%{MONTH}/%{YEAR}:%{TIME} %{INT}

# Shortcuts
QS %{QUOTEDSTRING}

# Log formats
SYSLOGBASE %{SYSLOGTIMESTAMP:timestamp} (?:%{SYSLOGFACILITY} )?%{SYSLOGHOST:logsource} %{SYSLOGPROG}:

# Log Levels
LOGLEVEL ([Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn?(?:ing)?|WARN?(?:ING)?|[Ee]rr?(?:or)?|ERR?(?:OR)?|[Cc]rit?(?:ical)?|CRIT?(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE|EMERG(?:ENCY)?|[Ee]merg(?:ency)?)
```

----------

## 其他

- [Grok in filebeat?](https://discuss.elastic.co/t/grok-in-filebeat/66259)
- [Do you grok Grok?](https://www.elastic.co/blog/do-you-grok-grok)
- [Elastic Search -> Ingest Node -> processor -> grok](https://discuss.elastic.co/t/elastic-search-ingest-node-processor-grok/64316)
- [plugins-filters-grok](https://www.elastic.co/guide/en/logstash/current/plugins-filters-grok.html)
