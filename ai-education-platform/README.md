# AI教育平台

基于Go-zero框架构建的智能化教育学习平台，服务于学生、教师、家长和教培机构。

## 项目特性

- **多角色用户系统**: 支持学生、教师、家长、教培机构四种角色
- **AI学习助手**: 智能知识点讲解、互动问答、个性化学习路径推荐
- **题库与练习系统**: 多学科题库、智能出题、错题本管理
- **课程管理**: 完整的课程创建、学习进度追踪体系
- **微服务架构**: 基于Go-zero的高性能微服务架构

## 技术栈

- **后端框架**: Go-zero
- **数据库**: MySQL 8.0+
- **缓存**: Redis
- **认证**: JWT
- **API文档**: Go-zero API定义
- **部署**: Docker (支持)

## 服务架构

```
ai-education-platform/
├── services/
│   ├── user/          # 用户服务
│   ├── course/        # 课程服务
│   ├── ai_assistant/  # AI助手服务
│   └── exam/          # 考试服务
├── common/            # 公共模块
├── docs/              # 文档
├── scripts/           # 脚本
└── sql/              # 数据库脚本
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- goctl工具 (可选)

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd ai-education-platform
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置数据库**
```bash
# 创建数据库
mysql -u root -p < sql/schema.sql
```

4. **配置服务**
```bash
# 复制配置文件
cp services/user/user/etc/user-api.yaml.example services/user/user/etc/user-api.yaml
# 编辑配置文件，修改数据库连接等信息
```

5. **启动服务**
```bash
# 启动用户服务
cd services/user/user
go run user.go -f etc/user-api.yaml

# 启动其他服务
cd services/course/course
go run course.go -f etc/course-api.yaml
```

### 服务端口

- 用户服务: 8881
- 课程服务: 8882
- AI助手服务: 8883
- 考试服务: 8884

## API文档

### 用户服务 API

#### 用户注册
```http
POST /api/v1/user/register
Content-Type: application/json

{
    "username": "student123",
    "password": "password123",
    "email": "student@example.com",
    "nickname": "小明",
    "role": "student",
    "grade": "三年级"
}
```

#### 用户登录
```http
POST /api/v1/user/login
Content-Type: application/json

{
    "username": "student123",
    "password": "password123"
}
```

#### 获取用户信息
```http
GET /api/v1/user/info
Authorization: Bearer <token>
```

### 课程服务 API

#### 获取课程列表
```http
GET /api/v1/courses?page=1&page_size=10&subject=math&grade=三年级
```

#### 课程报名
```http
POST /api/v1/course/123/enroll
Authorization: Bearer <token>
```

### AI助手服务 API

#### AI知识点讲解
```http
POST /api/v1/ai/explain
Authorization: Bearer <token>
Content-Type: application/json

{
    "subject": "math",
    "grade": "三年级",
    "topic": "分数的基本概念",
    "difficulty": 2,
    "language": "zh",
    "style": "detailed"
}
```

#### AI智能问答
```http
POST /api/v1/ai/chat
Authorization: Bearer <token>
Content-Type: application/json

{
    "question": "什么是分数？",
    "subject": "math",
    "grade": "三年级"
}
```

### 考试服务 API

#### 智能出题
```http
POST /api/v1/teacher/questions/generate
Authorization: Bearer <token>
Content-Type: application/json

{
    "subject": "math",
    "grade": "三年级",
    "question_types": ["single_choice", "fill_blank"],
    "difficulty": 2,
    "question_count": 10
}
```

#### 获取错题本
```http
GET /api/v1/wrong-book?subject=math&page=1&page_size=10
Authorization: Bearer <token>
```

## 数据库设计

### 主要表结构

- `users`: 用户信息表
- `courses`: 课程信息表
- `chapters`: 章节表
- `lessons`: 课时表
- `knowledge_points`: 知识点表
- `questions`: 题目表
- `exams`: 试卷表
- `student_progress`: 学习进度表
- `wrong_questions`: 错题记录表

完整的数据库设计请参考 `sql/schema.sql` 文件。

## 开发指南

### 项目结构

```
services/user/user/
├── user.go                    # 服务入口
├── etc/
│   └── user-api.yaml         # 配置文件
└── internal/
    ├── config/               # 配置定义
    ├── handler/              # HTTP处理器
    ├── logic/                # 业务逻辑
    ├── models/               # 数据模型
    ├── svc/                  # 服务上下文
    └── types/                # 类型定义
```

### 添加新服务

1. 在 `services/` 目录下创建新的服务目录
2. 定义API文件 (`.api`)
3. 使用goctl生成基础代码 (可选)
4. 实现业务逻辑
5. 配置服务启动

### 代码规范

- 遵循Go语言标准代码规范
- 使用统一的错误处理
- 添加必要的注释和文档
- 编写单元测试

## 部署说明

### Docker部署

```bash
# 构建镜像
docker build -t ai-education-platform/user-service .

# 运行容器
docker run -p 8881:8881 ai-education-platform/user-service
```

### 生产环境配置

1. 配置环境变量
2. 设置数据库连接池
3. 配置Redis集群
4. 设置负载均衡
5. 配置监控和日志

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码
4. 创建Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [Issues]
- 邮箱: admin@ai-education.com

---

**注意**: 这是一个完整的AI教育平台项目，包含用户管理、课程管理、AI助手、考试系统等核心功能。项目采用微服务架构，易于扩展和维护。