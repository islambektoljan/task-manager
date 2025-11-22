import React from 'react';
import TaskDetails from '../components/tasks/TaskDetails';

const TaskPage: React.FC = () => {
    return (
        <div className="min-h-screen bg-gray-100 py-8">
            <div className="container mx-auto px-4">
                <TaskDetails />
            </div>
        </div>
    );
};

export default TaskPage;