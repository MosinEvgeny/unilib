<template>
  <div id="app">
    <nav class="reader-nav" v-if="showReaderNav">
      <router-link to="/reader" class="nav-link">Каталог книг</router-link>
      <router-link to="/reader/history" class="nav-link">История заказов</router-link>
      <button @click="logout" class="nav-link">Выход</button>
    </nav>

    <router-view />
  </div>
</template>

<script>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

export default {
  setup() {
    const route = useRoute();
    const router = useRouter();

    const showReaderNav = computed(() => route.path.startsWith('/reader'));

    const logout = () => {
      localStorage.removeItem('reader_id');
      localStorage.removeItem('role');
      router.push('/');
    };
    return {
      showReaderNav,
      logout
    };
  }
};
</script>

<style>
.reader-nav {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background-color: #a3d22d;
}

.nav-link {
  padding: 0.8rem 1.6rem;
  border-radius: 4px;
  background-color: #fbd604;
  color: #fff;
  text-decoration: none;
}

.nav-link:hover {
  background-color: darken(#fbd604, 10%);
}

.router-link-exact-active {
  background-color: #e6330f;
}

.reader-nav {
  justify-content: center
}
</style>