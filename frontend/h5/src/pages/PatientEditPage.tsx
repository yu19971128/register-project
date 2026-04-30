import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { NavBar, Form, Input, Button, Toast, Radio } from 'antd-mobile'
import { patientApi } from '../api/client'

export default function PatientEditPage() {
  const { id } = useParams<{ id?: string }>()
  const navigate = useNavigate()
  const isEdit = Boolean(id)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (isEdit) {
      patientApi.get(Number(id)).then((p) => {
        form.setFieldsValue(p)
      }).catch((e: any) => {
        Toast.show({ content: e.message || '加载失败', position: 'bottom' })
      })
    }
  }, [id, isEdit, form])

  const onSubmit = async (values: any) => {
    setLoading(true)
    try {
      const payload = { ...values }
      if (payload.age !== undefined && payload.age !== '') {
        payload.age = Number(payload.age)
      }
      if (isEdit) {
        await patientApi.update(Number(id), payload)
        Toast.show({ content: '更新成功', position: 'bottom' })
      } else {
        await patientApi.create(payload)
        Toast.show({ content: '创建成功', position: 'bottom' })
      }
      navigate('/h5/patients')
    } catch (e: any) {
      Toast.show({ content: e.message || '保存失败', position: 'bottom' })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar onBack={() => navigate(-1)}>{isEdit ? '编辑就诊人' : '添加就诊人'}</NavBar>
      <div className="p-4">
        <Form
          form={form}
          onFinish={onSubmit}
          layout="horizontal"
          footer={
            <Button block color="primary" type="submit" loading={loading}>
              保存
            </Button>
          }
        >
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
            <Input placeholder="请输入姓名" />
          </Form.Item>
          <Form.Item name="id_card" label="身份证号" rules={[{ required: true, message: '请输入身份证号' }]}>
            <Input placeholder="请输入身份证号" />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ required: true, message: '请输入手机号' }]}>
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="gender" label="性别">
            <Radio.Group>
              <Radio value="male">男</Radio>
              <Radio value="female">女</Radio>
            </Radio.Group>
          </Form.Item>
          <Form.Item name="age" label="年龄">
            <Input type="number" placeholder="请输入年龄" />
          </Form.Item>
          <Form.Item name="address" label="住址">
            <Input placeholder="请输入住址" />
          </Form.Item>
        </Form>
      </div>
    </div>
  )
}
