// admin-pos project main.go
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"

	//	"github.com/jinzhu/gorm"
	//	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int64) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	dsn := ""
	if gin.Mode() == "release" { //prod
		dsn = "host=localhost user=orgchart password=##123211agc## dbname=orgapp port=5432 sslmode=disable TimeZone=America/Caracas"
		//dsn = "host=localhost user=postgres password=developer dbname=orgapp port=5432 sslmode=disable TimeZone=America/Caracas"
	} else { //local
		//dsn = "host=localhost user=postgres password=123456 dbname=orgapp port=5432 sslmode=disable TimeZone=America/Caracas"
		//dsn = "host=localhost user=postgres password=developer dbname=orgapp port=5432 sslmode=disable TimeZone=America/Lima"
		dsn = "host=localhost user=postgres password=forever dbname=orgapp port=5433 sslmode=disable TimeZone=America/Lima"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//	defer db.Close()

	// result := db.AutoMigrate(&Contact{}, &ContactRight{}, &Position{}, &Sheet{}, &Config{}, &Project{})
	// if result.Error != nil {
	// 	fmt.Print(result.Error)
	// }

	//db.LogMode(true)
	router := gin.Default()
	router.Use(CORSMiddleware())
	v1 := router.Group("/api")
	router.Static("/app", "./app")
	router.Static("/images", "./images")
	router.Static("/assets/", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
	})

	v1.GET("/position", func(c *gin.Context) {
		searchKey := c.Query("k")
		result, err := positionGetAll(db, searchKey)
		if err != "" {
			fmt.Println("Error: positionMovGetAll() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/position/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := positionGet(db, idParam)
		if err != "" {
			fmt.Println("Error: positionGet() " + err)
		}
		c.JSON(200, result)
	})
	v1.POST("/position", func(c *gin.Context) {
		var postData Position
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := positionCreate(db, postData); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})
	v1.PUT("/position/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var postData Position
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := positionUpdate(db, postData, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "update sucess"})
	})
	v1.DELETE("/position/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		if err2 := positionDelete(db, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delete success"})
	})

	v1.GET("/sheet", func(c *gin.Context) {
		searchKey := c.Query("k")
		result, err := sheetGetAll(db, searchKey)
		if err != "" {
			fmt.Println("Error: sheetMovGetAll() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/sheet/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := sheetGet(db, idParam)
		if err != "" {
			fmt.Println("Error: sheetGet() " + err)
		}
		c.JSON(200, result)
	})
	v1.POST("/sheet", func(c *gin.Context) {
		var postData Sheet
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := sheetCreate(db, postData); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})
	v1.PUT("/sheet/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		isNameUpdate := c.Query("isn")
		var postData Sheet
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err2 := sheetUpdate(db, postData, idParam, isNameUpdate); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "update sucess"})
	})
	v1.DELETE("/sheet/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		if err2 := sheetDelete(db, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delete success"})
	})

	//project
	v1.GET("/project", func(c *gin.Context) {
		searchKey := c.Query("k")
		result, err := projectGetAll(db, searchKey)
		if err != "" {
			fmt.Println("Error: sheetMovGetAll() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/project/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := projectGet(db, idParam)
		if err != "" {
			fmt.Println("Error: sheetGet() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/project/:id/sheets", func(c *gin.Context) {
		idParam := c.Param("id")
		var records []Sheet

		result := db.Where("project_id = ?", idParam).Find(&records) // find all
		if result.Error != nil {
			err := result.Error.Error() + " - sheetGetAll(), position"
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(200, records)

	})

	v1.POST("/project", func(c *gin.Context) {
		var postData Project
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := projectCreate(db, postData); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})
	v1.PUT("/project/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		isNameUpdate := c.Query("isn")

		var postData Project
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := projectUpdate(db, postData, idParam, isNameUpdate); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "update sucess"})
	})
	v1.DELETE("/project/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		if err2 := projectDelete(db, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delete success"})
	})

	//contactRight
	v1.GET("/contact-right", func(c *gin.Context) {
		contactID, _ := strconv.Atoi(c.Query("cid"))
		itemtype := c.Query("itype")
		optype := c.Query("otype")

		result, err := contactRightGetAll(db, contactID, itemtype, optype)
		if err != "" {
			fmt.Println("Error: contactRightGetAll() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/contact-right/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := contactRightGet(db, idParam)
		if err != "" {
			fmt.Println("Error: contactRightGet() " + err)
		}
		c.JSON(200, result)
	})

	v1.POST("/contact-right", func(c *gin.Context) {
		var postData ContactRight
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := contactRightCreate(db, postData); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})
	v1.PUT("/contact-right/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var postData ContactRight
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := contactRightUpdate(db, postData, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "update sucess"})
	})
	v1.DELETE("/contact-right/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		if err2 := contactRightDelete(db, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delete success"})
	})
	v1.GET("/contact", func(c *gin.Context) {

		result, err := contactGetAll(db)
		if err != "" {
			fmt.Println("Error: contactRightGetAll() " + err)
		}
		c.JSON(200, result)
	})
	v1.GET("/contact/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := contactRightGet(db, idParam)
		if err != "" {
			fmt.Println("Error: contactRightGet() " + err)
		}
		c.JSON(200, result)
	})

	v1.POST("/contact", func(c *gin.Context) {
		var postData Contact
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := contactCreate(db, postData); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})
	v1.PUT("/contact/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		var postData Contact
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err2 := contactUpdate(db, postData, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "update sucess"})
	})
	v1.DELETE("/contact/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		if err2 := contactDelete(db, idParam); err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delete success"})
	})
	v1.GET("/config/:key", func(c *gin.Context) {
		idParam := c.Param("key")

		var data = []Config{}
		var record Config
		//		var err string
		result := db.Where("key = ?", idParam).First(&record) // find all

		if result.Error != nil {
			//			err := result.Error.Error()
		}
		if record.ID != 0 {
			data = append(data, record)
		}

		c.JSON(200, data)
	})
	v1.POST("/config/:key", func(c *gin.Context) {
		idParam := c.Param("key")
		var postData Config
		if err := c.ShouldBindJSON(&postData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//		var data = []Config{}
		var record Config
		//		var err string
		result := db.Where("key = ?", idParam).First(&record) // find all
		if result.Error != nil {
			//			err := result.Error.Error()
		}
		if record.ID != 0 {
			record.Value = postData.Value
			db.Save(record)
		} else {
			result := db.Create(&Config{
				Key:   postData.Key,
				Value: postData.Value,
			})
			if result.Error != nil {
				fmt.Println("Error: ", result.Error.Error())
			}

		}

		c.JSON(http.StatusOK, gin.H{"status": "create sucess"})
	})

	//misc
	v1.GET("/export-tree", func(c *gin.Context) {
		str := c.Query("d")
		fileName := c.Query("n")
		if fileName == "" {
			fileName = "treeexport"
		}

		values := strings.Split(str, ";")
		lettersToLevels := [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		f := excelize.NewFile()
		style, err := f.NewStyle(`{"fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1}}`)
		if err != nil {
			fmt.Println(err)
		}

		i := 1
		for _, item := range values {
			vals := strings.Split(item, ":")
			if vals[0] != "" {
				letter, _ := strconv.Atoi(vals[0])
				f.SetCellValue("Sheet1", lettersToLevels[letter]+strconv.Itoa(i), vals[1])
				f.SetCellStyle("Sheet1", lettersToLevels[letter]+strconv.Itoa(i), lettersToLevels[letter]+strconv.Itoa(i), style)
				i += 1
			}
		}
		if err := f.SaveAs(filepath.Join("files", fileName+".xlsx")); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{"error": err})

		}
		c.File(filepath.Join("files", fileName+".xlsx"))
		fmt.Println(fileName + ".xlsx")

		c.JSON(http.StatusOK, gin.H{"ok": "ok"})

	})

	v1.POST("/upload", func(c *gin.Context) {
		// image upload
		//filename := RandStringBytesMaskImprSrc(16)
		file, header, err := c.Request.FormFile("upload")
		filename := header.Filename
		if err != nil {
			fmt.Println(err)
		} else {
			out, err := os.Create("./images/" + filename)

			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				fmt.Println(err)
			}
		}

		image := "/images/noimage.png"
		if file != nil {
			image = "/images/" + filename
		}
		c.JSON(http.StatusOK, gin.H{"image": image})

	})

	v1.POST("/upload-file", func(c *gin.Context) {
		// image upload
		//filename := RandStringBytesMaskImprSrc(16)
		file, header, err := c.Request.FormFile("upload")
		fileid := RandStringBytesMaskImprSrc(3)
		filename := fileid + "-" + header.Filename
		if err != nil {
			fmt.Println(err)
		} else {
			out, err := os.Create("./files/" + filename)

			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				fmt.Println(err)
			}
		}

		newfile := "/files/nofile"
		if file != nil {
			newfile = "/files/" + filename
		}
		c.JSON(http.StatusOK, gin.H{"file": newfile})

	})

	s := &http.Server{
		Addr:    ":8087",
		Handler: router,
		//ReadTimeout:    10 * 60, //
		//WriteTimeout:   10 * 60, //
		//MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
