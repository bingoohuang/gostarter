package ging

import (
	"net/http"

	"github.com/bingoohuang/gg/pkg/ginx"
	"github.com/bingoohuang/gg/pkg/jsoni"
	"github.com/bingoohuang/gg/pkg/jsoni/extra"
	"github.com/bingoohuang/gg/pkg/strcase"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// JsoniConfig tries to be 100% compatible with standard library behavior
var JsoniConfig = func() jsoni.API {
	c := jsoni.Config{EscapeHTML: true}.Froze()
	c.RegisterExtension(&extra.NamingStrategyExtension{Translate: strcase.ToCamelLower})
	c.RegisterTypeEncoderFunc("time.Time", jsoni.CreateTimeEncodeFn("2006-01-02 15:04:05"), nil)
	return c
}()

func JSON(g *gin.Context, data interface{}) {
	g.Render(http.StatusOK, &ginx.JSONRender{Data: data, JsoniAPI: JsoniConfig})
}
