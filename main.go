package main

import (
    "os"
    "log"
	"strconv"
    "net/http"
    "gin_sample/crypto"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
    gorm.Model
    Text   string
    Status string
}

type User struct {
    gorm.Model
    Username string `form:"username" binding:"required" gorm:"unique;not null"`
    Password string `form:"password" binding:"required"`
}

func dbInit() {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbInit）")
    }
    db.AutoMigrate(&Todo{})
    db.AutoMigrate(&User{})
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

// ユーザー登録処理
func createUser(username string, password string) []error {
    passwordEncrypt, _ := crypto.PasswordEncrypt(password)
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbDelete)")
    }
    defer db.Close()
    // Insert処理
    if err := db.Create(&User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
        return err
    }
    return nil
}

// ユーザーを一件取得
func getUser(username string) User {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbDelete)")
    }
    var user User
    db.First(&user, "username = ?", username)
    db.Close()
    return user
}

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        port = "3000"
    }
    router := gin.Default()
    router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets/css")

	dbInit()

	//Index
    router.GET("/", func(ctx *gin.Context) {
        todos := dbGetAll()
        ctx.HTML(http.StatusOK, "index.html", gin.H{
            "todos": todos,
        })
    })

	// New
	router.GET("/new", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "new.html", gin.H{})
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

    // ユーザー登録画面
    router.GET("/signup", func(c *gin.Context) {

        c.HTML(200, "signup.html", gin.H{})
    })

    // ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form User
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            c.Abort()
        } else {
            username := c.PostForm("username")
            password := c.PostForm("password")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := createUser(username, password); err != nil {
                c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // ユーザーログイン画面
    router.GET("/login", func(c *gin.Context) {

        c.HTML(200, "login.html", gin.H{})
    })

    // ユーザーログイン
    router.POST("/login", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        dbPassword := getUser(c.PostForm("username")).Password
        log.Println(dbPassword)
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("ログインできませんでした")
            c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            c.Abort()
        } else {
            log.Println("ログインできました")
            c.Redirect(302, "/")
        }
    })

    router.Run(":" + port)
}