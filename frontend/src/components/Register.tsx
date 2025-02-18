import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { authService } from '../services/authService';
import { useUser } from '../context/UserContext';
import { Toast } from './common/Toast';

interface FormErrors {
  name?: string;
  email?: string;
  password?: string;
}

export function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<FormErrors>({});
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { setUser } = useUser();

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};
    
    if (!name) newErrors.name = t('auth.errors.required');
    if (!email) newErrors.email = t('auth.errors.required');
    else if (!/\S+@\S+\.\S+/.test(email)) newErrors.email = t('auth.errors.invalidEmail');
    if (!password) newErrors.password = t('auth.errors.required');
    else if (password.length < 6) newErrors.password = t('auth.errors.passwordLength');

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;

    try {
      const response = await authService.register({ name, email, password });
      setUser(response.user);
      setToast({
        message: t('auth.success.registration'),
        type: 'success'
      });
      
      // Redirect after a short delay to show the success message
      setTimeout(() => {
        navigate('/dashboard');
      }, 1500);
    } catch (err: any) {
      setToast({
        message: err.response?.data?.message || t('auth.errors.registrationFailed'),
        type: 'error'
      });
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
        <div>
          <h2 className="text-center text-3xl font-extrabold text-gray-900">
            {t('auth.labels.register')}
          </h2>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="name" className="sr-only">
              {t('auth.labels.name')}
            </label>
            <input
              id="name"
              type="text"
              required
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.name ? 'border-red-500' : 'border-gray-300'
              } placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm`}
              placeholder={t('auth.labels.name')}
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            {errors.name && (
              <p className="mt-1 text-sm text-red-500">{errors.name}</p>
            )}
          </div>
          <div>
            <label htmlFor="email" className="sr-only">
              {t('auth.labels.email')}
            </label>
            <input
              id="email"
              type="email"
              required
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.email ? 'border-red-500' : 'border-gray-300'
              } placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm`}
              placeholder={t('auth.labels.email')}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-500">{errors.email}</p>
            )}
          </div>
          <div>
            <label htmlFor="password" className="sr-only">
              {t('auth.labels.password')}
            </label>
            <input
              id="password"
              type="password"
              required
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.password ? 'border-red-500' : 'border-gray-300'
              } placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm`}
              placeholder={t('auth.labels.password')}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-500">{errors.password}</p>
            )}
          </div>
          <div>
            <button
              type="submit"
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              {t('auth.labels.register')}
            </button>
          </div>
          <div className="text-sm text-center">
            <Link to="/login" className="text-indigo-600 hover:text-indigo-500">
              {t('auth.labels.haveAccount')}
            </Link>
          </div>
        </form>
      </div>
      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
} 