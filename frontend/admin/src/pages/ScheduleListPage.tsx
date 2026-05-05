import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Card, DatePicker, Tag, message, Select } from 'antd'
import { scheduleApi, type Schedule } from '../api/client'
import { DEPARTMENT_DOCTORS, DEPARTMENT_OPTIONS } from '../constants/departments'
import dayjs from 'dayjs'

export default function ScheduleListPage() {
  const [data, setData] = useState<Schedule[]>([])
  const [total, setTotal] = useState(0)
  const [date, setDate] = useState(dayjs().format('YYYY-MM-DD'))
  const [department, setDepartment] = useState('')
  const [doctorName, setDoctorName] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const doctorOptions = (department ? DEPARTMENT_DOCTORS[department] || [] : []).map(d => ({ label: d, value: d }))

  const load = async (page = 1, pageSize = 10, currentDept = department, currentDoc = doctorName) => {
    setLoading(true)
    try {
      const res = await scheduleApi.list({ date, department: currentDept || undefined, doctor_name: currentDoc || undefined, page, pageSize })
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
      load(1, 10, department, doctorName)
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
        <Select
          placeholder="全部科室"
          allowClear
          style={{ width: 140 }}
          value={department || undefined}
          onChange={(v) => {
            const d = v || ''
            setDepartment(d)
            setDoctorName('')
            load(1, 10, d, '')
          }}
        >
          {DEPARTMENT_OPTIONS.map(opt => (
            <Select.Option key={opt.value} value={opt.value}>{opt.label}</Select.Option>
          ))}
        </Select>
        <Select
          placeholder={department ? '全部医生' : '请先选择科室'}
          allowClear
          style={{ width: 140 }}
          disabled={!department}
          value={doctorName || undefined}
          onChange={(v) => {
            const doc = v || ''
            setDoctorName(doc)
            load(1, 10, department, doc)
          }}
        >
          {doctorOptions.map(opt => (
            <Select.Option key={opt.value} value={opt.value}>{opt.label}</Select.Option>
          ))}
        </Select>
        <Button onClick={() => { setDepartment(''); setDoctorName(''); load(1, 10, '', '') }}>重置</Button>
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
