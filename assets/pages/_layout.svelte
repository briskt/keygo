<script lang="ts">
import {authStatus} from 'data/api/auth'
import {user} from 'data/store/user'
import { AppDrawer } from 'components'
import {isAdmin} from 'data/types/user'
import * as routes from 'helpers/routes'

$: userIsAdmin = isAdmin($user)
$: userNotAdmin = !userIsAdmin || userIsAnonymous
$: userIsAnonymous = $authStatus.IsValid && !$authStatus.IsAuthenticated

$: menuItems = [
  // Admin menu items
  {
    url: routes.ADMIN,
    urlPattern: /\/admin$/,
    icon: 'admin_panel_settings',
    label: 'Admin',
    hide: !userIsAdmin,
  },

  // Non admin menu items
  {
    url: routes.PROFILE,
    urlPattern: /\/user$/,
    icon: 'account_circle',
    label: 'Profile',
    hide: userIsAnonymous,
  },


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
