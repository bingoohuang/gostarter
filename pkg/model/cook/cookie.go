package cook

import (
	"fmt"
	"log"

	"github.com/bingoohuang/gg/pkg/jsoni"
	"github.com/bingoohuang/gostarter/pkg/model"
	"github.com/bingoohuang/gostarter/pkg/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const UserInfo = "userInfo"

func Write(ctx *gin.Context, l model.Login) {
	session := sessions.Default(ctx)
	session.Set(UserInfo, string(util.PickA(jsoni.Marshal(l))))
	if err := session.Save(); err != nil {
		log.Printf("E! save session failed: %v", err)
	}
}

func Clear(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("E! save session failed: %v", err)
	}
}

func GetLogin(ctx *gin.Context) (*model.Login, error) {
	session := sessions.Default(ctx)
	ck := session.Get(UserInfo)
	if ck == nil {
		return nil, fmt.Errorf("cookie %s not found", UserInfo)
	}

	user := &model.Login{}
	if err := jsoni.Unmarshal([]byte(ck.(string)), user); err != nil {
		return nil, fmt.Errorf("parse cooke %s failed: %v", UserInfo, err)
	}

	return user, nil
}
