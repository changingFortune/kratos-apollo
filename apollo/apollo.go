package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/go-kratos/kratos/v2/config"
	"os"
	"path/filepath"
)

type apollo struct {
	client *agollo.Client
}

type Option func(*options)

type options struct {
	appid          string
	secret         string
	cluster        string
	endpoint       string
	namespace      string
	isBackupConfig bool
	backPath		string
}

// WithAppID with apollo config app id
func WithAppID(appID string) Option {
	return func(o *options) {
		o.appid = appID
	}
}

// WithCluster with apollo config cluster
func WithCluster(cluster string) Option {
	return func(o *options) {
		o.cluster = cluster
	}
}

// WithEndpoint with apollo config conf server ip
func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}

// WithEnableBackup with apollo config enable backup config
func WithEnableBackup() Option {
	return func(o *options) {
		o.isBackupConfig = true
	}
}

// WithDisableBackup with apollo config enable backup config
func WithDisableBackup() Option {
	return func(o *options) {
		o.isBackupConfig = false
	}
}

// WithSecret with apollo config app secret
func WithSecret(secret string) Option {
	return func(o *options) {
		o.secret = secret
	}
}

// WithBackPath with apollo config cach backpath
func WithBackPath(backPath string) Option {
	return func(o *options) {
		o.backPath = backPath
	}
}


func NewSource(opts ...Option) config.Source {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	if op.backPath == "" {
		op.backPath = "cache"
	}
	
	client, err := agollo.StartWithConfig(func() (*apolloConfig.AppConfig, error) {
		return &apolloConfig.AppConfig{
			AppID:          op.appid,
			Cluster:        op.cluster,
			NamespaceName:  op.namespace,
			IP:             op.endpoint,
			IsBackupConfig: op.isBackupConfig,
			Secret:         op.secret,
			BackupConfigPath: op.backPath,
		}, nil
	})
	if err != nil {
		panic(err)
	}
	
	// check have dir path
	absPath,_ :=filepath.Abs("")
	thePath := filepath.Join(absPath,op.backPath)
	if _,err:= os.Stat(thePath);os.IsNotExist(err){
		os.Mkdir(thePath,0777)
	}
	return &apollo{client}
}

func (e *apollo) load() []*config.KeyValue {
	kv := make([]*config.KeyValue, 0)
	e.client.GetDefaultConfigCache().Range(func(key, value interface{}) bool {
		kv = append(kv, &config.KeyValue{
			Key:   key.(string),
			Value: []byte(value.(string)),
		})
		return true
	})
	return kv
}

func (e *apollo) Load() (kv []*config.KeyValue, err error) {
	return e.load(), nil
}

func (e *apollo) Watch() (config.Watcher, error) {
	w, err := NewWatcher(e)
	if err != nil {
		return nil, err
	}
	return w, nil
}
