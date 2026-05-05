import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { NavBar, Card, Button, Toast } from 'antd-mobile'
import { registrationApi, type TicketResult } from '../api/client'

export default function RegisterTicketPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const orderNo = searchParams.get('order_no') || ''
  const [ticket, setTicket] = useState<TicketResult | null>(null)

  useEffect(() => {
    if (!orderNo) return
    loadTicket()
  }, [orderNo])

  const loadTicket = async () => {
    try {
      const t = await registrationApi.getTicket(orderNo)
      setTicket(t)
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', icon: 'fail' })
    }
  }

  if (!ticket) {
    return (
      <div className="min-h-screen bg-gray-50">
        <NavBar onBack={() => navigate(-1)}>挂号成功</NavBar>
        <div className="p-4 text-center">加载中...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <NavBar onBack={() => navigate(-1)} right={<Button fill="none" size="mini" onClick={() => navigate('/h5/patients')}>完成</Button>}>挂号成功</NavBar>

      <div className="text-center py-6">
        <div className="text-green-600 text-lg font-medium mb-2">✅ 挂号成功！</div>
      </div>

      <Card className="mx-4">
        <div className="flex flex-col items-center py-4">
          <div className="w-32 h-32 bg-gray-900 flex items-center justify-center text-white text-xs mb-3">
            {/* QRCode placeholder */}
            <div className="text-center leading-tight">
              {Array.from({ length: 8 }).map((_, i) => (
                <div key={i} className="flex justify-center">
                  {Array.from({ length: 8 }).map((_, j) => (
                    <span key={j} className={`inline-block w-3 h-3 ${Math.random() > 0.5 ? 'bg-white' : 'bg-transparent'}`} />
                  ))}
                </div>
              ))}
            </div>
          </div>
          <div className="text-sm text-gray-500">订单号：{ticket.order_no}</div>
        </div>
      </Card>

      <Card className="mx-4 mt-4" title="就诊信息">
        <div className="space-y-2 text-sm">
          <div className="flex justify-between"><span className="text-gray-500">姓名</span><span>{ticket.patient_name}</span></div>
          <div className="flex justify-between"><span className="text-gray-500">科室</span><span>{ticket.department}</span></div>
          <div className="flex justify-between"><span className="text-gray-500">医生</span><span>{ticket.doctor_name}</span></div>
          <div className="flex justify-between"><span className="text-gray-500">时间</span><span>{ticket.date} {ticket.start_time}-{ticket.end_time}</span></div>
          <div className="flex justify-between"><span className="text-gray-500">地点</span><span>{ticket.location}</span></div>
        </div>
      </Card>

      <Card className="mx-4 mt-4" title="温馨提示">
        <div className="space-y-1 text-sm text-gray-600">
          {ticket.notice.map((n, i) => (
            <div key={i}>{i + 1}. {n}</div>
          ))}
        </div>
      </Card>

      <div className="fixed bottom-0 left-0 right-0 p-4 bg-white border-t">
        <Button block color="primary" size="large" onClick={() => navigate('/h5/orders')}>
          查看挂号记录
        </Button>
      </div>
    </div>
  )
}
