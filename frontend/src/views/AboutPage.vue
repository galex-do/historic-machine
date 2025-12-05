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
          <li>{{ t('aboutUseCase5') }}</li>
        </ul>
      </section>

      <section class="open-source-section">
        <div class="open-source-notice">
          <span class="os-icon">📖</span>
          <div class="os-content">
            <p>{{ t('aboutOpenSource') }}</p>
            <div class="os-links">
              <a 
                href="https://github.com/galex-do/historic-machine" 
                target="_blank" 
                rel="noopener noreferrer"
                class="github-link"
              >
                <span class="github-icon">🔗</span>
                GitHub
              </a>
              <a 
                v-if="contactEmail"
                :href="`mailto:${contactEmail}`"
                class="contact-link"
              >
                <span class="contact-icon">✉️</span>
                {{ t('aboutContact') }}
              </a>
            </div>
          </div>
        </div>
      </section>

      <section v-if="supportCredentials.length > 0" class="support-section">
        <h2>{{ t('aboutSupportTitle') }}</h2>
        <p class="support-description">{{ t('aboutSupportDescription') }}</p>
        <div class="support-options">
          <div 
            v-for="credential in supportCredentials" 
            :key="credential.name"
            class="support-item"
          >
            <span class="support-name">{{ credential.name }}:</span>
            <a 
              v-if="credential.is_url" 
              :href="credential.value" 
              target="_blank" 
              rel="noopener noreferrer"
              class="support-link"
            >
              {{ t('aboutSupportLinkText') }}
            </a>
            <span v-else class="support-value">{{ credential.value }}</span>
          </div>
        </div>
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
import { ref, onMounted } from 'vue'
import { useLocale } from '@/composables/useLocale.js'

export default {
  name: 'AboutPage',
  setup() {
    const { t } = useLocale()
    const supportCredentials = ref([])
    const contactEmail = ref('')

    const fetchSupportCredentials = async () => {
      try {
        const response = await fetch('/api/support')
        if (response.ok) {
          const data = await response.json()
          supportCredentials.value = data || []
        }
      } catch (error) {
        console.error('Failed to fetch support credentials:', error)
      }
    }

    const fetchConfig = async () => {
      try {
        const response = await fetch('/api/config')
        if (response.ok) {
          const data = await response.json()
          contactEmail.value = data.contact_email || ''
        }
      } catch (error) {
        console.error('Failed to fetch config:', error)
      }
    }

    onMounted(() => {
      fetchSupportCredentials()
      fetchConfig()
    })

    return { t, supportCredentials, contactEmail }
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

.open-source-section {
  margin-bottom: 2rem;
}

.open-source-notice {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.25rem;
  background: linear-gradient(135deg, #f0f4ff 0%, #e8f4f8 100%);
  border-radius: 12px;
  border-left: 4px solid #667eea;
}

.os-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.os-content {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.os-content p {
  margin: 0;
  color: #4a5568;
}

.os-links {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.github-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #24292e;
  color: white;
  text-decoration: none;
  border-radius: 6px;
  font-weight: 500;
  font-size: 0.9rem;
  transition: background 0.2s ease;
}

.github-link:hover {
  background: #3a3f44;
}

.github-icon {
  font-size: 0.9rem;
}

.contact-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  border-radius: 6px;
  font-weight: 500;
  font-size: 0.9rem;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.contact-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.contact-icon {
  font-size: 0.9rem;
}

.support-section {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: #fffbeb;
  border-radius: 12px;
  border: 1px solid #fcd34d;
}

.support-section h2 {
  color: #92400e;
  margin-bottom: 0.75rem;
}

.support-description {
  color: #78716c;
  margin-bottom: 1rem;
  font-size: 0.95rem;
}

.support-options {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.support-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #fcd34d;
}

.support-name {
  font-weight: 600;
  color: #78350f;
  min-width: 120px;
}

.support-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.support-link:hover {
  text-decoration: underline;
}

.support-value {
  color: #4a5568;
  font-family: monospace;
  font-size: 0.9rem;
  word-break: break-all;
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
