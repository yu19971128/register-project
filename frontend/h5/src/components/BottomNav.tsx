import { useLocation, useNavigate } from 'react-router-dom'
import { TabBar } from 'antd-mobile'
import { CalendarOutline, UserOutline } from 'antd-mobile-icons'

export default function BottomNav({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate()
  const location = useLocation()

  const tabs = [
    { key: '/h5/register', title: '当日挂号', icon: <CalendarOutline /> },
    { key: '/h5/me', title: '个人中心', icon: <UserOutline /> },
  ]

  const activeKey = tabs.some(t => location.pathname === t.key)
    ? location.pathname
    : '/h5/register'

  return (
    <div className="flex flex-col min-h-screen">
      <div className="flex-1 pb-16">{children}</div>
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t z-50">
        <TabBar activeKey={activeKey} onChange={(key) => navigate(key)}>
          {tabs.map((item) => (
            <TabBar.Item key={item.key} icon={item.icon} title={item.title} />
          ))}
        </TabBar>
      </div>
    </div>
  )
}
