# Index Pattern V.S. Index Template

## Index Pattern

![Kibana 未导入 Index Pattern 前](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Kibana%20%E6%9C%AA%E5%AF%BC%E5%85%A5%20Index%20Pattern%20%E5%89%8D.png "Kibana 未导入 Index Pattern 前")

- 每一种加载到 Elasticsearch 中的数据集都具有一种 index pattern ；
- **index pattern** 就是能够匹配多个 index 的、带有可选通配符（wildcards）的字符串；

    举例：针对常规的日志处理场景，一种比较典型的 index 名字将具有 `YYYY.MM.DD` 格式的日期，相应的，针对 5 月的 index pattern 看起来类似 `logstash-2015.05*` ；

- 如果数据源提供的数据内容确定包含 time-series 数据（例如将 Logstash 作为数据源时），则在点击 **Add New** 定义 index 时，请务必勾选上 `Index contains time-based events` 复选框，并从 `Time-field name` 下拉列表中选择 `@timestamp` 字段（field）；

> **注意**：当你定义一个 index pattern 时，匹配该 pattern 的 indices **必须**事先已经存在于 Elasticsearch 之中；并且这些 indices **必须**包含了数据内容；

![Configure an Index Pattern](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Configure%20an%20Index%20Pattern.jpeg "Configure an Index Pattern")


- 为了使用 Kibana ，你必须通过配置一个或多个 index patterns 告诉它你想要使用哪些 Elasticsearch indices ；
- Kibana 会查找匹配指定 pattern 的 index 名字；pattern 中出现的星号 (*) 用于匹配零个或多个字符； 
- index pattern 也可以简单的设置为单独一个 index 的名字；
- 若想在 index 名字中使用 event time ，可以在 pattern 中使用 [] 将静态文本括起来，之后再指定日期格式；例如 `[logstash-]YYYY.MM.DD` 能够匹配 `logstash-2015.01.31` 和 `logstash-2015-02-01` ；
- 在你浏览 Discover 页面时，设置成默认的 index pattern 将会被自动加载；Kibana 会在默认 pattern 左侧显示一个星星；你所创建的首个 pattern 自动被指定为默认 pattern ；
- 当你添加了一个 index mapping 后，Kibana 会自动扫描（scans）匹配 pattern 的 indices 以变显示（新的）index fields 列表；你也可以重新加载 index fields 列表以便获取任意新加的 fields ；


参考：

- [Defining Your Index Patterns](https://www.elastic.co/guide/en/kibana/current/tutorial-define-index.html)
- [Index Patterns](https://www.elastic.co/guide/en/kibana/current/index-patterns.html)


## Index Template

- Index templates 允许你定义 templates ，以便在新 indices 被创建时自动被应用；
- templates 内容由 settings 和 mappings 构成，以及一个简单的 **pattern template** 用于控制当前 template 是否被应用到新 index 上；
    > pattern template 其实就是创建模版时在 JSON 内容中设置的 `"template": "<regex>"` ；
- Templates 仅在 index 创建时被应用；变更一个 template 不会对已经存在 indices 造成影响；
- Index templates 提供了 C 风格的 `/* */` 块注释；注释内容允许出现在 JSON 文档的任何位置，除了 JSON 的初始 '{' 之上；
- 多个 index templates 可能同时匹配到同一个 index 上，在这种情况下，settings 和 mappings 内容会**合并**为目标 index 的最终配置；而合并的顺序可以基于 `order` 参数进行控制，order 值越小越优先被应用，而值越大越后被应用（实现**覆盖**效果）；
- Templates 可以选择性添加一个 `version` 字段指定版本号信息（其值为任意整数），以便简化外部系统对 template 进行管理；若要重置一个版本号，只需简单的使用未指定版本号的 template 进行替换；


参考：

- [Index Templates](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-templates.html)

