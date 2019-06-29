package httpserver

func (s *Server) RouterSetting() {

	s.Engine.GET("/hc", HC)

	youbike := s.Engine.Group("/youbike")
	{
		youbike.GET("/topic-one/:text", SearchStationByText)
		youbike.GET("/topic-two", SearchSpaceRank)
		youbike.GET("/topic-three", SearchMaxAndMinFreeSpaceByTime)
	}
}
