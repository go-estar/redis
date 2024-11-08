package redis

import (
	"context"
	"encoding/json"
	"github.com/go-estar/id-generator/snowflake"
	"github.com/go-estar/local-time"
	"github.com/go-estar/redis/script"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var client = New(&Config{
	Addr:     "127.0.0.1:6379",
	Password: "GBkrIO9bkOcWrdsC",
})

func TestRandom(t *testing.T) {
	var s = redis.NewScript(`
math.randomseed(ARGV[1])
return math.random(3)
`)
	for i := 0; i < 10; i++ {
		go func(i int) {
			var seed = localTime.Now().UnixNano()
			result, err := s.Run(context.Background(), client, []string{}, seed).Result()
			if err != nil {
				t.Log(i, "failed", err)
				return
			}
			t.Log(i, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestQuotaApply(t *testing.T) {
	var keyFields = make([][]interface{}, 0)
	keyFields = append(keyFields, []interface{}{"1-1-1", "C-1", 1, 2})
	keyFields = append(keyFields, []interface{}{"1-2-1", "C-1", 1, 1})

	keyFieldData, err := json.Marshal(keyFields)
	if err != nil {
		t.Fatal(err)
	}

	var orderIdWorker = snowflake.New(0)
	for i := 0; i < 100; i++ {
		go func(i int) {
			orderId := orderIdWorker.String()
			result, err := script.QuotaApply.Run(context.Background(), client, []string{"test-quota-apply", "test-quota"}, orderId, string(keyFieldData)).Result()
			if err != nil {
				t.Log(orderId, "failed", err)
				return
			}
			t.Log(orderId, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestQuotaRelease(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			orderId := 334
			result, err := script.QuotaRelease.Run(context.Background(), client, []string{"promotion-quota-apply", "promotion-quota", "promotion-quota-release"}, orderId).Result()
			if err != nil {
				t.Log(orderId, "failed", err)
				return
			}
			t.Log(orderId, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestCalculateApply(t *testing.T) {
	var (
		weightsStr = `[["12-12:E",14,1,10],["12-12:E",15,1,1]]`
	)
	var orderIdWorker = snowflake.New(0)
	for i := 0; i < 10; i++ {
		go func(i int) {
			var seed = localTime.Now().UnixNano()
			orderId := orderIdWorker.String()
			result, err := script.CalculateApply.Run(context.Background(), client, []string{"promotion-random-apply", "promotion-random"}, orderId, weightsStr, seed).Result()
			if err != nil {
				t.Log(orderId, "failed", err)
				return
			}
			t.Log(orderId, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestCouponApply(t *testing.T) {
	var keyFields = make([]interface{}, 0)
	keyFields = append(keyFields, "1", "2", "3")

	keyFieldData, err := json.Marshal(keyFields)
	if err != nil {
		t.Fatal(err)
	}

	var orderIdWorker = snowflake.New(0)
	for i := 0; i < 10; i++ {
		go func(i int) {
			orderId := orderIdWorker.String()
			result, err := script.CouponApply.Run(context.Background(), client, []string{"test-coupon-apply", "test-coupon"}, orderId, string(keyFieldData)).Result()
			if err != nil {
				t.Log(orderId, "failed", err)
				return
			}
			t.Log(orderId, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestCouponRelease(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			orderId := 1396420543510478842
			result, err := script.CouponRelease.Run(context.Background(), client, []string{"test-coupon-apply", "test-coupon", "test-coupon-release"}, orderId).Result()
			if err != nil {
				t.Log(orderId, "failed", err)
				return
			}
			t.Log(orderId, "succeed", result)
		}(i)
	}

	time.Sleep(time.Second * 5)
	t.Log("done")
}

func TestBatchGet(t *testing.T) {
	result, err := client.BatchGet(context.Background(), "test-01:1", "0758", "10759")
	if err != nil {
		t.Log("failed", err)
		return
	}
	t.Log("succeed", result)
}

func TestBatchLock(t *testing.T) {
	result, err := client.BatchLock(context.Background(), []string{"test-1", "test-3"}, 120)
	if err != nil {
		t.Log("failed", err)
		return
	}
	t.Log("succeed", result)
}

func TestBatchUnlock(t *testing.T) {
	result := client.BatchUnlock(context.Background(), []string{"test-1", "test-2", "test-3", "test-4", "test-5"})
	t.Log("result", result)
}

func TestBalanceCharge(t *testing.T) {
	result, err := script.BalanceCharge.Run(context.Background(), client, []string{"test-balance", "test-balance-history"},
		1, localTime.Now().UnixMillisecond(), 1000).Float64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestBalanceChargeRevoke(t *testing.T) {
	result, err := script.BalanceChargeRevoke.Run(context.Background(), client, []string{"test-balance", "test-balance-history"},
		1, "1730795081978", 1000).Float64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestBalanceConsume(t *testing.T) {
	result, err := script.BalanceConsume.Run(context.Background(), client, []string{"test-balance", "test-balance-history"},
		1, localTime.Now().UnixMillisecond(), -20).Float64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestBalanceConsumeRevoke(t *testing.T) {
	result, err := script.BalanceConsumeRevoke.Run(context.Background(), client, []string{"test-balance", "test-balance-history"},
		1, "1730795098131", -20).Float64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}
