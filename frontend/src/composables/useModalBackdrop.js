import { ref, computed } from 'vue'

const open_modals = ref(new Set())

export function useModalBackdrop() {
  const register_modal = (id) => {
    const s = new Set(open_modals.value)
    s.add(id)
    open_modals.value = s
  }

  const unregister_modal = (id) => {
    const s = new Set(open_modals.value)
    s.delete(id)
    open_modals.value = s
  }

  const any_modal_open = computed(() => open_modals.value.size > 0)

  return { register_modal, unregister_modal, any_modal_open }
}
