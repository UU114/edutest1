-- AI教育平台数据库结构设计
-- 创建数据库
CREATE DATABASE IF NOT EXISTS ai_education_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE ai_education_platform;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码(加密存储)',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    nickname VARCHAR(50) NOT NULL COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像URL',
    role ENUM('student', 'teacher', 'parent', 'institution') NOT NULL COMMENT '用户角色',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0-禁用,1-正常,2-待审核',
    real_name VARCHAR(50) COMMENT '真实姓名',
    gender TINYINT DEFAULT 0 COMMENT '性别:0-未知,1-男,2-女',
    birthday DATE COMMENT '生日',
    school VARCHAR(100) COMMENT '学校',
    grade VARCHAR(20) COMMENT '年级',
    class VARCHAR(20) COMMENT '班级',
    parent_id BIGINT COMMENT '家长ID',
    institution_id BIGINT COMMENT '机构ID',
    bio TEXT COMMENT '个人简介',
    last_login_at BIGINT COMMENT '最后登录时间',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_status (status),
    INDEX idx_parent_id (parent_id),
    INDEX idx_institution_id (institution_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 课程表
CREATE TABLE IF NOT EXISTS courses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '课程标题',
    description TEXT NOT NULL COMMENT '课程描述',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    difficulty TINYINT NOT NULL DEFAULT 2 COMMENT '难度:1-简单,2-中等,3-困难',
    cover_image VARCHAR(255) COMMENT '封面图片URL',
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '价格',
    duration INT NOT NULL COMMENT '课程时长(分钟)',
    teacher_id BIGINT NOT NULL COMMENT '教师ID',
    institution_id BIGINT NOT NULL COMMENT '机构ID',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态:0-草稿,1-发布,2-下架',
    student_count INT NOT NULL DEFAULT 0 COMMENT '学生数量',
    rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分',
    objectives TEXT COMMENT '学习目标',
    prerequisites TEXT COMMENT '先修要求',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_teacher_id (teacher_id),
    INDEX idx_institution_id (institution_id),
    INDEX idx_status (status),
    INDEX idx_rating (rating),
    FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (institution_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 章节表
CREATE TABLE IF NOT EXISTS chapters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL COMMENT '课程ID',
    title VARCHAR(200) NOT NULL COMMENT '章节标题',
    description TEXT NOT NULL COMMENT '章节描述',
    `order` INT NOT NULL COMMENT '排序',
    duration INT NOT NULL COMMENT '章节时长(分钟)',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_course_id (course_id),
    INDEX idx_order (`order`),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节表';

-- 课时表
CREATE TABLE IF NOT EXISTS lessons (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    chapter_id BIGINT NOT NULL COMMENT '章节ID',
    title VARCHAR(200) NOT NULL COMMENT '课时标题',
    type ENUM('video', 'audio', 'document', 'exercise') NOT NULL COMMENT '课时类型',
    content TEXT NOT NULL COMMENT '课时内容',
    duration INT NOT NULL COMMENT '课时时长(分钟)',
    `order` INT NOT NULL COMMENT '排序',
    is_free BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否免费',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_chapter_id (chapter_id),
    INDEX idx_order (`order`),
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课时表';

-- 知识点表
CREATE TABLE IF NOT EXISTS knowledge_points (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL COMMENT '课程ID',
    chapter_id BIGINT COMMENT '章节ID',
    title VARCHAR(200) NOT NULL COMMENT '知识点标题',
    description TEXT NOT NULL COMMENT '知识点描述',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    difficulty TINYINT NOT NULL DEFAULT 2 COMMENT '难度:1-简单,2-中等,3-困难',
    keywords JSON COMMENT '关键词',
    prerequisites JSON COMMENT '前置知识点ID数组',
    created_at BIGINT NOT NULL,
    INDEX idx_course_id (course_id),
    INDEX idx_chapter_id (chapter_id),
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_difficulty (difficulty),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识点表';

-- 题目表
CREATE TABLE IF NOT EXISTS questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(500) NOT NULL COMMENT '题目标题',
    type ENUM('multiple_choice', 'single_choice', 'fill_blank', 'essay', 'calculation', 'judgment') NOT NULL COMMENT '题目类型',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    difficulty TINYINT NOT NULL DEFAULT 3 COMMENT '难度:1-5',
    content TEXT NOT NULL COMMENT '题目内容',
    options JSON COMMENT '选择题选项',
    correct_answer TEXT NOT NULL COMMENT '正确答案',
    analysis TEXT COMMENT '解析',
    knowledge_points JSON COMMENT '关联知识点ID数组',
    tags JSON COMMENT '题目标签',
    creator_id BIGINT NOT NULL COMMENT '创建者ID',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0-草稿,1-发布,2-禁用',
    usage_count INT NOT NULL DEFAULT 0 COMMENT '使用次数',
    correct_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '正确率',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_type (type),
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_difficulty (difficulty),
    INDEX idx_creator_id (creator_id),
    INDEX idx_status (status),
    INDEX idx_usage_count (usage_count),
    FULLTEXT INDEX ft_content (content, title),
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目表';

-- 试卷表
CREATE TABLE IF NOT EXISTS exams (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '试卷标题',
    description TEXT NOT NULL COMMENT '试卷描述',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    duration INT NOT NULL COMMENT '考试时长(分钟)',
    total_score DECIMAL(10,2) NOT NULL COMMENT '总分',
    pass_score DECIMAL(10,2) NOT NULL COMMENT '及格分',
    question_count INT NOT NULL COMMENT '题目数量',
    creator_id BIGINT NOT NULL COMMENT '创建者ID',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态:0-草稿,1-发布,2-禁用',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_creator_id (creator_id),
    INDEX idx_status (status),
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='试卷表';

-- 试卷题目关联表
CREATE TABLE IF NOT EXISTS exam_questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    exam_id BIGINT NOT NULL COMMENT '试卷ID',
    question_id BIGINT NOT NULL COMMENT '题目ID',
    score DECIMAL(5,2) NOT NULL COMMENT '分值',
    `order` INT NOT NULL COMMENT '题目顺序',
    created_at BIGINT NOT NULL,
    INDEX idx_exam_id (exam_id),
    INDEX idx_question_id (question_id),
    INDEX idx_order (`order`),
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='试卷题目关联表';

-- 学生学习进度表
CREATE TABLE IF NOT EXISTS student_progress (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL COMMENT '课程ID',
    user_id BIGINT NOT NULL COMMENT '学生ID',
    completed_lessons INT NOT NULL DEFAULT 0 COMMENT '已完成课时数',
    total_lessons INT NOT NULL COMMENT '总课时数',
    progress DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT '进度百分比',
    last_lesson_id BIGINT COMMENT '最后学习课时ID',
    study_time INT NOT NULL DEFAULT 0 COMMENT '学习时长(分钟)',
    started_at BIGINT NOT NULL COMMENT '开始时间',
    completed_at BIGINT COMMENT '完成时间',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    UNIQUE KEY uk_course_user (course_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_progress (progress),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生学习进度表';

-- 课程报名表
CREATE TABLE IF NOT EXISTS course_enrollments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL COMMENT '课程ID',
    user_id BIGINT NOT NULL COMMENT '学生ID',
    status ENUM('pending', 'active', 'completed', 'cancelled') NOT NULL DEFAULT 'pending' COMMENT '状态',
    enrolled_at BIGINT NOT NULL COMMENT '报名时间',
    completed_at BIGINT COMMENT '完成时间',
    created_at BIGINT NOT NULL,
    UNIQUE KEY uk_course_user (course_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程报名表';

-- 考试记录表
CREATE TABLE IF NOT EXISTS exam_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    exam_id BIGINT NOT NULL COMMENT '考试ID',
    user_id BIGINT NOT NULL COMMENT '学生ID',
    score DECIMAL(10,2) COMMENT '得分',
    total_score DECIMAL(10,2) NOT NULL COMMENT '总分',
    status ENUM('in_progress', 'completed', 'timeout', 'submitted') NOT NULL DEFAULT 'in_progress' COMMENT '考试状态',
    start_time BIGINT NOT NULL COMMENT '开始时间',
    end_time BIGINT COMMENT '结束时间',
    time_used INT NOT NULL DEFAULT 0 COMMENT '用时(秒)',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_exam_id (exam_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='考试记录表';

-- 考试答案表
CREATE TABLE IF NOT EXISTS exam_answers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    exam_record_id BIGINT NOT NULL COMMENT '考试记录ID',
    question_id BIGINT NOT NULL COMMENT '题目ID',
    answer TEXT NOT NULL COMMENT '学生答案',
    is_correct BOOLEAN COMMENT '是否正确',
    score DECIMAL(5,2) COMMENT '得分',
    time_spent INT NOT NULL DEFAULT 0 COMMENT '答题用时(秒)',
    created_at BIGINT NOT NULL,
    INDEX idx_exam_record_id (exam_record_id),
    INDEX idx_question_id (question_id),
    FOREIGN KEY (exam_record_id) REFERENCES exam_records(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='考试答案表';

-- 错题记录表
CREATE TABLE IF NOT EXISTS wrong_questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '学生ID',
    question_id BIGINT NOT NULL COMMENT '题目ID',
    student_answer TEXT NOT NULL COMMENT '学生答案',
    correct_answer TEXT NOT NULL COMMENT '正确答案',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    wrong_count INT NOT NULL DEFAULT 1 COMMENT '错误次数',
    last_wrong_time BIGINT NOT NULL COMMENT '最后错误时间',
    mastered BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已掌握',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_question_id (question_id),
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_mastered (mastered),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='错题记录表';

-- 练习记录表
CREATE TABLE IF NOT EXISTS practice_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '学生ID',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    question_count INT NOT NULL COMMENT '题目数量',
    correct_count INT NOT NULL COMMENT '正确数量',
    score DECIMAL(5,2) NOT NULL COMMENT '得分',
    time_used INT NOT NULL COMMENT '用时(秒)',
    difficulty VARCHAR(20) NOT NULL COMMENT '难度',
    created_at BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='练习记录表';

-- AI对话记录表
CREATE TABLE IF NOT EXISTS ai_chat_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    session_id VARCHAR(100) NOT NULL COMMENT '会话ID',
    role ENUM('user', 'assistant') NOT NULL COMMENT '角色',
    content TEXT NOT NULL COMMENT '消息内容',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') COMMENT '学科',
    grade VARCHAR(20) COMMENT '年级',
    created_at BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI对话记录表';

-- AI讲解记录表
CREATE TABLE IF NOT EXISTS ai_explanations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    explanation_id VARCHAR(100) NOT NULL COMMENT '讲解ID',
    subject ENUM('math', 'chinese', 'english', 'physics', 'chemistry', 'biology', 'history', 'geography', 'politics') NOT NULL COMMENT '学科',
    grade VARCHAR(20) NOT NULL COMMENT '年级',
    topic VARCHAR(200) NOT NULL COMMENT '主题',
    difficulty TINYINT NOT NULL COMMENT '难度',
    content TEXT NOT NULL COMMENT '讲解内容',
    summary TEXT COMMENT '摘要',
    key_points JSON COMMENT '关键点',
    estimated_time INT COMMENT '预估学习时间',
    created_at BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_subject (subject),
    INDEX idx_grade (grade),
    INDEX idx_topic (topic),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI讲解记录表';

-- 课程标签表
CREATE TABLE IF NOT EXISTS course_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL COMMENT '课程ID',
    tag VARCHAR(50) NOT NULL COMMENT '标签',
    created_at BIGINT NOT NULL,
    INDEX idx_course_id (course_id),
    INDEX idx_tag (tag),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程标签表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_config (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT NOT NULL COMMENT '配置值',
    description VARCHAR(500) COMMENT '配置描述',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 初始化系统配置数据
INSERT INTO system_config (config_key, config_value, description, created_at, updated_at) VALUES
('site_name', 'AI教育平台', '网站名称', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('site_description', '智能化教育学习平台', '网站描述', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('max_file_size', '10485760', '最大文件上传大小(字节)', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('default_avatar', '/static/images/default-avatar.png', '默认头像', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('ai_service_timeout', '30', 'AI服务超时时间(秒)', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('max_daily_ai_requests', '100', '每日AI请求最大次数', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 创建管理员账户 (密码: admin123)
INSERT INTO users (username, password, email, nickname, role, status, created_at, updated_at) VALUES
('admin', '$2a$10$rOZXp7mGXmHWK7vJtxB7uO5D3Q7J8Y.K.KYzLkHqJ2Q7Y.K.YzLkHq', 'admin@ai-education.com', '系统管理员', 'institution', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());