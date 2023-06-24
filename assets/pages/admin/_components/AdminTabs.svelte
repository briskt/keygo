<script lang="ts">
  import type { TabData } from 'data/types/tab'
  import { goto, isActive } from '@roxi/routify'
  import { TabBar } from '@silintl/ui-components'

  const Scroller = TabBar.Scroller
  const Tab = TabBar.Tab

  let activeTabIndex = 0
  let tabs: TabData[]

  $: tabs = [
    {
      label: 'Admin Home',
      tabUrl: '/admin',
      visible: true,
    },
    {
      label: 'Users',
      tabUrl: '/admin/users',
      visible: true,
    },
  ]
  $: tabs.forEach((tab, tabIndex) => {
    if ($isActive(tab.tabUrl, {}, true)) {
      activeTabIndex = tabIndex
    }
  })
</script>

<TabBar class="mb-1" tab={activeTabIndex}>
  <Scroller>
    {#each tabs as { label, tabUrl, visible } (tabUrl)}
      <Tab {label} on:click={() => $goto(tabUrl)} />
    {/each}
  </Scroller>
</TabBar>
