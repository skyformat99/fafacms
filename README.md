# 花花CMS(FaFa CMS): 支持分布式部署的管理系统

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/fafacms.svg?style=social&label=Forks)](https://github.com/hunterhug/fafacms/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/fafacms.svg?style=social&label=Stars)](https://github.com/hunterhug/fafacms/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/fafacms)](https://goreportcard.com/report/github.com/hunterhug/fafacms)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms/issues)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu) 
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

[English](/README_EN.md)

## 项目说明

此项目代号为 `fafacms`。花花拼音 `fafa`，名称来源于广东话发发，花花的谐音，听起来有诙谐，娱乐等感觉，是一个使用 `Golang` 开发的前后端分离 --> 内容管理系统(CMS)。
目前基本的后端API以及配套文档都已经完成，大家可以从这个项目中学会 `Golang` 相关的开发技能，包括数据库操作，图片存储，邮件发送，容器部署等知识，五脏俱全的典型企业级应用，配置对象存储后可分布式多副本部署。

## 产品概述

1. 用户注册，填入相应信息如QQ，微博，邮箱，自我介绍，头像等，然后收到注册邮件，点击进行激活。未激活用户登陆后会显示未激活，无法使用平台。激活后用户可以登录后台编辑内容。用户注册后不提供注销功能。用户如果违禁被拉进黑名单不允许任何操作。
2. 用户超级管理员高级权限控制，需要由管理员为用户分配用户组，用户组下有若干超级管理员路由资源，路由资源均为特殊路由，如更改其他用户密码，查看所有用户文章，用户信息，拉黑违禁用户等路由，如果用户不进入特殊资源路由，正常使用后台，即只能操作自己的资源，否则需要具备相应的组权限。该功能为用户无感知隐藏功能。
3. 用户信息一般操作，用户登录后台，进入后台后可以随时退出登录以及补充注册时的用户信息，修改密码等。用户忘记密码可以通过邮件找回。用户昵称一个月只能修改两次，且全局唯一。
4. 内容编辑，用户可以创建内容节点，节点下可以有子节点，但最多两层，节点间实现了拖曳排序的功能，智能无比，在节点下可以新建文章，可以更新内容，设置隐藏文章，文章置顶，设置文章密码等，文章设计了特殊的发布机制和历史版本功能，文章内容先保存在预发布字段，点击发布按钮才真正刷新进正式字段，每次更新内容时可以将草稿保存进历史，每次发布时，会相应保存进发布历史，可以从历史内容版本中恢复等。同时可以对文章进行拖曳排序。文章实现二次删除，被删除时会移到回收站，可以从回收站恢复或彻底删除。
5. 首页阅读和内容评论，其他用户可以浏览其他用户文章并进行评论，内容所有者可以设置关闭或者开启评论，评论相对智能仿QQ音乐，评论可以由评论所有者删除。其他用户也可以为内容或者内容的某条评论点赞或者取消点赞，详细记录登陆用户点赞等情况，防止多次点赞。其他用户可以举报文章和评论。可以设置举报阈值自动将内容或评论违禁。
6. 文件存储功能：用户头像，节点背景图，文章背景图等内部图片均需要通过上传接口保存进数据库，禁止使用不安全外部图片链接，图片存储在本地或者云对象存储服务中。文件有相应的列出，分类打标签等API功能。
7. 可以关闭用户注册，将用户加入黑名单，将内容封禁等。
8. 互动消息站内信，如评论被点赞，内容被点赞，内容被评论，评论被评论。系统通知站内信，内容被违禁，评论被违禁，管理员通知。
9. 关注好友，好友发布内容站内信通知。
10. 私信。

需求时刻迭代，最新更新参见[产品更新记录](/log.md)，待处理事宜参见[待做清单](/todo.md)。详细设计，约束请参考实际[API文档](https://github.com/hunterhug/fafadoc)，及可用产品:[前端-完善过程中](https://github.com/hunterhug/fafafront)
。

## 写给前端人员

不关心部署，只想参与前端UI开发的看这里，请拥有一台类Unix机器，安装 `Docker`，`Docker-compose` 后一键部署。

```
# Linux使用install.sh
# Mac请使用install_mac.sh
cd install
chmod 777 install.sh
sudo ./install.sh
```

打开浏览器: `IP:8080` 进行开发，超级管理员账户密码：`admin/admin`

## 写给后端人员

关心部署细节，见 [给后端开发详细的部署说明](/install/README.md)，强烈建议阅读。

## 支持

微信支持:

![](/doc/support/weixin.jpg)

支付宝支持:

![](/doc/support/alipay.png)

## CopyRight

All right reserved. Attribution-NonCommercial-NoDerivatives 4.0 International.You can use it for education only but can't make profits for any companies and individuals!
