import { useState } from 'react'
import { Button, Form, Input, message } from 'antd'
import { useNavigate } from 'react-router-dom'

interface LoginForm {
  username: string
  password: string
}

export default function LoginPage() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const onFinish = async (values: LoginForm) => {
    setLoading(true)
    try {
      const res = await fetch('/api/v1/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      })
      const json = await res.json()
      if (json.code !== 200) {
        throw new Error(json.message || '登录失败')
      }
      localStorage.setItem('admin_token', json.data.token)
      message.success('登录成功')
      window.location.replace('/patients')
    } catch (err: any) {
      message.error(err.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="min-h-screen flex items-center justify-center relative overflow-hidden"
      style={{
        backgroundImage:
          'linear-gradient(135deg, #e0f2fe 0%, #bae6fd 50%, #7dd3fc 100%)',
      }}
    >
      {/* Medical-themed background decorations */}
      <div className="absolute inset-0 opacity-10 pointer-events-none">
        <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <pattern
              id="medical-cross"
              x="0"
              y="0"
              width="120"
              height="120"
              patternUnits="userSpaceOnUse"
            >
              <path
                d="M50 20 h20 v30 h30 v20 h-30 v30 h-20 v-30 h-30 v-20 h30 z"
                fill="#0284c7"
              />
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#medical-cross)" />
        </svg>
      </div>

      {/* Floating medical icons */}
      <div className="absolute top-10 left-10 opacity-20 pointer-events-none">
        <svg width="120" height="120" viewBox="0 0 24 24" fill="#0369a1">
          <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-4v4h-4v-4H6v-4h4V6h4v4h4v4z" />
        </svg>
      </div>
      <div className="absolute bottom-20 right-10 opacity-15 pointer-events-none">
        <svg width="180" height="180" viewBox="0 0 24 24" fill="#0284c7">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" />
        </svg>
      </div>
      <div className="absolute top-1/3 right-1/4 opacity-10 pointer-events-none">
        <svg width="100" height="100" viewBox="0 0 24 24" fill="#0ea5e9">
          <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-4v4h-4v-4H6v-4h4V6h4v4h4v4z" />
        </svg>
      </div>

      {/* Heartbeat line decoration */}
      <div className="absolute bottom-0 left-0 right-0 h-32 opacity-10 pointer-events-none">
        <svg width="100%" height="100%" preserveAspectRatio="none">
          <path
            d="M0,64 L200,64 L240,20 L280,108 L320,64 L400,64 L440,64 L480,30 L520,100 L560,64 L800,64 L840,64 L880,25 L920,110 L960,64 L1200,64 L1240,64 L1280,35 L1320,105 L1360,64 L1600,64 L1640,64 L1680,30 L1720,108 L1760,64 L2000,64"
            fill="none"
            stroke="#0369a1"
            strokeWidth="3"
          />
        </svg>
      </div>

      {/* Login card */}
      <div className="relative z-10 w-full max-w-md mx-4">
        <div className="bg-white/80 backdrop-blur-md p-8 rounded-2xl shadow-xl border border-white/50">
          {/* Logo / Title area */}
          <div className="flex items-center justify-center mb-6">
            <div className="w-14 h-14 bg-sky-500 rounded-xl flex items-center justify-center shadow-lg mr-3">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="white">
                <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-4v4h-4v-4H6v-4h4V6h4v4h4v4z" />
              </svg>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-slate-800">
                挂号系统
              </h1>
              <p className="text-sm text-slate-500">管理后台</p>
            </div>
          </div>

          <Form<LoginForm>
            layout="vertical"
            onFinish={onFinish}
            initialValues={{ username: 'admin', password: 'admin123' }}
          >
            <Form.Item
              label="用户名"
              name="username"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input size="large" placeholder="请输入用户名" />
            </Form.Item>
            <Form.Item
              label="密码"
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password size="large" placeholder="请输入密码" />
            </Form.Item>
            <Form.Item className="mb-2">
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
                size="large"
                className="bg-sky-500 hover:bg-sky-600"
              >
                登录
              </Button>
            </Form.Item>
          </Form>

          <div className="text-center text-xs text-slate-400 mt-4">
            默认账号：admin / admin123
          </div>
        </div>
      </div>
    </div>
  )
}
