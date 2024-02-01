package utils

import (
	"errors"
	"sync"
	"time"

	urlcache "github.com/lflxp/djangolang/model"

	cache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ClusterCachePrefix          = "cluster-cache-"
	TokenCachePrefix            = "token-cache-"
	UserCachePrefix             = "user-cache-"
	AlertIndicatorCache         = "alert-indicator"
	PromAlertsCache             = "prom-alerts"
	AlertSilenceCache           = "alert-silence"
	AlertTenantClusterCache     = "alert-tenant-clusters"
	AlertClustersCache          = "alert-clusters"
	AlertMessagesCache          = "alert-message"
	AlertClusterCache           = "alert-cluster"
	AlertNodeCache              = "alert-node"
	AlertAppCache               = "alert-app"
	LicAllClusterVcpusCache     = "lic-all-cluster-vcpus"
	LicProductCache             = "lic-product-cache"
	ClusterPrometheusCacheKey   = "ClusterPrometheus__"
	ClusterAlertmanagerCacheKey = "ClusterAlertmanager__"
	ClusterSkywalkingCacheKey   = "ClusterSkywalking__"
)

var (
	// 统计缓存个数,分别是cache和redis
	cacheCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Cache_Count",
			Help: "Count how many value in each cache",
		},
		[]string{
			"source",
		},
	)

	// 有过期时间的cache
	localCacheWithTTL *cache.Cache
	once              sync.Once
	onceUrl           sync.Once
)

func init() {
	prometheus.MustRegister(cacheCount)
}

func NewCacheCliWithTTL() *cache.Cache {
	once.Do(func() {
		localCacheWithTTL = cache.New(10*time.Minute, 10*time.Minute)
	})
	return localCacheWithTTL
}

// 设置URL缓存
func CacheUrlSet(value urlcache.UrlCache) error {
	onceUrl.Do(func() {
		_, isExist := NewCacheCliWithTTL().Get("UrlCache")
		if !isExist {
			NewCacheCliWithTTL().Set("UrlCache", map[string]urlcache.UrlCache{}, cache.NoExpiration)
		}
	})

	err := value.Vaild()
	if err != nil {
		// 校验不通过不允许启动服务
		panic(err)
	}

	data, isExist := NewCacheCliWithTTL().Get("UrlCache")
	if !isExist {
		return errors.New("UrlCache 不存在")
	}

	data.(map[string]urlcache.UrlCache)[value.Name] = value
	NewCacheCliWithTTL().Set("UrlCache", data, cache.NoExpiration)
	return nil
}

func CacheUrlSetList(values []urlcache.UrlCache) error {
	onceUrl.Do(func() {
		_, isExist := NewCacheCliWithTTL().Get("UrlCache")
		if !isExist {
			NewCacheCliWithTTL().Set("UrlCache", map[string]urlcache.UrlCache{}, cache.NoExpiration)
		}
	})

	for _, value := range values {
		err := value.Vaild()
		if err != nil {
			// 校验不通过不允许启动服务
			panic(err)
		}

		data, isExist := NewCacheCliWithTTL().Get("UrlCache")
		if !isExist {
			return errors.New("UrlCache 不存在")
		}

		data.(map[string]urlcache.UrlCache)[value.Name] = value
		NewCacheCliWithTTL().Set("UrlCache", data, cache.NoExpiration)
	}
	return nil
}

// 获取URL缓存
func CacheUrlGet(key string) (urlcache.UrlCache, bool) {
	data, isExist := NewCacheCliWithTTL().Get("UrlCache")
	if !isExist {
		return urlcache.UrlCache{}, isExist
	}

	if v, ok := data.(map[string]urlcache.UrlCache)[key]; ok {
		return v, true
	} else {
		return urlcache.UrlCache{}, false
	}
}

// 获取URL全量数据
func CacheUrlAll() map[string]urlcache.UrlCache {
	data, _ := NewCacheCliWithTTL().Get("UrlCache")
	return data.(map[string]urlcache.UrlCache)
}

type ValueWrapper struct {
	Value     interface{}
	CacheTime time.Time
	Duration  time.Duration
}
