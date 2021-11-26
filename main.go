package main

import (
	"strconv"
    "net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
    gorm.Model
    Text   string
    Status string
}

func dbInit() {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbInit）")
    }
    db.AutoMigrate(&Todo{})
    defer db.Close()
}

//DB追加
func dbInsert(text string, status string) {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbInsert)")
    }
    db.Create(&Todo{Text: text, Status: status})
    defer db.Close()
}

//DB全取得
func dbGetAll() []Todo {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！(dbGetAll())")
    }
    var todos []Todo
    db.Order("created_at desc").Find(&todos)
    db.Close()
    return todos
}

//DB一つ取得
func dbGetOne(id int) Todo {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！(dbGetOne())")
    }
    var todo Todo
    db.First(&todo, id)
    db.Close()
    return todo
}

//DB更新
func dbUpdate(id int, text string, status string) {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbUpdate)")
    }
    var todo Todo
    db.First(&todo, id)
    todo.Text = text
    todo.Status = status
    db.Save(&todo)
    db.Close()
}

//DB削除
func dbDelete(id int) {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbDelete)")
    }
    var todo Todo
    db.First(&todo, id)
    db.Delete(&todo)
    db.Close()
}

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets/css")

	dbInit()

	//Index
    router.GET("/", func(ctx *gin.Context) {
        todos := dbGetAll()
        ctx.HTML(200, "index.html", gin.H{
            "todos": todos,
        })
    })

	// New
	router.GET("/new", func(ctx *gin.Context) {
		ctx.HTML(200, "new.html", gin.H{})
	})

	//Create
    router.POST("/post", func(ctx *gin.Context) {
        text := ctx.PostForm("text")
        status := ctx.PostForm("status")
        dbInsert(text, status)
        ctx.Redirect(302, "/")
    })

	// Show
	router.GET("/show/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		ctx.HTML(200, "show.html", gin.H{"todo": todo})
	})

	//Edit
    router.GET("/edit/:id", func(ctx *gin.Context) {
        n := ctx.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic(err)
        }
        todo := dbGetOne(id)
        ctx.HTML(200, "edit.html", gin.H{"todo": todo})
    })

	//Update
    router.POST("/update/:id", func(ctx *gin.Context) {
        n := ctx.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        text := ctx.PostForm("text")
        status := ctx.PostForm("status")
        dbUpdate(id, text, status)
        ctx.Redirect(302, "/")
    })

	//削除確認
    router.GET("/delete_check/:id", func(ctx *gin.Context) {
        n := ctx.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        todo := dbGetOne(id)
        ctx.HTML(200, "delete.html", gin.H{"todo": todo})
    })

	//Delete
    router.POST("/delete/:id", func(ctx *gin.Context) {
        n := ctx.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        dbDelete(id)
        ctx.Redirect(302, "/")

    })

    // router.Run()
    http.ListenAndServe(":8080", router)
}