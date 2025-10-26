CREATE TABLE IF NOT EXISTS task_schema.tasks (
                                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                 title VARCHAR(255) NOT NULL,
                                                 description TEXT,
                                                 status VARCHAR(50) DEFAULT 'pending',
                                                 priority VARCHAR(50) DEFAULT 'medium',
                                                 due_date TIMESTAMP,
                                                 created_by UUID NOT NULL,
                                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON task_schema.tasks(created_by);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON task_schema.tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON task_schema.tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON task_schema.tasks(priority);

-- Комментарии к таблице и колонкам
COMMENT ON TABLE task_schema.tasks IS 'Stores user tasks with status and priority tracking';
COMMENT ON COLUMN task_schema.tasks.status IS 'Task status: pending, in_progress, completed, cancelled';
COMMENT ON COLUMN task_schema.tasks.priority IS 'Task priority: low, medium, high, urgent';