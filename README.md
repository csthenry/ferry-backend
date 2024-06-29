<center><h1>Ferry (neo)</h1></center>

## 基于Gin + Vue + Element UI前后端分离的工单系统

**流程中心**

通过灵活的配置流程、模版等数据，非常快速方便的生成工单流程，通过对流程进行任务绑定，实现流程中的钩子操作，目前支持绑定邮件来通知处理，当然为兼容更多的通知方式，也可以自己写任务脚本来进行任务通知，可根据自己的需求定制。

兼容了多种处理情况，包括串行处理、并行处理以及根据条件判断进行节点跳转。

可通过变量设置处理人，例如：直接负责人、部门负责人、HRBP等变量数据。

**系统管理**

基于casbin的RBAC权限控制，借鉴了go-admin项目的前端权限管理，可以在页面对API、菜单、页面按钮等操作，进行灵活且简单的配置。

## 功能介绍

<!-- wp:paragraph -->
<p>下面对本系统的功能做一个简单介绍。</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>工单系统相关功能：</p>
<!-- /wp:paragraph -->

<!-- wp:list -->
<ul><li>工单提交申请</li><li>工单统计</li><li>多维度工单列表，包括（我创建的、我相关的、我待办的、所有工单）</li><li>自定义流程</li><li>自定义模版</li><li>任务钩子</li><li>任务管理</li><li>催办</li><li>转交</li><li>手动结单</li><li>加签</li><li>多维度处理人，包括（个人，变量(创建者、创建者负责人)）</li><li>排他网关，即根据条件判断进行工单跳转</li><li>并行网关，即多个节点同时进行审批处理</li><li>通知提醒（目前仅支持邮件）</li><li>流程分类管理</li></ul>
<!-- /wp:list -->

<!-- wp:paragraph -->
<p>权限管理相关功能，使用casbin实现接口权限控制：</p>
<!-- /wp:paragraph -->

<!-- wp:list -->
<ul><li>用户、角色、岗位的增删查改，批量删除，多条件搜索</li><li>角色、岗位数据导出Excel</li><li>重置用户密码</li><li>维护个人信息，上传管理头像，修改当前账户密码</li><li>部门的增删查改</li><li>菜单目录、跳转、按钮及API接口的增删查改</li><li>登录日志管理</li><li>左菜单权限控制</li><li>页面按钮权限控制</li><li>API接口权限控制</li></ul>
<!-- /wp:list -->

## 交流

加群条件是需给项目一个star，不需要您费多大的功夫与力气，一个小小的star是作者能维护下去的动力。

## 特别感谢
[go-amdin # 不错的后台开发框架](https://github.com/go-admin-team/go-admin)

[vue-element-admin # 不错的前端模版框架](https://github.com/PanJiaChen/vue-element-admin)

[vue-form-making # 表单设计器，开源版本比较简单，如果有能力的话可以自己进行二次开发 ](https://github.com/GavinZhuLei/vue-form-making.git)

[wfd-vue # 流程设计器](https://github.com/guozhaolong/wfd-vue)

[machinery # 任务队列](https://github.com/RichardKnop/machinery.git)

## 鸣谢

特别感谢 [JetBrains](https://www.jetbrains.com/?from=ferry) 为本开源项目提供免费的 [IntelliJ GoLand](https://www.jetbrains.com/go/?from=ferry-web) 授权

## License

开源不易，请尊重作者的付出，感谢。

在此处声明，本系统目前不建议商业产品使用，因本系统使用的`流程设计器`未设置开源协议，`表单设计器`是LGPL v3的协议。

因此避免纠纷，不建议商业产品使用，若执意使用，请联系原作者获得授权。

再次声明，若是未联系作者直接将本系统使用于商业产品，出现的商业纠纷，本系统概不承担，感谢。

[LGPL-3.0](https://github.com/csthenry/ferry-web/blob/master/LICENSE)

Copyright (C) 2024 Henry 基于 ferry@lanyulei 二次开发
