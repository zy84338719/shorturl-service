package shorturl

import (
	"errors"
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"shorturl/internal/model/short"
	"strconv"
	"strings"
	"time"
)

func hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func hash32(data []byte) uint32 {
	return murmur3.Sum32(data)
}

func getShortCode(key, url string, ct int /*冲突标记*/) string {
	//加密字符传前的混合 KEY
	// 要使用生成 URL 的字符
	// 对传入网址进行 MMH3 加密
	h32 := hash32([]byte(key + url))
	if ct != 0 {
		h32 = hash32([]byte(key + url + string(ct)))
	}

	// 把加密字符按照 8 位一组 16 进制与 0x3FFFFFFF 进行位与运算
	lHexLong := 0x3FFFFFFF & h32
	outChars := ""
	//循环获得每组6位的字符串
	for j := 0; j < 6; j++ {
		// 把得到的值与 0x0000003E 这是62 进行mod运算，取得字符数组 chars 索引(具体需要看chars数组的长度   以防下标溢出，注意起点为0)
		index := lHexLong % 0x0000003E
		// 把取得的字符相加
		outChars += chars[index]
		// 每次循环按位右移 5 位
		lHexLong = lHexLong >> 5
	}
	return outChars
}

func getTempShortCode(key, url string, ct int /*冲突标记*/) string {
	return getShortCode(key, url, ct) + createShortExtra()
}

func judgePermanentShort(code string) bool {
	if len(code) > 6 {
		return false
	}
	return true
}

func getExpire(exprieInMinutes int64, exprieTime string) (exprie *time.Time, err error) {
	exp := time.Now()
	if exprieInMinutes != 0 {
		exp = exp.Add(time.Minute * time.Duration(exprieInMinutes))
	} else if exprieInMinutes == -1 {
		return nil, nil
	} else {
		exp, err = time.Parse("2006-01-02 15:04:05", exprieTime)
		if err != nil {
			return nil, err
		}
	}
	return &exp, nil
}

// 这块可以做分库分表设计
func createShortExtra() string {
	x := getDays() % 62
	return chars[x]
}

func getDays() int64 {
	yearMonthDay := time.Date(2023, 01, 01, 0, 0, 0, 0, time.Local)
	return (time.Now().Unix() - yearMonthDay.Unix()) / 3600 / 24
}

var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h",
	"i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
	"u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H",
	"I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z"}

func switchShortDao(db *gorm.DB, cache *redis.Redis, daoName string) short.ShortDaoInterface {
	switch daoName {
	case "url":
		return short.NewUrlDao(db, cache)
	case "data":
		return short.NewDataDao(db, cache)
	case "url_permanent":
		return short.NewPermanentUrlDao(db, cache)
	case "data_permanent":
		return short.NewPermanentDataDao(db, cache)
	default:
		return short.NewUrlDao(db, cache)
	}
}

func checkExprieTime(exprieTime *time.Time, status int) (err error) {
	if exprieTime == nil || status != 2 {
		return nil
	}
	if status == 2 || time.Now().Sub(*exprieTime) > 0 {
		err = errors.New("已经失效")
		return
	}
	return
}

func TimerCheckStatus(db *gorm.DB, cache *redis.Redis) {
	m := []string{
		"url",
		"data",
	}
	// 每小时检查一次 status = 1 的数据是否过期
	for range time.NewTicker(time.Hour).C {
		count, err := switchShortDao(db, cache, "").GetExpireCount()
		if err != nil {
			continue
		}
		for _, k := range m {
			for i := 0; i < int(count/1000)+1; i++ {
				urls, err := switchShortDao(db, cache, k).ExpireList(i*1000, 1000)
				if err != nil {
					continue
				}
				for _, url := range urls {
					_ = switchShortDao(db, cache, k).UpdateByCodeStatus(url.Code, 2)
				}
			}
		}
	}
}

func TimerUpdateCount(db *gorm.DB, cache *redis.Redis) {

	for {
		t := time.Now()
		next := t.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 1, 0, 0, next.Location())
		timer := time.NewTimer(next.Sub(t))
		<-timer.C

		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local)
		yesterday := date.Add(-time.Hour)

		smembers, err := cache.Smembers(getSetKey(yesterday))
		if err != nil {
			continue
		}
		for _, smember := range smembers {
			str, err := cache.Get(smember)
			if err != nil {
				continue
			}
			count, err := strconv.Atoi(str)
			if err != nil {
				continue
			}
			split := strings.Split(smember, "_")
			nameDao := split[0]
			if judgePermanentShort(smember) {
				nameDao += "_permanent"
			}
			_ = switchShortDao(db, cache, nameDao).UpdateByCodeCount(split[1], count)
		}

	}

}

func addCount(cache *redis.Redis, code string, typeName string) {
	now := time.Now()
	_, _ = cache.Incr(getKey(now, code, typeName))
	_ = cache.Expire(getKey(now, code, typeName), 3600*24)
	_, _ = cache.Sadd(getSetKey(now), typeName+"_"+code)
	_ = cache.Expire(getSetKey(now), 3600*24)
}

func getSetKey(t time.Time) string {
	return fmt.Sprintf("short:set:%d-%d-%d", t.Year(), t.Month(), t.Day())
}

func getKey(t time.Time, code string, typeName string) string {
	return fmt.Sprintf("%s_%s_%d-%d-%d", typeName, code, t.Year(), t.Month(), t.Day())
}
