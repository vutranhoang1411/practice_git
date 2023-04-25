package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/vutranhoang1411/SimpleBank/db/sqlc"
	"github.com/vutranhoang1411/SimpleBank/util"
)

func newTestServer(store db.Store)(*Server,error){
	config:=util.Config{
		KeyString: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server,err:=NewServer(config,store)
	return server,err
}
func TestMain(m *testing.M){
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}