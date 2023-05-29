<script lang="ts">
import {getAuthStatus} from 'data/api/auth'
import {user} from 'data/store/user'
import { AppDrawer } from 'components'
import {isAdmin} from 'data/types/user'
import * as routes from 'helpers/routes'
import {afterUpdate} from 'svelte'

let userIsAnonymous: boolean

afterUpdate(async () => {
  const status = await getAuthStatus()
  userIsAnonymous = !status.IsAuthenticated
})

$: userIsAdmin = isAdmin($user)
$: userNotAdmin = !userIsAdmin || userIsAnonymous

$: menuItems = [
  {},

  // Admin menu items
  {
    url: routes.DASHBOARD,
    urlPattern: /\/home$/,
    icon: 'home',
    label: 'Home',
    hide: !userIsAdmin,
  },

  // Non admin menu items


  // Menu items for anonymous users
  {
    url: routes.ROOT,
    icon: 'person',
    label: 'Login',
    hide: !userIsAnonymous,
  },
  {
    url: routes.TERMS_OF_USE,
    icon: 'article',
    label: 'Terms of Service',
    hide: !userIsAnonymous,
  },
  {
    url: routes.PRIVACY_POLICY,
    icon: 'policy',
    label: 'Privacy Policy',
    hide: !userIsAnonymous,
  },
]
</script>

<AppDrawer {menuItems}>
  <slot />
</AppDrawer>
