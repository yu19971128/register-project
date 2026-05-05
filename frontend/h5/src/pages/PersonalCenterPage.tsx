import { useNavigate } from 'react-router-dom'
import { NavBar, List, Button, Dialog, Input, Toast } from 'antd-mobile'
import { UserOutline, FileOutline, RightOutline } from 'antd-mobile-icons'

export default function PersonalCenterPage() {
  const navigate = useNavigate()
  const visitorPhone = localStorage.getItem('visitor_phone') || '-'

  const handleSwitchPhone = () => {
    let phone = ''
    const dialog = Dialog.show({
      title: '切换手机号',
      content: (
        <Input
          placeholder="请输入手机号"
          onChange={(v) => { phone = v }}
          maxLength={11}
        />
      ),
      closeOnAction: false,
      actions: [
        [
          { key: 'cancel', text: '取消', onClick: () => dialog.close() },
          {
            key: 'confirm',
            text: '确认',
            primary: true,
            onClick: () => {
              const trimmed = phone.trim()
              if (!trimmed || !/^1[3-9]\d{9}$/.test(trimmed)) {
                Toast.show({ content: '请输入正确的手机号', icon: 'fail' })
                return
              }
              localStorage.setItem('visitor_phone', trimmed)
              dialog.close()
              window.location.reload()
            },
          },
        ],
      ],
    })
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
        <Button block color="danger" fill="outline" onClick={handleSwitchPhone}>
          切换手机号
        </Button>
      </div>
    </div>
  )
}
