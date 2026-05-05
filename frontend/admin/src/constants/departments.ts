export const DEPARTMENT_DOCTORS: Record<string, string[]> = {
  '泌尿外科': ['张医生', '李医生'],
  '内科': ['王医生', '刘医生'],
  '外科': ['陈医生'],
  '儿科': ['赵医生'],
  '妇科': ['孙医生'],
  '眼科': ['周医生'],
  '口腔科': ['吴医生'],
  '皮肤科': ['郑医生'],
}

export const DEPARTMENT_OPTIONS = Object.keys(DEPARTMENT_DOCTORS).map(d => ({ label: d, value: d }))
