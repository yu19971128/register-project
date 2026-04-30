# 数据流：Orchestrator 自循环主流程

```mermaid
flowchart TD
    User[用户输入需求] --> Init[阶段0: Orchestrator 启动]
    Init --> Gate1{Gatekeeper<br/>文档完整性检查}
    Gate1 -- 缺失 --> Halt[暂停并引导补全]
    Gate1 -- 完整 --> S1[阶段1: 需求澄清]
    S1 --> S2[阶段2: 架构与模块拆分]
    S2 --> S3[阶段3: 接口与数据库设计]
    S3 --> S4[阶段4: 前端设计<br/>条件触发]
    S4 --> S5[阶段5: 审查与计划生成]
    S5 --> S6[阶段6: TDD 实现]
    S6 --> S7[阶段7: 质量检查与测试自检]
    S7 --> S8[阶段8: 发布与经验沉淀]
    S8 --> Done[更新 ROADMAP & Task Card]
    S1 -.->|偏离检测| Gate1
    S2 -.->|偏离检测| Gate1
    S6 -.->|偏离检测| Gate1
    style Gate1 fill:#f9f,stroke:#333,stroke-width:2px
    style Halt fill:#faa,stroke:#333,stroke-width:2px
    style Done fill:#afa,stroke:#333,stroke-width:2px
```

## 说明

1. **阶段0**：Orchestrator 读取 `AGENTS.md` + `docs/00-index.md`，确认文档链完整
2. **阶段1-5**：设计阶段，产出并更新 GSD 文档树
3. **阶段6**：实现阶段，按依赖图并行派发 subagent，Red-Green 配对执行
4. **阶段7**：并行运行 `health` 和 `qa`，门禁不通过则自动修复（有上限）
5. **阶段8**：调用 `ship` + `learn`，闭环结束
6. **偏离检测**：任何阶段若发现实现与文档约束不符，回退到 Gatekeeper 进行校正
