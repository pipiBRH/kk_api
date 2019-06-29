package httpserver

func (s *Server) RouterSetting() {
	youbike := s.Engine.Group("/youbike")
	{
		youbike.GET("/topic-one/:text", SearchStationByText)
		youbike.GET("topic-two", SearchSpaceRank)
		youbike.GET("topic-three", SearchMaxAndMinFreeSpaceByTime)
	}
}
