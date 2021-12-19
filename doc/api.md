# API

## User management

### Create

```shell
❯ http POST :8080/v1/user nick_name=vink phone_number='+1 123 456 789 00' -v                                                                                                                             05:15:06
POST /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 58
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "nick_name": "vink",
    "phone_number": "+1 123 456 789 00"
}

HTTP/1.1 200 OK
Content-Length: 66
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:15:24 GMT

{
    "user": {
        "id": 0,
        "nick_name": "vink",
        "phone_number": "112345678900"
    }
}
```

### Update

```shell
http PUT :8080/v1/user id:=0 phone_number='+1 123 456 789 42' -v                                                                                                                                       05:16:49
PUT /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 46
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "id": 0,
    "phone_number": "+1 123 456 789 42"
}

HTTP/1.1 200 OK
Content-Length: 66
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:18:25 GMT

{
    "user": {
        "id": 0,
        "nick_name": "vink",
        "phone_number": "112345678942"
    }
}

```

### Search

```shell
 http :8080/v1/users nick_name==vink -v                                                                                                                                                                 05:19:19
GET /v1/users?nick_name=vink HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 69
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:19:32 GMT

{
    "users": [
        {
            "id": 0,
            "nick_name": "vink",
            "phone_number": "112345678942"
        }
    ]
}
```

### Delete

```shell
❯ http DELETE :8080/v1/user id:=0 -v                                                                                                                                                                     05:19:32
DELETE /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 9
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "id": 0
}

HTTP/1.1 200 OK
Content-Length: 66
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:20:22 GMT

{
    "user": {
        "id": 0,
        "nick_name": "vink",
        "phone_number": "112345678942"
    }
}

```

## Contacts

### Connect

```shell
http POST :8080/v1/contact from:=0 to:=1 -v                                                                                                                                                            05:22:10
POST /v1/contact HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 20
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "from": 0,
    "to": 1
}

HTTP/1.1 200 OK
Content-Length: 2
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:22:29 GMT

{}
```

### Lookup

```shell
❯ http :8080/v1/contact/1 -v                                                                                                                                                                             05:22:30
GET /v1/contact/1 HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 72
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:23:15 GMT

{
    "Contacts": [
        {
            "id": 0,
            "nick_name": "vink",
            "phone_number": "112345678900"
        }
    ]
}
```

### Reverse lookup

```shell
❯ http :8080/v1/contact/1/reverse -v                                                                                                                                                                     05:24:55
GET /v1/contact/1/reverse HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 133
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:24:59 GMT

{
    "Contacts": [
        {
            "id": 2,
            "nick_name": "zulgrar",
            "phone_number": "112345678988"
        },
        {
            "id": 0,
            "nick_name": "vink",
            "phone_number": "112345678900"
        }
    ]
}

```

### Suggest

```shell
❯ http :8080/v1/contact/1/suggest -v                                                                                                                                                                     05:25:13
GET /v1/contact/1/suggest HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 72
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:26:11 GMT

{
    "Contacts": [
        {
            "id": 1,
            "nick_name": "taek",
            "phone_number": "112345678942"
        }
    ]
}
```


## Dev

### Stats

```shell
❯ http :8080/v1/dev/stats -v                                                                                                                                                                             05:26:41
GET /v1/dev/stats HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 46
Content-Type: text/plain; charset=utf-8
Date: Sun, 19 Dec 2021 05:26:41 GMT

{
    "connections": 3,
    "user_connected": 3,
    "users": 3
}
```

### Generate

```shell
POST /v1/dev/generate HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 42
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "connection": 50,
    "population": 10000000
}

HTTP/1.1 202 Accepted
Content-Length: 0
Date: Sun, 19 Dec 2021 05:27:29 GMT

```