package view

import (
	"fmt"
	"net/http"

	"github.com/GanePrivate/go-rest-API/api/controller"
	"github.com/gin-gonic/gin"
)

type File struct {
	Name string `uri:"name" binding:"required"`
}

func StartServer() {
	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")
	files := v1.Group("/files")

	// ファイルを受け取るコード
	files.POST("/", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// フォームから保存するパスの情報を取得
		filePath := c.PostForm("filePath")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 受け取ったファイルを保存する
		n, err := controller.Upload(file, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "Uploaded successfully",
			"name":    fmt.Sprintf("%s", n),
		})
	})

	// ファイル名を取得してそのデータを返すコード
	files.GET("/:name/", func(c *gin.Context) {
		var f File
		if err := c.ShouldBindUri(&f); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		m, cn, err := controller.Download(f.Name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+f.Name)
		c.Data(http.StatusOK, m, cn)
	})

	// GETでファイル一覧を取得する
	v1.GET("/list/", func(c *gin.Context) {
		names, err := controller.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "Listing successfully",
			"files":   names,
		})
	})

	_ = router.Run(":8085")
}
