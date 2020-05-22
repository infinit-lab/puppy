# log 接口

## GET /api/1/operate-log

获取操作日志

### Header

```json5
{
  "Authorization": "string" //Token
}
```

### Param

```json5
{
  "startTime": "string", //required
  "stopTime": "string", //required
  "username": "string", //optional
  "processId": "int", //optional
  "rows": "int", //required
  "offset": "int", //optional
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": [
    {
      "id": "int",
      "username": "string",
      "ip": "string",
      "operate": "string",
      "processId": "int",
      "processName": "string",
      "time": "string"
    }
  ]
}
```

## GET /api/1/login-log

获取登录日志

### Header

```json5
{
  "Authorization": "string", //Token
}
```

### Params

```json5
{
  "startTime": "string", //required
  "stopTime": "string", //required
  "rows": "int", //required
  "offset": "int", //optional
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": [
    {
      "id": "int",
      "username": "string",
      "ip": "string",
      "isLogin": "bool",
      "time": "string"
    } 
  ]
}
```
