package server

import "api/internals/modules/health"

func (s *Server)mountModules(){
	health.New(s.db).Register(s.router)
}