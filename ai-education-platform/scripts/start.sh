#!/bin/bash

# AI教育平台启动脚本

set -e

echo "=== AI教育平台启动脚本 ==="

# 检查依赖
check_dependencies() {
    echo "检查依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "错误: Go未安装"
        exit 1
    fi
    
    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        echo "警告: MySQL未找到，请确保数据库已配置"
    fi
    
    # 检查Redis
    if ! command -v redis-cli &> /dev/null; then
        echo "警告: Redis未找到，请确保Redis已配置"
    fi
    
    echo "依赖检查完成"
}

# 安装依赖
install_dependencies() {
    echo "安装Go依赖..."
    go mod tidy
    echo "依赖安装完成"
}

# 初始化数据库
init_database() {
    echo "初始化数据库..."
    
    # 检查数据库连接
    if mysql -u root -p -e "USE ai_education_platform;" 2>/dev/null; then
        echo "数据库已存在，跳过初始化"
    else
        echo "创建数据库..."
        mysql -u root -p < sql/schema.sql
        echo "数据库初始化完成"
    fi
}

# 构建服务
build_services() {
    echo "构建服务..."
    
    # 构建用户服务
    cd services/user/user
    go build -o user-service user.go
    echo "用户服务构建完成"
    
    # 构建课程服务
    cd ../../course/course
    go build -o course-service course.go
    echo "课程服务构建完成"
    
    # 构建AI助手服务
    cd ../../ai_assistant/ai_assistant
    go build -o ai-service ai_assistant.go
    echo "AI助手服务构建完成"
    
    # 构建考试服务
    cd ../../exam/exam
    go build -o exam-service exam.go
    echo "考试服务构建完成"
    
    cd ../../../..
}

# 启动服务
start_services() {
    echo "启动服务..."
    
    # 启动用户服务
    cd services/user/user
    ./user-service -f etc/user-api.yaml &
    echo "用户服务已启动 (端口: 8881)"
    
    # 启动课程服务
    cd ../../course/course
    ./course-service -f etc/course-api.yaml &
    echo "课程服务已启动 (端口: 8882)"
    
    # 启动AI助手服务
    cd ../../ai_assistant/ai_assistant
    ./ai-service -f etc/ai_assistant-api.yaml &
    echo "AI助手服务已启动 (端口: 8883)"
    
    # 启动考试服务
    cd ../../exam/exam
    ./exam-service -f etc/exam-api.yaml &
    echo "考试服务已启动 (端口: 8884)"
    
    cd ../../../..
    
    echo "所有服务已启动"
    echo "用户服务: http://localhost:8881"
    echo "课程服务: http://localhost:8882"
    echo "AI助手服务: http://localhost:8883"
    echo "考试服务: http://localhost:8884"
}

# 停止服务
stop_services() {
    echo "停止服务..."
    
    # 停止所有相关进程
    pkill -f "user-service" || true
    pkill -f "course-service" || true
    pkill -f "ai-service" || true
    pkill -f "exam-service" || true
    
    echo "服务已停止"
}

# 检查服务状态
check_status() {
    echo "检查服务状态..."
    
    services=("user-service:8881" "course-service:8882" "ai-service:8883" "exam-service:8884")
    
    for service in "${services[@]}"; do
        name=$(echo $service | cut -d':' -f1)
        port=$(echo $service | cut -d':' -f2)
        
        if lsof -i :$port > /dev/null; then
            echo "✓ $name 运行中 (端口: $port)"
        else
            echo "✗ $name 未运行"
        fi
    done
}

# 显示帮助信息
show_help() {
    echo "AI教育平台启动脚本"
    echo ""
    echo "用法:"
    echo "  ./start.sh [命令]"
    echo ""
    echo "命令:"
    echo "  install     安装依赖和初始化数据库"
    echo "  build       构建所有服务"
    echo "  start       启动所有服务"
    echo "  stop        停止所有服务"
    echo "  restart     重启所有服务"
    echo "  status      检查服务状态"
    echo "  help        显示帮助信息"
    echo ""
    echo "示例:"
    echo "  ./start.sh install    # 首次运行"
    echo "  ./start.sh start      # 启动服务"
    echo "  ./start.sh status     # 检查状态"
}

# 主函数
main() {
    case "${1:-}" in
        "install")
            check_dependencies
            install_dependencies
            init_database
            build_services
            echo "安装完成！使用 './start.sh start' 启动服务"
            ;;
        "build")
            build_services
            echo "构建完成！"
            ;;
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            sleep 2
            start_services
            ;;
        "status")
            check_status
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        "")
            echo "错误: 请指定命令"
            echo "使用 './start.sh help' 查看帮助"
            exit 1
            ;;
        *)
            echo "错误: 未知命令 '$1'"
            echo "使用 './start.sh help' 查看帮助"
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"