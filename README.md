## zero-api

小巧的 api

## 主要文件

- app.go 应用，负责整合`Server`和`Router`
- context.go 一次请求的上下文
- cookie.go `Cookie`相关
- group.go 组路由，即一组由相同前缀组成的路由
- option.go 应用的选项设置
- route.go 每一个`Method(GET/POST/PUT...)`的路由由一个`Route`存储
- router.go 路由管理器
- server.go `HTTP Server`

## 路由

静态路由

- 示例: `/blog/list`

动态路由

- 格式: `:param`
- 示例: `/blog/list/:id`
  - `/blog/list/1001` 匹配, id="1001"
  - `/blog/list/p1001` 匹配, id="p1001"
  - `/blog/list` 不匹配
  - `/blog/list/1001/add` 不匹配

动态路由，带正则表达式

- 格式: `:param(regexp)`
- 示例: `/blog/list/:id(^\d+$)`
  - `/blog/list/1001` 匹配，id="1001"
  - `/blog/list/p1001` 不匹配
  - `/blog/list` 不匹配
  - `/blog/list/1001/add` 不匹配

动态路由，带验证函数

- 格式: `:param|validator...|`，验证函数必须包裹在`|`内
- 备注: 框架自带常用验证函数，也可以自行定义
- 示例: `/blog/list/:id|isNum|less4|`，id 为数字且小于 4 位数
  - `/blog/list/1001` 匹配，id="1001"
  - `/blog/list/1000001` 不匹配
  - `/blog/list/p101` 不匹配

动态路由，混合各种类型

- 格式: `:param(regexp)|validator...|`
- 示例: `/blog/:id(^\d+$)|less4|`
  - `/blog/100` 匹配
  - `/blog/1001` 不匹配

## 中间件

共有三种，添加方式如下

- 应用级别中间件，作用在所有路由中
- 组路由级别中间件，作用在该组路由中
- 路由级别中间件，作用在当前路由中
