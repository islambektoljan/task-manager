import React, { useEffect } from 'react';
import { useTasks } from '../../contexts/TaskContext';
import TaskItem from './TaskItem.tsx';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';

const TaskList: React.FC = () => {
    const { tasks, loading, error, fetchTasks, clearError } = useTasks();

    useEffect(() => {
        fetchTasks();
    }, [fetchTasks]);

    if (loading) {
        return <LoadingSpinner />;
    }

    if (error) {
        return <ErrorMessage message={error} onClose={clearError} />;
    }

    return (
        <div className="space-y-4">
            <h2 className="text-2xl font-bold">Tasks</h2>
            {tasks.length === 0 ? (
                <p className="text-gray-500">No tasks found. Create your first task!</p>
            ) : (
                tasks.map(task => <TaskItem key={task.id} task={task} />)
            )}
        </div>
    );
};

export default TaskList;