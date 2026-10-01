<script setup>
import { ref, onMounted } from 'vue'
import BannerSwiper from './BannerSwiper.vue'

const config = ref(null)
const partners = ref([])
const bankIcons = ref([])
const socialIcons = ref([])
const banners = ref([])
const lang = ref('en')
const loading = ref(true)
const loadError = ref(false)

async function loadData() {
  try {
    const [cfgRes, partnersRes, iconsRes, bannersRes] = await Promise.all([
      fetch('/api/config'),
      fetch('/api/partners'),
      fetch('/api/payment-icons'),
      fetch('/api/banners')
    ])
    config.value = await cfgRes.json()
    partners.value = await partnersRes.json()
    const icons = await iconsRes.json()
    bankIcons.value = icons.filter(i => i.category === 'bank')
    socialIcons.value = icons.filter(i => i.category === 'social')
    banners.value = await bannersRes.json()
    applyBackground(config.value)
  } catch (e) {
    loadError.value = true
  } finally {
    loading.value = false
  }
}

// Background color/image is site-wide, so it's applied straight to
// <body> rather than a component inside .wrap.
function applyBackground(cfg) {
  if (!cfg) return
  if (cfg.background_color) {
    document.body.style.backgroundColor = cfg.background_color
  }
  if (cfg.background_image_url) {
    document.body.style.backgroundImage = `url('${cfg.background_image_url}')`
    document.body.style.backgroundSize = 'cover'
    document.body.style.backgroundPosition = 'center'
    document.body.style.backgroundAttachment = 'fixed'
    document.body.style.backgroundRepeat = 'no-repeat'
  } else {
    document.body.style.backgroundImage = ''
  }
}

onMounted(loadData)
</script>

<template>
  <div v-if="loading" class="loading-note">Loading directory…</div>

  <div v-else-if="loadError" class="loading-note">
    Couldn't reach the API. Make sure the Go backend is running on :8080.
  </div>

  <div v-else >
    <!-- Header -->
    <header class="site-header">
      <img class="site-logo" :src="config.logo_url" alt="" />
      <div class="site-name">{{ config.site_name }}</div>
    </header>

    <!-- Social row -->
    <div class="wrap">
      <div class="social-row">
        <a class="social-item" :href="config.telegram_url" target="_blank" rel="noopener">
          <div class="social-label">តេលេក្រាមផ្លូវការ</div>
          <div class="social-icon">
            <img v-if="config.telegram_icon_url" :src="config.telegram_icon_url" alt="" />
            <span v-else>✈</span>
          </div>
        </a>
        <a class="social-item" :href="config.facebook_url" target="_blank" rel="noopener">
          <div class="social-label">ហ្វេសប៊ុកផ្លូវការ</div>
          <div class="social-icon">
            <img v-if="config.facebook_icon_url" :src="config.facebook_icon_url" alt="" />
            <span v-else>f</span>
          </div>
        </a>
      </div>

      <!-- Hero banner swiper -->
      <BannerSwiper :banners="banners" />

      <!-- Partner grid -->
      <h2 class="section-title">ហ្គេមពេញនិយម</h2>
      <div class="partner-grid">
        <a
          v-for="p in partners"
          :key="p.id"
          class="partner-card"
          :href="p.link"
          target="_blank"
          rel="noopener"
        >
          <div class="partner-media">
            <img
              v-if="p.media_image_url"
              class="partner-media-img"
              :src="p.media_image_url"
              alt=""
            />
            <div v-else class="partner-media-img partner-media-fallback"></div>
            <img class="partner-logo" :src="p.logo_url" alt="" />
            <div class="partner-name">{{ p.name }}</div>
          </div>
          <div class="partner-detail">{{ p.detail }}</div>
        </a>
      </div>
    </div>

    <!-- Footer -->
    <footer>
      <div class="wrap">
        <div class="footer-logo-row">
          <img class="site-logo" :src="config.logo_url" alt="" />
          <div class="site-name">{{ config.site_name }}</div>
          <div class="footer-note-marquee" v-if="config.footer_note">
            <div class="footer-note-track">{{ config.footer_note }}</div>
          </div>
        </div>

        <div class="icon-groups">
          <div v-if="bankIcons.length">
            <div class="icon-group-label">ធនាគារ</div>
            <div class="icon-row">
              <div class="icon-chip" v-for="icon in bankIcons" :key="icon.id" :title="icon.name">
                <img :src="icon.icon_url" :alt="icon.name" />
              </div>
            </div>
          </div>
          <div v-if="socialIcons.length">
            <div class="icon-group-label">បណ្ដាញសង្គម</div>
            <div class="icon-row">
              <div class="icon-chip" v-for="icon in socialIcons" :key="icon.id" :title="icon.name">
                <img :src="icon.icon_url" :alt="icon.name" />
              </div>
            </div>
          </div>
        </div>

        <div class="copyright">
          © {{ new Date().getFullYear() }} {{ config.site_name }} រក្សាសិទ្ធិគ្រប់យ៉ាង | 18+
        </div>
      </div>
    </footer>
  </div>
</template>