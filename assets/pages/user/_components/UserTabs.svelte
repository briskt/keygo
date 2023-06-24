<script lang="ts">
  import { user } from 'data/store/user'
  import type { TabData } from 'data/types/tab'
  import { goto, isActive } from '@roxi/routify'
  import { TabBar } from '@silintl/ui-components'

  const Scroller = TabBar.Scroller
  const Tab = TabBar.Tab

  let activeTabIndex = 0
  let tabs: TabData[]

  $: tabs = [
    {
      label: 'View',
      tabUrl: '/user',
      visible: true,
    },
    {
      label: 'Edit',
      tabUrl: '/user/edit',
      visible: true,
    },
  ]
  $: tabs.forEach((tab, tabIndex) => {
    if ($isActive(tab.tabUrl, {}, true)) {
      activeTabIndex = tabIndex
    }
  })
</script>

{#if $user.ID}
  <TabBar class="mb-1" tab={activeTabIndex}>
    <Scroller>
      {#each tabs as { label, tabUrl, visible } (tabUrl)}
        <Tab {label} on:click={() => $goto(tabUrl)} />
      {/each}
    </Scroller>
  </TabBar>
{/if}
