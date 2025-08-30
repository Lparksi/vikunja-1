<template>
  <div class="not-found">
    <div class="not-found-content">
      <h1>404</h1>
      <h2>Page Not Found</h2>
      <p>The page you're looking for doesn't exist.</p>
      <t-button theme="primary" @click="$router.push('/')">
        Go Home
      </t-button>
    </div>
  </div>
</template>

<script setup lang="ts">
// No script needed for this simple component
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f5f5f5;
}

.not-found-content {
  text-align: center;
  padding: 40px;
}

.not-found-content h1 {
  font-size: 72px;
  font-weight: bold;
  color: #0052d9;
  margin-bottom: 16px;
}

.not-found-content h2 {
  font-size: 24px;
  color: #333;
  margin-bottom: 8px;
}

.not-found-content p {
  color: #666;
  margin-bottom: 24px;
}
</style>