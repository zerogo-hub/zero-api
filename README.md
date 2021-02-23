## zero-web

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
- 示例: `/blog/list/:id|isNum|less4|`
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
