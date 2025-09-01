<template>
  <div class="admin-layout">
    <div class="admin-header">
      <h1>Administration Panel</h1>
      <nav class="admin-nav">
        <router-link to="/admin/events" class="admin-nav-link">Events</router-link>
        <router-link to="/admin/date-template-groups" class="admin-nav-link">Date Template Groups</router-link>
        <router-link to="/admin/date-templates" class="admin-nav-link">Date Templates</router-link>
      </nav>
    </div>
    
    <div class="admin-content">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminLayout'
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f8f9fa;
}

.admin-header {
  background: white;
  border-bottom: 1px solid #e2e8f0;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.admin-header h1 {
  margin: 0;
  color: #2d3748;
  font-size: 1.5rem;
}

.admin-nav {
  display: flex;
  gap: 1rem;
}

.admin-nav-link {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  transition: background 0.2s;
}

.admin-nav-link:hover {
  background: #5a67d8;
}

.admin-nav-link.router-link-active {
  background: #4c51bf;
}

.admin-content {
  padding: 2rem;
}
</style>