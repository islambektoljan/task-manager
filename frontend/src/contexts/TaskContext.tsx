import React, { createContext, useContext, useReducer, useCallback } from 'react';
import type { TaskState, Task, CreateTaskRequest, UpdateTaskRequest } from '../types/task';
import { taskService } from '../services/taskService';

type TaskAction =
    | { type: 'TASKS_LOADING' }
    | { type: 'TASKS_SUCCESS'; payload: Task[] }
    | { type: 'TASK_SUCCESS'; payload: Task }
    | { type: 'TASKS_ERROR'; payload: string }
    | { type: 'TASK_CREATE'; payload: Task }
    | { type: 'TASK_UPDATE'; payload: Task }
    | { type: 'TASK_DELETE'; payload: string }
    | { type: 'CLEAR_ERROR' }
    | { type: 'CLEAR_CURRENT_TASK' };

const initialState: TaskState = {
    tasks: [],
    currentTask: null,
    loading: false,
    error: null,
};

const taskReducer = (state: TaskState, action: TaskAction): TaskState => {
    switch (action.type) {
        case 'TASKS_LOADING':
            return { ...state, loading: true, error: null };
        case 'TASKS_SUCCESS':
            return { ...state, loading: false, tasks: action.payload };
        case 'TASK_SUCCESS':
            return { ...state, loading: false, currentTask: action.payload };
        case 'TASKS_ERROR':
            return { ...state, loading: false, error: action.payload };
        case 'TASK_CREATE':
            return { ...state, tasks: [...state.tasks, action.payload] };
        case 'TASK_UPDATE':
            return {
                ...state,
                tasks: state.tasks.map(task =>
                    task.id === action.payload.id ? action.payload : task
                ),
                currentTask: action.payload,
            };
        case 'TASK_DELETE':
            return {
                ...state,
                tasks: state.tasks.filter(task => task.id !== action.payload),
                currentTask: state.currentTask?.id === action.payload ? null : state.currentTask,
            };
        case 'CLEAR_ERROR':
            return { ...state, error: null };
        case 'CLEAR_CURRENT_TASK':
            return { ...state, currentTask: null };
        default:
            return state;
    }
};

interface TaskContextType extends TaskState {
    fetchTasks: () => Promise<void>;
    fetchTask: (id: string) => Promise<void>;
    createTask: (task: CreateTaskRequest) => Promise<void>;
    updateTask: (id: string, task: UpdateTaskRequest) => Promise<void>;
    deleteTask: (id: string) => Promise<void>;
    updateTaskStatus: (id: string, status: string) => Promise<void>;
    clearError: () => void;
    clearCurrentTask: () => void;
}

const TaskContext = createContext<TaskContextType | undefined>(undefined);

export const TaskProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [state, dispatch] = useReducer(taskReducer, initialState);

    const fetchTasks = useCallback(async () => {
        dispatch({ type: 'TASKS_LOADING' });
        try {
            const response = await taskService.getTasks();
            if (response.success && response.data) {
                dispatch({ type: 'TASKS_SUCCESS', payload: response.data });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to fetch tasks' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to fetch tasks',
            });
        }
    }, []);

    const fetchTask = useCallback(async (id: string) => {
        dispatch({ type: 'TASKS_LOADING' });
        try {
            const response = await taskService.getTask(id);
            if (response.success && response.data) {
                dispatch({ type: 'TASK_SUCCESS', payload: response.data });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to fetch task' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to fetch task',
            });
        }
    }, []);

    const createTask = useCallback(async (taskData: CreateTaskRequest) => {
        dispatch({ type: 'TASKS_LOADING' });
        try {
            const response = await taskService.createTask(taskData);
            if (response.success && response.data) {
                dispatch({ type: 'TASK_CREATE', payload: response.data });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to create task' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to create task',
            });
        }
    }, []);

    const updateTask = useCallback(async (id: string, taskData: UpdateTaskRequest) => {
        dispatch({ type: 'TASKS_LOADING' });
        try {
            const response = await taskService.updateTask(id, taskData);
            if (response.success && response.data) {
                dispatch({ type: 'TASK_UPDATE', payload: response.data });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to update task' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to update task',
            });
        }
    }, []);

    const deleteTask = useCallback(async (id: string) => {
        dispatch({ type: 'TASKS_LOADING' });
        try {
            const response = await taskService.deleteTask(id);
            if (response.success) {
                dispatch({ type: 'TASK_DELETE', payload: id });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to delete task' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to delete task',
            });
        }
    }, []);

    const updateTaskStatus = useCallback(async (id: string, status: string) => {
        try {
            const response = await taskService.updateTaskStatus(id, status);
            if (response.success && response.data) {
                dispatch({ type: 'TASK_UPDATE', payload: response.data });
            } else {
                dispatch({ type: 'TASKS_ERROR', payload: response.error || 'Failed to update task status' });
            }
        } catch (error: any) {
            dispatch({
                type: 'TASKS_ERROR',
                payload: error.response?.data?.error || 'Failed to update task status',
            });
        }
    }, []);

    const clearError = useCallback(() => {
        dispatch({ type: 'CLEAR_ERROR' });
    }, []);

    const clearCurrentTask = useCallback(() => {
        dispatch({ type: 'CLEAR_CURRENT_TASK' });
    }, []);

    return (
        <TaskContext.Provider
            value={{
                ...state,
                fetchTasks,
                fetchTask,
                createTask,
                updateTask,
                deleteTask,
                updateTaskStatus,
                clearError,
                clearCurrentTask,
            }}
        >
            {children}
        </TaskContext.Provider>
    );
};

export const useTasks = () => {
    const context = useContext(TaskContext);
    if (context === undefined) {
        throw new Error('useTasks must be used within a TaskProvider');
    }
    return context;
};