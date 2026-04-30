import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Input, Button, Space, Card, message } from 'antd'
import { patientApi, type Patient } from '../api/client'

export default function PatientListPage() {
  const [data, setData] = useState<Patient[]>([])
  const [total, setTotal] = useState(0)
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const load = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const res = await patientApi.list(keyword, page, pageSize)
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

  const handleDelete = async (id: number) => {
    try {
      await patientApi.remove(id)
      message.success('删除成功')
      load()
    } catch (e: any) {
      message.error(e.message || '删除失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 80 },
    { title: '姓名', dataIndex: 'name' },
    { title: '身份证号', dataIndex: 'id_card' },
    { title: '手机号', dataIndex: 'phone' },
    { title: '性别', dataIndex: 'gender', render: (v: string) => (v === 'male' ? '男' : v === 'female' ? '女' : '未知') },
    { title: '年龄', dataIndex: 'age' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Patient) => (
        <Space>
          <Button type="link" onClick={() => navigate(`/patients/${record.id}`)}>详情</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="就诊人管理">
      <Space className="mb-4">
        <Input.Search
          placeholder="搜索姓名/手机号/身份证号"
          allowClear
          onSearch={(v) => { setKeyword(v); load(1, 10) }}
        />
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
