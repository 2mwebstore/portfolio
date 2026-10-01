<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { clearSession, getUsername } from '../api'
import ConfigEditor from './ConfigEditor.vue'
import PartnersManager from './PartnersManager.vue'
import IconsManager from './IconsManager.vue'
import BannersManager from './BannersManager.vue'
import PasswordForm from './PasswordForm.vue'

const router = useRouter()
const tab = ref('config')

function logout() {
  clearSession()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="admin-shell">
    <header class="admin-header">
      <div class="admin-header-inner">
        <div class="admin-title">Directory Admin</div>
        <div class="admin-header-right">
          <span class="admin-user">{{ getUsername() }}</span>
          <a class="view-site-link" href="/" target="_blank">View site ↗</a>
          <button class="logout-btn" @click="logout">Log out</button>
        </div>
      </div>
    </header>

    <div class="admin-body">
      <nav class="admin-tabs">
        <button :class="{ active: tab === 'config' }" @click="tab = 'config'">Site Config</button>
        <button :class="{ active: tab === 'banners' }" @click="tab = 'banners'">Banners</button>
        <button :class="{ active: tab === 'partners' }" @click="tab = 'partners'">Partners</button>
        <button :class="{ active: tab === 'icons' }" @click="tab = 'icons'">Bank / Social Icons</button>
        <button :class="{ active: tab === 'password' }" @click="tab = 'password'">Password</button>
      </nav>

      <main class="admin-content">
        <ConfigEditor v-if="tab === 'config'" />
        <BannersManager v-else-if="tab === 'banners'" />
        <PartnersManager v-else-if="tab === 'partners'" />
        <IconsManager v-else-if="tab === 'icons'" />
        <PasswordForm v-else-if="tab === 'password'" />
      </main>
    </div>
  </div>
</template>

<style scoped>
.admin-shell{ min-height: 100vh; background: #f4f6f8; }
.admin-header{ background: #16294f; color: #fff; }
.admin-header-inner{
  max-width: 1040px; margin: 0 auto; padding: 16px 24px;
  display: flex; align-items: center; justify-content: space-between;
}
.admin-title{ font-weight: 700; font-size: 16px; }
.admin-header-right{ display:flex; align-items:center; gap: 16px; font-size: 13px; }
.admin-user{ color: rgba(255,255,255,0.7); }
.view-site-link{ color: #f2a768; }
.logout-btn{
  background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.25);
  color:#fff; padding: 6px 12px; border-radius: 6px; cursor:pointer; font-size: 13px;
}

.admin-body{ max-width: 1040px; margin: 0 auto; padding: 24px; display:flex; gap: 24px; }
.admin-tabs{ display:flex; flex-direction: column; gap: 6px; width: 190px; flex-shrink:0; }
.admin-tabs button{
  text-align:left; padding: 10px 14px; border-radius: 8px; border: 1px solid transparent;
  background: transparent; color: #16294f; font-size: 14px; cursor:pointer;
}
.admin-tabs button.active{ background: #fff; border-color: #e3e7ec; font-weight: 600; }
.admin-content{ flex:1; min-width: 0; }

@media (max-width: 680px){
  .admin-body{ flex-direction: column; }
  .admin-tabs{ width: 100%; flex-direction: row; flex-wrap: wrap; }
}
</style>
