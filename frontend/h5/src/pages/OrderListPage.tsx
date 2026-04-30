import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { List, NavBar, Tag, Toast } from 'antd-mobile'
import { orderApi, type Order } from '../api/client'

const statusMap: Record<string, { text: string; color: string }> = {
  confirmed: { text: '待就诊', color: 'primary' },
  cancelled: { text: '已退号', color: 'default' },
  completed: { text: '已完成', color: 'success' },
}

export default function OrderListPage() {
  const [orders, setOrders] = useState<Order[]>([])
  const navigate = useNavigate()

  const load = async () => {
    try {
      const res = await orderApi.list()
      setOrders(res.list)
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', position: 'bottom' })
    }
  }

  useEffect(() => {
    load()
  }, [])

  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar back={null}>挂号记录</NavBar>
      <List>
        {orders.map((o) => (
          <List.Item
            key={o.id}
            title={
              <div className="flex items-center gap-2">
                <span>{o.patient_name}</span>
                <Tag color={statusMap[o.status]?.color || 'default'}>
                  {statusMap[o.status]?.text || o.status}
                </Tag>
              </div>
            }
            description={`${o.department} · ${o.doctor_name} · ${o.date} ${o.start_time}-${o.end_time}`}
            onClick={() => navigate(`/h5/orders/${o.id}`)}
          />
        ))}
      </List>
      {orders.length === 0 && (
        <div className="p-8 text-center text-gray-400">暂无挂号记录</div>
      )}
    </div>
  )
}
