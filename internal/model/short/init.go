package short

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yi-nology/common/biz/lock"
	"github.com/yi-nology/sdk/conf"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type ShortDaoInterface interface {
	UpdateByCodeCount(code string, count int) (err error)
	UpdateByCodeStatus(code string, status int64) (err error)
	DeleteByCode(code string) (err error)
	GetList(page, pageSize int, shortCode string, createBy string) (urls []Short, count int64, err error)
	UpdateByCode(code string, url Short) (err error)
	Create(url Short) (err error)
	GetShortUrl(code string) (url *Short, err error)
	ExpireList(page, pageSize int) (urls []Short, err error)
	GetExpireCount() (count int64, err error)
}

type Short struct {
	Id       int64          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Code     string         `gorm:"column:code;NOT NULL;COMMENT:'code 短码'" json:"code"`
	Data     string         `gorm:"column:data;NOT NULL;COMMENT:'数据'" json:"data"`
	Exprie   *time.Time     `gorm:"column:exprie;DEFAULT NULL;COMMENT:'过期时间'" json:"exprie"`
	Status   int            `gorm:"column:status;NOT NULL" json:"status"`
	Count    int            `gorm:"column:count;NOT NULL" json:"count"`
	CreateBy string         `gorm:"column:create_by;DEFAULT NULL" json:"createBy"`
	CreateAt *time.Time     `gorm:"column:create_at;autoCreateTime:nano;" json:"createAt"`
	UpdateAt *time.Time     `gorm:"column:update_at;DEFAULT NULL ON UPDATE current_timestamp()" json:"UpdateAt"`
	DeleteAt gorm.DeletedAt `gorm:"column:delete_at;DEFAULT NULL" json:"-"`
}

type shortDao struct {
	ctx   context.Context
	db    *gorm.DB
	model any
	cache *redis.Client
	bid   string
	log   logx.Logger
}

func newShortDao(ctx context.Context, db *gorm.DB, model any, cache *redis.Client, bid string, log logx.Logger) *shortDao {
	return &shortDao{ctx: ctx, db: db, model: model, cache: cache, bid: bid, log: log}
}

// GetShortUrl 获取短链接
func (d *shortDao) GetShortUrl(code string) (url *Short, err error) {
	//u := &Short{}
	lockCode := lock.NewLock(*conf.RedisClient, d.log)
	lockCode.Up(context.Background(), d.key(d.bid, code), 5)
	defer lockCode.Down(d.ctx, d.key(d.bid, code))

	data, err := d.cache.Get(context.Background(), d.bid+"_"+code).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil || len(data) == 0 {
		err = d.db.Model(d.model).Table(d.bid).Where("code = ?", code).First(&url).Error
		if err != nil {
			return nil, err
		}
		data, err = jsonx.MarshalToString(url)
		if err != nil {
			return nil, err
		}
		err = d.cache.SetEx(d.ctx, d.bid+"_"+code, data, 360*time.Second).Err()
		if err != nil {
			return nil, err
		}
	}
	err = jsonx.UnmarshalFromString(data, &url)
	return url, nil
}

// Create 创建数据
func (d *shortDao) Create(url Short) (err error) {
	lockCode := lock.NewLock(*conf.RedisClient, d.log)
	lockCode.Up(context.Background(), d.key(d.bid, url.Code), 5)
	defer lockCode.Down(context.Background(), d.key(d.bid, url.Code))

	err = d.db.Model(d.model).Table(d.bid).Create(&url).Error
	if err != nil {
		return err
	}
	return
}

// UpdateByCode 更新数据
func (d *shortDao) UpdateByCode(code string, data Short) (err error) {
	lockCode := lock.NewLock(*conf.RedisClient, d.log)
	lockCode.Up(context.Background(), d.key(d.bid, code), 5)
	defer lockCode.Down(context.Background(), d.key(d.bid, code))

	err = d.db.Model(d.model).Table(d.bid).Where("code = ?", code).Updates(data).Error
	if err != nil {
		return err
	}
	_, _ = d.cache.Del(d.ctx, d.bid+"_"+code).Result()
	return
}

// GetList 获取列表
func (d *shortDao) GetList(page, pageSize int, shortCode string, createBy string) (urls []Short, count int64, err error) {
	tx := d.db.Model(d.model).Table(d.bid)
	err = tx.Where("code like ? and create_by = ?", "%"+shortCode+"%", createBy).Offset((page - 1) * pageSize).Limit(pageSize).Order("create_at DESC").Find(&urls).Error
	if err != nil {
		return nil, 0, err
	}
	err = tx.Count(&count).Error
	return
}

func (d *shortDao) GetExpireCount() (count int64, err error) {
	err = d.db.Model(d.model).Table(d.bid).Where("status = 1 and exprie < current_timestamp()").Count(&count).Error
	return
}

func (d *shortDao) ExpireList(page, pageSize int) (urls []Short, err error) {
	tx := d.db.Model(d.model).Table(d.bid).Where("status = 1 and exprie < current_timestamp()")

	err = tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&urls).Error
	if err != nil {
		return nil, err
	}

	return urls, err
}

// UpdateByCodeCount 更新count数据
func (d *shortDao) UpdateByCodeCount(code string, count int) (err error) {

	err = d.db.Model(d.model).Table(d.bid).Where("code = ?", code).Select("count").UpdateColumn("count", gorm.Expr("count + ?", count)).Error
	_, _ = d.cache.Del(d.ctx, d.bid+"_"+code).Result()
	return
}

// UpdateByCodeStatus 更新状态
func (d *shortDao) UpdateByCodeStatus(code string, status int64) (err error) {
	err = d.db.Model(d.model).Table(d.bid).Where("code = ?", code).Update("status", status).Error
	if err != nil {
		return err
	}
	_, _ = d.cache.Del(d.ctx, d.bid+"_"+code).Result()
	return
}

// DeleteByCode 逻辑删除
func (d *shortDao) DeleteByCode(code string) (err error) {
	err = d.db.Model(d.model).Table(d.bid).Where("code = ?", code).Update("delete_at", time.Now()).Error
	if err != nil {
		return err
	}
	_, _ = d.cache.Del(d.ctx, d.bid+"_"+code).Result()
	return
}

func (d *shortDao) key(bid, code string) string {
	return fmt.Sprintf("short:%s:%s", bid, code)
}
