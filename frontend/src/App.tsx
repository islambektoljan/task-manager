import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { TaskProvider } from './contexts/TaskContext';
import PrivateRoute from './components/common/PrivateRoute';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import TaskPage from './pages/TaskPage';
import './App.css';

const App: React.FC = () => {
    return (
        <AuthProvider>
            <TaskProvider>
                <Router>
                    <div className="App">
                        <Routes>
                            <Route path="/login" element={<Login />} />
                            <Route path="/register" element={<Register />} />
                            <Route
                                path="/dashboard"
                                element={
                                    <PrivateRoute>
                                        <Dashboard />
                                    </PrivateRoute>
                                }
                            />
                            <Route
                                path="/tasks/:id"
                                element={
                                    <PrivateRoute>
                                        <TaskPage />
                                    </PrivateRoute>
                                }
                            />
                            <Route path="/" element={<Navigate to="/dashboard" replace />} />
                            <Route path="*" element={<Navigate to="/dashboard" replace />} />
                        </Routes>
                    </div>
                </Router>
            </TaskProvider>
        </AuthProvider>
    );
};

export default App;