package rollout

import "testing"

func TestAllPlayersHitAtFullPercent(t *testing.T) {
	book := New(1, 2)
	result, err := book.Apply(1, "flag", 100, []string{"a", "b", "c"})
	if err != nil || len(result.Hit) != 3 {
		t.Fatalf("百分之百灰度应该全员命中，得到 %+v err=%v", result, err)
	}
}

func TestNoPlayersHitAtZeroPercent(t *testing.T) {
	book := New(1, 1)
	result, _ := book.Apply(1, "flag", -1, []string{"a", "b"})
	if len(result.Hit) != 0 {
		t.Fatalf("零灰度不该有人命中，得到 %+v", result)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(1, 1).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
