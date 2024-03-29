package infrastructure

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	yaml "gopkg.in/yaml.v2"
)

type Routing struct {
	Gin          *gin.Engine
	FaceList     FaceList
	AbsolutePath string
}

type FaceList struct {
	Faces []string `yaml:"faces"`
	Zura  []string `yaml:"zura"`
}

func NewRouting() *Routing {
	c, _ := NewConfig()
	//facelistをyamlとして読み込んでRoutingへ格納
	f, _ := os.ReadFile("./dist/assets/facelist.yaml") //TODO:pathを合わせる
	var f2 FaceList
	err := yaml.UnmarshalStrict(f, &f2)
	if err != nil {
		panic(err)
	}

	r := &Routing{
		Gin:          gin.Default(),
		FaceList:     f2,
		AbsolutePath: c.AbsolutePath,
	}
	r.loadTemplates()
	r.setRouting()
	return r
}

func (r *Routing) loadTemplates() {
	r.Gin.Use(favicon.New("./dist/assets/favicon.ico"))
	r.Gin.Static("/assets", r.AbsolutePath+"/dist/assets")
	r.Gin.LoadHTMLGlob(r.AbsolutePath + "/app/interfaces/presenters/*")
}

func (r *Routing) setRouting() {
	var zura = r.getZura()
	const DEPLOY = "https://zura-chan-zura.com"

	r.Gin.GET("/", func(c *gin.Context) {
		face := r.getFace()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": zura + "💓",
			"text":  zura,
			"face":  face,
			"href": "https://twitter.com/intent/tweet" +
				"?url=" + "\n\n" + DEPLOY +
				"&text=" + zura + face,
		})
	})

	r.Gin.HEAD("/", func(c *gin.Context) {
		face := r.getFace()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": zura + "💓",
			"text":  zura,
			"face":  face,
			"href": "https://twitter.com/share" +
				"?url=" + "\n\n" + DEPLOY +
				"&text=" + zura + face,
		})
	})

	r.Gin.POST("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})

	r.Gin.PUT("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})

	r.Gin.DELETE("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})

	r.Gin.Handle(http.MethodConnect, "/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})

	r.Gin.OPTIONS("/*any", func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"text": "ずらちゃん is member of Aqours.",
		})
	})

	r.Gin.PATCH("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})

	r.Gin.Handle(http.MethodTrace, "/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})
}

func (r *Routing) Run() error {
	port := "8703"
	return r.Gin.Run(":" + port)
}

func (r *Routing) getFace() string {
	rand.Seed(time.Now().UnixNano())
	return r.FaceList.Faces[rand.Intn(len(r.FaceList.Faces))]
}

func (r *Routing) getZura() string {
	rand.Seed(time.Now().UnixNano())
	return r.FaceList.Zura[rand.Intn(len(r.FaceList.Zura))]
}
