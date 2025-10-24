import Link from 'next/link'

import { Badge } from '@/components/ui/badge'
import { buttonVariants } from '@/components/ui/button-variants'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

export default function Page() {
  return (
    <div className="grid gap-8">
      <section className="grid items-stretch gap-4 md:grid-cols-2">
        <Card className="flex h-full flex-col">
          <CardHeader className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div className="space-y-1">
              <CardTitle>Plane ↔ CNB 集成配置</CardTitle>
              <CardDescription>建立仓库与项目之间的映射关系，配置同步方向与标签过滤策略。</CardDescription>
            </div>
            <Badge variant="info">Mappings</Badge>
          </CardHeader>
          <CardContent className="flex-1 space-y-3 text-sm text-muted-foreground">
            <p>确保 Plane 项目准确接收来自 CNB 的 Issue / PR 信号，并反向回写状态。</p>
            <ul className="space-y-1">
              <li>• 维护仓库 ↔ 项目与状态映射，避免错误同步。</li>
              <li>• 配置标签选择器与同步方向，控制流量与优先级。</li>
              <li>• 支持快速刷新现有映射与按项目过滤视图。</li>
            </ul>
          </CardContent>
          <CardFooter className="flex flex-col gap-3 text-sm text-muted-foreground md:flex-row md:items-center md:justify-between">
            <span>工作区、项目、同步策略均可在一处集中管理。</span>
            <Link href="/mappings" className={buttonVariants({ variant: 'primary' })}>
              打开映射面板
            </Link>
          </CardFooter>
        </Card>

        <Card className="flex h-full flex-col">
          <CardHeader className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
            <div className="space-y-1">
              <CardTitle>用户身份映射管理</CardTitle>
              <CardDescription>维护 Plane 成员与 CNB 账号的绑定，确保活动准确归属。</CardDescription>
            </div>
            <Badge variant="success">Users</Badge>
          </CardHeader>
          <CardContent className="flex-1 space-y-3 text-sm text-muted-foreground">
            <p>映射缺失时，评论会以机器人身份写入。通过后台面板可以：</p>
            <ul className="space-y-1">
              <li>• 按 Plane 或 CNB 用户 ID 检索现有映射。</li>
              <li>• 批量录入待绑定的账号对，支持显示名备注。</li>
              <li>• 查看最近更新与连接时间，排查权限问题。</li>
            </ul>
          </CardContent>
          <CardFooter className="flex flex-col gap-3 text-sm text-muted-foreground md:flex-row md:items-center md:justify-between">
            <span>同步前请确认关键成员均已绑定，避免消息匿名化。</span>
            <Link href="/users" className={buttonVariants({ variant: 'secondary' })}>
              管理用户映射
            </Link>
          </CardFooter>
        </Card>
      </section>

      <section className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>运行与环境</CardTitle>
            <CardDescription>常用运维入口与本地开发说明。</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-muted-foreground">
            <div>
              <span className="font-medium text-foreground">后端健康检查：</span>
              <code className="ml-2">GET /healthz</code>
            </div>
            <div>
              <span className="font-medium text-foreground">前端环境变量：</span>
              <ul className="ml-4 mt-1 list-disc space-y-1">
                <li>
                  <code>NEXT_PUBLIC_API_BASE</code> — 浏览器直连后端；未配置时使用 Next API 代理。
                </li>
                <li>
                  <code>API_BASE</code> — 服务端 Route Handler 调用后端地址。
                </li>
              </ul>
            </div>
            <div>
              <span className="font-medium text-foreground">开发流程：</span>
              <ol className="ml-4 mt-1 list-decimal space-y-1">
                <li>
                  <code>npm ci</code> 安装依赖。
                </li>
                <li>
                  <code>API_BASE=http://localhost:8080 npm run dev</code> 启动面板。
                </li>
                <li>确保后端与数据库连接正常，刷新页面即可看到实时配置。</li>
              </ol>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>需要手动管理的配置</CardTitle>
            <CardDescription>所有影响同步策略的关键表均在后台提供可视化入口。</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-muted-foreground">
            <ul className="space-y-2">
              <li>
                <span className="font-medium text-foreground">repo_project_mappings</span> — 仓库 ↔ 项目绑定、同步方向、标签选择器。
              </li>
              <li>
                <span className="font-medium text-foreground">pr_state_mappings</span> — Pull Request 状态与 Plane Issue 状态对齐。
              </li>
              <li>
                <span className="font-medium text-foreground">label_mappings</span> — CNB 标签到 Plane 标签的映射表。
              </li>
              <li>
                <span className="font-medium text-foreground">issue_links</span> — Plane Issue ↔ CNB Issue 的人工纠错映射。
              </li>
              <li>
                <span className="font-medium text-foreground">thread_links</span> — 飞书线程 ↔ Plane Issue 的同步绑定。
              </li>
              <li>
                <span className="font-medium text-foreground">user_mappings</span> — Plane 成员与外部账号的身份绑定。
              </li>
            </ul>
            <p>如需扩展新的映射表，可在后端新增 API 并复用 Coss UI 组件快速构建管理界面。</p>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
