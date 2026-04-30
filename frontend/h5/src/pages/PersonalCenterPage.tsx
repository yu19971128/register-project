import { useNavigate } from 'react-router-dom'
import { NavBar, List, Button } from 'antd-mobile'
import { UserOutline, FileOutline, RightOutline } from 'antd-mobile-icons'

export default function PersonalCenterPage() {
  const navigate = useNavigate()
  const visitorPhone = localStorage.getItem('visitor_phone') || '-'

  const handleLogout = () => {
    localStorage.removeItem('visitor_phone')
    window.location.reload()
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar back={null}>个人中心</NavBar>

      <div className="bg-blue-500 text-white p-6">
        <div className="text-2xl font-bold">您好</div>
        <div className="text-sm mt-1 opacity-80">手机号：{visitorPhone}</div>
      </div>

      <List className="mt-4">
        <List.Item
          prefix={<UserOutline />}
          onClick={() => navigate('/h5/patients')}
          arrow={<RightOutline />}
        >
          就诊人管理
        </List.Item>
        <List.Item
          prefix={<FileOutline />}
          onClick={() => navigate('/h5/orders')}
          arrow={<RightOutline />}
        >
          挂号记录
        </List.Item>
      </List>

      <div className="p-4 mt-8">
        <Button block color="danger" fill="outline" onClick={handleLogout}>
          切换手机号
        </Button>
      </div>
    </div>
  )
}
