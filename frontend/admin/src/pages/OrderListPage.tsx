import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Input, Select, Button, Space, Card, Tag, message } from 'antd'
import { orderApi, type Order } from '../api/client'

const statusMap: Record<string, { text: string; color: string }> = {
  confirmed: { text: '待就诊', color: 'blue' },
  cancelled: { text: '已退号', color: 'default' },
  completed: { text: '已完成', color: 'green' },
}

export default function OrderListPage() {
  const [data, setData] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [keyword, setKeyword] = useState('')
  const [status, setStatus] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const load = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await orderApi.list({ keyword, status, page, pageSize })
      setData(res.list)
      setTotal(res.total)
    } catch (e: any) {
      message.error(e.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const handleCancel = async (id: number) => {
    try {
      await orderApi.cancel(id, '管理员退号')
      message.success('退号成功')
      load()
    } catch (e: any) {
      message.error(e.message || '退号失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '订单号', dataIndex: 'order_no' },
    { title: '就诊人', dataIndex: 'patient_name' },
    { title: '科室', dataIndex: 'department' },
    { title: '医生', dataIndex: 'doctor_name' },
    {
      title: '就诊时间',
      render: (_: any, record: Order) => `${record.date} ${record.start_time}-${record.end_time}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      render: (v: string) => (
        <Tag color={statusMap[v]?.color || 'default'}>{statusMap[v]?.text || v}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Space>
          <Button type="link" onClick={() => navigate(`/orders/${record.id}`)}>详情</Button>
          {record.status === 'confirmed' && (
            <Button type="link" danger onClick={() => handleCancel(record.id)}>退号</Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <Card title="挂号订单管理">
      <Space className="mb-4">
        <Input.Search
          placeholder="搜索订单号/就诊人"
          allowClear
          onSearch={(v) => { setKeyword(v); load(1, 10) }}
        />
        <Select
          placeholder="状态筛选"
          allowClear
          style={{ width: 120 }}
          onChange={(v) => { setStatus(v || ''); load(1, 10) }}
        >
          <Select.Option value="confirmed">待就诊</Select.Option>
          <Select.Option value="cancelled">已退号</Select.Option>
          <Select.Option value="completed">已完成</Select.Option>
        </Select>
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={{ total, showTotal: (t) => `共 ${t} 条`, onChange: load }}
      />
    </Card>
  )
}
