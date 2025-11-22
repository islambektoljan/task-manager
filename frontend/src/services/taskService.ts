import { api } from './api';
import type {Task, CreateTaskRequest, UpdateTaskRequest} from '../types/task.ts';
import type {ApiResponse} from '../types/common';

export const taskService = {
    async getTasks(): Promise<ApiResponse<Task[]>> {
        const response = await api.get('/tasks');
        return response.data;
    },

    async getTask(id: string): Promise<ApiResponse<Task>> {
        const response = await api.get(`/tasks/${id}`);
        return response.data;
    },

    async createTask(task: CreateTaskRequest): Promise<ApiResponse<Task>> {
        const response = await api.post('/tasks', task);
        return response.data;
    },

    async updateTask(id: string, task: UpdateTaskRequest): Promise<ApiResponse<Task>> {
        const response = await api.put(`/tasks/${id}`, task);
        return response.data;
    },

    async deleteTask(id: string): Promise<ApiResponse<void>> {
        const response = await api.delete(`/tasks/${id}`);
        return response.data;
    },

    async updateTaskStatus(id: string, status: string): Promise<ApiResponse<Task>> {
        const response = await api.patch(`/tasks/${id}/status`, { status });
        return response.data;
    }
};