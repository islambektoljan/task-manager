import React from 'react';

interface ErrorMessageProps {
    message: string;
    onClose?: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, onClose }) => {
    return (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong className="font-bold">Error: </strong>
            <span className="block sm:inline">{message}</span>
            {onClose && (
                <button
                    className="absolute top-0 right-0 px-4 py-3"
                    onClick={onClose}
                >
                    <span className="text-red-500">×</span>
                </button>
            )}
        </div>
    );
};

export default ErrorMessage;