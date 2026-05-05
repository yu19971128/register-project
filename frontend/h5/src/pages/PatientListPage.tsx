import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { List, NavBar, Button, Dialog, Toast, Empty } from 'antd-mobile'
import { patientApi, type Patient } from '../api/client'

export default function PatientListPage() {
  const [patients, setPatients] = useState<Patient[]>([])
  const navigate = useNavigate()

  const load = async () => {
    try {
      const res = await patientApi.list()
      setPatients(res.list || [])
    } catch (e: any) {
      Toast.show({ content: e.message || '加载失败', position: 'bottom' })
    }
  }

  useEffect(() => {
    load()
  }, [])

  const handleDelete = async (p: Patient) => {
    try {
      await patientApi.remove(p.id)
      Toast.show({ content: '删除成功', position: 'bottom' })
      load()
    } catch (e: any) {
      console.log('err')
      Toast.show({ content: e?.message || '删除失败', position: 'bottom' })
    }
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar back={null}>就诊人管理</NavBar>
   
      {!patients?.length ? (
        <div className="flex flex-col items-center justify-center pt-20">
          <Empty description="暂无就诊人" />
          {/* <Button color="primary" className="mt-6 w-48" onClick={() => navigate('/h5/patients/edit')}>
            添加就诊人
          </Button> */}
        </div>
      ) : (
        <List>
          {patients.map((p) => (
            <List.Item
              key={p.id}
              title={
                <span onClick={() => navigate(`/h5/patients/edit/${p.id}`)}>
                  {p.name}
                </span>
              }
              description={`${p.phone} · ${p.gender === 'male' ? '男' : p.gender === 'female' ? '女' : '未知'} · ${p.age ?? '-'}岁`}
              extra={
                <Button size="mini" color="danger" onClick={() => handleDelete(p)}>
                  删除
                </Button>
              }
            />
          ))}
        </List>
      )}
      <div className="p-4">
        <Button color="primary" block onClick={() => navigate('/h5/patients/edit')}>
          添加就诊人
        </Button>
      </div>
    </div>
  )
}
