package main

import (
	"database/sql" // 导入 Go 的标准库，用于操作数据库
	"log"          // 导入日志库，用于记录日志信息
	"net/http"     // 导入 HTTP 库，用于处理 HTTP 请求和响应
	"strconv"      // 导入字符串转换库，用于处理字符串和数字之间的转换

	"github.com/gin-gonic/gin"         // 导入 Gin 框架，用于构建 Web 服务
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动，支持与 MySQL 数据库交互
)

// Student 结构体定义学生信息，用于表示数据库中的学生记录
type Student struct {
	ID    int    `json:"id"`    // 学生 ID，唯一标识
	Name  string `json:"name"`  // 学生姓名
	Age   int    `json:"age"`   // 学生年龄
	Email string `json:"email"` // 学生邮箱
	Grade string `json:"grade"` // 学生年级
}

// 定义一个全局变量 db，用于保存数据库连接
var db *sql.DB

// 初始化数据库连接
func initDB() {
	var err error // 定义一个变量用于保存错误信息

	// 数据库连接字符串
	dsn := "root:02020202@tcp(192.168.178.1:3306)/StudentService?charset=utf8mb4&parseTime=True&loc=Local"

	// 打开数据库连接，返回一个数据库对象
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		// 如果连接失败，记录错误日志并退出程序
		log.Fatal("数据库连接失败:", err)
	}

	// 测试数据库连接是否正常
	if err = db.Ping(); err != nil {
		// 如果测试失败，记录错误日志并退出程序
		log.Fatal("数据库连接测试失败:", err)
	}

	// 如果连接成功，记录日志信息
	log.Println("数据库连接成功")
}

// ListStudents 获取所有学生信息
func ListStudents(c *gin.Context) {
	// 执行 SQL 查询，获取所有学生信息
	rows, err := db.Query("SELECT id, name, age, email, grade FROM students ORDER BY id")
	if err != nil {
		// 如果查询失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close() // 确保查询结果关闭，释放资源

	// 定义一个切片用于保存学生信息
	var students []Student
	for rows.Next() {
		var student Student // 定义一个变量用于保存单个学生信息
		// 将查询结果赋值给 student 变量
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.Grade)
		if err != nil {
			// 如果解析数据失败，返回 HTTP 500 错误和错误信息
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "数据解析失败",
				"error":   err.Error(),
			})
			return
		}
		// 将学生信息添加到切片中
		students = append(students, student)
	}

	// 返回 HTTP 200 响应和学生信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    students,
	})
}

// CreateStudent 创建新学生
func CreateStudent(c *gin.Context) {
	var newStudent Student // 定义一个变量用于保存新学生信息

	// 解析请求体中的 JSON 数据，并赋值给 newStudent
	if err := c.ShouldBindJSON(&newStudent); err != nil {
		// 如果解析失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 执行 SQL 插入操作，将新学生信息保存到数据库
	result, err := db.Exec("INSERT INTO students (name, age, email, grade) VALUES (?, ?, ?, ?)",
		newStudent.Name, newStudent.Age, newStudent.Email, newStudent.Grade)
	if err != nil {
		// 如果插入失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建学生失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取新插入记录的 ID
	id, err := result.LastInsertId()
	if err != nil {
		// 如果获取 ID 失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取ID失败",
			"error":   err.Error(),
		})
		return
	}

	// 将新学生的 ID 设置为返回值
	newStudent.ID = int(id)

	// 返回 HTTP 201 响应和新学生信息
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "学生创建成功",
		"data":    newStudent,
	})
}

// GetStudent 根据ID获取单个学生信息
func GetStudent(c *gin.Context) {
	idParam := c.Param("id")         // 获取 URL 参数中的学生 ID
	id, err := strconv.Atoi(idParam) // 将学生 ID 转换为整数
	if err != nil {
		// 如果转换失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的学生ID",
		})
		return
	}

	var student Student // 定义一个变量用于保存学生信息
	// 执行 SQL 查询，根据学生 ID 获取学生信息
	err = db.QueryRow("SELECT id, name, age, email, grade FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.Grade)
	if err == sql.ErrNoRows {
		// 如果学生不存在，返回 HTTP 404 错误和错误信息
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生未找到",
		})
		return
	} else if err != nil {
		// 如果查询失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回 HTTP 200 响应和学生信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    student,
	})
}

// UpdateStudent 更新学生信息
func UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")         // 获取 URL 参数中的学生 ID
	id, err := strconv.Atoi(idParam) // 将学生 ID 转换为整数
	if err != nil {
		// 如果转换失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的学生ID",
		})
		return
	}

	var updatedStudent Student // 定义一个变量用于保存更新后的学生信息
	// 解析请求体中的 JSON 数据，并赋值给 updatedStudent
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		// 如果解析失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 执行 SQL 更新操作，根据学生 ID 更新学生信息
	result, err := db.Exec("UPDATE students SET name = ?, age = ?, email = ?, grade = ? WHERE id = ?",
		updatedStudent.Name, updatedStudent.Age, updatedStudent.Email, updatedStudent.Grade, id)
	if err != nil {
		// 如果更新失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 如果检查失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查更新结果失败",
			"error":   err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		// 如果没有行被更新，返回 HTTP 404 错误和错误信息
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生未找到",
		})
		return
	}

	// 将更新后的学生 ID 设置为返回值
	updatedStudent.ID = id

	// 返回 HTTP 200 响应和更新后的学生信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "学生信息更新成功",
		"data":    updatedStudent,
	})
}

// DeleteStudent 删除学生
func DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")         // 获取 URL 参数中的学生 ID
	id, err := strconv.Atoi(idParam) // 将学生 ID 转换为整数
	if err != nil {
		// 如果转换失败，返回 HTTP 400 错误和错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的学生ID",
		})
		return
	}

	// 执行 SQL 删除操作，根据学生 ID 删除学生信息
	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		// 如果删除失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查是否有行被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 如果检查失败，返回 HTTP 500 错误和错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查删除结果失败",
			"error":   err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		// 如果没有行被删除，返回 HTTP 404 错误和错误信息
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学生未找到",
		})
		return
	}

	// 返回 HTTP 200 响应，表示学生删除成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "学生删除成功",
	})
}

// 主函数，程序入口
func main() {
	// 初始化数据库连接，确保服务启动时能够连接到 MySQL 数据库
	initDB()
	// 使用 defer 关键字确保程序退出时关闭数据库连接，避免资源泄漏
	defer db.Close()

	// 创建一个默认的 Gin 路由器，包含日志记录和错误恢复中间件
	router := gin.Default()

	// 定义一个 API 路由组，所有路由都以 "/api/v1" 开头，方便管理和扩展
	api := router.Group("/api/v1")
	{
		// 注册 GET 路由，用于获取所有学生信息，调用 ListStudents 函数处理请求
		api.GET("/students", ListStudents)
		// 注册 POST 路由，用于创建新学生，调用 CreateStudent 函数处理请求
		api.POST("/students", CreateStudent)
		// 注册 GET 路由，用于根据学生 ID 获取单个学生信息，调用 GetStudent 函数处理请求
		api.GET("/students/:id", GetStudent)
		// 注册 PUT 路由，用于根据学生 ID 更新学生信息，调用 UpdateStudent 函数处理请求
		api.PUT("/students/:id", UpdateStudent)
		// 注册 DELETE 路由，用于根据学生 ID 删除学生信息，调用 DeleteStudent 函数处理请求
		api.DELETE("/students/:id", DeleteStudent)
	}

	// 输出日志，提示服务器已经启动并监听 8080 端口
	log.Println("服务器启动在 :8080")
	// 启动 HTTP 服务器，监听 8080 端口，等待客户端请求
	router.Run(":8080")
}
