// Package rollout 把热更开关按百分比灰度给玩家，结果必须与并发度无关。
// 缺陷：命中判定掺进了 worker 编号、版本不校验、结果按完成顺序返回而不是玩家名升序、
// 汇总时对每个玩家重扫已有结果。
package rollout

import (
	"errors"
	"hash/fnv"
	"sort"
	"sync"
)

// ErrStaleVersion 传入版本不比当前新。
var ErrStaleVersion = errors.New("版本不比当前新")

// Result 是一次灰度结果。
type Result struct {
	Version int
	Hit     []string
	Miss    int
}

// Rollout 是灰度台账。
type Rollout struct {
	version int
	workers int
	scanned int
}

// New 建台账：version 是当前开关版本，workers 是并发度。
func New(version, workers int) *Rollout {
	if workers < 1 {
		workers = 1
	}
	return &Rollout{version: version, workers: workers}
}

// Apply 按百分比把玩家分批命中，返回结果。缺陷：哈希掺了 worker 编号、不校验版本、结果未排序、汇总平方级。
func (r *Rollout) Apply(version int, flag string, percent int, players []string) (Result, error) {
	type chunk struct {
		index   int
		players []string
	}
	chunks := []chunk{}
	size := (len(players) + r.workers - 1) / r.workers
	for start := 0; start < len(players); start += size {
		end := start + size
		if end > len(players) {
			end = len(players)
		}
		chunks = append(chunks, chunk{index: len(chunks), players: players[start:end]})
	}
	out := make([][]string, len(chunks))
	var wait sync.WaitGroup
	for _, item := range chunks {
		wait.Add(1)
		go func(item chunk) {
			defer wait.Done()
			hit := []string{}
			for _, player := range item.players {
				hasher := fnv.New32a()
				hasher.Write([]byte(flag + "|" + player))
				hasher.Write([]byte{byte(item.index)})
				if int(hasher.Sum32()%100) <= percent {
					hit = append(hit, player)
				}
			}
			out[item.index] = hit
		}(item)
	}
	wait.Wait()
	result := Result{Version: r.version}
	r.version = version
	for _, part := range out {
		result.Hit = append(result.Hit, part...)
	}
	result.Miss = len(players) - len(result.Hit)
	r.scanned += len(players) * len(result.Hit)
	return result, nil
}

// Sorted 返回排序后的命中名单（口径要求结果升序）。
func (r Result) Sorted() []string {
	out := append([]string{}, r.Hit...)
	sort.Strings(out)
	return out
}

// Scanned 是累计的比较次数（规模观测）。
func (r *Rollout) Scanned() int { return r.scanned }
