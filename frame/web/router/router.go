package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	htError "github.com/qq754174349/ht/ht-frame/common/error"
	"github.com/qq754174349/ht/ht-frame/common/result"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web/config"
)

type Registrar struct {
	*group
	engine *gin.Engine
}

type group struct {
	group     *gin.RouterGroup
	timeout   time.Duration
	registrar *Registrar
}

func New(engine *gin.Engine) *Registrar {
	timeout := config.Get().Web.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	registrar := &Registrar{
		engine: engine,
	}

	registrar.group = &group{
		group:     &engine.RouterGroup,
		timeout:   timeout,
		registrar: registrar,
	}

	return registrar
}

func wrapHandlers(timeout time.Duration, handlers ...gin.HandlerFunc) []gin.HandlerFunc {
	resp := make([]gin.HandlerFunc, 0, len(handlers))

	for _, handler := range handlers {
		resp = append(resp, func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
			defer cancel()
			c.Request = c.Request.WithContext(ctx)

			var err error

			done := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf("panic: %v", r)
						switch x := r.(type) {
						case string:
							err = errors.New(x)
						case error:
							err = x
						default:
							err = errors.New("unknown error")
						}
					}
					close(done)
				}()
				handler(c)
			}()

			select {
			case <-done:
				if err != nil {
					result.FailDefault(c)
					return
				}
			case <-ctx.Done():
				result.FailWithHttpCode(c, http.StatusGatewayTimeout, htError.FAILURE.Code, htError.FAILURE.Msg)
				return
			}
		})
	}

	return resp
}

// Handle 注册任意HTTP方法的路由
func (g *group) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) *group {
	if g.timeout > 0 {
		g.group.Handle(httpMethod, relativePath, wrapHandlers(g.timeout*time.Second, handlers...)...)
	} else {
		g.group.Handle(httpMethod, relativePath, handlers...)
	}
	return g
}

// GET 注册GET方法路由
func (g *group) GET(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodGet, relativePath, handlers...)
}

// POST 注册POST方法路由
func (g *group) POST(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodPost, relativePath, handlers...)
}

// DELETE 注册DELETE方法路由
func (g *group) DELETE(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodDelete, relativePath, handlers...)
}

// PATCH 注册PATCH方法路由
func (g *group) PATCH(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodPatch, relativePath, handlers...)
}

// PUT 注册PUT方法路由
func (g *group) PUT(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodPut, relativePath, handlers...)
}

// OPTIONS 注册OPTIONS方法路由
func (g *group) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodOptions, relativePath, handlers...)
}

// HEAD 注册HEAD方法路由
func (g *group) HEAD(relativePath string, handlers ...gin.HandlerFunc) *group {
	return g.Handle(http.MethodHead, relativePath, handlers...)
}

// Any 注册所有HTTP方法的路由
func (g *group) Any(relativePath string, handlers ...gin.HandlerFunc) *group {
	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
	for _, method := range methods {
		g.Handle(method, relativePath, handlers...)
	}
	return g
}

// Match 注册匹配指定方法的路由
func (g *group) Match(methods []string, relativePath string, handlers ...gin.HandlerFunc) *group {
	for _, method := range methods {
		g.Handle(method, relativePath, handlers...)
	}
	return g
}

// StaticFile 注册静态文件路由
func (g *group) StaticFile(relativePath, filepath string) *group {
	g.group.StaticFile(relativePath, filepath)
	return g
}

// StaticFileFS 使用文件系统注册静态文件路由
func (g *group) StaticFileFS(relativePath, filepath string, fs http.FileSystem) *group {
	g.group.StaticFileFS(relativePath, filepath, fs)
	return g
}

// Static 注册静态文件目录路由
func (g *group) Static(relativePath, root string) *group {
	g.group.Static(relativePath, root)
	return g
}

// StaticFS 使用文件系统注册静态文件目录路由
func (g *group) StaticFS(relativePath string, fs http.FileSystem) *group {
	g.group.StaticFS(relativePath, fs)
	return g
}

// Group 创建子路由组
func (g *group) Group(relativePath string, handlers ...gin.HandlerFunc) *group {
	routerGroup := g.group.Group(relativePath, handlers...)
	return &group{
		group:     routerGroup,
		timeout:   g.timeout,
		registrar: g.registrar,
	}
}

// GroupWithTimeout 创建子路由组并设置超时时间, 0则没有时限
func (g *group) GroupWithTimeout(relativePath string, timeout time.Duration, handlers ...gin.HandlerFunc) *group {
	routerGroup := g.group.Group(relativePath, handlers...)
	return &group{
		group:     routerGroup,
		timeout:   timeout,
		registrar: g.registrar,
	}
}
