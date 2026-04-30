import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Card, Descriptions, Button, message } from 'antd'
import { patientApi, type Patient } from '../api/client'

export default function PatientDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [patient, setPatient] = useState<Patient | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    patientApi.get(Number(id))
      .then((p) => setPatient(p))
      .catch((e: any) => message.error(e.message || '加载失败'))
      .finally(() => setLoading(false))
  }, [id])

  if (!patient) {
    return <Card loading={loading}>加载中...</Card>
  }

  return (
    <Card
      title="就诊人详情"
      extra={<Button onClick={() => navigate('/patients')}>返回</Button>}
    >
      <Descriptions bordered column={2}>
        <Descriptions.Item label="ID">{patient.id}</Descriptions.Item>
        <Descriptions.Item label="姓名">{patient.name}</Descriptions.Item>
        <Descriptions.Item label="身份证号">{patient.id_card}</Descriptions.Item>
        <Descriptions.Item label="手机号">{patient.phone}</Descriptions.Item>
        <Descriptions.Item label="性别">
          {patient.gender === 'male' ? '男' : patient.gender === 'female' ? '女' : '未知'}
        </Descriptions.Item>
        <Descriptions.Item label="年龄">{patient.age ?? '-'}</Descriptions.Item>
        <Descriptions.Item label="住址">{patient.address || '-'}</Descriptions.Item>
        <Descriptions.Item label="创建时间">{patient.created_at || '-'}</Descriptions.Item>
        <Descriptions.Item label="更新时间">{patient.updated_at || '-'}</Descriptions.Item>
      </Descriptions>
    </Card>
  )
}
