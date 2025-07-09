package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListStudents 获取所有学生信息
// 从数据库中查询所有学生记录，并返回 JSON 格式的响应
func ListStudents(c *gin.Context) {
	// 执行 SQL 查询，获取所有学生信息
	rows, err := DB.Query("SELECT id, name, age, email, grade FROM students ORDER BY id")
	if err != nil {
		// 如果查询失败，返回 HTTP 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	defer rows.Close() // 确保查询结果集被正确关闭

	var students []Student
	// 遍历查询结果，将每条记录添加到 students 列表中
	for rows.Next() {
		var student Student
		rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.Grade)
		students = append(students, student)
	}
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"data": students})
}

// CreateStudent 创建新学生
// 从请求体中解析 JSON 数据，插入到数据库中，并返回新创建的学生信息
func CreateStudent(c *gin.Context) {
	var newStudent Student
	// 解析请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&newStudent); err != nil {
		// 如果解析失败，返回 HTTP 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
		return
	}

	// 执行 SQL 插入操作，将新学生信息保存到数据库
	result, err := DB.Exec("INSERT INTO students (name, age, email, grade) VALUES (?, ?, ?, ?)",
		newStudent.Name, newStudent.Age, newStudent.Email, newStudent.Grade)
	if err != nil {
		// 如果插入失败，返回 HTTP 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建学生失败"})
		return
	}

	// 获取新插入记录的 ID，并设置到 newStudent 对象中
	id, _ := result.LastInsertId()
	newStudent.ID = int(id)
	// 返回新创建的学生信息
	c.JSON(http.StatusOK, gin.H{"data": newStudent})
}

// GetStudent 根据 ID 获取单个学生信息
// 从数据库中查询指定 ID 的学生记录，并返回 JSON 格式的响应
func GetStudent(c *gin.Context) {
	// 从 URL 参数中获取学生 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 如果 ID 无效，返回 HTTP 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的学生ID"})
		return
	}

	var student Student
	// 执行 SQL 查询，获取指定 ID 的学生信息
	err = DB.QueryRow("SELECT id, name, age, email, grade FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.Grade)
	if err != nil {
		// 如果学生记录未找到，返回 HTTP 404 错误
		c.JSON(http.StatusNotFound, gin.H{"message": "学生未找到"})
		return
	}
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"data": student})
}

// UpdateStudent 更新学生信息
// 根据请求体中的 JSON 数据和指定的 ID，更新数据库中的学生记录
func UpdateStudent(c *gin.Context) {
	// 从 URL 参数中获取学生 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 如果 ID 无效，返回 HTTP 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的学生ID"})
		return
	}

	var updatedStudent Student
	// 解析请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		// 如果解析失败，返回 HTTP 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
		return
	}

	// 执行 SQL 更新操作，将学生信息更新到数据库
	result, err := DB.Exec("UPDATE students SET name = ?, age = ?, email = ?, grade = ? WHERE id = ?",
		updatedStudent.Name, updatedStudent.Age, updatedStudent.Email, updatedStudent.Grade, id)
	if err != nil {
		// 如果更新失败，返回 HTTP 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		// 如果没有记录被更新，返回 HTTP 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	// 返回更新后的学生信息
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": updatedStudent})
}

// DeleteStudent 删除学生
// 根据指定的 ID 删除数据库中的学生记录
func DeleteStudent(c *gin.Context) {
	// 从 URL 参数中获取学生 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 如果 ID 无效，返回 HTTP 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的学生ID"})
		return
	}

	// 执行 SQL 删除操作
	result, err := DB.Exec("DELETE FROM students WHERE id = ?", id)
	// 检查是否有记录被删除
	rowAffected, err := result.RowsAffected()
	if err != nil || rowAffected == 0 {
		// 如果没有记录被删除，返回 HTTP 500 错误
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	// 返回删除成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "学生删除成功"})
}
