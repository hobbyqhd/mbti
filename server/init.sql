-- 创建数据库
CREATE DATABASE IF NOT EXISTS mbti;
USE mbti;

-- 创建题目表
-- 使用7级量表评分：
-- 7-强烈同意
-- 6-同意
-- 5-略微同意
-- 4-中立
-- 3-略微反对
-- 2-反对
-- 1-强烈反对
CREATE TABLE IF NOT EXISTS questions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    question TEXT NOT NULL,
    dimension VARCHAR(2) NOT NULL,
    direction TINYINT NOT NULL COMMENT '1表示正向题目，-1表示反向题目',
    type VARCHAR(10) NOT NULL COMMENT '题目类型：simple-简单版，detailed-详细版，full-完整版',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
INSERT INTO questions (question, dimension, direction, type) VALUES
-- E-I维度（外向-内向）题目
-- 简单版（10题）
('我喜欢在社交场合中主动与他人交谈', 'EI', 1, 'simple'),
('我享受成为人群关注的焦点', 'EI', 1, 'simple'),
('我更喜欢和朋友一起度过周末', 'EI', 1, 'simple'),
('我能很快适应并融入新的环境', 'EI', 1, 'simple'),
('我更喜欢团队合作而不是独自工作', 'EI', 1, 'simple'),
('我在会议中经常主动发言', 'EI', 1, 'simple'),
('我喜欢认识新朋友', 'EI', 1, 'simple'),
('我享受组织和领导团队活动', 'EI', 1, 'simple'),
('我在公共场合表达自己很自在', 'EI', 1, 'simple'),
('我喜欢在开放的办公环境中工作', 'EI', 1, 'simple'),

-- 详细版新增（10题）
('我经常在社交媒体上分享生活', 'EI', 1, 'detailed'),
('我更喜欢结伴旅行', 'EI', 1, 'detailed'),
('我倾向于在团队中担任领导角色', 'EI', 1, 'detailed'),
('我喜欢参加社交活动和聚会', 'EI', 1, 'detailed'),
('我在陌生环境中容易建立新的人际关系', 'EI', 1, 'detailed'),
('我喜欢与他人分享自己的想法和感受', 'EI', 1, 'detailed'),
('我在团队讨论中经常提出建议', 'EI', 1, 'detailed'),
('我享受成为话题的中心', 'EI', 1, 'detailed'),
('我更喜欢在热闹的环境中工作', 'EI', 1, 'detailed'),
('我经常主动发起社交活动', 'EI', 1, 'detailed'),

-- 超详细版新增（10题）
('我在公共演讲时感到自在', 'EI', 1, 'full'),
('我喜欢参与小组讨论和头脑风暴', 'EI', 1, 'full'),
('我倾向于在社交场合中展示自己', 'EI', 1, 'full'),
('我享受与不同背景的人交流', 'EI', 1, 'full'),
('我在团队项目中表现更出色', 'EI', 1, 'full'),
('我喜欢在工作中与他人合作', 'EI', 1, 'full'),
('我经常在社交网络上发表观点', 'EI', 1, 'full'),
('我更愿意在团队中解决问题', 'EI', 1, 'full'),
('我喜欢参加社交网络活动', 'EI', 1, 'full'),
('我在群体活动中感到energized', 'EI', 1, 'full'),

-- S-N维度（感觉-直觉）题目
-- 简单版（10题）
('我更关注具体的事实和细节', 'SN', -1, 'simple'),
('我喜欢按部就班地学习新知识', 'SN', -1, 'simple'),
('我更相信实践经验', 'SN', -1, 'simple'),
('我注重计划的具体执行步骤', 'SN', -1, 'simple'),
('我倾向于使用已经证实有效的方法', 'SN', -1, 'simple'),
('阅读时我更注意具体的描述', 'SN', -1, 'simple'),
('我更依赖数据和事实做决策', 'SN', -1, 'simple'),
('我喜欢有明确规范的工作', 'SN', -1, 'simple'),
('我是一个务实的人', 'SN', -1, 'simple'),
('我更关注当下的实际情况', 'SN', -1, 'simple'),

-- 详细版新增（10题）
('我喜欢稳定有序的工作环境', 'SN', -1, 'detailed'),
('我更注重方案的可行性', 'SN', -1, 'detailed'),
('我倾向于循序渐进地解决问题', 'SN', -1, 'detailed'),
('我更看重实际的工作成果', 'SN', -1, 'detailed'),
('我喜欢按照既定流程工作', 'SN', -1, 'detailed'),
('我更相信亲身经历', 'SN', -1, 'detailed'),
('我注重细节的准确性', 'SN', -1, 'detailed'),
('我更愿意遵循传统的方法', 'SN', -1, 'detailed'),
('我喜欢具体的工作指导', 'SN', -1, 'detailed'),
('我更关注实际的应用价值', 'SN', -1, 'detailed'),

-- 超详细版新增（10题）
('我更重视实际操作经验', 'SN', -1, 'full'),
('我喜欢处理具体的问题', 'SN', -1, 'full'),
('我更看重现实的可能性', 'SN', -1, 'full'),
('我倾向于保持稳定的工作方式', 'SN', -1, 'full'),
('我更相信已经验证的方法', 'SN', -1, 'full'),
('我注重实际的工作效果', 'SN', -1, 'full'),
('我喜欢按照计划行事', 'SN', -1, 'full'),
('我更关注具体的执行细节', 'SN', -1, 'full'),
('我倾向于选择传统可靠的方案', 'SN', -1, 'full'),
('我更看重实际的工作经验', 'SN', -1, 'full'),

-- T-F维度（思维-情感）题目
-- 简单版（10题）
('我做决定时更依赖逻辑分析', 'TF', 1, 'simple'),
('我处理矛盾时会客观分析利弊', 'TF', 1, 'simple'),
('我评价他人时更看重能力和成就', 'TF', 1, 'simple'),
('我面对批评时会理性分析', 'TF', 1, 'simple'),
('我的决策主要基于客观事实', 'TF', 1, 'simple'),
('我更重视结果的效率', 'TF', 1, 'simple'),
('我更关注任务的完成度', 'TF', 1, 'simple'),
('我在讨论中强调逻辑严谨', 'TF', 1, 'simple'),
('我做项目时以目标达成为重', 'TF', 1, 'simple'),
('我评估方案时更看重成本效益', 'TF', 1, 'simple'),

-- 详细版新增（10题）
('我更重视客观的事实依据', 'TF', 1, 'detailed'),
('我制定规则时强调公平公正', 'TF', 1, 'detailed'),
('我更看重工作的专业性', 'TF', 1, 'detailed'),
('我倾向于用数据支持观点', 'TF', 1, 'detailed'),
('我更重视逻辑的连贯性', 'TF', 1, 'detailed'),
('我在评价时保持客观中立', 'TF', 1, 'detailed'),
('我更关注问题的本质', 'TF', 1, 'detailed'),
('我喜欢分析问题的因果关系', 'TF', 1, 'detailed'),
('我更重视决策的合理性', 'TF', 1, 'detailed'),
('我倾向于用理性方式解决问题', 'TF', 1, 'detailed'),

-- 超详细版新增（10题）
('我更看重论证的严密性', 'TF', 1, 'full'),
('我注重决策的客观依据', 'TF', 1, 'full'),
('我倾向于理性分析情况', 'TF', 1, 'full'),
('我更重视效率和成果', 'TF', 1, 'full'),
('我在工作中强调专业性', 'TF', 1, 'full'),
('我更看重逻辑的正确性', 'TF', 1, 'full'),
('我喜欢用数据说明问题', 'TF', 1, 'full'),
('我更重视客观的评判标准', 'TF', 1, 'full'),
('我倾向于理性决策', 'TF', 1, 'full'),
('我注重分析问题的逻辑性', 'TF', 1, 'full'),

-- J-P维度（判断-知觉）题目
-- 简单版（10题）
('我习惯提前计划和安排', 'JP', 1, 'simple'),
('我喜欢按步就班地完成任务', 'JP', 1, 'simple'),
('我会提前完成截止日期的工作', 'JP', 1, 'simple'),
('我习惯将物品分类存放', 'JP', 1, 'simple'),
('我喜欢详细规划行程', 'JP', 1, 'simple'),
('我会立即处理收到的文件', 'JP', 1, 'simple'),
('我倾向于快速做出决定', 'JP', 1, 'simple'),
('我喜欢保持环境整洁有序', 'JP', 1, 'simple'),
('我严格按计划执行项目', 'JP', 1, 'simple'),
('我喜欢固定的日程安排', 'JP', 1, 'simple'),

-- 详细版新增（10题）
('我习惯制定详细的工作计划', 'JP', 1, 'detailed'),
('我倾向于提前做好准备工作', 'JP', 1, 'detailed'),
('我喜欢按照既定计划行事', 'JP', 1, 'detailed'),
('我会定期整理工作环境', 'JP', 1, 'detailed'),
('我习惯将任务分类处理', 'JP', 1, 'detailed'),
('我更喜欢有条理的工作方式', 'JP', 1, 'detailed'),
('我会为项目设定明确的时间表', 'JP', 1, 'detailed'),
('我倾向于及时完成任务', 'JP', 1, 'detailed'),
('我喜欢按计划进行时间管理', 'JP', 1, 'detailed'),
('我习惯将工作按优先级排序', 'JP', 1, 'detailed'),

-- 超详细版新增（10题）
('我喜欢制定长期发展规划', 'JP', 1, 'full'),
('我会为每个目标设定具体步骤', 'JP', 1, 'full'),
('我习惯建立系统的工作流程', 'JP', 1, 'full'),
('我倾向于提前准备应对方案', 'JP', 1, 'full'),
('我重视工作中的条理性', 'JP', 1, 'full'),
('我会定期检查计划执行情况', 'JP', 1, 'full'),
('我喜欢有序地处理事务', 'JP', 1, 'full'),
('我习惯为团队制定明确规则', 'JP', 1, 'full'),
('我倾向于系统地解决问题', 'JP', 1, 'full'),
('我重视工作的计划性和条理性', 'JP', 1, 'full');