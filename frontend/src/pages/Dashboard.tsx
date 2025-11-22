import React, { useState, useEffect } from 'react';
import { useTasks } from '../contexts/TaskContext';
import { useAuth } from '../contexts/AuthContext';
import Header from '../components/common/Header';
import TaskList from '../components/tasks/TaskList';
import TaskForm from '../components/tasks/TaskForm';
import LoadingSpinner from '../components/common/LoadingSpinner';
import ErrorMessage from '../components/common/ErrorMessage';

const Dashboard: React.FC = () => {
    const { user } = useAuth();
    const { tasks, loading, error, fetchTasks, clearError } = useTasks();
    const [showTaskForm, setShowTaskForm] = useState(false);

    useEffect(() => {
        fetchTasks();
    }, []); // Убрана зависимость fetchTasks

    const handleTaskCreated = () => {
        setShowTaskForm(false);
        fetchTasks(); // Явный вызов после создания задачи
    };

    return (
        <div className="min-h-screen bg-gray-50">
            <Header />

            <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <div className="flex justify-between items-center">
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                            <p className="text-gray-600 mt-2">
                                {user?.role === 'admin'
                                    ? 'You can view and manage all tasks as an administrator'
                                    : 'Manage your tasks and track your progress'
                                }
                            </p>
                        </div>
                        <button
                            onClick={() => setShowTaskForm(!showTaskForm)}
                            className="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 font-medium"
                        >
                            {showTaskForm ? 'Cancel' : '+ New Task'}
                        </button>
                    </div>

                    {user?.role === 'admin' && (
                        <div className="mt-4 bg-blue-50 border border-blue-200 rounded-lg p-4">
                            <div className="flex items-center">
                                <div className="flex-shrink-0">
                                    <svg className="h-5 w-5 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                                    </svg>
                                </div>
                                <div className="ml-3">
                                    <p className="text-sm text-blue-700">
                                        <strong>Admin Mode:</strong> You are viewing all tasks in the system.
                                    </p>
                                </div>
                            </div>
                        </div>
                    )}
                </div>

                {error && (
                    <div className="mb-6">  
                        <ErrorMessage message={error} onClose={clearError} />
                    </div>
                )}

                {showTaskForm && (
                    <div className="mb-8">
                        <TaskForm onSuccess={handleTaskCreated} onCancel={() => setShowTaskForm(false)} />
                    </div>
                )}

                <div className="bg-white rounded-lg shadow">
                    <div className="px-6 py-4 border-b border-gray-200">
                        <h2 className="text-xl font-semibold text-gray-800">
                            Tasks ({tasks.length})
                        </h2>
                    </div>
                    <div className="p-6">
                        {loading ? <LoadingSpinner /> : <TaskList />}
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Dashboard;