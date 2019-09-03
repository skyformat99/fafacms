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

## 产品概述

1. 用户注册，填入相应信息如QQ，微博，邮箱，自我介绍，头像等，然后收到注册邮件，点击进行激活。未激活用户登陆后会显示未激活，无法使用平台。激活后用户可以登录后台编辑内容。用户注册后不提供注销功能。
2. 用户根高级权限控制，需要由管理员为用户分配用户组，用户组下有若干路由资源，路由资源均为特殊路由，如激活用户，更改其他用户密码，查看所有用户文章，用户信息等路由，如果用户不进入特殊资源路由，正常使用后台，否则需要具备相应的组权限。该功能为用户无感知隐藏功能。
3. 用户信息一般操作，用户登录后台，进入后台后可以随时退出登录以及补充注册时的用户信息，修改密码等。用户忘记密码可以通过邮件找回。
4. 内容编辑，用户可以创建内容节点，节点下可以有子节点，但最多两层，节点间实现了拖曳排序的功能，智能无比，在节点下可以新建文章，可以更新内容，设置隐藏文章，文章置顶，设置文章密码等，文章设计了特殊的历史版本功能，可以从历史版本恢复，并且可以对文章进行拖曳排序，以及拖曳移动到另外的节点目录。文章被删除可以回到回收站，回收站可以恢复。
5. 首页阅读和内容评论，其他用户可以浏览其他用户文章并进行评论，所有者可以设置自动审核，或者手动审核，通过的评论会被显示，评论有堆楼功能。评论可以由所有者删除，删除的评论及其子评论均会消失。其他用户也可以为内容或者内容的某条评论点赞或者反对，详细记录登陆用户点赞等情况，防止多次点赞。
6. 图片存储：用户头像，节点背景图，文章背景图等内部图片均需要通过上传接口保存进数据库，禁止使用不安全图片链接，图片存储在本地或者云对象存储服务中。
7. 内容编辑器使用markdown，插入图片时调用图片接口，抽取数据库已上传图片供编辑者选择，在此可以上传本地图片，并为图片打标签等。
8. 可以关闭用户注册，将用户加入黑名单，将内容封禁等。

新功能准备参见：[待做清单](/todo.md)。
其他详细设计，以及约束请参考实际可用产品及[文档-完善过程中](https://github.com/hunterhug/fafadoc)，[前端-完善过程中](https://github.com/hunterhug/fafafront)
。 

## 写给后端人员

见 [给后端开发详细的部署说明](/install/README.md)，强烈建议阅读。

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

## Update

### 201908

1. 可选择不将内容记录进历史。
2. 改为单点登录，任何用户登录会挤掉其他端的用户。
3. 用户访问临时令牌加长到7天有效。
4. 数据库切为utf8mb4完善文本保存。
5. 自动创建数据库，不必再手动创建。
6. 启动时打印版本号
7. 完善邮件发送逻辑，验证码有效期为5分钟，激活码重置不需要之前的激活码。
8. 完善用户注册，忘记密码，用户，用户组操作等前端API文档。
9. 完成用户注册，登录，忘记密码，注销前端UI。

### 201904-06

1. 支持本地文件模式，也支持阿里云对象存储文件上下载
2. 注册用户，用户忘记密码发邮件功能
3. 取消了cookie功能，让大前端自己实现，只实现自定义的用户session redis功能，改为token API模式，需授权API均检测该token，token和用户信息均保存在redis，缓存击穿时再从mysql加载
4. 实现自动创建数据库，数据库表，填充管理员账户。
5. 自动将admin URL置于变量内存中，避免管理员API权限过滤时，频繁查找数据库
6. 文章节点，文章等有菜单排序的功能，均支持强大拖拽排序
7. 实现自动docker hub打包镜像发布，提供从数据库到后端的一键部署脚本。


## 支持

微信支持:

![](/doc/support/weixin.jpg)

支付宝支持:

![](/doc/support/alipay.png)
