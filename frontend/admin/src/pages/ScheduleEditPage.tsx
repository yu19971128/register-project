import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Card, Form, Select, DatePicker, TimePicker, Button, InputNumber, message, Space } from 'antd'
import { scheduleApi } from '../api/client'
import { DEPARTMENT_DOCTORS, DEPARTMENT_OPTIONS } from '../constants/departments'
import dayjs from 'dayjs'

interface FormValues {
  date: dayjs.Dayjs
  department: string
  doctor_name: string
  timeRange: [dayjs.Dayjs, dayjs.Dayjs]
  total_quota: number
}

export default function ScheduleEditPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const id = searchParams.get('id')
  const [form] = Form.useForm<FormValues>()
  const [loading, setLoading] = useState(false)
  const isEdit = !!id
  const department = Form.useWatch('department', form)
  const doctorOptions = (DEPARTMENT_DOCTORS[department] || []).map(d => ({ label: d, value: d }))

  useEffect(() => {
    if (!id) return
    setLoading(true)
    scheduleApi.get(Number(id))
      .then((s) => {
        form.setFieldsValue({
          date: dayjs(s.date),
          department: s.department,
          doctor_name: s.doctor_name,
          timeRange: [dayjs(s.start_time, 'HH:mm'), dayjs(s.end_time, 'HH:mm')],
          total_quota: s.total_quota,
        })
      })
      .catch((e: any) => message.error(e.message || '加载失败'))
      .finally(() => setLoading(false))
  }, [id])

  const onFinish = async (values: FormValues) => {
    setLoading(true)
    try {
      const payload = {
        date: values.date.format('YYYY-MM-DD'),
        department: values.department,
        doctor_name: values.doctor_name,
        start_time: values.timeRange[0].format('HH:mm'),
        end_time: values.timeRange[1].format('HH:mm'),
        total_quota: values.total_quota,
      }
      if (isEdit) {
        await scheduleApi.update(Number(id), payload)
        message.success('更新成功')
      } else {
        await scheduleApi.create(payload)
        message.success('创建成功')
      }
      navigate('/schedules')
    } catch (e: any) {
      message.error(e.message || '保存失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title={isEdit ? '编辑号源' : '添加号源'} loading={loading}>
      <Form form={form} layout="vertical" onFinish={onFinish} style={{ maxWidth: 480 }}>
        <Form.Item name="date" label="出诊日期" rules={[{ required: true, message: '请选择出诊日期' }]}>
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="department" label="科室" rules={[{ required: true, message: '请选择科室' }]}>
          <Select
            placeholder="请选择科室"
            options={DEPARTMENT_OPTIONS}
            onChange={() => form.setFieldValue('doctor_name', undefined)}
          />
        </Form.Item>
        <Form.Item name="doctor_name" label="医生" rules={[{ required: true, message: '请选择医生' }]}>
          <Select
            placeholder={department ? '请选择医生' : '请先选择科室'}
            options={doctorOptions}
            disabled={!department}
          />
        </Form.Item>
        <Form.Item name="timeRange" label="时间段" rules={[{ required: true, message: '请选择时间段' }]}>
          <TimePicker.RangePicker format="HH:mm" style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item name="total_quota" label="总号数" rules={[{ required: true, message: '请输入总号数' }]}>
          <InputNumber min={1} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
            <Button onClick={() => navigate('/schedules')}>取消</Button>
          </Space>
        </Form.Item>
      </Form>
    </Card>
  )
}
