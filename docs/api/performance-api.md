# Performance 接口

## GET /api/1/performance/cpu

获取cpu利用率

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
  "data": "int" //CPU利用率
}
```

## GET /api/1/performance/mem

获取内存利用率

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
    "rate": "int",
    "total": "int",
    "avail": "int"
  }
}
```