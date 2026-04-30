import { Layout, Menu, Button } from 'antd'
import { UserOutlined, MedicineBoxOutlined, FileTextOutlined, LogoutOutlined } from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'

const { Header, Sider, Content } = Layout

export default function BasicLayout({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate()
  const location = useLocation()

  const menuItems = [
    { key: '/patients', icon: <UserOutlined />, label: '就诊人管理' },
    { key: '/schedules', icon: <MedicineBoxOutlined />, label: '号源管理' },
    { key: '/orders', icon: <FileTextOutlined />, label: '挂号订单' },
  ]

  const handleLogout = () => {
    localStorage.removeItem('admin_token')
    navigate('/login')
  }

  return (
    <Layout className="min-h-screen">
      <Header className="bg-white shadow-sm flex items-center justify-between px-6">
        <h1 className="text-lg font-semibold m-0">诊所挂号系统</h1>
        <Button type="link" icon={<LogoutOutlined />} onClick={handleLogout}>
          退出登录
        </Button>
      </Header>
      <Layout>
        <Sider className="bg-white" width={200}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Content className="p-6">{children}</Content>
      </Layout>
    </Layout>
  )
}
