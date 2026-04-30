import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { NavBar, Card, Button, Dialog, Toast, Tag } from 'antd-mobile'
import { orderApi, type OrderDetail } from '../api/client'

const statusMap: Record<string, { text: string; color: string }> = {
  confirmed: { text: '待就诊', color: 'primary' },
  cancelled: { text: '已退号', color: 'default' },
  completed: { text: '已完成', color: 'success' },
}

export default function OrderDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<OrderDetail | null>(null)

  const load = async () => {
    try {
      const res = await orderApi.get(Number(id))
      setOrder(res)
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', position: 'bottom' })
    }
  }

  useEffect(() => {
    load()
  }, [id])

  const handleCancel = () => {
    if (!order) return
    Dialog.confirm({
      content: '确定要退号吗？退号后号源将释放给其他患者。',
      onConfirm: async () => {
        try {
          await orderApi.cancel(order.id, '个人原因')
          Toast.show({ content: '退号成功', position: 'bottom' })
          load()
        } catch (e: any) {
          Toast.show({ content: e.message || '退号失败', position: 'bottom' })
        }
      },
    })
  }

  if (!order) return null

  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar onBack={() => navigate(-1)}>订单详情</NavBar>
      <div className="p-4 space-y-4">
        <Card title="就诊信息">
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-500">订单号</span>
              <span>{order.order_no}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">状态</span>
              <Tag color={statusMap[order.status]?.color || 'default'}>
                {statusMap[order.status]?.text || order.status}
              </Tag>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">就诊人</span>
              <span>{order.patient.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">科室</span>
              <span>{order.schedule.department}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">医生</span>
              <span>{order.schedule.doctor_name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">时间</span>
              <span>{order.schedule.date} {order.schedule.start_time}-{order.schedule.end_time}</span>
            </div>
          </div>
        </Card>

        {order.status === 'confirmed' && (
          <Button color="danger" block onClick={handleCancel}>
            申请退号
          </Button>
        )}
      </div>
    </div>
  )
}
