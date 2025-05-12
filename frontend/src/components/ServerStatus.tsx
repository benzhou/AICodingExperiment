import React, { useState, useEffect } from 'react';
import { Alert, Button } from 'antd';
import api from '../services/api';

const ServerStatus: React.FC = () => {
  const [status, setStatus] = useState<string>('Checking...');
  const [error, setError] = useState<string | null>(null);

  const checkServerStatus = async () => {
    try {
      setStatus('Checking...');
      setError(null);
      
      const token = localStorage.getItem('token');
      
      // Log token info for debugging (truncated for security)
      if (token) {
        console.log('Token available, length:', token.length);
        console.log('Token prefix:', token.substring(0, 15) + '...');
      } else {
        console.warn('No token found in localStorage');
        setStatus('No Auth');
        setError('You need to log in first');
        return;
      }
      
      // Try to reach the backend using the health endpoint which doesn't require auth
      const healthResponse = await api.get('/health');
      
      if (healthResponse.data && healthResponse.data.status === 'ok') {
        console.log('Server health check response:', healthResponse.data);
        setStatus('Connected');
      } else {
        setStatus('Error');
        setError('Unexpected response from server');
      }
    } catch (err: any) {
      console.error('Server connection error:', err);
      if (err.response) {
        // The request was made and the server responded with a status code outside of 2xx
        setStatus(`Error ${err.response.status}`);
        setError(err.response.data?.message || 'Authentication error');
      } else if (err.request) {
        // The request was made but no response was received
        setStatus('No Response');
        setError('The server is not responding. Is the backend running?');
      } else {
        // Something happened in setting up the request that triggered an Error
        setStatus('Request Failed');
        setError(err.message);
      }
    }
  };

  useEffect(() => {
    checkServerStatus();
  }, []);

  return (
    <div className="server-status my-2 p-2 border rounded">
      <div className="flex items-center">
        <div className="mr-2">
          <span className={`inline-block w-3 h-3 rounded-full ${
            status === 'Connected' ? 'bg-green-500' :
            status === 'Checking...' ? 'bg-yellow-500' : 'bg-red-500'
          }`}></span>
        </div>
        <div>
          <strong>Backend Server:</strong> {status}
        </div>
        <Button 
          size="small" 
          className="ml-2" 
          onClick={checkServerStatus}
        >
          Check
        </Button>
      </div>
      
      {error && (
        <Alert 
          type="error" 
          message={error} 
          className="mt-2"
        />
      )}
    </div>
  );
};

export default ServerStatus; 