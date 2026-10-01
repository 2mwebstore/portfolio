<script setup>
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Autoplay, Pagination } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/pagination'

const props = defineProps({
  banners: { type: Array, default: () => [] },
  intervalMs: { type: Number, default: 4500 },
})
</script>

<template>
  <section class="hero">
    <Swiper
      v-if="banners.length"
      class="banner-swiper"
      :modules="[Autoplay, Pagination]"
      :loop="banners.length > 1"
      :autoplay="{ delay: intervalMs, disableOnInteraction: false, pauseOnMouseEnter: true }"
      :pagination="banners.length > 1 ? { clickable: true } : false"
    >
      <SwiperSlide v-for="b in banners" :key="b.id">
        <a
          class="banner-slide"
          :href="b.link || '#'"
          :target="b.link && b.link !== '#' ? '_blank' : undefined"
          rel="noopener"
        >
          <img :src="b.image_url" alt="" />
        </a>
      </SwiperSlide>
    </Swiper>

    <div v-else class="hero-banner">
      <div class="hero-fallback">
        <h2>Welcome</h2>
        <p>Add banners from the admin panel to fill this space.</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.banner-swiper {
  --swiper-pagination-color: #e8c369;
  --swiper-pagination-bullet-inactive-color: #ffffff;
  --swiper-pagination-bullet-inactive-opacity: 0.55;
  --swiper-pagination-bullet-size: 8px;
  --swiper-pagination-bullet-horizontal-gap: 4px;

  border-radius: 14px;
  overflow: hidden;
  border: 2px solid #e8c369;
  box-shadow: 0px 0px 10px #e8c369;
}
.banner-slide {
  display: block;
  line-height: 0;
}
.banner-slide img {
  width: 100%;
  max-height: 300px;
  height: 100%;
  display: block;
}

@media (max-width: 600px) {
  .banner-slide img { max-height: 140px; }
}
</style>