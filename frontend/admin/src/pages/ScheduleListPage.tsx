import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Card, DatePicker, Tag, message } from 'antd'
import { scheduleApi, type Schedule } from '../api/client'
import dayjs from 'dayjs'

export default function ScheduleListPage() {
  const [data, setData] = useState<Schedule[]>([])
  const [total, setTotal] = useState(0)
  const [date, setDate] = useState(dayjs().format('YYYY-MM-DD'))
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const load = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await scheduleApi.list(date, page, pageSize)
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
  }, [date])

  const handleDelete = async (id: number) => {
    try {
      await scheduleApi.remove(id)
      message.success('删除成功')
      load()
    } catch (e: any) {
      message.error(e.message || '删除失败')
    }
  }

  const statusTag = (status: string) => {
    const map: Record<string, { color: string; text: string }> = {
      available: { color: 'success', text: '可约' },
      full: { color: 'error', text: '已满' },
      stopped: { color: 'default', text: '已停诊' },
    }
    const s = map[status] || { color: 'default', text: status }
    return <Tag color={s.color}>{s.text}</Tag>
  }

  const columns = [
    { title: '科室', dataIndex: 'department' },
    { title: '医生', dataIndex: 'doctor_name' },
    { title: '时间段', render: (_: any, r: Schedule) => `${r.start_time} - ${r.end_time}` },
    { title: '总号数', dataIndex: 'total_quota', width: 90 },
    { title: '余量', dataIndex: 'remaining', width: 80 },
    { title: '状态', dataIndex: 'status', render: (v: string) => statusTag(v) },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Schedule) => (
        <Space>
          <Button type="link" onClick={() => navigate(`/schedules/edit?id=${record.id}`)}>编辑</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="号源管理">
      <Space className="mb-4" wrap>
        <DatePicker value={dayjs(date)} onChange={(d) => setDate(d ? d.format('YYYY-MM-DD') : dayjs().format('YYYY-MM-DD'))} />
        <Button type="primary" onClick={() => navigate('/schedules/edit')}>+ 新增号源</Button>
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
