# Net接口

## GET /api/1/net/interface

获取网卡列表

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
      "index": "int",
      "name": "string",
      "description": "string",
      "mac": "string",
      "ip": "string",
      "mask": "string",
      "gateway": "string"
    }
  ]
}
```

## PUT /api/1/net/interface/{index}

设置网卡信息

### Header

```json5
{
  "Authorization": "string"
}
```

### Request

```json5
{
  "index": "int",
  "name": "string",
  "description": "string",
  "mac": "string",
  "ip": "string",
  "mask": "string",
  "gateway": "string"
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string"
}
```