import { IconSettings } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Tabs } from "@/components/ui/tabs"
import { AccessPanel } from "@/features/settings/access-panel"
import { GroupsPanel } from "@/features/settings/groups-panel"
import { OptionsPanel } from "@/features/options/options-panel"

export type SettingsTab = "options" | "groups" | "role-access" | "user-access"

export interface SettingsPageProps {
  tab: SettingsTab
  onTabChange: (tab: SettingsTab) => void
  onCreateGroup: () => void
  onEditGroup: (uuid: string) => void
}

/** 系统设置页：平台参数、群集分组与访问授权。 */
export function SettingsPage({ tab, onTabChange, onCreateGroup, onEditGroup }: SettingsPageProps) {

  return (
    <>
      <PageHeader
        title="系统设置"
        description="平台参数、群集分组与访问授权"
        action={
          <span className="flex items-center gap-2 text-sm text-muted-foreground">
            <IconSettings size={16} stroke={1.5} />
            仅管理员可见
          </span>
        }
      />

      <Tabs
        className="px-1"
        value={tab}
        onChange={(key) => onTabChange(key as SettingsTab)}
        items={[
          { key: "options", label: "平台参数" },
          { key: "groups", label: "群集分组" },
          { key: "role-access", label: "角色授权" },
          { key: "user-access", label: "用户授权" },
        ]}
      />

      {tab === "options" ? (
        <SectionCard
          title="平台参数"
          description="只读项由服务端固定，可写项能在这里改值"
          bodyClassName="p-0"
        >
          <OptionsPanel showHeader={false} />
        </SectionCard>
      ) : null}

      {tab === "groups" ? (
        <GroupsPanel onCreate={onCreateGroup} onEdit={onEditGroup} />
      ) : null}

      {tab === "role-access" ? <AccessPanel scope="role" /> : null}
      {tab === "user-access" ? <AccessPanel scope="user" /> : null}
    </>
  )
}
