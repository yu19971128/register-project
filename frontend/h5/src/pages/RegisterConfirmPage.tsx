import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { NavBar, Card, Button, Radio, Toast } from 'antd-mobile'
import { patientApi, scheduleApi, registrationApi, type Patient, type Schedule } from '../api/client'

export default function RegisterConfirmPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const scheduleId = Number(searchParams.get('schedule_id'))
  const [schedule, setSchedule] = useState<Schedule | null>(null)
  const [patients, setPatients] = useState<Patient[]>([])
  const [selectedPatientId, setSelectedPatientId] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [scheduleId])

  const loadData = async () => {
    try {
      const [s, p] = await Promise.all([
        scheduleApi.get(scheduleId),
        patientApi.list(),
      ])
      setSchedule(s)
      setPatients(p.list)
      if (p.list.length > 0) {
        setSelectedPatientId(p.list[0].id)
      }
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', icon: 'fail' })
    }
  }

  const handleSubmit = async () => {
    if (!selectedPatientId) {
      Toast.show({ content: '请选择就诊人', icon: 'fail' })
      return
    }
    setLoading(true)
    try {
      const phone = localStorage.getItem('visitor_phone') || ''
      const result = await registrationApi.submit(scheduleId, selectedPatientId, phone)
      Toast.show({ content: '挂号成功', icon: 'success' })
      navigate(`/h5/register/ticket?order_no=${result.order_no}`)
    } catch (e: any) {
      Toast.show({ content: e.message || '挂号失败', icon: 'fail' })
    } finally {
      setLoading(false)
    }
  }

  if (!schedule) return <div className="p-4 text-center">加载中...</div>

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <NavBar onBack={() => navigate(-1)}>确认挂号</NavBar>
      <Card className="m-4">
        <div className="text-sm text-gray-500 mb-1">号源信息</div>
        <div className="font-medium text-lg">{schedule.department}</div>
        <div className="text-sm text-gray-600">{schedule.doctor_name}</div>
        <div className="text-sm text-gray-600">{schedule.date} {schedule.start_time}-{schedule.end_time}</div>
        <div className="text-sm text-green-600 mt-1">挂号费：免费</div>
      </Card>

      <div className="px-4 mt-4">
        <div className="text-sm font-medium mb-2">选择就诊人 *</div>
        <Radio.Group value={selectedPatientId} onChange={(v) => setSelectedPatientId(v as number)}>
          <div className="space-y-3">
            {patients.map(p => (
              <Card key={p.id} className={selectedPatientId === p.id ? 'border-blue-500 border-l-4' : ''}>
                <Radio value={p.id} className="w-full">
                  <div className="flex items-center gap-2 ml-2">
                    <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 text-sm">
                      {p.name.charAt(0)}
                    </div>
                    <div>
                      <div className="font-medium">{p.name} {p.gender === 'male' ? '男' : p.gender === 'female' ? '女' : ''} {p.age}岁</div>
                      <div className="text-xs text-gray-400">{p.phone}</div>
                    </div>
                  </div>
                </Radio>
              </Card>
            ))}
          </div>
        </Radio.Group>
        <div className="mt-4 text-center">
          <Button fill="none" onClick={() => navigate('/h5/patients/edit')}>
            + 添加新就诊人
          </Button>
        </div>
      </div>

      <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t">
        <Button block color="primary" size="large" loading={loading} onClick={handleSubmit}>
          确认挂号
        </Button>
      </div>
    </div>
  )
}
