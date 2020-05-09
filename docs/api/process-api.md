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
      "enable": "bool"
    } 
  ]
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
