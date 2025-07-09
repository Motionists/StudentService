package main

type Student struct {
	ID    int    `json:"id"`    // 学生 ID
	Name  string `json:"name"`  // 学生姓名
	Age   int    `json:"age"`   // 学生年龄
	Email string `json:"email"` // 学生邮箱
	Grade string `json:"grade"` // 学生年级
}
