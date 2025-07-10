# 学生管理系统

## 系统概述

学生管理系统是一个基于 Go 语言和 Gin 框架开发的 RESTful API 服务，提供完整的学生信息管理功能，包括创建、查询、更新和删除学生记录。系统采用 JWT 认证机制保护 API 访问安全。

## 项目结构

```
StudentService/
├── main.go          # 程序入口点，设置路由和启动服务器
├── models.go        # 定义数据模型
├── db.go            # 数据库连接和初始化
├── jwt.go           # JWT 认证相关功能
├── student_api.go   # 学生管理 API 实现
└── README.md        # 项目文档
```

## 系统架构

### API 流程图

```
┌─────────┐     ┌──────────────┐     ┌────────────┐     ┌─────────────┐
│         │     │              │     │            │     │             │
│ 客户端  │────▶│ 用户认证     │────▶│ JWT 中间件 │────▶│ 学生管理 API │
│         │     │ /login       │     │            │     │             │
└─────────┘     └──────────────┘     └────────────┘     └──────┬──────┘
     ▲                                                        │
     │                                                        │
     │                                                        ▼
     │                                                 ┌─────────────┐
     └─────────────────────────────────────────────────┤ MySQL 数据库 │
                        返回数据                        └─────────────┘
```

### 数据模型

```
┌─────────────────┐
│     Student     │
├─────────────────┤
│ id: int         │
│ name: string    │
│ age: int        │
│ email: string   │
│ grade: string   │
│ created_at: time│
│ updated_at: time│
└─────────────────┘
```

## API 文档

### 1. 认证 API

#### 登录获取 Token

- **URL**: `/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "admin",
    "password": "password"
  }
  ```
- **响应**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **状态码**:
  - `200 OK`: 登录成功
  - `401 Unauthorized`: 用户名或密码错误

### 2. 学生管理 API

> **注意**: 以下所有 API 都需要在请求头中添加 JWT Token:
> ```
> Authorization: Bearer <token>
> ```

#### 2.1 获取所有学生

- **URL**: `/api/v1/students`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "张三",
        "age": 20,
        "email": "zhangsan@example.com",
        "grade": "大一"
      },
      {
        "id": 2,
        "name": "李四",
        "age": 21,
        "email": "lisi@example.com",
        "grade": "大二"
      }
    ]
  }
  ```
- **状态码**:
  - `200 OK`: 查询成功
  - `401 Unauthorized`: 未授权
  - `500 Internal Server Error`: 服务器错误

#### 2.2 创建新学生

- **URL**: `/api/v1/students`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "name": "王五",
    "age": 19,
    "email": "wangwu@example.com",
    "grade": "大一"
  }
  ```
- **响应**:
  ```json
  {
    "data": {
      "id": 3,
      "name": "王五",
      "age": 19,
      "email": "wangwu@example.com",
      "grade": "大一"
    }
  }
  ```
- **状态码**:
  - `200 OK`: 创建成功
  - `400 Bad Request`: 请求参数错误
  - `401 Unauthorized`: 未授权
  - `500 Internal Server Error`: 服务器错误

#### 2.3 获取单个学生

- **URL**: `/api/v1/students/:id`
- **方法**: `GET`
- **参数**: 
  - `id`: 学生ID (路径参数)
- **响应**:
  ```json
  {
    "data": {
      "id": 1,
      "name": "张三",
      "age": 20,
      "email": "zhangsan@example.com",
      "grade": "大一"
    }
  }
  ```
- **状态码**:
  - `200 OK`: 查询成功
  - `400 Bad Request`: ID 参数无效
  - `401 Unauthorized`: 未授权
  - `404 Not Found`: 学生不存在
  - `500 Internal Server Error`: 服务器错误

#### 2.4 更新学生信息

- **URL**: `/api/v1/students/:id`
- **方法**: `PUT`
- **参数**: 
  - `id`: 学生ID (路径参数)
- **请求体**:
  ```json
  {
    "name": "张三",
    "age": 22,
    "email": "zhangsan@example.com",
    "grade": "大三"
  }
  ```
- **响应**:
  ```json
  {
    "message": "更新成功",
    "data": {
      "id": 1,
      "name": "张三",
      "age": 22,
      "email": "zhangsan@example.com",
      "grade": "大三"
    }
  }
  ```
- **状态码**:
  - `200 OK`: 更新成功
  - `400 Bad Request`: ID 参数或请求体无效
  - `401 Unauthorized`: 未授权
  - `404 Not Found`: 学生不存在
  - `500 Internal Server Error`: 服务器错误

#### 2.5 删除学生

- **URL**: `/api/v1/students/:id`
- **方法**: `DELETE`
- **参数**: 
  - `id`: 学生ID (路径参数)
- **响应**:
  ```json
  {
    "message": "学生删除成功"
  }
  ```
- **状态码**:
  - `200 OK`: 删除成功
  - `400 Bad Request`: ID 参数无效
  - `401 Unauthorized`: 未授权
  - `404 Not Found`: 学生不存在
  - `500 Internal Server Error`: 服务器错误

## 错误响应格式

所有错误响应都采用以下统一格式:

```json
{
  "message": "错误描述"
}
```

## 认证与安全

### JWT 认证流程

1. 客户端提交用户名和密码到 `/login` 端点
2. 服务器验证凭据，生成 JWT token 并返回
3. 客户端在后续请求中，将 token 添加到请求头中
4. 服务器验证 token 是否有效，并允许或拒绝请求

## 部署要求

- Go 1.13 或更高版本
- MySQL 5.7 或更高版本
- 至少 512MB RAM

## 数据库初始化

执行以下 SQL 语句创建必要的表结构:

```sql
CREATE TABLE students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    email VARCHAR(255) NOT NULL,
    grade VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```