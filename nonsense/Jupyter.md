# Jupyter

> 官网地址：[这里](http://jupyter.org/)

## The Jupyter Notebook

The Jupyter Notebook is a web application that allows you to create and share documents that contain live code, equations, visualizations and explanatory text. Uses include: data cleaning and transformation, numerical simulation, statistical modeling, machine learning and much more.

### 特点

- Language of choice
- Share notebooks
- Interactive widgets
- Big data integration


### Jupyter Architecture

The Jupyter Notebook is based on a set of open standards for interactive computing. Think HTML and CSS for interactive computing on the web. These open standards can be leveraged by third party developers to build customized applications with embedded interactive computing.


#### Notebook Document Format

Jupyter Notebooks are an open document format based on JSON. They contain a complete record of the user's sessions and embed code, narrative text, equations and rich output.

#### Interactive Computing Protocol

The Notebook communicates with computational Kernels using the Interactive Computing Protocol, an open network protocol based on JSON data over ZMQ and WebSockets.

#### The Kernel

Kernels are processes that run interactive code in a particular programming language and return output to the user. Kernels also respond to tab completion and introspection requests.


----------


## 在线试用

[这里](https://try.jupyter.org/)

## 软件安装

方式一

We recommend using the `Anaconda` distribution to install Python and Jupyter. 

For new users, we highly recommend installing `Anaconda`. `Anaconda` conveniently installs Python, the Jupyter Notebook, and other commonly used packages for scientific computing and data science.

方式二

As an existing Python user, you may wish to install Jupyter using Python’s package manager, pip, instead of Anaconda.

```shell
pip install --upgrade pip
pip install jupyter
```

## 关于 Kernel

- 在 Jupyter Notebook 中每一种编程语言都对应了一种 kernels ；
- Jupyter Notebook 中会预安装 IPython kernel ；若想使用其他语言，则需要安装对应的 kernel ，详见[这里](https://github.com/jupyter/jupyter/wiki/Jupyter-kernels)；


## 运行

- 默认启动

```shell
jupyter notebook
```

- 指定端口启动

```shell
jupyter notebook --port 9999
```

- 启动时不默认打开浏览器页面

```shell
jupyter notebook --no-browser
```

- 帮助信息获取

```shell
jupyter notebook --help
```

## Migrating from IPython Notebook

如果你想知道 IPython 和 Jupyter 之间的恩怨纠葛，请移步[这里](https://blog.jupyter.org/2015/04/15/the-big-split/)；

简言之

> Jupyter is the new home of language-agnostic projects that began as part of IPython, such as the notebook.





## 创建 web 版 PPT

```shell
jupyter nbconvert xxx.ipynb --to slides --post serve
```


There is a nice guide here: [Presentation slides with Jupyter Notebook](http://echorand.me/presentation-slides-with-jupyter-notebook.html#.WIRfjbZ94p8)


```
jupyter-nbconvert --to slides slides.ipynb --reveal-prefix=reveal.js
```

注：reveal.js 内容看下面；

> `jupyter nbconvert` 等价于 `jupyter-nbconvert`


## jupyter-nbconvert

```shell
# jupyter-nbconvert -h

This application is used to convert notebook files (*.ipynb) to various other
formats.

WARNING: THE COMMANDLINE INTERFACE MAY CHANGE IN FUTURE RELEASES.

Options
-------

Arguments that take values are actually convenience aliases to full
Configurables, whose aliases are listed on the help line. For more information
on full configurables, see '--help-all'.

--execute
    Execute the notebook prior to export.
--allow-errors
    Continue notebook execution even if one of the cells throws an error and include the error message in the cell output (the default behaviour is to abort conversion). This flag is only relevant if '--execute' was specified, too.
--stdout
    Write notebook output to stdout instead of files.
--stdin
    read a single notebook file from stdin. Write the resulting notebook with default basename 'notebook.*'
--inplace
    Run nbconvert in place, overwriting the existing notebook (only
    relevant when converting to notebook format)
-y
    Answer yes to any questions instead of prompting.
--debug
    set log level to logging.DEBUG (maximize logging output)
--generate-config
    generate default config file
--nbformat=<Enum> (NotebookExporter.nbformat_version)
    Default: 4
    Choices: [1, 2, 3, 4]
    The nbformat version to write. Use this to downgrade notebooks.
--output-dir=<Unicode> (FilesWriter.build_directory)
    Default: ''
    Directory to write output(s) to. Defaults to output to the directory of each
    notebook. To recover previous default behaviour (outputting to the current
    working directory) use . as the flag value.
--writer=<DottedObjectName> (NbConvertApp.writer_class)
    Default: 'FilesWriter'
    Writer class used to write the  results of the conversion
--log-level=<Enum> (Application.log_level)
    Default: 30
    Choices: (0, 10, 20, 30, 40, 50, 'DEBUG', 'INFO', 'WARN', 'ERROR', 'CRITICAL')
    Set the log level by value or name.
--reveal-prefix=<Unicode> (SlidesExporter.reveal_url_prefix)
    Default: u''
    The URL prefix for reveal.js. This can be a a relative URL for a local copy
    of reveal.js, or point to a CDN.
    For speaker notes to work, a local reveal.js prefix must be used.
--to=<Unicode> (NbConvertApp.export_format)
    Default: 'html'
    The export format to be used, either one of the built-in formats, or a
    dotted object name that represents the import path for an `Exporter` class
--template=<Unicode> (TemplateExporter.template_file)
    Default: u''
    Name of the template file to use
--output=<Unicode> (NbConvertApp.output_base)
    Default: ''
    overwrite base name use for output files. can only be used when converting
    one notebook at a time.
--post=<DottedOrNone> (NbConvertApp.postprocessor_class)
    Default: u''
    PostProcessor class used to write the results of the conversion
--config=<Unicode> (JupyterApp.config_file)
    Default: u''
    Full path of a config file.

To see all available configurables, use `--help-all`

Examples
--------

    The simplest way to use nbconvert is

    > jupyter nbconvert mynotebook.ipynb

    which will convert mynotebook.ipynb to the default format (probably HTML).

    You can specify the export format with `--to`.
    Options include ['asciidoc', 'custom', 'html', 'latex', 'markdown', 'notebook', 'pdf', 'python', 'rst', 'script', 'slides']

    > jupyter nbconvert --to latex mynotebook.ipynb

    Both HTML and LaTeX support multiple output templates. LaTeX includes
    'base', 'article' and 'report'.  HTML includes 'basic' and 'full'. You
    can specify the flavor of the format used.

    > jupyter nbconvert --to html --template basic mynotebook.ipynb

    You can also pipe the output to stdout, rather than a file

    > jupyter nbconvert mynotebook.ipynb --stdout

    PDF is generated via latex

    > jupyter nbconvert mynotebook.ipynb --to pdf

    You can get (and serve) a Reveal.js-powered slideshow

    > jupyter nbconvert myslides.ipynb --to slides --post serve

    Multiple notebooks can be given at the command line in a couple of
    different ways:

    > jupyter nbconvert notebook*.ipynb
    > jupyter nbconvert notebook1.ipynb notebook2.ipynb

    or you can specify the notebooks list in a config file, containing::

        c.NbConvertApp.notebooks = ["my_notebook.ipynb"]

    > jupyter nbconvert --config mycfg.py

```


## jupyter-notebook

该命令用于创建 "The Jupyter HTML Notebook" ；


----------

```shell
# jupyter-notebook -h
```

> 以下内容为 jupyter-notebook 命令的帮助说明；

通过 `jupyter-notebook` 命令能够启动一个基于 Tornado 的 HTML Notebook Server 并为
HTML5/Javascript Notebook client 提供服务；

### 子命令

子命令调用方式：`jupyter-notebook cmd [args]`
子命令帮助信息：`jupyter-notebook cmd -h`

**`list`**
    List currently running notebook servers.

### 选项

> Arguments that take values are actually convenience aliases to full
Configurables, whose aliases are listed on the help line. For more information
on full configurables, see '--help-all'.

**`--script`**

    DEPRECATED, IGNORED
    
**`--pylab`**

    DISABLED: use %pylab or %matplotlib in the notebook to enable matplotlib.
    
**`--debug`**

    set log level to logging.DEBUG (maximize logging output)
    
**`--no-browser`**

    Don't open the notebook in a browser after startup.
    
**`-y`**

    Answer yes to any questions instead of prompting.
    
**`--no-mathjax`**

    Disable MathJax
    
    MathJax is the javascript library Jupyter uses to render math/LaTeX. It is
    very large, so you may want to disable it if you have a slow internet
    connection, or for offline use of the notebook.

    When disabled, equations etc. will appear as their untransformed TeX source.
    
**`--no-script`**

    DEPRECATED, IGNORED
    
**`--generate-config`**

    generate default config file
    
**`--certfile=<Unicode>`** (NotebookApp.certfile)

    Default: u''
    The full path to an SSL/TLS certificate file.
    
**`--ip=<Unicode>`** (NotebookApp.ip)

    Default: 'localhost'
    The IP address the notebook server will listen on.
    
**`--pylab=<Unicode>`** (NotebookApp.pylab)

    Default: 'disabled'
    DISABLED: use %pylab or %matplotlib in the notebook to enable matplotlib.
    
**`--log-level=<Enum>`** (Application.log_level)

    Default: 30
    Choices: (0, 10, 20, 30, 40, 50, 'DEBUG', 'INFO', 'WARN', 'ERROR', 'CRITICAL')
    Set the log level by value or name.
    
**`--port-retries=<Integer>`** (NotebookApp.port_retries)

    Default: 50
    The number of additional ports to try if the specified port is not
    available.
    
**`--notebook-dir=<Unicode>`** (NotebookApp.notebook_dir)

    Default: u''
    The directory to use for notebooks and kernels.
    
**`--keyfile=<Unicode>`** (NotebookApp.keyfile)

    Default: u''
    The full path to a private key file for usage with SSL/TLS.
    
**`--client-ca=<Unicode>`** (NotebookApp.client_ca)

    Default: u''
    The full path to a certificate authority certificate for SSL/TLS client
    authentication.
    
**`--config=<Unicode>`** (JupyterApp.config_file)

    Default: u''
    Full path of a config file.
    
**`--port=<Integer>`** (NotebookApp.port)

    Default: 8888
    The port the notebook server will listen on.
    
**`--transport=<CaselessStrEnum>`** (KernelManager.transport)

    Default: 'tcp'
    Choices: [u'tcp', u'ipc']
    
**`--browser=<Unicode>`** (NotebookApp.browser)

    Default: u''
    Specify what command to use to invoke a web browser when opening the
    notebook. If not specified, the default browser will be determined by the
    `webbrowser` standard library module, which allows setting of the BROWSER
    environment variable to override it.

To see all available configurables, use `--help-all`

### 实例

```
    jupyter notebook                       # start the notebook
    jupyter notebook --certfile=mycert.pem # use SSL/TLS certificate
```


----------


# reveal.js

作品展示：

- [作品 1](http://lab.hakim.se/reveal-js/#/)
- [作品 2](https://knewone.com/about/startup.html#/)

教程：[这里](http://tchen.me/posts/2012-12-26-reveal-js-support-for-octopress.html)
其他：[这里](http://hackerzhang.com/post/tips/make-presention-using-reveal-js)