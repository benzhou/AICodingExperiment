// A simple test script to verify the login API endpoint
const axios = require('axios');

// Test admin login
async function testLogin() {
  const API_URL = 'http://localhost:8080/api/v1';
  
  const loginData = {
    email: 'admin@example.com',
    password: 'securepassword123'
  };
  
  try {
    console.log('Testing login with:', loginData);
    console.log('Request URL:', `${API_URL}/auth/login`);
    
    // Make the API request with detailed logging
    const response = await axios.post(`${API_URL}/auth/login`, loginData, {
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    console.log('Login successful!');
    console.log('Status:', response.status);
    console.log('Response data:', response.data);
    
    return true;
  } catch (error) {
    console.error('Login failed!');
    
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('Status:', error.response.status);
      console.error('Response data:', error.response.data);
      console.error('Response headers:', error.response.headers);
    } else if (error.request) {
      // The request was made but no response was received
      console.error('No response received from server');
      console.error('Request details:', error.request);
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('Error setting up request:', error.message);
    }
    
    return false;
  }
}

// Test with different endpoint format
async function testLoginAlternative() {
  try {
    const API_URL = 'http://localhost:8080';
    
    const loginData = {
      email: 'admin@example.com',
      password: 'securepassword123'
    };
    
    console.log('\nTrying alternative API URL format...');
    console.log('Request URL:', `${API_URL}/api/v1/auth/login`);
    
    const response = await axios.post(`${API_URL}/api/v1/auth/login`, loginData, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    });
    
    console.log('Alternative login successful!');
    console.log('Status:', response.status);
    console.log('Response data:', response.data);
    
    return true;
  } catch (error) {
    console.error('Alternative login failed!');
    
    if (error.response) {
      console.error('Status:', error.response.status);
      console.error('Response data:', error.response.data);
    } else if (error.request) {
      console.error('No response received');
    } else {
      console.error('Error:', error.message);
    }
    
    return false;
  }
}

// Check for server connection issues
async function testServerConnection() {
  try {
    console.log('\nTesting basic server connection...');
    const response = await axios.get('http://localhost:8080/health');
    console.log('Server is up! Status:', response.status);
    return true;
  } catch (error) {
    console.error('Server connection test failed.');
    if (error.code === 'ECONNREFUSED') {
      console.error('Connection refused - server may be down');
    } else {
      console.error('Error:', error.message);
    }
    return false;
  }
}

// Run the tests
async function runTests() {
  console.log('=== Login API Testing Tool ===');
  
  let serverUp = await testServerConnection();
  if (!serverUp) {
    console.error('\n❌ Server appears to be offline. Please start the backend server.');
    return;
  }
  
  let success = await testLogin();
  if (!success) {
    await testLoginAlternative();
  }
  
  console.log('\n=== Debugging tips if login fails ===');
  console.log('1. Check if the server is running (we already tested this)');
  console.log('2. Verify the API endpoint URL is correct');
  console.log('3. Ensure the login credentials match what\'s in the database');
  console.log('4. Check Content-Type header is application/json');
  console.log('5. Look for CORS issues in browser console');
  console.log('6. Examine request format in browser network tab');
  console.log('7. Check backend server logs for detailed error messages');
  console.log('8. Verify database connection in the backend');
}

runTests().catch(err => {
  console.error('Test runner error:', err);
}); 