import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { NavBar, Tabs, Card, Button, Badge, Toast, Modal, Input } from 'antd-mobile'
import { scheduleApi, type Schedule } from '../api/client'
import dayjs from 'dayjs'

export default function RegisterPage() {
  const navigate = useNavigate()
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [departments, setDepartments] = useState<string[]>(['全部'])
  const [activeDept, setActiveDept] = useState('全部')
  const [visitorPhone, setVisitorPhone] = useState('')
  const [phoneModalVisible, setPhoneModalVisible] = useState(false)
  const today = dayjs().format('YYYY-MM-DD')

  useEffect(() => {
    const vp = localStorage.getItem('visitor_phone') || ''
    if (!vp) {
      setPhoneModalVisible(true)
    } else {
      setVisitorPhone(vp)
      loadSchedules()
    }
  }, [])

  const loadSchedules = async () => {
    try {
      const res = await scheduleApi.list(today)
      setSchedules(res.list || [])
      const depts = Array.from(new Set((res.list || []).map(s => s.department)))
      setDepartments(['全部', ...depts])
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', icon: 'fail' })
    }
  }

  const handleSavePhone = () => {
    if (!visitorPhone || !/^1[3-9]\d{9}$/.test(visitorPhone)) {
      Toast.show({ content: '请输入正确的手机号', icon: 'fail' })
      return
    }
    localStorage.setItem('visitor_phone', visitorPhone)
    setPhoneModalVisible(false)
    loadSchedules()
  }

  const filtered = activeDept === '全部' ? schedules : schedules.filter(s => s.department === activeDept)

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar back={null}>当天挂号</NavBar>
      <div className="px-4 py-2 text-sm text-gray-500">📅 {today} (今天)</div>
      <Tabs activeKey={activeDept} onChange={setActiveDept} className="bg-white">
        {departments.map(d => (
          <Tabs.Tab title={d} key={d} />
        ))}
      </Tabs>
      <div className="p-4 space-y-3">
        {filtered.map(s => (
          <Card key={s.id} className={s.remaining === 0 ? 'bg-gray-100 opacity-60' : 'bg-white border-l-4 border-blue-500'}>
            <div className="flex justify-between items-center">
              <div>
                <div className="font-medium text-base">{s.department}</div>
                <div className="text-sm text-gray-500">{s.doctor_name} · {s.start_time}-{s.end_time}</div>
                <div className="text-sm mt-1">余号 {s.remaining}/{s.total_quota}</div>
              </div>
              {s.remaining > 0 ? (
                <Button color="primary" size="small" onClick={() => navigate(`/h5/register/confirm?schedule_id=${s.id}`)}>
                  立即挂号
                </Button>
              ) : (
                <Badge content="已满" style={{ '--color': '#999' }} />
              )}
            </div>
          </Card>
        ))}
        {filtered.length === 0 && (
          <div className="text-center text-gray-400 py-10">暂无号源</div>
        )}
      </div>

      <Modal
        visible={phoneModalVisible}
        title="请输入您的手机号"
        content={
          <Input
            placeholder="请输入手机号"
            value={visitorPhone}
            onChange={(v) => setVisitorPhone(v)}
            maxLength={11}
          />
        }
        closeOnAction
        onClose={() => {}}
        actions={[
          {
            key: 'confirm',
            text: '确定',
            primary: true,
            onClick: handleSavePhone,
          },
        ]}
      />
    </div>
  )
}
