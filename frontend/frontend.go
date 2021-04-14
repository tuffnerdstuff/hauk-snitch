package frontend

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tuffnerdstuff/hauk-snitch/mapper"
)

type Server struct {
	sessions map[string]mapper.TopicSession
	config   Config
}

func New(config Config) *Server {
	return &Server{sessions: make(map[string]mapper.TopicSession), config: config}
}

func (t *Server) Run(newSessions <-chan mapper.TopicSession) {
	go t.trackSessions(newSessions)
	router := gin.Default()
	route := router.Group("/")

	if !t.config.IsAnonymous {
		route = router.Group("/", gin.BasicAuth(t.getAuthorizedUsers()))
	}

	route.GET("/", t.listSessions)

	router.Run(fmt.Sprintf(":%d", t.config.Port))

}

func (t *Server) getAuthorizedUsers() gin.Accounts {
	return gin.Accounts{t.config.User: t.config.Password}
}

func (t *Server) listSessions(c *gin.Context) {
	sessions := t.sessions
	for _, session := range sessions {
		fmt.Fprint(c.Writer, "<html><body>")
		fmt.Fprintf(c.Writer, "<a href=\"%s\">%s</a><br>", session.URL, session.Topic)
		fmt.Fprint(c.Writer, "</body></html>")
	}
}

func (t *Server) trackSessions(newSessions <-chan mapper.TopicSession) {
	log.Println("Tracking Sessions")
	for newSession := range newSessions {
		t.sessions[newSession.Topic] = newSession
	}
}
