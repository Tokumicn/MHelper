package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()

	vp.SetConfigName("config")
	// 设置配置文件路径
	for _, conf := range configs {
		if conf != "" {
			vp.AddConfigPath(conf)
		}
	}
	vp.SetConfigType("yaml") // 设置配置文件格式
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}
	s.WatchSettingChange() // 开启配置文件变化监听
	return s, nil
}

// WatchSettingChange 监听配置文件变化并重新加载配置文件
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			s.ReloadAllSection()
		})
	}()
}

// ReloadAllSection 重新加载配置文件细则
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadSection 读取配置细则项
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	// 不存在的才更新
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}
