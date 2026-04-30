import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Card, Descriptions, Button, Tag, message } from 'antd'
import { orderApi, type OrderDetail } from '../api/client'

const statusMap: Record<string, { text: string; color: string }> = {
  confirmed: { text: '待就诊', color: 'blue' },
  cancelled: { text: '已退号', color: 'default' },
  completed: { text: '已完成', color: 'green' },
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
      message.error(e.message || '加载失败')
    }
  }

  useEffect(() => {
    load()
  }, [id])

  const handleCancel = async () => {
    if (!order) return
    try {
      await orderApi.cancel(order.id, '管理员退号')
      message.success('退号成功')
      load()
    } catch (e: any) {
      message.error(e.message || '退号失败')
    }
  }

  if (!order) return null

  return (
    <Card
      title="订单详情"
      extra={
        <Button onClick={() => navigate('/orders')}>返回列表</Button>
      }
    >
      <Descriptions bordered column={2}>
        <Descriptions.Item label="订单号">{order.order_no}</Descriptions.Item>
        <Descriptions.Item label="状态">
          <Tag color={statusMap[order.status]?.color || 'default'}>
            {statusMap[order.status]?.text || order.status}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="就诊人">{order.patient.name}</Descriptions.Item>
        <Descriptions.Item label="性别">
          {order.patient.gender === 'male' ? '男' : order.patient.gender === 'female' ? '女' : '未知'}
        </Descriptions.Item>
        <Descriptions.Item label="年龄">{order.patient.age}</Descriptions.Item>
        <Descriptions.Item label="访客手机号">{order.visitor_phone}</Descriptions.Item>
        <Descriptions.Item label="科室">{order.schedule.department}</Descriptions.Item>
        <Descriptions.Item label="医生">{order.schedule.doctor_name}</Descriptions.Item>
        <Descriptions.Item label="就诊日期">{order.schedule.date}</Descriptions.Item>
        <Descriptions.Item label="就诊时间">
          {order.schedule.start_time} - {order.schedule.end_time}
        </Descriptions.Item>
        <Descriptions.Item label="创建时间">{order.created_at}</Descriptions.Item>
        <Descriptions.Item label="更新时间">{order.updated_at}</Descriptions.Item>
      </Descriptions>

      {order.status === 'confirmed' && (
        <div className="mt-4">
          <Button danger onClick={handleCancel}>退号</Button>
        </div>
      )}
    </Card>
  )
}
