# Binance SDK Examples

## 设置环境变量

首先设置你的API密钥：

```bash
export BINANCE_API_KEY="your-api-key"
export BINANCE_SECRET_KEY="your-secret-key"
```

## 可用示例

### 1. 账户信息查询 (account_info.go)

查询账户所有资产余额和提现配额：

```bash
go run examples/account_info.go
```

功能：
- 显示账户VIP等级
- 显示24小时提现限额（总额度/已使用/剩余）
- 列出所有资产余额
- 显示保存的提现地址

### 2. 小额提现 (withdrawal.go)

交互式小额提现工具：

```bash
go run examples/withdrawal.go
```

功能：
- 自动启用快速提现（币安内部转账免费）
- 显示可提现的币种
- 动态加载每个币种支持的网络
- 显示最小提现金额和手续费
- 使用withdrawOrderId跟踪订单状态
- 实时监控提现进度

## 注意事项

1. **API权限**：确保你的API密钥有以下权限：
   - 读取权限（查询余额）
   - 交易权限（启用快速提现）
   - 提现权限（执行提现）

2. **安全提醒**：
   - 不要将API密钥提交到代码仓库
   - 建议先用小额测试
   - 内部转账（转到其他币安账户）是免费且即时的

3. **提现状态码**：
   - 0: 邮件已发送
   - 1: 已取消
   - 2: 等待批准
   - 3: 已拒绝
   - 4: 处理中
   - 5: 失败
   - 6: 完成