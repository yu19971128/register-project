# api-conventions

## URL 设计
- 资源名使用复数
- 动作通过 HTTP Method 表达

## 请求规范
- Content-Type: `application/json`
- 时间格式：ISO 8601

## 响应规范
```json
{
  "code": 200,
  "data": {},
  "message": "ok"
}
```
