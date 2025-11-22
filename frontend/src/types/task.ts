export interface Task {
    id: string;
    title: string;
    description: string;
    status: 'pending' | 'in progress' | 'completed' | 'cancelled';
    priority: 'low' | 'medium' | 'high' | 'urgent';
    due_date?: string;
    created_by: string;
    created_at: string;
    updated_at: string;
}

export interface CreateTaskRequest {
    title: string;
    description: string;
    status?: string;
    priority?: string;
    due_date?: string;
}

export interface UpdateTaskRequest {
    title?: string;
    description?: string;
    status?: string;
    priority?: string;
    due_date?: string;
}

export interface TaskState {
    tasks: Task[];
    currentTask: Task | null;
    loading: boolean;
    error: string | null;
}