# Pholcus

[![GitHub release](https://img.shields.io/github/release/andeya/pholcus.svg?style=flat-square)](https://github.com/andeya/pholcus/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/andeya/pholcus.svg)](https://pkg.go.dev/github.com/andeya/pholcus)
[![Go Report Card](https://goreportcard.com/badge/github.com/andeya/pholcus?style=flat-square)](https://goreportcard.com/report/andeya/pholcus)
[![GitHub issues](https://img.shields.io/github/issues/andeya/pholcus.svg?style=flat-square)](https://github.com/andeya/pholcus/issues?q=is%3Aopen+is%3Aissue)
[![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/andeya/pholcus.svg?style=flat-square)](https://github.com/andeya/pholcus/issues?q=is%3Aissue+is%3Aclosed)

Pholcus（幽灵蛛）是一款纯 Go 语言编写的支持分布式的高并发爬虫软件，仅用于编程学习与研究。

它支持单机、服务端、客户端三种运行模式，拥有 Web、GUI、命令行三种操作界面；规则简单灵活、批量任务并发、输出方式丰富（mysql/mongodb/kafka/csv/excel 等）；另外它还支持横纵向两种抓取模式，支持模拟登录和任务暂停、取消等一系列高级功能。

![Pholcus](https://github.com/andeya/pholcus/raw/master/doc/icon.png)

## 免责声明

**本软件仅用于学术研究，使用者需遵守其所在地的相关法律法规，请勿用于非法用途！！
如在中国大陆频频爆出爬虫开发者涉诉与违规的[新闻](https://github.com/HiddenStrawberry/Crawler_Illegal_Cases_In_China)。**

**郑重声明：因违法违规使用造成的一切后果，使用者自行承担！！**

## 爬虫原理

![模块结构](https://github.com/andeya/pholcus/raw/master/doc/module.png)

![项目架构](https://github.com/andeya/pholcus/raw/master/doc/project.png)

![分布式架构](https://github.com/andeya/pholcus/raw/master/doc/distribute.png)

## 框架特点

- 为具备一定 Go 或 JS 编程基础的用户提供只需关注规则定制、功能完备的重量级爬虫工具
- 支持单机、服务端、客户端三种运行模式
- GUI（仅 Windows）、Web、Cmd 三种操作界面，可通过参数控制打开方式
- 支持状态控制，如暂停、恢复、停止等
- 可控制采集量与并发协程数
- 支持多采集任务并发执行
- 支持代理 IP 列表，可控制更换频率
- 支持采集过程随机停歇，模拟人工行为
- 根据规则需求，提供自定义配置输入接口
- 支持 mysql、mongodb、kafka、csv、excel、原文件下载共六种输出方式
- 支持分批输出，且每批数量可控
- 支持静态 Go 和动态 JS 两种采集规则，支持横纵向两种抓取模式，且有大量 Demo
- 持久化成功记录，便于自动去重
- 序列化失败请求，支持反序列化自动重载处理
- 采用 surfer 高并发下载器，支持 GET/POST/HEAD 方法及 http/https 协议，同时支持固定 UserAgent 自动保存 cookie 与随机大量 UserAgent 禁用 cookie 两种模式，高度模拟浏览器行为，可实现模拟登录等功能
- 服务器/客户端模式采用 Teleport 高并发 SocketAPI 框架，全双工长连接通信，内部数据传输格式为 JSON

## 快速开始

### 获取源码

```bash
git clone https://github.com/andeya/pholcus.git
cd pholcus
```

### 创建项目

参考 `sample/pholcus_web.go`：

```go
package main

import (
    "github.com/andeya/pholcus/exec"
    _ "github.com/andeya/pholcus/rules" // 公开维护的 spider 规则库
    // _ "yourproject/rules_pte"         // 也可以自由添加自己的规则库
)

func main() {
    // 设置运行时默认操作界面，并开始运行
    // 可通过 -a_ui 参数指定为 "web"、"gui" 或 "cmd"
    // 其中 "gui" 仅支持 Windows 系统
    exec.DefaultRun("web")
}
```

### 编译运行

```bash
# 编译（非 Windows 平台会自动排除 GUI 包）
go build -o pholcus ./sample/

# 查看可选参数
./pholcus -h
```

Windows 下隐藏 cmd 窗口的编译方法：

```bash
go build -ldflags="-H=windowsgui -linkmode=internal" -o pholcus.exe ./sample/
```

![命令行帮助](https://github.com/andeya/pholcus/raw/master/doc/help.jpg)

> *Web 版操作界面*

![Web 界面](https://github.com/andeya/pholcus/raw/master/doc/webshow_1.png)

> *GUI 版操作界面（仅 Windows）*

![GUI 界面](https://github.com/andeya/pholcus/raw/master/doc/guishow_0.jpg)

> *Cmd 版运行参数设置示例*

```bash
pholcus -_ui=cmd -a_mode=0 -c_spider=3,8 -a_outtype=csv -a_thread=20 \
    -a_dockercap=5000 -a_pause=300 -a_proxyminute=0 \
    -a_keyins="<pholcus><golang>" -a_limit=10 -a_success=true -a_failure=true
```

> **注意：** Mac 下如使用代理 IP 功能，请务必获取 root 用户权限，否则无法通过 `ping` 获取可用代理。

## 运行时目录结构

```
├── pholcus                    可执行文件
└── pholcus_pkg/               运行时文件目录
    ├── config.ini             配置文件
    ├── proxy.lib              代理 IP 列表文件
    ├── spiders/               动态规则目录
    │   └── xxx.pholcus.xml    动态规则文件
    ├── phantomjs              PhantomJS 程序文件
    ├── text_out/              文本数据文件输出目录
    ├── file_out/              文件结果输出目录
    ├── logs/                  日志目录
    ├── history/               历史记录目录
    └── cache/                 临时缓存目录
```

## 动态规则示例

特点：动态加载规则，无需重新编译软件，书写简单，添加自由，适用于轻量级的采集项目。

将以下内容保存为 `pholcus_pkg/spiders/example.pholcus.xml`：

```xml
<Spider>
    <Name>HTML动态规则示例</Name>
    <Description>HTML动态规则示例 [Auto Page] [http://xxx.xxx.xxx]</Description>
    <Pausetime>300</Pausetime>
    <EnableLimit>false</EnableLimit>
    <EnableCookie>true</EnableCookie>
    <EnableKeyin>false</EnableKeyin>
    <NotDefaultField>false</NotDefaultField>
    <Namespace>
        <Script></Script>
    </Namespace>
    <SubNamespace>
        <Script></Script>
    </SubNamespace>
    <Root>
        <Script param="ctx">
        console.log("Root");
        ctx.JsAddQueue({
            Url: "http://xxx.xxx.xxx",
            Rule: "登录页"
        });
        </Script>
    </Root>
    <Rule name="登录页">
        <AidFunc>
            <Script param="ctx,aid">
            </Script>
        </AidFunc>
        <ParseFunc>
            <Script param="ctx">
            console.log(ctx.GetRuleName());
            ctx.JsAddQueue({
                Url: "http://xxx.xxx.xxx",
                Rule: "登录后",
                Method: "POST",
                PostData: "username=user@example.com&amp;password=pass&amp;submit=login"
            });
            </Script>
        </ParseFunc>
    </Rule>
    <Rule name="登录后">
        <ParseFunc>
            <Script param="ctx">
            console.log(ctx.GetRuleName());
            ctx.Output({
                "全部": ctx.GetText()
            });
            ctx.JsAddQueue({
                Url: "http://accounts.xxx.xxx/member",
                Rule: "个人中心",
                Header: {
                    "Referer": [ctx.GetUrl()]
                }
            });
            </Script>
        </ParseFunc>
    </Rule>
    <Rule name="个人中心">
        <ParseFunc>
            <Script param="ctx">
            console.log("个人中心: " + ctx.GetRuleName());
            ctx.Output({
                "全部": ctx.GetText()
            });
            </Script>
        </ParseFunc>
    </Rule>
</Spider>
```

## 静态规则示例

特点：随软件一同编译，定制性更强，效率更高，适用于重量级的采集项目。

在 `rules/` 目录下新建 Go 文件（参考 `rules/chinanews/chinanews.go`）：

```go
package rules

import (
    "net/http"

    "github.com/andeya/pholcus/app/downloader/request"
    spider "github.com/andeya/pholcus/app/spider"
)

func init() {
    mySpider.Register()
}

var mySpider = &spider.Spider{
    Name:         "静态规则示例",
    Description:  "静态规则示例 [Auto Page] [http://xxx.xxx.xxx]",
    EnableCookie: true,
    RuleTree: &spider.RuleTree{
        Root: func(ctx *spider.Context) {
            ctx.AddQueue(&request.Request{
                Url:  "http://xxx.xxx.xxx",
                Rule: "登录页",
            })
        },
        Trunk: map[string]*spider.Rule{
            "登录页": {
                ParseFunc: func(ctx *spider.Context) {
                    ctx.AddQueue(&request.Request{
                        Url:      "http://xxx.xxx.xxx",
                        Rule:     "登录后",
                        Method:   "POST",
                        PostData: "username=user@example.com&password=pass&submit=login",
                    })
                },
            },
            "登录后": {
                ParseFunc: func(ctx *spider.Context) {
                    ctx.Output(map[int]interface{}{
                        0: ctx.GetText(),
                    })
                    ctx.AddQueue(&request.Request{
                        Url:    "http://accounts.xxx.xxx/member",
                        Rule:   "个人中心",
                        Header: http.Header{"Referer": []string{ctx.GetUrl()}},
                    })
                },
            },
            "个人中心": {
                ParseFunc: func(ctx *spider.Context) {
                    ctx.Output(map[int]interface{}{
                        0: ctx.GetText(),
                    })
                },
            },
        },
    },
}
```

## 代理 IP

代理 IP 写在 `pholcus_pkg/proxy.lib` 文件中，格式如下（一行一个）：

```
http://183.141.168.95:3128
https://60.13.146.92:8088
http://59.59.4.22:8090
```

在操作界面选择"代理 IP 更换频率"或命令行设置 `-a_proxyminute` 参数即可启用。

> **注意：** Mac 下如使用代理 IP 功能，请务必获取 root 用户权限，否则无法通过 `ping` 获取可用代理。

## FAQ

**请求队列中，重复的 URL 是否会自动去重？**

URL 默认情况下是去重的，但可以通过设置 `Request.Reloadable = true` 忽略重复。

**URL 指向的页面内容若有更新，框架是否有判断机制？**

框架无法直接判断页面内容更新，但用户可以在规则中自定义支持。

**请求成功是依据 HTTP 状态码判断？**

不是判断状态码，而是判断服务器有无响应流返回。即 404 页面同样属于请求成功。

**请求失败后的重新请求机制？**

每个 URL 尝试下载指定次数后，若依然失败，则将该请求追加到一个类似 defer 性质的特殊队列中。在当前任务正常结束后，将自动添加至下载队列再次下载。如果依然失败，则保存至失败历史记录。下次执行该条爬虫规则时，可通过选择继承历史失败记录，把这些失败请求自动加入重试队列。
