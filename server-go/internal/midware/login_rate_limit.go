package midware

import (
	"net/http"
	"sync"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"

	"github.com/gin-gonic/gin"
)

// ========== v0.7.0 新增：登录限流 + IP 锁定 ==========
//
// 策略：
//   1. 同一 IP 在「短窗口」内登录失败次数超过阈值 → 锁定一段时间
//   2. 锁定期间所有登录尝试直接拒绝，不消耗验证资源
//   3. 登录成功后立即清空该 IP 的失败计数
//
// 仅针对 POST /api/v1/user/login/account。不持久化（进程内内存即可），避免 SQLite 压力。

const (
	// LoginMaxFailPerWindow 登录失败阈值
	LoginMaxFailPerWindow = 5
	// LoginFailWindow 滑动窗口大小
	LoginFailWindow = 5 * time.Minute
	// LoginLockDuration 达到阈值后的锁定时长
	LoginLockDuration = 15 * time.Minute
	// LoginFailRecordMax 内存中最多保留的 IP 数（简易 LRU 防内存爆炸）
	LoginFailRecordMax = 10000
	// LoginFailGCInterval 过期清理间隔
	LoginFailGCInterval = 10 * time.Minute
)

type loginAttempt struct {
	failures  []time.Time // 失败时间点（滑动窗口）
	lockUntil time.Time   // 锁定解除时间（零值表示未锁定）
}

var (
	loginAttemptMap  = make(map[string]*loginAttempt)
	loginAttemptLock sync.Mutex
	loginGCOnce      sync.Once
)

// initLoginRateLimitGC 启动后台 GC（懒启动）
func initLoginRateLimitGC() {
	loginGCOnce.Do(func() {
		go func() {
			t := time.NewTicker(LoginFailGCInterval)
			defer t.Stop()
			for range t.C {
				gcLoginAttempts()
			}
		}()
	})
}

func gcLoginAttempts() {
	loginAttemptLock.Lock()
	defer loginAttemptLock.Unlock()

	now := time.Now()
	cutoff := now.Add(-LoginFailWindow)
	for ip, a := range loginAttemptMap {
		// 清掉窗口外的失败记录
		fresh := a.failures[:0]
		for _, ts := range a.failures {
			if ts.After(cutoff) {
				fresh = append(fresh, ts)
			}
		}
		a.failures = fresh
		// 锁定已解除 && 无近期失败 → 删除记录
		if len(a.failures) == 0 && (a.lockUntil.IsZero() || a.lockUntil.Before(now)) {
			delete(loginAttemptMap, ip)
		}
	}

	// 若记录数仍超上限，粗暴清理（实际使用中极少触发）
	if len(loginAttemptMap) > LoginFailRecordMax {
		ylog.Warnf("LoginRateLimit", "loginAttemptMap size %d exceeds limit, force clear", len(loginAttemptMap))
		loginAttemptMap = make(map[string]*loginAttempt)
	}
}

// LoginRateLimitCheck 判断某 IP 是否被锁定；未锁定返回 nil，否则返回锁定剩余时间
func LoginRateLimitCheck(ip string) (locked bool, remaining time.Duration) {
	loginAttemptLock.Lock()
	defer loginAttemptLock.Unlock()

	a, ok := loginAttemptMap[ip]
	if !ok {
		return false, 0
	}
	now := time.Now()
	if !a.lockUntil.IsZero() && a.lockUntil.After(now) {
		return true, a.lockUntil.Sub(now)
	}
	return false, 0
}

// RecordLoginFailure 记录一次登录失败，若达到阈值则锁定
func RecordLoginFailure(ip string) {
	if ip == "" {
		return
	}
	loginAttemptLock.Lock()
	defer loginAttemptLock.Unlock()

	now := time.Now()
	cutoff := now.Add(-LoginFailWindow)

	a, ok := loginAttemptMap[ip]
	if !ok {
		a = &loginAttempt{}
		loginAttemptMap[ip] = a
	}

	// 清理窗口外失败
	fresh := a.failures[:0]
	for _, ts := range a.failures {
		if ts.After(cutoff) {
			fresh = append(fresh, ts)
		}
	}
	a.failures = append(fresh, now)

	if len(a.failures) >= LoginMaxFailPerWindow {
		a.lockUntil = now.Add(LoginLockDuration)
		ylog.Warnf("LoginRateLimit", "IP %s locked until %s after %d failures",
			ip, a.lockUntil.Format(time.RFC3339), len(a.failures))
	}
}

// RecordLoginSuccess 登录成功时清理该 IP 的失败计数
func RecordLoginSuccess(ip string) {
	if ip == "" {
		return
	}
	loginAttemptLock.Lock()
	defer loginAttemptLock.Unlock()
	delete(loginAttemptMap, ip)
}

// LoginRateLimit 中间件：锁定 IP 直接拒绝；未锁定则放行，由 handler 在登录失败时调用 RecordLoginFailure
func LoginRateLimit() gin.HandlerFunc {
	initLoginRateLimitGC()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if locked, remaining := LoginRateLimitCheck(ip); locked {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"errorCode": common.AuthFailedErrorCode,
				"msg":       "too many failed attempts, please try again later",
				"data": gin.H{
					"lockedRemainingSeconds": int(remaining.Seconds()),
				},
			})
			return
		}
		c.Next()
	}
}
