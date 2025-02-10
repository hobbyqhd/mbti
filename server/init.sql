-- 创建数据库
CREATE DATABASE IF NOT EXISTS mbti;
USE mbti;

-- 创建题目表
CREATE TABLE IF NOT EXISTS questions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    question TEXT NOT NULL,
    options JSON NOT NULL
);

-- 创建结果表
CREATE TABLE IF NOT EXISTS results (
    id VARCHAR(50) PRIMARY KEY,
    type VARCHAR(4) NOT NULL,
    dimensions JSON NOT NULL,
    report TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入MBTI测试题目
INSERT INTO questions (question, options) VALUES
('在社交场合，你通常会...', '["主动与他人交谈","等待他人来与你交谈"]'),
('你更倾向于...', '["关注现实和具体的细节","关注概念和可能性"]'),
('做决定时，你更看重...', '["逻辑和客观分析","个人价值和感受"]'),
('你更喜欢...', '["按计划行事","随机应变"]'),
-- 继续插入其他89个问题...
('你更喜欢的工作环境是...', '["有条理和结构化的","灵活和适应性强的"]');