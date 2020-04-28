# Account 接口

## PUT /api/1/password/{username}

修改密码

### Header

```json5
{
  "Authorization": "string" //Token
}
```

### Request

```json5
{
  "origin": "string", //原始密码
  "new": "string", //新密码
}
```

### Response

```json5
{
  "result": "bool",
  "error": "string",
}
```