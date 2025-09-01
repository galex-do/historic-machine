<template>
  <div id="app">
    <!-- Header -->
    <header class="app-header">
      <div class="header-left">
        <router-link to="/" class="home-link">Historical Events Mapping</router-link>
      </div>
      <div class="header-right">
        <div class="admin-dropdown">
          <button @click="show_admin_dropdown = !show_admin_dropdown" class="admin-button">
            Admin ▼
          </button>
          <div v-show="show_admin_dropdown" class="dropdown-menu">
            <router-link to="/admin/events" class="dropdown-item">Events</router-link>
            <router-link to="/admin/date-template-groups" class="dropdown-item">Date Template Groups</router-link>
            <router-link to="/admin/date-templates" class="dropdown-item">Date Templates</router-link>
          </div>
        </div>
      </div>
    </header>
    
    <!-- Router View -->
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      show_admin_dropdown: false
    }
  },
  
  mounted() {
    // Close dropdown when clicking outside
    document.addEventListener('click', this.handle_outside_click)
  },
  
  beforeUnmount() {
    document.removeEventListener('click', this.handle_outside_click)
  },
  
  methods: {
    handle_outside_click(event) {
      const dropdown = event.target.closest('.admin-dropdown')
      if (!dropdown) {
        this.show_admin_dropdown = false
      }
    }
  }
}
</script>

<style>
/* Global Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: #f8f9fa;
}

#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header Styles */
.app-header {
  background: white;
  border-bottom: 1px solid #e2e8f0;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 80px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left .home-link {
  font-size: 1.5rem;
  font-weight: 600;
  color: #2d3748;
  text-decoration: none;
  transition: color 0.2s;
}

.header-left .home-link:hover {
  color: #667eea;
}

.admin-dropdown {
  position: relative;
}

.admin-button {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background 0.2s;
}

.admin-button:hover {
  background: #5a67d8;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 200px;
  z-index: 1000;
  margin-top: 0.5rem;
}

.dropdown-item {
  display: block;
  padding: 0.75rem 1rem;
  color: #4a5568;
  text-decoration: none;
  border-bottom: 1px solid #f7fafc;
  transition: background 0.2s;
}

.dropdown-item:hover {
  background: #f7fafc;
  color: #667eea;
}

.dropdown-item:last-child {
  border-bottom: none;
}
</style>