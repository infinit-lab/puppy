# License接口

## Get /api/1/fingerprint

获取机器指纹

### Header

```json5
{
  "Authorization": "string" //Token
}
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": "string"
}
```

## Get /api/1/license

获取授权证书状态

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
  "data": {
    "status": "int",
    "isForever": "bool",
    "validDateTime": "string",
    "validDuration": "int", //second
  }
}
```

## Put /api/1/license

上传授权证书

### Request

授权证书文件

### Response

```json5
{
  "result": "bool",
  "error": "string"
}
```
