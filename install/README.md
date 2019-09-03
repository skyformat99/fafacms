# 给后端开发详细的部署说明

最快速部署请安装`docker`,`docker-compose`, 然后直接一键脚本：

```
git clone github.com/hunterhug/fafacms
cd fafacms

# Linux使用install.sh
# Mac请使用install_mac.sh
cd install

# 修改邮箱地址等
vim config.json

chmod 777 install.sh
sudo ./install.sh
```

主要集成了`mysql5.7.27`，`phpmyadmin:edge-4.9`和`redis:5.0.5`，端口分别为`3306`，`8000`，`6379`，

`MYSQL`账号密码：`root/123456789`,`Redis`密码：`123456789`，打开`IP:8000`登录数据库进行查看。

持久卷将会挂载在 `/data/mydocker` 中。具体配置和挂载卷可修改`docker-compose.yaml`和`config.json`文件。

运行后，请打开`IP:8080`进行API对接，超级管理员账户密码：`admin/admin`

# 详细说明

## 后端部署(常规)

获取代码:

```
go get -v github.com/hunterhug/fafacms
```

代码就会保存在`Golang GOPATH`目录下.

运行:

```
fafacms -config=./config.json
```

其中`config.json`说明如下（具体参考实际配置）:

```
{
  "DefaultConfig": {
    "WebPort": ":8080",                 # 程序运行端口(可改)
    "StorageOss": false,                # 文件是否保存在对象存储，默认否，true时 OssConfig 有效
    "StoragePath": "./data/storage",    # 文件保存在本地地址(可改)
    "LogPath": "./data/log/fafacms_log.log",        # 日志保存地址(可改)
    "LogDebug": true,   					        # 打开调试(建议保持为true)
    "CloseRegister": false                          # 是否关闭注册功能
  },
  "OssConfig": {
    "Endpoint": "oss-cn-qingdao.aliyuncs.com",      # 对象存储配置（区域，桶和密钥对）
    "BucketName": "syoss",
    "AccessKeyId": "",
    "AccessKeySecret": ""
  },
  "DbConfig": {
    "DriverName": "mysql",      # 关系型数据库驱动(不能改，等功能拓展可支持不同驱动)
    "Name": "fafa",             # 关系型数据库名字(可改)
    "Host": "127.0.0.1",        # 关系型数据库地址(可改)
    "User": "root",             # 关系型数据库用户(可改)
    "Pass": "123456789",        # 关系型数据库密码(可改)
    "Port": "3306",             # 关系型数据库端口(可改)
    "MaxIdleConns": 20,         # 关系型数据库池闲置连接数(默认保持)
    "MaxOpenConns": 20,         # 关系型数据库池打开连接数(默认保持)
    "DebugToFile": true,        # SQL调试是否输出到文件(默认保持)
    "DebugToFileName": "./data/log/fafacms_db.log",     # SQL调试输出文件路径(默认保持)
    "Debug": true                                       # SQL调试(默认保持)
  },
  "Email": {
    "Host": "smtp-mail.outlook.com",    # 忘记密码，激活用户时发邮件服务器
    "Port": 587,                        # 邮件服务器端口
    "Email": "gdccmcm14@live.com",      # 邮箱账号
    "Password": "",                     # 邮箱密码
    "Subject": "FaFa CMS Code",         # 邮箱发送主题
    "Body": "%s Code is <br/> <p style='text-align:center'>%s</p> <br/>Valid in 5 minute."  # 邮箱内容，两个占位符，第二个%s为验证码，第一个是字符串功能。
  },
  "SessionConfig": {
    "RedisHost": "127.0.0.1:6379",  # Redis地址(可改)
    "RedisMaxIdle": 64,             # (默认保持)
    "RedisMaxActive": 0,            # (默认保持)
    "RedisIdleTimeout": 120,        # (默认保持)
    "RedisDB": 0,                   # Redis默认连接数据库(默认保持)
    "RedisPass": "123456789"        # Redis密码(可为空,可改)
  }
}
```

具体命令参数如下：

```
  -auth_skip_debug
        Auth skip debug
  -config string
        config file (default "./config.json")
  -email_debug
        Email debug
  -history_record
        Content history record
  -init_db
        create db table (default true)
```

正常启动如下：

```
./fafacms config=/root/fafacms/config.json -history_record=true -init_db=false
```

表示文章内容开启历史记录功能，并且关闭数据库数据填充（第二次启动时可设置为false）。

## 后端部署(Docker)

你也可以使用`docker`进行部署, 构建镜像(Docker版本必须大于17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

先新建数据卷, 并且移动配置并修改:

```
mkdir /root/fafacms
cp docker_config.json /root/fafacms/config.json
```

启动容器:

```
sudo docker run -d --name fafacms -p 8080:8080 -v /root/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.json -history_record=true -init_db=true" hunterhug/fafacms

sudo docker logs -f --tail 10 fafacms
```

其中`/root/fafacms`是挂载的持久化卷, 配置`config.json`放置在该文件夹下.

开发中`Debug`:

```
sudo docker run -d --name fafacms -p 8080:8080 -v /root/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.json -email_debug=true -auth_skip_debug=true -history_record=true -init_db=true" hunterhug/fafacms
```