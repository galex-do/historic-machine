<template>
  <div class="about-page">
    <div class="about-container">
      <h1>{{ t('aboutTitle') }}</h1>
      
      <section class="about-section">
        <p>{{ t('aboutProjectDescription') }}</p>
      </section>

      <section class="about-section">
        <h2>{{ t('aboutFeaturesTitle') }}</h2>
        <ul class="features-list">
          <li>{{ t('aboutFeature1') }}</li>
          <li>{{ t('aboutFeature2') }}</li>
          <li>{{ t('aboutFeature3') }}</li>
          <li>{{ t('aboutFeature4') }}</li>
          <li>{{ t('aboutFeature5') }}</li>
        </ul>
      </section>

      <section class="about-section">
        <h2>{{ t('aboutUseCasesTitle') }}</h2>
        <ul class="use-cases-list">
          <li>{{ t('aboutUseCase1') }}</li>
          <li>{{ t('aboutUseCase2') }}</li>
          <li>{{ t('aboutUseCase3') }}</li>
          <li>{{ t('aboutUseCase4') }}</li>
        </ul>
      </section>

      <div class="back-link">
        <router-link to="/" class="btn-back">
          <span>🗺️</span>
          {{ t('backToMap') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'AboutPage',
  setup() {
    const { t } = useLocale()
    return { t }
  }
}
</script>

<style scoped>
.about-page {
  min-height: calc(100vh - 80px);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 2rem;
}

.about-container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  border-radius: 16px;
  padding: 2.5rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

h1 {
  color: #2d3748;
  font-size: 2rem;
  margin-bottom: 2rem;
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.about-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.about-section:last-of-type {
  border-bottom: none;
}

h2 {
  color: #4a5568;
  font-size: 1.25rem;
  margin-bottom: 1rem;
}

p {
  color: #718096;
  line-height: 1.7;
  font-size: 1rem;
}

.features-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.features-list li {
  color: #718096;
  padding: 0.5rem 0;
  padding-left: 1.5rem;
  position: relative;
  line-height: 1.6;
}

.features-list li::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: #48bb78;
  font-weight: bold;
}

.use-cases-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.use-cases-list li {
  color: #718096;
  padding: 0.75rem 0;
  padding-left: 2rem;
  position: relative;
  line-height: 1.6;
  border-bottom: 1px solid #edf2f7;
}

.use-cases-list li:last-child {
  border-bottom: none;
}

.use-cases-list li::before {
  content: '💡';
  position: absolute;
  left: 0;
  font-size: 1rem;
}

.back-link {
  text-align: center;
  margin-top: 2rem;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-back:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

@media (max-width: 768px) {
  .about-page {
    padding: 1rem;
  }
  
  .about-container {
    padding: 1.5rem;
  }
  
  h1 {
    font-size: 1.5rem;
  }
  
  .tech-stack {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
