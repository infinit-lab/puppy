# Token 接口

## POST /api/1/token 

创建Token-1

### Request

```json
{
  "username": "string",
  "password": "string"
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
  "data": "string" //Token
}
```

## DELETE /api/1/token/{token}

续约Token-1

### Response

```json5
{
  "result": "bool",
  "error": "string",
}
```
