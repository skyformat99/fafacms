# FaFa CMS

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/fafacms.svg?style=social&label=Forks)](https://github.com/hunterhug/fafacms/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/fafacms.svg?style=social&label=Stars)](https://github.com/hunterhug/fafacms/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/fafacms)](https://goreportcard.com/report/github.com/hunterhug/fafacms)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms/issues)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu) 
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

## Project description

`fafa` -- means `flower` in Cantonese.

A content management system written in go, which frontend and backend is highly splited. Support multi-users, post blogs, view blogs. We hope to bring a practices that are generally desired across the industry.

Backend returns JSON API.You can use any mainstream frameworks to develop frontend. This project framework is scalable.

Dependencies:

1. [Gin](https://github.com/gin-gonic/gin)
2. [Go Validator](https://github.com/go-playground/validator)

...

Structure:

```
├── config.json 
├── core    	# backend files
│   ├── config      
│   ├── flog        
│   ├── controllers 
│   ├── model       
│   ├── router     
│   ├── server      
│   └── util        
├── main.go 	# entrance
└── web  		# frontend files
```

API DOC：[https://github.com/hunterhug/fafadoc](https://github.com/hunterhug/fafadoc)

## Product description

1. User registration with personal information(qq, weibo, email and other profile data), then receive the email and activate the account. Inactivated users could not use the platform. No close functions so far.
2. Privilege management. Administrator could set up user group and routing access, like activate user, change other users' password, check all articles. This is invisible function.
3. User could manage their profile. If user forget password, they can use their e-mail to reset password. 
4. Content management. User could create two layers tags at most. They can create new posts, hide or update posts, recover post from specific version, drag to different tags.
5. Comment management. User could comment on other's article, they can delete their own comments, or vote for others' comment.
6. Image storage. All image needs to upload to database via interface. 
7. Using Markdown as content editor. Users can choose existing images on database or upload local images.
8. User registration could be disabled. Blocklist of user and content implemented.

More details and restrictions are available in production.

## Instruction

### Easy Install

Install docker, docker-compose, Then:

```
cd install
chmod 777 install.sh
sudo ./install.sh
```

Account/password：admin/admin, url: http://IP:8080.

### Backend deployment(normal)

get codes:

```
go get -v github.com/hunterhug/fafacms
```

Then the repository will be downloaded and could be found inside `Golang GOPATH`.

run:

```
fafacms -config=./config.json
```

description of`config.json`:

```
{
  "DefaultConfig": {
    "WebPort": ":8080", 				    	# Port for project(optional)
    "StoragePath": "./data/storage",  # Path for file saving(optional)
    "LogPath": "./data/log/fafacms_log.log", 	# Log saving path(optional)
    "LogDebug": true   					        # Debug(default)
  },
  "DbConfig": {
    "DriverName": "mysql",  	# Relational DB driver(default)
    "Name": "fafa", 					# DB name(optional)
    "Host": "127.0.0.1", 			# DB host(optional)
    "User": "root", 					# DB user(optional)
    "Pass": "123456789", 			# DB password(optional)
    "Port": "3306", 					# DB port(optional)
    "MaxIdleConns": 20, 			# Max Idle connections(default)
    "MaxOpenConns": 20, 			# Max Idle connections(default)
    "DebugToFile": true, 			# Debug output files(default)
    "DebugToFileName": "./data/log/fafacms_db.log", # SQL output file path(default)
    "Debug": true 										# sql Debug(default)
  },
  "SessionConfig": {
    "RedisHost": "127.0.0.1:6379", 		# RedisHost(optional)
    "RedisMaxIdle": 64, 							# (default)
    "RedisMaxActive": 0, 							# (default)
    "RedisIdleTimeout": 120, 					# (default)
    "RedisDB": 0, 										# Redis connect database(default)
    "RedisPass": "123456789"   				# Redis password(optional, optional)
  }
}
```

### Backend deployment(Docker)

We can also use `docker` to deploy, construct the image(Docker version must later than 17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

Make Dir and add config file:

```
mkdir /root/fafacms
cp docker_config.json /root/fafacms/config.json
```

Initialize container:

```
sudo docker run -d --name fafacms -p 8080:8080 -v /root/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.json" hunterhug/fafacms

sudo docker logs -f --tail 10 fafacms
```

`/root/fafacms` is persistent volume, please put `config.json` under the folder.

## Frontend Web

doing...