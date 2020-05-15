# Process 接口

## GET /api/1/process

获取进程列表

### Header 

```json5
{
  "Authorization": "string", //Token
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
      "name": "string",
      "path": "string",
      "dir": "string",
      "config": "string",
      "enable": "bool",
      "startTime": "string", //yyyy-MM-dd hh:mm:ss
    } 
  ]
}
```

## GET /api/1/process/{processId}

获取单个进程信息

### Header

```json5
{
  "Authorization": "string", //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": {
  "id": "int",
  "name": "string",
  "path": "string",
  "dir": "string",
  "config": "string",
  "enable": "bool",
  "startTime": "string", //yyyy-MM-dd hh:mm:ss
  } 
}
```

## PUT /api/1/process/{processId}/operation

操作进程

### Header

```json5
{
  "Authorization": "string", //Token
}
```

### Request

```json5
{
  "operation": "string" //start, stop, restart, enable, disable
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string"
}
```

## GET /api/1/process/{processId}/status

获取状态列表

### Header

```json5
{
  "Authorization": "string", //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": [
    {
      "processId": "int",
      "type": "string",
      "value": "string"
    } 
  ]
}
```

## GET /api/1/process/{processId}/status/{statusType}

获取状态

### Header

```json5
{
  "Authorization": "string", //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": {
    "processId": "int",
    "type": "string",
    "value": "string"
  }
}
```

## GET /api/1/status/{statusType}

获取状态列表

### Header

```json5
{
  "Authorization": "string" //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": [
    {
      "processId": "int",
      "type": "string",
      "value": "string"
    } 
  ]
}
```

## GET /api/1/process/statistic

获取进程统计

### Header 

```json5
{
  "Authorization": "string", //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": {
    "total": "int",
    "running": "int",
    "stopped": "int",
    "disable": "int"
  }
}
```

## PUT /api/1/process/{processId}/update-file/{fileId}

更新进程

### Header

```json5
{
  "Authorization": "string", //Token
  "Content-Type": "string", //application/x-zip-compressed
  "File-Name": "string"
}
```

### Request

升级压缩包

### Response

```json5
{
  "result": "bool",
  "error": "message"
}
```

## GET /api/1/process/{processId}/config-file

获取配置文件

### Header

```json5
{
  "Authorization": "string" //Token
}
```

### Response

```json5
{
  "result": "bool",
  "error": "message",
  "data": "string" //Base64
}
```

## PUT /api/1/process/{processId}/config-file

更新配置文件

### Header

```json5
{
  "Authorization": "string" //Token
}
```

### Request

```json5
{
  "content": "string" //Base64
}
```

### Response

```json5
{
  "result": "bool",
  "error": "message",
}
```
