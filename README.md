# hotfix-rollout

热更开关灰度：按百分比把玩家分批命中，并发执行，结果必须与并发度无关，只用 Go 标准库。

```
go run ./cmd/hotfix-rollout --sample=stable
go run ./cmd/hotfix-rollout --sample=version
go run ./cmd/hotfix-rollout --sample=sorted
go run ./cmd/hotfix-rollout --sample=work
go test ./...
go vet ./...
```

## 口径（README 为准）

- **版本单调**：只有传入版本严格大于当前版本的开关才生效，否则返回 `ErrStaleVersion` 且不改状态。
- **命中判定**：`fnv1a(开关名 + "|" + 玩家名) % 100` 严格小于百分比才算命中，正好等于不算。
- **与并发度无关**：同一输入无论 workers 取多少，命中集合与 Miss 完全一致。
- **顺序**：命中名单按玩家名升序返回，与 goroutine 完成顺序无关。
- **不变量**：命中数加 Miss 等于玩家总数；同一输入重复执行结果相同。
- **代价**：`Apply` 是 O(玩家数) 摊还，`Scanned` 不随玩家数平方增长；10 万玩家每秒 10 万次灰度。

## 输出契约（不改格式）

```
{"same":..}
{"first_ok":..,"stale_ok":..}
{"sorted":..}
{"scanned":..,"hit":..}
```

## 目录

```
rollout/rollout.go               版本、哈希桶与并发汇总
cmd/hotfix-rollout/main.go       命令行入口
rollout/rollout_test.go          单元用例
```
