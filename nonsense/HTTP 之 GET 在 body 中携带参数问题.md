# HTTP 之 GET 在 body 中携带参数问题

在测试 Elasticsearch 时遇到下面这种测试命令：

```
curl -u elastic:changeme -XGET http://localhost:9200/_search\?pretty -d '{ "from": 30, "size":10}'
```

可以看到这里指定的 `-XGET` ，但同时通过 `-d` 参数指定了数据内容：

> 在 `man 1 curl` 中可以看到
> 
> **`-d, --data <data>`**
> 
> (HTTP) Sends the specified data in a **POST** request to the HTTP server, in the same way that a browser does when a user has filled in an HTML form and presses the submit button. This will cause `curl` to pass the data to the server using the **content-type** `application/x-www-form-urlencoded`.
> 
> `-d, --data` is the same as `--data-ascii`. `--data-raw` is almost the same but does not have a special interpretation of the `@` character. To post data purely binary, you should instead use the `--data-binary` option. To **URL-encode** the value of a form field you may use  `--data-urlencode`.
If any of these options is used more than once on the same command line, the data pieces specified will be merged together with a separating `&`-symbol. Thus, using '`-d name=daniel -d skill=lousy`' would generate a post chunk that looks like '`name=daniel&skill=lousy`'.
> 
> If you start the data with the letter `@`, the rest should be a **file name** to read the data from, or `-` if you want `curl` to read the data from `stdin`. Multiple files can also be specified. Posting data from a file named 'foobar' would thus be done with `--data @foobar`. When `--data` is told to read from a file like that, carriage returns and newlines will be stripped out. If you don't want the `@` character to have a special interpretation use `--data-raw` instead.

这里就引出了一个矛盾点：在 `GET` 请求中通过 `-d` 发送数据是什么鬼？！

抓包发现确实可以这么使用：

![HTTP GET with data in body](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/HTTP%20GET%20with%20data%20in%20body.png "HTTP GET with data in body")

一般来讲，任何语言（特别是 javascript）实现的 HTTP 库都不允许 `GET` 请求中携带交互数据。因此，很多用户会惊讶于 `GET` 请求中居然会允许这么做。

真实情况是：在规定 HTTP 语义及内容的 [RFC 7231](http://tools.ietf.org/html/rfc7231#page-24) 中，并未限制 `GET` 请求中是否允许携带交互数据！所以，有些 HTTP 服务允许这种行为，而另外一些（特别是缓存代理）则不允许这种行为。

Elasticsearch 的作者们倾向于使用 `GET` 提交查询请求，因为他们觉得这个 verb 相比 `POST` 来说，能更好的描述这种行为。然而，因为携带交互数据的 GET 请求并不被广泛支持，所以 search API 同样支持 `POST` 请求，即

```
POST /_search
{
  "from": 30,
  "size": 10
}
```
