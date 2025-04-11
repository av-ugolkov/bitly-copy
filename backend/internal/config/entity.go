package config

import (
	"fmt"
	"time"
)

type Server struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

func (s *Server) Addr() string {
	return fmt.Sprintf(":%s", s.Port)
}

type Redis struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Name        string        `yaml:"name"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

func (r *Redis) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
