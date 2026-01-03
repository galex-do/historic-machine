<template>
  <div v-if="isOpen" class="event_detail_modal_overlay" @click.self="closeModal">
    <div class="event_detail_modal">
      <div class="event_detail_header">
        <div class="event_title_row">
          <span class="event_icon">{{ getEventEmoji(event.lens_type) }}</span>
          <a 
            v-if="event.source"
            :href="event.source" 
            target="_blank" 
            rel="noopener noreferrer"
            class="event_name_link"
          >{{ event.name }}</a>
          <span v-else class="event_name">{{ event.name }}</span>
        </div>
        <button class="close_btn" @click="closeModal" :title="t('close')">×</button>
      </div>

      <div class="event_detail_content">
        <div class="event_date">{{ formatEventDisplayDate(event.event_date, event.era) }}</div>
        <div v-if="event.description" class="event_description">{{ event.description }}</div>
        <div v-if="event.tags && event.tags.length > 0" class="event_tags">
          <span
            v-for="tag in event.tags"
            :key="tag.id"
            class="event_tag"
            :style="{ color: tag.color || '#6366f1' }"
            @click.stop="handleTagClick(tag)"
            :title="tag.description"
          >#{{ tag.name }}</span>
        </div>

        <div class="event_actions">
          <button class="focus_btn" @click="handleFocusEvent" :title="t('focusOnMap')">
            ⌖ {{ t('focusOnMap') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { watch, ref, onMounted, onUnmounted } from 'vue'
import { useLocale } from '@/composables/useLocale.js'
import { getEventEmoji } from '@/utils/event-utils.js'

export default {
  name: 'EventDetailModal',
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    event: {
      type: Object,
      default: () => ({})
    }
  },
  emits: ['close', 'focus-event', 'tag-clicked'],
  setup(props, { emit }) {
    const { t, formatEventDisplayDate } = useLocale()
    const previouslyFocusedElement = ref(null)

    const closeModal = () => {
      emit('close')
    }

    const handleFocusEvent = () => {
      emit('focus-event', props.event)
      closeModal()
    }

    const handleTagClick = (tag) => {
      emit('tag-clicked', tag)
      closeModal()
    }

    const handleKeydown = (e) => {
      if (e.key === 'Escape' && props.isOpen) {
        closeModal()
      }
    }

    watch(() => props.isOpen, (newVal) => {
      if (newVal) {
        previouslyFocusedElement.value = document.activeElement
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = ''
        if (previouslyFocusedElement.value) {
          previouslyFocusedElement.value.focus()
        }
      }
    })

    onMounted(() => {
      document.addEventListener('keydown', handleKeydown)
    })

    onUnmounted(() => {
      document.removeEventListener('keydown', handleKeydown)
      document.body.style.overflow = ''
    })

    return {
      t,
      formatEventDisplayDate,
      getEventEmoji,
      closeModal,
      handleFocusEvent,
      handleTagClick
    }
  }
}
</script>

<style scoped>
.event_detail_modal_overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.event_detail_modal {
  background: white;
  border-radius: 0.75rem;
  max-width: 500px;
  width: 100%;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.event_detail_header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #e2e8f0;
  gap: 1rem;
}

.event_title_row {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  flex: 1;
  min-width: 0;
}

.event_icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.event_name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
  word-break: break-word;
}

.event_name_link {
  font-size: 1.1rem;
  font-weight: 600;
  color: #3b82f6;
  text-decoration: none;
  word-break: break-word;
}

.event_name_link:hover {
  text-decoration: underline;
}

.close_btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #64748b;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  flex-shrink: 0;
}

.close_btn:hover {
  color: #334155;
}

.event_detail_content {
  padding: 1.25rem;
  overflow-y: auto;
}

.event_date {
  font-size: 0.8rem;
  font-weight: 600;
  color: #64748b;
  margin-bottom: 0.75rem;
}

.event_description {
  font-size: 0.9rem;
  line-height: 1.6;
  color: #334155;
  white-space: pre-wrap;
  word-break: break-word;
}

.event_tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.event_tag {
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.event_tag:hover {
  opacity: 0.7;
}

.event_actions {
  margin-top: 1.25rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
}

.focus_btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.focus_btn:hover {
  background: #e2e8f0;
  color: #3b82f6;
}
</style>
